package xpay

import (
	"context"
)

const (
	uriStartUploadGoods  = "/xpay/start_upload_goods"
	uriQueryUploadGoods  = "/xpay/query_upload_goods"
	uriStartPublishGoods = "/xpay/start_publish_goods"
	uriQueryPublishGoods = "/xpay/query_publish_goods"
)

// UploadItem 道具信息。
type UploadItem struct {
	ID      string `json:"id"`       // 道具 ID，≤20 字符，仅字母/数字/_/-
	Name    string `json:"name"`     // 道具名称，≤20 字符
	Price   int64  `json:"price"`    // 单价，单位分，>0
	Remark  string `json:"remark"`   // 备注，≤1024 字符
	ItemURL string `json:"item_url"` // 图片 url，支持 jpg/png
}

// ---------- 批量上传道具 ----------

// StartUploadGoodsReq 批量上传道具请求（一次仅支持 1 个道具）。
type StartUploadGoodsReq struct {
	UploadItem []UploadItem `json:"upload_item"` // 上传商品列表，长度 1
	Env        int          `json:"env"`         // 由 SDK 自动填充
}

// StartUploadGoods 启动批量上传道具任务，用 QueryUploadGoods 查询结果。
func (c *Client) StartUploadGoods(ctx context.Context, req *StartUploadGoodsReq) error {
	return c.post(ctx, uriStartUploadGoods, req, requestOptions{needPaySig: true}, nil)
}

// 道具上传状态枚举（QueryUploadGoodsResp.UploadItem[].UploadStatus）。
const (
	GoodsUploadProcessing = 0 // 上传中
	GoodsUploadExist      = 1 // id 已存在
	GoodsUploadSuccess    = 2 // 成功
	GoodsUploadFailed     = 3 // 失败
)

// UploadedItem 上传结果中的道具信息。
type UploadedItem struct {
	UploadItem
	ItemURL      string `json:"item_url"`      // 微信转存后的图片 url
	UploadStatus int    `json:"upload_status"` // 见 GoodsUpload*
	ErrMsg       string `json:"errmsg"`        // 失败原因
}

// QueryUploadGoodsResp 查询上传道具任务响应。
type QueryUploadGoodsResp struct {
	baseResp
	Status     int            `json:"status"`      // 0-无任务 1-运行中 2-失败或部分失败 3-成功
	UploadItem []UploadedItem `json:"upload_item"` // 道具列表
}

// QueryUploadGoods 查询批量上传道具任务结果。
func (c *Client) QueryUploadGoods(ctx context.Context) (*QueryUploadGoodsResp, error) {
	req := struct {
		Env int `json:"env"`
	}{}
	var out QueryUploadGoodsResp
	if err := c.post(ctx, uriQueryUploadGoods, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 启动批量发布道具 ----------

// PublishItem 发布道具信息。
type PublishItem struct {
	ID string `json:"id"` // 道具 ID（开发环境上传时的 id）
}

// StartPublishGoodsReq 批量发布道具请求（一次仅支持 1 个道具）。
type StartPublishGoodsReq struct {
	PublishItem []PublishItem `json:"publish_item"` // 发布列表，长度 1
	Env         int           `json:"env"`          // 由 SDK 自动填充
}

// StartPublishGoods 启动批量发布道具任务，用 QueryPublishGoods 查询结果。
func (c *Client) StartPublishGoods(ctx context.Context, req *StartPublishGoodsReq) error {
	return c.post(ctx, uriStartPublishGoods, req, requestOptions{needPaySig: true}, nil)
}

// PublishedItem 发布结果中的道具信息。
type PublishedItem struct {
	ID            string `json:"id"`             // 道具 ID
	PublishStatus int    `json:"publish_status"` // 见 GoodsUpload*
	ErrMsg        string `json:"errmsg"`         // 失败原因
}

// QueryPublishGoodsResp 查询发布道具任务响应。
type QueryPublishGoodsResp struct {
	baseResp
	Status      int             `json:"status"`       // 0-无任务 1-运行中 2-失败或部分失败 3-成功
	PublishItem []PublishedItem `json:"publish_item"` // 道具列表
}

// QueryPublishGoods 查询批量发布道具任务结果。
func (c *Client) QueryPublishGoods(ctx context.Context) (*QueryPublishGoodsResp, error) {
	req := struct {
		Env int `json:"env"`
	}{}
	var out QueryPublishGoodsResp
	if err := c.post(ctx, uriQueryPublishGoods, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
