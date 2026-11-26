package xpay

import (
	"context"
)

const (
	uriQuerySubscribeContract  = "/xpay/query_subscribe_contract"
	uriSendSubscribePrePayment = "/xpay/send_subscribe_pre_payment"
	uriSubmitSubscribePayOrder = "/xpay/submit_subscribe_pay_order"
	uriCancelSubscribeContract = "/xpay/cancel_subscribe_contract"
)

// 签约状态枚举（QuerySubscribeContractResp.AuthorizationState）。
const (
	SubscribeStateSigned     = "SIGNED"     // 签约生效中
	SubscribeStateTerminated = "TERMINATED" // 已解约（终态）
	SubscribeStateUnbind     = "UNBINDUSER" // 从未签约
)

// ---------- 查询签约关系 ----------

// QuerySubscribeContractReq 查询签约关系请求。
type QuerySubscribeContractReq struct {
	Openid          string `json:"openid"`            // 用户 openid
	ProductID       string `json:"product_id"`        // 道具 id，须为订阅制道具
	OutContractCode string `json:"out_contract_code"` // 签约协议号
}

// QuerySubscribeContractResp 查询签约关系响应。
type QuerySubscribeContractResp struct {
	baseResp
	AuthorizationState string `json:"authorization_state"` // 见 SubscribeState*
}

// QuerySubscribeContract 查询签约关系。
func (c *Client) QuerySubscribeContract(ctx context.Context, req *QuerySubscribeContractReq) (*QuerySubscribeContractResp, error) {
	var out QuerySubscribeContractResp
	if err := c.post(ctx, uriQuerySubscribeContract, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 预通知扣款 ----------

// SendSubscribePrePaymentReq 预通知扣款请求。
type SendSubscribePrePaymentReq struct {
	Openid          string `json:"openid"`            // 用户 openid
	DeductPrice     int64  `json:"deduct_price"`      // 扣款金额，单位分，范围 [100, 道具价格]
	ProductID       string `json:"product_id"`        // 订阅制道具 id
	OutContractCode string `json:"out_contract_code"` // 签约协议号
}

// SendSubscribePrePayment 预通知扣款。
func (c *Client) SendSubscribePrePayment(ctx context.Context, req *SendSubscribePrePaymentReq) error {
	return c.post(ctx, uriSendSubscribePrePayment, req, requestOptions{needPaySig: true}, nil)
}

// ---------- 发起订阅扣款 ----------

// SubmitSubscribePayOrderReq 发起订阅扣款请求。
type SubmitSubscribePayOrderReq struct {
	Openid       string `json:"openid"`        // 用户 openid
	OfferID      string `json:"offer_id"`      // 米大师应用 id（小程序后台「虚拟支付-支付基础配置」）
	BuyQuantity  int64  `json:"buy_quantity"`  // 填 1
	CurrencyType string `json:"currency_type"` // 填 CNY
	ProductID    string `json:"product_id"`    // 订阅道具 ID
	DeductPrice  int64  `json:"deduct_price"`  // 扣款金额，单位分，范围 [1, 道具价格]
	OrderID      string `json:"order_id"`      // 业务订单号，8-32 字符，仅数字/字母/_/-，只能使用一次
	Attach       string `json:"attach"`        // 透传数据，发货通知时原样透传给开发者
	// 注意：该接口 env 仅支持 0（正式环境），沙箱环境不可用。
}

// SubmitSubscribePayOrder 发起订阅扣款。
//
// 成功仅代表受理，扣款结果经 xpay_goods_deliver_notify 消息通知；失败可在 t～t+6 重试，
// 同用户同道具每小时限 1 次。
func (c *Client) SubmitSubscribePayOrder(ctx context.Context, req *SubmitSubscribePayOrderReq) error {
	return c.post(ctx, uriSubmitSubscribePayOrder, req, requestOptions{needPaySig: true}, nil)
}

// ---------- 商家解约 ----------

// CancelSubscribeContractReq 商家解约请求。
type CancelSubscribeContractReq struct {
	Openid            string `json:"openid"`             // 用户 openid
	TerminationReason string `json:"termination_reason"` // 解约原因
	ProductID         string `json:"product_id"`         // 订阅制道具 id
	OutContractCode   string `json:"out_contract_code"`  // 签约协议号
}

// CancelSubscribeContract 商家解约。
func (c *Client) CancelSubscribeContract(ctx context.Context, req *CancelSubscribeContractReq) error {
	return c.post(ctx, uriCancelSubscribeContract, req, requestOptions{needPaySig: true}, nil)
}
