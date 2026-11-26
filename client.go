package xpay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

// 环境常量，对应各接口请求体中的 env 字段。
const (
	EnvProduction = 0 // 正式环境
	EnvSandbox    = 1 // 沙箱环境
)

const apiBaseURL = "https://api.weixin.qq.com"

// Config SDK 配置。
type Config struct {
	AppKey     string       // 正式环境的虚拟支付 AppKey（小程序后台「虚拟支付-支付基础配置」）
	SandboxKey string       // 沙箱环境的虚拟支付 AppKey，仅 Env=EnvSandbox 时必填
	Env        int          // 默认环境：EnvProduction(0) 正式、EnvSandbox(1) 沙箱，默认正式
	HTTPClient *http.Client // 自定义 HTTP 客户端（可选），默认 15 秒超时

	// TokenProvider 外部 access_token 获取函数（必填）。
	// SDK 不负责获取 access_token（如小程序 AppSecret 换 stable_token），
	// 统一由外部传入（如公司内部统一 token 服务），SDK 每次请求实时获取，不做缓存。
	TokenProvider func(ctx context.Context) (string, error)
}

// Client 微信小程序虚拟支付（xpay）客户端。
//
// Client 并发安全；access_token 每次请求时由外部 TokenProvider 实时获取，
// 各接口请求体中的 env 字段由 SDK 按 Config.Env 自动填充，无需调用方设置。
type Client struct {
	cfg  Config
	http *http.Client
}

// New 创建虚拟支付客户端。
func New(cfg Config) (*Client, error) {
	if cfg.TokenProvider == nil {
		return nil, errors.New("xpay: TokenProvider 不能为空（access_token 统一由外部传入）")
	}
	if cfg.AppKey == "" {
		return nil, errors.New("xpay: AppKey 不能为空（获取自小程序后台「虚拟支付-支付基础配置」）")
	}
	if cfg.Env != EnvProduction && cfg.Env != EnvSandbox {
		return nil, fmt.Errorf("xpay: 非法 Env 值 %d（仅支持 0-正式 / 1-沙箱）", cfg.Env)
	}
	if cfg.Env == EnvSandbox && cfg.SandboxKey == "" {
		return nil, errors.New("xpay: 沙箱环境必须配置 SandboxKey")
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{cfg: cfg, http: cfg.HTTPClient}, nil
}

// appKey 返回与当前环境匹配的 AppKey。
func (c *Client) appKey() string {
	if c.cfg.Env == EnvSandbox {
		return c.cfg.SandboxKey
	}
	return c.cfg.AppKey
}

// requestOptions 单次请求的签名选项。
type requestOptions struct {
	needPaySig  bool   // 是否附加支付签名 pay_sig
	needUserSig bool   // 是否附加用户态签名 signature
	sessionKey  string // 用户态签名密钥
}

// post 发送 JSON POST 请求并解析响应。
//
// body 中的 Env 字段（若存在）会先被自动填充；参与 pay_sig 签名的
// 请求体与实际发送的请求体保证完全一致。
func (c *Client) post(ctx context.Context, uri string, body any, opt requestOptions, out any) error {
	injectEnv(body, c.cfg.Env)

	raw, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("xpay: 序列化请求体失败: %w", err)
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	q := uri + "?access_token=" + token
	if opt.needPaySig {
		q += "&pay_sig=" + PaySig(c.appKey(), uri, string(raw))
	}
	if opt.needUserSig {
		if opt.sessionKey == "" {
			return errors.New("xpay: 该接口需要用户态签名，sessionKey 不能为空")
		}
		q += "&signature=" + UserSignature(opt.sessionKey, string(raw))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiBaseURL+q, bytesReader(raw))
	if err != nil {
		return fmt.Errorf("xpay: 构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("xpay: 请求 %s 失败: %w", uri, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("xpay: 读取响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("xpay: %s http_status=%d body=%s", uri, resp.StatusCode, truncate(data, 512))
	}

	// 先解出公共错误码
	var base baseResp
	if err := json.Unmarshal(data, &base); err != nil {
		return fmt.Errorf("xpay: 解析响应失败: %w body=%s", err, truncate(data, 512))
	}
	if base.ErrCode != 0 {
		return &APIError{Code: base.ErrCode, Msg: base.ErrMsg, URI: uri}
	}

	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("xpay: 解析响应失败: %w body=%s", err, truncate(data, 512))
		}
	}
	return nil
}

// injectEnv 若 body 是含 Env int 字段的结构体指针，则自动填充环境值。
func injectEnv(body any, env int) {
	if body == nil {
		return
	}
	v := reflect.ValueOf(body)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	f := v.FieldByName("Env")
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.Int {
		f.SetInt(int64(env))
	}
}

// getAccessToken 通过外部 TokenProvider 实时获取 access_token（SDK 不做缓存）。
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	token, err := c.cfg.TokenProvider(ctx)
	if err != nil {
		return "", fmt.Errorf("xpay: TokenProvider 获取 access_token 失败: %w", err)
	}
	if token == "" {
		return "", errors.New("xpay: TokenProvider 返回的 access_token 为空")
	}
	return token, nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "...(truncated)"
}
