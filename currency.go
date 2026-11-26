package xpay

import (
	"context"
)

// xpay 接口路径（不带 query，参与 pay_sig 签名）。
const (
	uriQueryUserBalance  = "/xpay/query_user_balance"
	uriCurrencyPay       = "/xpay/currency_pay"
	uriCancelCurrencyPay = "/xpay/cancel_currency_pay"
	uriPresentCurrency   = "/xpay/present_currency"
)

// ---------- 查询代币余额 ----------

// QueryUserBalanceReq 查询代币余额请求（需用户态签名）。
type QueryUserBalanceReq struct {
	Openid string `json:"openid"`  // 用户 openid
	UserIP string `json:"user_ip"` // 用户 ip，如 1.1.1.1
}

// UserBalance 用户代币余额信息。
type UserBalance struct {
	baseResp
	Balance        int64 `json:"balance"`         // 代币总余额（含赠送部分）
	PresentBalance int64 `json:"present_balance"` // 赠送账户的代币余额
	SumSave        int64 `json:"sum_save"`        // 累计有价货币充值数量
	SumPresent     int64 `json:"sum_present"`     // 累计赠送无价货币数量
	SumBalance     int64 `json:"sum_balance"`     // 历史总增加的代币金额
	SumCost        int64 `json:"sum_cost"`        // 历史总消耗代币金额
	FirstSaveFlag  bool  `json:"first_save_flag"` // 是否满足首充活动标记
}

// QueryUserBalance 查询用户代币余额。
//
// sessionKey 为用户当前有效 session_key（auth.code2Session 获取）。
func (c *Client) QueryUserBalance(ctx context.Context, sessionKey string, req *QueryUserBalanceReq) (*UserBalance, error) {
	var out UserBalance
	if err := c.post(ctx, uriQueryUserBalance, req,
		requestOptions{needPaySig: true, needUserSig: true, sessionKey: sessionKey}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 扣减代币 ----------

// CurrencyPayReq 扣减代币请求（代币支付，需用户态签名）。
type CurrencyPayReq struct {
	Openid  string `json:"openid"`             // 用户 openid
	UserIP  string `json:"user_ip"`            // 用户 ip
	Amount  int64  `json:"amount"`             // 支付的代币数量
	OrderID string `json:"order_id"`           // 订单号
	PayItem string `json:"pay_item,omitempty"` // 物品信息，记录到账户流水，如 [{"productid":"物品id","unit_price":1,"quantity":2}]
	Remark  string `json:"remark,omitempty"`   // 备注
}

// CurrencyPayResp 扣减代币响应。
type CurrencyPayResp struct {
	baseResp
	OrderID           string `json:"order_id"`            // 订单号
	Balance           int64  `json:"balance"`             // 支付后总余额
	UsedPresentAmount int64  `json:"used_present_amount"` // 本次使用的赠送代币数量
}

// CurrencyPay 扣减代币（一般用于代币支付）。
func (c *Client) CurrencyPay(ctx context.Context, sessionKey string, req *CurrencyPayReq) (*CurrencyPayResp, error) {
	var out CurrencyPayResp
	if err := c.post(ctx, uriCurrencyPay, req,
		requestOptions{needPaySig: true, needUserSig: true, sessionKey: sessionKey}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 代币支付退款 ----------

// CancelCurrencyPayReq 代币支付退款请求（currency_pay 的逆操作，需用户态签名）。
type CancelCurrencyPayReq struct {
	Openid     string `json:"openid"`       // 用户 openid
	UserIP     string `json:"user_ip"`      // 用户 ip
	PayOrderID string `json:"pay_order_id"` // 代币支付时传入的 order_id
	OrderID    string `json:"order_id"`     // 本次退款单的单号
	Amount     int64  `json:"amount"`       // 退款代币数量
}

// CancelCurrencyPayResp 代币支付退款响应。
type CancelCurrencyPayResp struct {
	baseResp
	OrderID string `json:"order_id"` // 退款订单号
}

// CancelCurrencyPay 代币支付退款。
func (c *Client) CancelCurrencyPay(ctx context.Context, sessionKey string, req *CancelCurrencyPayReq) (*CancelCurrencyPayResp, error) {
	var out CancelCurrencyPayResp
	if err := c.post(ctx, uriCancelCurrencyPay, req,
		requestOptions{needPaySig: true, needUserSig: true, sessionKey: sessionKey}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 代币赠送 ----------

// PresentCurrencyReq 代币赠送请求。
type PresentCurrencyReq struct {
	Openid  string `json:"openid"`   // 用户 openid
	OrderID string `json:"order_id"` // 赠送单号
	Amount  int64  `json:"amount"`   // 赠送代币数量
}

// PresentCurrencyResp 代币赠送响应。
type PresentCurrencyResp struct {
	baseResp
	Balance        int64  `json:"balance"`         // 赠送后用户的代币余额
	OrderID        string `json:"order_id"`        // 赠送单号
	PresentBalance int64  `json:"present_balance"` // 用户收到的总赠送金额
}

// PresentCurrency 代币赠送。
//
// 不支持按单号查赠送单，若返回 IsDuplicateOperation（268490004）表示之前已赠送成功，可视为成功。
func (c *Client) PresentCurrency(ctx context.Context, req *PresentCurrencyReq) (*PresentCurrencyResp, error) {
	var out PresentCurrencyResp
	if err := c.post(ctx, uriPresentCurrency, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
