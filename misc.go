package xpay

import (
	"context"
)

const (
	uriQueryPunishmentReasons = "/xpay/query_punishment_reasons"
	uriQueryOrderByIosOrder   = "/xpay/query_order_id_by_ios_order"
)

// ---------- 商户被管控原因查询 ----------

// RecoverySpecification 管控恢复说明。
type RecoverySpecification struct {
	LimitationCaseID         string `json:"limitation_case_id"`         // 管控单据号（对应管控通知 business_code）
	LimitationReasonType     string `json:"limitation_reason_type"`     // 管控原因类型
	LimitationReason         string `json:"limitation_reason"`          // 管控原因
	LimitationReasonDescribe string `json:"limitation_reason_describe"` // 管控原因描述
	RelateLimitations        string `json:"relate_limitations"`         // 受影响能力
	OtherRelateLimitations   string `json:"other_relate_limitations"`   // 其他受影响能力
	RecoverWay               string `json:"recover_way"`                // 恢复方式
	RecoverWayParam          string `json:"recover_way_param"`          // 恢复方式参数
	RecoverHelpURL           string `json:"recover_help_url"`           // 恢复帮助链接
	LimitationActionType     string `json:"limitation_action_type"`     // 管控动作类型
	LimitationStartDate      string `json:"limitation_start_date"`      // 管控开始日期
	LimitationDate           string `json:"limitation_date"`            // 管控期限
}

// PunishmentReasons 商户被管控原因。
type PunishmentReasons struct {
	baseResp
	AppID                  string                  `json:"appid"`                   // 小程序 appid
	Nickname               string                  `json:"nickname"`                // 小程序昵称
	MerchantCode           string                  `json:"merchant_code"`           // 商户号
	LimitedFunctions       []string                `json:"limited_functions"`       // 被管控能力列表
	OtherLimitedFunctions  string                  `json:"other_limited_functions"` // 其他被管控能力描述
	RecoverySpecifications []RecoverySpecification `json:"recovery_specifications"` // 恢复说明列表
}

// QueryPunishmentReasons 查询商户被微信支付管控原因。
func (c *Client) QueryPunishmentReasons(ctx context.Context) (*PunishmentReasons, error) {
	var out PunishmentReasons
	if err := c.post(ctx, uriQueryPunishmentReasons, struct{}{},
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 查询苹果订单对应单号 ----------

// IosOrderInfo 苹果订单对应的商户订单信息。
type IosOrderInfo struct {
	OutTradeNo            string `json:"out_trade_no"`            // 商家订单号
	TransactionID         string `json:"transaction_id"`          // 平台交易单号
	WebOrderLineItemID    string `json:"web_order_line_item_id"`  // 苹果订阅支付单号（可选）
	OriginalTransactionID string `json:"original_transaction_id"` // 原始平台交易单号
	PaidTime              string `json:"paid_time"`               // 支付时间，秒级时间戳字符串
	PaidFee               string `json:"paid_fee"`                // 实付金额，单位分
	PaidCurrency          string `json:"paid_currency"`           // 实付币种
}

// QueryOrderByIosOrderResp 查询苹果订单对应单号响应。
type QueryOrderByIosOrderResp struct {
	baseResp
	OrderInfoList []IosOrderInfo `json:"order_info_list"` // 订单信息列表
}

// QueryOrderByIosOrder 通过用户侧 Apple 支付凭证号（iosorderid）查询对应的商家交易单号及订单信息。
//
// 注意：受苹果服务限制，禁止轮询和大规模批量调用。
func (c *Client) QueryOrderByIosOrder(ctx context.Context, iosOrderID string) (*QueryOrderByIosOrderResp, error) {
	req := struct {
		IosOrderID string `json:"ios_order_id"`
	}{IosOrderID: iosOrderID}
	var out QueryOrderByIosOrderResp
	if err := c.post(ctx, uriQueryOrderByIosOrder, req,
		requestOptions{needPaySig: true}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
