package xpay

import (
	"context"
)

const (
	uriGetComplaintList      = "/xpay/get_complaint_list"
	uriGetComplaintDetail    = "/xpay/get_complaint_detail"
	uriGetNegotiationHistory = "/xpay/get_negotiation_history"
	uriResponseComplaint     = "/xpay/response_complaint"
	uriCompleteComplaint     = "/xpay/complete_complaint"
	uriUploadVpFile          = "/xpay/upload_vp_file"
	uriGetUploadFileSign     = "/xpay/get_upload_file_sign"
)

// 投诉状态枚举（Complaint.ComplaintState）。
const (
	ComplaintStatePending    = "PENDING"    // 待处理
	ComplaintStateProcessing = "PROCESSING" // 处理中
	ComplaintStateProcessed  = "PROCESSED"  // 已处理
)

// 投诉问题类型枚举（Complaint.ProblemType）。
const (
	ComplaintProblemRefund       = "REFUND"           // 退款
	ComplaintProblemNotAvailable = "SERVICE_NOT_WORK" // 服务无法使用
	ComplaintProblemOthers       = "OTHERS"           // 其他
)

// ComplaintMedia 投诉相关的媒体文件。
type ComplaintMedia struct {
	MediaType int      `json:"media_type"` // 媒体类型
	MediaURL  []string `json:"media_url"`  // 媒体 url 列表
}

// ComplaintOrderInfo 投诉涉及的订单信息。
type ComplaintOrderInfo struct {
	TransactionID string `json:"transaction_id"`   // 平台交易单号
	OutTradeNo    string `json:"out_trade_no"`     // 商户单号
	Amount        int64  `json:"amount"`           // 订单金额，单位分
	WxaOutTradeNo string `json:"wxa_out_trade_no"` // 小程序商户单号
	WxOrderID     string `json:"wx_order_id"`      // 微信内部单号
}

// ServiceOrderInfo 投诉涉及的服务订单信息。
type ServiceOrderInfo struct {
	OrderID    string `json:"order_id"`     // 订单 id
	OutOrderNo string `json:"out_order_no"` // 商户订单号
	State      string `json:"state"`        // DOING-进行中 REVOKED-已取消 WAITPAY-待支付 DONE-已完成
}

// Complaint 投诉信息。
type Complaint struct {
	ComplaintID           string               `json:"complaint_id"`            // 投诉 id
	ComplaintTime         int64                `json:"complaint_time"`          // 投诉时间，Unix 秒
	ComplaintDetail       string               `json:"complaint_detail"`        // 投诉内容
	ComplaintState        string               `json:"complaint_state"`         // 见 ComplaintState*
	PayerPhone            string               `json:"payer_phone"`             // 投诉人手机号（脱敏）
	PayerOpenid           string               `json:"payer_openid"`            // 投诉人 openid
	ComplaintOrderInfo    []ComplaintOrderInfo `json:"complaint_order_info"`    // 涉及订单
	ComplaintFullRefunded bool                 `json:"complaint_full_refunded"` // 是否全部全额退款
	IncomingUserResponse  bool                 `json:"incoming_user_response"`  // 是否有待回复留言
	UserComplaintTimes    int                  `json:"user_complaint_times"`    // 用户投诉次数
	ComplaintMediaList    []ComplaintMedia     `json:"complaint_media_list"`    // 投诉媒体文件
	ProblemDescription    string               `json:"problem_description"`     // 问题描述
	ProblemType           string               `json:"problem_type"`            // 见 ComplaintProblem*
	ApplyRefundAmount     int64                `json:"apply_refund_amount"`     // 申请退款金额，单位分
	UserTagList           []string             `json:"user_tag_list"`           // 用户标签，如 TRUSTED、HIGH_RISK
	ServiceOrderInfo      []ServiceOrderInfo   `json:"service_order_info"`      // 服务订单信息
}

// ---------- 获取投诉列表 ----------

// GetComplaintListReq 获取投诉列表请求。
type GetComplaintListReq struct {
	BeginDate string `json:"begin_date"` // 开始时间，yyyy-mm-dd
	EndDate   string `json:"end_date"`   // 结束时间，yyyy-mm-dd
	Offset    int    `json:"offset"`     // 偏移，从 0 开始
	Limit     int    `json:"limit"`      // 最多返回条数
	Env       int    `json:"env"`        // 由 SDK 自动填充
}

// GetComplaintListResp 获取投诉列表响应。
type GetComplaintListResp struct {
	baseResp
	Total      int         `json:"total"`      // 投诉总数
	Complaints []Complaint `json:"complaints"` // 投诉列表
}

// GetComplaintList 获取投诉列表。
func (c *Client) GetComplaintList(ctx context.Context, req *GetComplaintListReq) (*GetComplaintListResp, error) {
	var out GetComplaintListResp
	if err := c.post(ctx, uriGetComplaintList, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 获取投诉详情 ----------

// GetComplaintDetailReq 获取投诉详情请求。
type GetComplaintDetailReq struct {
	ComplaintID string `json:"complaint_id"` // 投诉 id
	Env         int    `json:"env"`          // 由 SDK 自动填充
}

// GetComplaintDetailResp 获取投诉详情响应。
type GetComplaintDetailResp struct {
	baseResp
	Complaint *Complaint `json:"complaint"` // 投诉详情
}

// GetComplaintDetail 获取投诉详情。
func (c *Client) GetComplaintDetail(ctx context.Context, req *GetComplaintDetailReq) (*GetComplaintDetailResp, error) {
	var out GetComplaintDetailResp
	if err := c.post(ctx, uriGetComplaintDetail, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 获取协商历史 ----------

// NegotiationHistory 投诉协商历史。
type NegotiationHistory struct {
	LogID              string           `json:"log_id"`               // 操作流水号
	Operator           string           `json:"operator"`             // 操作人
	OperateTime        string           `json:"operate_time"`         // 操作时间
	OperateType        string           `json:"operate_type"`         // 操作类型，如 USER_CREATE_COMPLAINT、MERCHANT_RESPONSE 等
	OperateDetails     string           `json:"operate_details"`      // 操作内容
	ComplaintMediaList []ComplaintMedia `json:"complaint_media_list"` // 媒体文件
}

// GetNegotiationHistoryReq 获取协商历史请求。
type GetNegotiationHistoryReq struct {
	ComplaintID string `json:"complaint_id"` // 投诉 id
	Offset      int    `json:"offset"`       // 偏移，从 0 开始
	Limit       int    `json:"limit"`        // 最多返回条数
	Env         int    `json:"env"`          // 由 SDK 自动填充
}

// GetNegotiationHistoryResp 获取协商历史响应。
type GetNegotiationHistoryResp struct {
	baseResp
	Total   int                  `json:"total"`   // 总条数
	History []NegotiationHistory `json:"history"` // 协商历史列表
}

// GetNegotiationHistory 获取投诉协商历史。
func (c *Client) GetNegotiationHistory(ctx context.Context, req *GetNegotiationHistoryReq) (*GetNegotiationHistoryResp, error) {
	var out GetNegotiationHistoryResp
	if err := c.post(ctx, uriGetNegotiationHistory, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 回复用户 ----------

// ResponseComplaintReq 回复用户请求。
type ResponseComplaintReq struct {
	ComplaintID     string   `json:"complaint_id"`     // 投诉 id
	ResponseContent string   `json:"response_content"` // 回复内容
	ResponseImages  []string `json:"response_images"`  // 图片 file_id 列表（UploadVpFile 返回）
	Env             int      `json:"env"`              // 由 SDK 自动填充
}

// ResponseComplaint 回复用户。
func (c *Client) ResponseComplaint(ctx context.Context, req *ResponseComplaintReq) error {
	return c.post(ctx, uriResponseComplaint, req, requestOptions{needPaySig: true}, nil)
}

// ---------- 完成投诉处理 ----------

// CompleteComplaintReq 完成投诉处理请求。
type CompleteComplaintReq struct {
	ComplaintID string `json:"complaint_id"` // 投诉 id
	Env         int    `json:"env"`          // 由 SDK 自动填充
}

// CompleteComplaint 完成投诉处理。
func (c *Client) CompleteComplaint(ctx context.Context, req *CompleteComplaintReq) error {
	return c.post(ctx, uriCompleteComplaint, req, requestOptions{needPaySig: true}, nil)
}

// ---------- 上传媒体文件 ----------

// UploadVpFileReq 上传媒体文件请求（投诉相关图片、凭证）。
type UploadVpFileReq struct {
	Base64Img string `json:"base64_img,omitempty"` // base64 图片，≤1M
	ImgURL    string `json:"img_url,omitempty"`    // 可直接下载（不能 302）的图片 url，≤2M，优先使用
	FileName  string `json:"file_name"`            // 图片名称
	Env       int    `json:"env"`                  // 由 SDK 自动填充
}

// UploadVpFileResp 上传媒体文件响应。
type UploadVpFileResp struct {
	baseResp
	FileID string `json:"file_id"` // 文件 id
}

// UploadVpFile 上传媒体文件（如图片，凭证等）。
func (c *Client) UploadVpFile(ctx context.Context, req *UploadVpFileReq) (*UploadVpFileResp, error) {
	var out UploadVpFileResp
	if err := c.post(ctx, uriUploadVpFile, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 获取投诉图片签名头部 ----------

// GetUploadFileSignReq 获取微信支付反馈投诉图片的签名头部请求。
type GetUploadFileSignReq struct {
	WxpayURL    string `json:"wxpay_url"`    // 微信支付图片地址，如 https://api.mch.weixin.qq.com/v3/merchant-service/images/{xxx}
	ConvertCOS  bool   `json:"convert_cos"`  // 是否转存 COS
	ComplaintID string `json:"complaint_id"` // 对应投诉 id
	Env         int    `json:"env"`          // 由 SDK 自动填充
}

// GetUploadFileSignResp 获取投诉图片签名头部响应。
type GetUploadFileSignResp struct {
	baseResp
	Sign   string `json:"sign"`    // 下载图片所需 Authorization 头值（还需 Accept: application/json 和非空 User-Agent）
	CosURL string `json:"cos_url"` // ConvertCOS=true 时返回转存 url，30 分钟有效
}

// GetUploadFileSign 获取微信支付反馈投诉图片的签名头部。
func (c *Client) GetUploadFileSign(ctx context.Context, req *GetUploadFileSignReq) (*GetUploadFileSignResp, error) {
	var out GetUploadFileSignResp
	if err := c.post(ctx, uriGetUploadFileSign, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
