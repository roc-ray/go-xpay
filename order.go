package xpay

import (
	"context"
)

const (
	uriQueryOrder         = "/xpay/query_order"
	uriRefundOrder        = "/xpay/refund_order"
	uriNotifyProvideGoods = "/xpay/notify_provide_goods"
)

// ---------- 查询创建的订单 ----------

// QueryOrderReq 查询订单请求（现金单，非代币单）。
type QueryOrderReq struct {
	Openid    string `json:"openid"`                // 用户 openid
	OrderID   string `json:"order_id,omitempty"`    // 商家订单号（与 WxOrderID 二选一）
	WxOrderID string `json:"wx_order_id,omitempty"` // 微信内部单号（与 OrderID 二选一）
}

// 订单状态枚举（Order.Status）。
const (
	OrderStatusInit            = 0  // 订单初始化（不可用于支付）
	OrderStatusCreated         = 1  // 订单创建成功
	OrderStatusPaid            = 2  // 已支付，待发货
	OrderStatusProviding       = 3  // 发货中
	OrderStatusProvided        = 4  // 已发货
	OrderStatusRefunded        = 5  // 退款：商户侧扣款完成
	OrderStatusClosed          = 6  // 已关闭
	OrderStatusRefundFailed    = 7  // 退款失败
	OrderStatusUserRefunded    = 8  // 退款：用户侧收款完成
	OrderStatusAdFundsRecycled = 9  // 回收广告金完成
	OrderStatusSplitReturned   = 10 // 分账回退完成
)

// 订单类型枚举（Order.OrderType）。
const (
	OrderTypeNormal    = 0 // 普通虚拟支付
	OrderTypeRefund    = 1 // 普通退款
	OrderTypeIosPay    = 7 // 苹果 iOS 支付
	OrderTypeIosRefund = 8 // 苹果 iOS 退款
)

// Order 虚拟支付订单信息。
type Order struct {
	OrderID        string `json:"order_id"`         // 商家订单号
	CreateTime     int64  `json:"create_time"`      // 创建时间，Unix 秒
	UpdateTime     int64  `json:"update_time"`      // 更新时间，Unix 秒
	Status         int    `json:"status"`           // 订单状态，见 OrderStatus*
	BizType        int    `json:"biz_type"`         // 业务类型，0-短剧
	OrderFee       int64  `json:"order_fee"`        // 订单金额，单位分
	CouponFee      int64  `json:"coupon_fee"`       // 优惠金额，单位分
	PaidFee        int64  `json:"paid_fee"`         // 用户支付金额，单位分
	OrderType      int    `json:"order_type"`       // 订单类型，见 OrderType*
	RefundFee      int64  `json:"refund_fee"`       // 退款单时表示退款金额，单位分
	PaidTime       int64  `json:"paid_time"`        // 支付/退款时间，Unix 秒
	ProvideTime    int64  `json:"provide_time"`     // 发货时间，Unix 秒
	BizMeta        string `json:"biz_meta"`         // 下单时商家自定义数据
	EnvType        int    `json:"env_type"`         // 环境类型，1-现网 2-沙箱
	Token          string `json:"token"`            // 下单时米大师返回的 token
	LeftFee        int64  `json:"left_fee"`         // 支付单经退款后剩余金额，单位分
	WxOrderID      string `json:"wx_order_id"`      // 微信内部单号
	ChannelOrderID string `json:"channel_order_id"` // 商户单号（用户微信支付详情页可见）
	WxpayOrderID   string `json:"wxpay_order_id"`   // 微信支付交易单号（用户微信支付详情页可见）
	SettTime       int64  `json:"sett_time"`        // 结算时间 Unix 秒，>0 表示结算成功
	SettState      int    `json:"sett_state"`       // 结算状态：0-未开始 1-结算中 2-成功 3-待结算
	PlatformFeeFen int64  `json:"platform_fee_fen"` // 虚拟支付技术服务费，分；结算成功后返回
	CpsFeeFen      int64  `json:"cps_fee_fen"`      // 公众号、视频号 CPS 服务费，分；结算成功后返回
}

// QueryOrderResp 查询订单响应。
type QueryOrderResp struct {
	baseResp
	Order *Order `json:"order"` // 订单信息
}

// QueryOrder 查询创建的订单（现金单，非代币单）。
func (c *Client) QueryOrder(ctx context.Context, req *QueryOrderReq) (*QueryOrderResp, error) {
	var out QueryOrderResp
	if err := c.post(ctx, uriQueryOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 启动订单退款任务 ----------

// 退款原因枚举（RefundOrderReq.RefundReason）。
const (
	RefundReasonNone         = "0" // 暂无描述
	RefundReasonProductIssue = "1" // 产品问题，影响使用或效果不佳
	RefundReasonAfterSale    = "2" // 售后问题，无法满足需求
	RefundReasonUserWish     = "3" // 意愿问题，用户主动退款
	RefundReasonPrice        = "4" // 价格问题
	RefundReasonOther        = "5" // 其他原因
)

// 退款来源枚举（RefundOrderReq.ReqFrom）。
const (
	RefundFromManualService = "1" // 人工客服退款
	RefundFromUser          = "2" // 用户自己发起
	RefundFromOther         = "3" // 其它
)

// RefundOrderReq 发起订单退款请求（可对支付时间 365 天以内的订单退款）。
type RefundOrderReq struct {
	Openid        string `json:"openid"`                // 下单时的用户 openid
	OrderID       string `json:"order_id,omitempty"`    // 下单时的商家单号（与 WxOrderID 二选一）
	WxOrderID     string `json:"wx_order_id,omitempty"` // 微信侧单号（与 OrderID 二选一）
	RefundOrderID string `json:"refund_order_id"`       // 本次退款单号，8-32 位，仅字母/数字/_/-
	LeftFee       int64  `json:"left_fee"`              // 当前单剩余可退金额，分，通过 QueryOrder 查询
	RefundFee     int64  `json:"refund_fee"`            // 本次退款金额，分，范围 (0, LeftFee]
	BizMeta       string `json:"biz_meta"`              // 商家自定义数据，0-1024 字符，查询时原样返回
	RefundReason  string `json:"refund_reason"`         // 退款原因，见 RefundReason*
	ReqFrom       string `json:"req_from"`              // 退款来源，见 RefundFrom*
	Env           int    `json:"env"`                   // 由 SDK 自动填充
}

// RefundOrderResp 发起退款响应。
type RefundOrderResp struct {
	baseResp
	RefundOrderID   string `json:"refund_order_id"`    // 退款单号
	RefundWxOrderID string `json:"refund_wx_order_id"` // 退款单的微信侧单号
	PayOrderID      string `json:"pay_order_id"`       // 对应的支付单商家单号
	PayWxOrderID    string `json:"pay_wx_order_id"`    // 对应的支付单微信侧单号
}

// RefundOrder 启动订单退款任务。
//
// 仅启动退款任务，需调用 QueryOrder 轮询退款单状态至退款完成（status=5/8）。
func (c *Client) RefundOrder(ctx context.Context, req *RefundOrderReq) (*RefundOrderResp, error) {
	var out RefundOrderResp
	if err := c.post(ctx, uriRefundOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 通知已发货完成 ----------

// NotifyProvideGoodsReq 通知发货请求（仅现金单）。
type NotifyProvideGoodsReq struct {
	OrderID   string `json:"order_id"`    // 下单时传的商家单号
	WxOrderID string `json:"wx_order_id"` // 微信内部单号（与 OrderID 二选一）
	Env       int    `json:"env"`         // 由 SDK 自动填充
}

// NotifyProvideGoods 通知已发货完成。
//
// 正常情况下通过 xpay_goods_deliver_notify 消息推送返回成功即可，本接口用于推送不成功时手动置为已发货。
func (c *Client) NotifyProvideGoods(ctx context.Context, req *NotifyProvideGoodsReq) error {
	return c.post(ctx, uriNotifyProvideGoods, req, requestOptions{}, nil)
}
