package xpay

import (
	"context"
)

// 广告金相关接口均不需要支付签名。
const (
	uriQueryTransferAccount = "/xpay/query_transfer_account"
	uriQueryAdverFunds      = "/xpay/query_adver_funds"
	uriCreateFundsBill      = "/xpay/create_funds_bill"
	uriBindTransferAccount  = "/xpay/bind_transfer_accout"
	uriQueryFundsBill       = "/xpay/query_funds_bill"
	uriQueryRecoverBill     = "/xpay/query_recover_bill"
)

// 广告金类型枚举。
const (
	FundTypeGeneral  = 0 // 通用赠送
	FundTypeAd       = 1 // 广告激励
	FundTypeDirected = 2 // 定向激励
)

// TransferAccount 广告金充值账户。
type TransferAccount struct {
	TransferAccountName       string `json:"transfer_account_name"`        // 充值账户名称
	TransferAccountUID        int64  `json:"transfer_account_uid"`         // 充值账户 uid
	TransferAccountAgencyID   int64  `json:"transfer_account_agency_id"`   // 服务商账号 id
	TransferAccountAgencyName string `json:"transfer_account_agency_name"` // 服务商账号名称
	State                     int    `json:"state"`                        // 0-待审核 1-通过 2-驳回
	BindResult                int    `json:"bind_result"`                  // 1-绑定成功 2-绑定失败
	ErrorMsg                  string `json:"error_msg"`                    // 错误信息
}

// QueryTransferAccountResp 查询广告金充值账户响应。
type QueryTransferAccountResp struct {
	baseResp
	AcctList []TransferAccount `json:"acct_list"` // 充值账户列表
}

// QueryTransferAccount 查询广告金充值账户。
func (c *Client) QueryTransferAccount(ctx context.Context) (*QueryTransferAccountResp, error) {
	req := struct {
		Env int `json:"env"`
	}{}
	var out QueryTransferAccountResp
	if err := c.post(ctx, uriQueryTransferAccount, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 查询广告金发放记录 ----------

// AdverFundsFilter 广告金发放记录筛选条件。
type AdverFundsFilter struct {
	SettleBegin int64 `json:"settle_begin,omitempty"` // 结算周期开始，Unix 秒
	SettleEnd   int64 `json:"settle_end,omitempty"`   // 结算周期结束，Unix 秒
	FundType    int   `json:"fund_type,omitempty"`    // 见 FundType*
}

// QueryAdverFundsReq 查询广告金发放记录请求。
type QueryAdverFundsReq struct {
	Page     int               `json:"page,omitempty"`      // 页码，≥1
	PageSize int               `json:"page_size,omitempty"` // 每页条数
	Filter   *AdverFundsFilter `json:"filter,omitempty"`
	Env      int               `json:"env"` // 由 SDK 自动填充
}

// AdverFund 广告金发放记录。
type AdverFund struct {
	SettleBegin  int64  `json:"settle_begin"`  // 结算周期开始，Unix 秒
	SettleEnd    int64  `json:"settle_end"`    // 结算周期结束，Unix 秒
	TotalAmount  int64  `json:"total_amount"`  // 发放金额，单位分
	RemainAmount int64  `json:"remain_amount"` // 剩余可用金额，单位分
	ExpireTime   int64  `json:"expire_time"`   // 过期时间，Unix 秒
	FundType     int    `json:"fund_type"`     // 见 FundType*
	FundID       string `json:"fund_id"`       // 发放 ID
}

// QueryAdverFundsResp 查询广告金发放记录响应。
type QueryAdverFundsResp struct {
	baseResp
	TotalPage      int         `json:"total_page"`       // 总页数
	AdverFundsList []AdverFund `json:"adver_funds_list"` // 发放记录列表
}

// QueryAdverFunds 查询广告金发放记录。
func (c *Client) QueryAdverFunds(ctx context.Context, req *QueryAdverFundsReq) (*QueryAdverFundsResp, error) {
	var out QueryAdverFundsResp
	if err := c.post(ctx, uriQueryAdverFunds, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 充值广告金 ----------

// CreateFundsBillReq 充值广告金请求。
type CreateFundsBillReq struct {
	TransferAmount          int64  `json:"transfer_amount"`            // 充值金额，单位分
	TransferAccountUID      int64  `json:"transfer_account_uid"`       // 充值账户 uid
	TransferAccountName     string `json:"transfer_account_name"`      // 充值账户名称
	TransferAccountAgencyID int64  `json:"transfer_account_agency_id"` // 服务商账号 id
	RequestID               string `json:"request_id"`                 // 唯一请求 id，≤1024 字符，重复视为重复请求
	SettleBegin             int64  `json:"settle_begin"`               // 所用广告金结算周期开始，Unix 秒
	SettleEnd               int64  `json:"settle_end"`                 // 结算周期结束，Unix 秒
	AuthorizeAdvertise      int    `json:"authorize_advertise"`        // 0-否 1-是
	FundType                int    `json:"fund_type"`                  // 见 FundType*
	Env                     int    `json:"env"`                        // 由 SDK 自动填充
}

// CreateFundsBillResp 充值广告金响应。
type CreateFundsBillResp struct {
	baseResp
	BillID string `json:"bill_id"` // 充值单 id
}

// CreateFundsBill 充值广告金。返回 IsDuplicateOperation 表示之前的充值已成功。
func (c *Client) CreateFundsBill(ctx context.Context, req *CreateFundsBillReq) (*CreateFundsBillResp, error) {
	var out CreateFundsBillResp
	if err := c.post(ctx, uriCreateFundsBill, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 绑定广告金充值账户 ----------

// BindTransferAccountReq 绑定广告金充值账户请求。
type BindTransferAccountReq struct {
	TransferAccountUID     int64  `json:"transfer_account_uid,omitempty"`      // 充值账户 uid
	TransferAccountOrgName string `json:"transfer_account_org_name,omitempty"` // 充值账户主体名称
	Env                    int    `json:"env,omitempty"`                       // 由 SDK 自动填充
}

// BindTransferAccount 绑定广告金充值账户。
func (c *Client) BindTransferAccount(ctx context.Context, req *BindTransferAccountReq) error {
	return c.post(ctx, uriBindTransferAccount, req, requestOptions{}, nil)
}

// ---------- 查询广告金充值记录 ----------

// FundsBillFilter 广告金充值记录筛选条件。
type FundsBillFilter struct {
	OperTimeBegin int64  `json:"oper_time_begin"`      // 充值开始时间，Unix 秒
	OperTimeEnd   int64  `json:"oper_time_end"`        // 充值结束时间，Unix 秒
	BillID        string `json:"bill_id,omitempty"`    // 充值单 ID
	RequestID     string `json:"request_id,omitempty"` // CreateFundsBill 传入的 request_id
}

// QueryFundsBillReq 查询广告金充值记录请求。
type QueryFundsBillReq struct {
	Page     int              `json:"page"`      // 页码，≥1
	PageSize int              `json:"page_size"` // 每页条数
	Filter   *FundsBillFilter `json:"filter"`
	Env      int              `json:"env"` // 由 SDK 自动填充
}

// FundsBill 广告金充值记录。
type FundsBill struct {
	BillID              string `json:"bill_id"`               // 充值单 ID
	OperTime            int64  `json:"oper_time"`             // 充值时间，Unix 秒
	SettleBegin         int64  `json:"settle_begin"`          // 对应结算周期开始，Unix 秒
	SettleEnd           int64  `json:"settle_end"`            // 对应结算周期结束，Unix 秒
	FundID              string `json:"fund_id"`               // 对应广告金 ID
	TransferAccountName string `json:"transfer_account_name"` // 充值账户名称
	TransferAccountUID  int64  `json:"transfer_account_uid"`  // 充值账户 uid
	TransferAmount      int64  `json:"transfer_amount"`       // 充值金额，单位分
	Status              int    `json:"status"`                // 0-充值中 1-成功 2-失败
	RequestID           string `json:"request_id"`            // 原请求 id
}

// QueryFundsBillResp 查询广告金充值记录响应。
type QueryFundsBillResp struct {
	baseResp
	TotalPage int         `json:"total_page"` // 总页数
	BillList  []FundsBill `json:"bill_list"`  // 充值记录列表
}

// QueryFundsBill 查询广告金充值记录。
func (c *Client) QueryFundsBill(ctx context.Context, req *QueryFundsBillReq) (*QueryFundsBillResp, error) {
	var out QueryFundsBillResp
	if err := c.post(ctx, uriQueryFundsBill, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 查询广告金回收记录 ----------

// RecoverBillFilter 广告金回收记录筛选条件。
type RecoverBillFilter struct {
	RecoverTimeBegin int64  `json:"recover_time_begin"` // 回收开始时间，Unix 秒
	RecoverTimeEnd   int64  `json:"recover_time_end"`   // 回收结束时间，Unix 秒
	BillID           string `json:"bill_id,omitempty"`  // 回收单 ID
}

// QueryRecoverBillReq 查询广告金回收记录请求。
type QueryRecoverBillReq struct {
	Page     int                `json:"page"`      // 页码，≥1
	PageSize int                `json:"page_size"` // 每页条数
	Filter   *RecoverBillFilter `json:"filter"`
	Env      int                `json:"env"` // 由 SDK 自动填充
}

// RecoverBill 广告金回收记录。
type RecoverBill struct {
	BillID             string   `json:"bill_id"`              // 回收单 ID
	RecoverTime        int64    `json:"recover_time"`         // 回收时间，Unix 秒
	SettleBegin        int64    `json:"settle_begin"`         // 结算周期开始，Unix 秒
	SettleEnd          int64    `json:"settle_end"`           // 结算周期结束，Unix 秒
	FundID             string   `json:"fund_id"`              // 对应发放广告金 ID
	RecoverAccountName string   `json:"recover_account_name"` // 回收账户
	RecoverAmount      int64    `json:"recover_amount"`       // 回收金额，单位分
	RefundOrderList    []string `json:"refund_order_list"`    // 对应退款订单 id
}

// QueryRecoverBillResp 查询广告金回收记录响应。
type QueryRecoverBillResp struct {
	baseResp
	TotalPage int           `json:"total_page"` // 总页数
	BillList  []RecoverBill `json:"bill_list"`  // 回收记录列表
}

// QueryRecoverBill 查询广告金回收记录。
func (c *Client) QueryRecoverBill(ctx context.Context, req *QueryRecoverBillReq) (*QueryRecoverBillResp, error) {
	var out QueryRecoverBillResp
	if err := c.post(ctx, uriQueryRecoverBill, req, requestOptions{}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
