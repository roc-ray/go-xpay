package xpay

import (
	"bytes"
	"errors"
	"fmt"
)

// baseResp 微信接口公共响应字段，各响应结构体内嵌以平铺 JSON 字段。
type baseResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// APIError 微信接口返回的业务错误（errcode != 0 时产生）。
type APIError struct {
	Code int    `json:"errcode"`
	Msg  string `json:"errmsg"`
	URI  string `json:"-"` // 出错的接口路径
}

func (e *APIError) Error() string {
	return fmt.Sprintf("xpay: %s errcode=%d errmsg=%s", e.URI, e.Code, e.Msg)
}

// 虚拟支付常用错误码（完整列表见官方文档）。
const (
	CodeSystemError          = -1        // 系统错误
	CodeInvalidOpenid        = 268490001 // openid 错误
	CodeInvalidParam         = 268490002 // 请求参数字段错误
	CodeSignatureError       = 268490003 // 签名错误
	CodeDuplicateOperation   = 268490004 // 重复操作（之前的操作已成功）
	CodeAlreadyRefunded      = 268490005 // 订单已退款，不支持再退款
	CodeInsufficientBalance  = 268490006 // 代币余额不足
	CodeSensitiveContent     = 268490007 // 图片或文字存在敏感内容
	CodeCurrencyNotPublished = 268490008 // 代币未发布
	CodeSessionKeyExpired    = 268490009 // session_key 不存在或已过期
	CodeDataGenerating       = 268490011 // 数据生成中，稍后重试
	CodeBatchTaskRunning     = 268490012 // 批量任务运行中
	CodeRateLimited          = 268490015 // 频率限制
	CodeLeftFeeMismatch      = 268490016 // left_fee 与实际不符，请通过 QueryOrder 确认
)

// AsAPIError 将 err 断言为 *APIError，断言失败时 ok 为 false。
func AsAPIError(err error) (apiErr *APIError, ok bool) {
	ok = errors.As(err, &apiErr)
	return
}

// IsDuplicateOperation 判断是否为「重复操作」（268490004），
// 赠送/代币支付/广告金充值接口返回该错误表示之前的操作已成功，可视为成功。
func IsDuplicateOperation(err error) bool {
	e, ok := AsAPIError(err)
	return ok && e.Code == CodeDuplicateOperation
}

// IsSessionKeyExpired 判断用户 session_key 是否已过期（268490009），
// 返回 true 时应引导用户重新登录获取新 session_key。
func IsSessionKeyExpired(err error) bool {
	e, ok := AsAPIError(err)
	return ok && e.Code == CodeSessionKeyExpired
}

// IsInsufficientBalance 判断代币余额是否不足（268490006）。
func IsInsufficientBalance(err error) bool {
	e, ok := AsAPIError(err)
	return ok && e.Code == CodeInsufficientBalance
}

// IsSignatureInvalid 判断签名是否错误（268490003）。
func IsSignatureInvalid(err error) bool {
	e, ok := AsAPIError(err)
	return ok && e.Code == CodeSignatureError
}

func bytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
