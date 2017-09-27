package stripe

import (
	"encoding/json"
)

// FeeRefundParams is the set of parameters that can be used when refunding a fee.
// For more details see https://stripe.com/docs/api#fee_refund.
type FeeRefundParams struct {
	Params         `form:"*"`
	Amount         uint64 `form:"amount"`
	ApplicationFee string `form:"-"` // Included in the URL
}

// FeeRefundListParams is the set of parameters that can be used when listing fee refunds.
// For more details see https://stripe.com/docs/api#list_fee_refunds.
type FeeRefundListParams struct {
	ListParams     `form:"*"`
	ApplicationFee string `form:"-"` // Included in the URL
}

// FeeRefund is the resource representing a Stripe fee refund.
// For more details see https://stripe.com/docs/api#fee_refunds.
type FeeRefund struct {
	Amount             uint64              `json:"amount"`
	BalanceTransaction *BalanceTransaction `json:"balance_transaction"`
	Created            int64               `json:"created"`
	Currency           Currency            `json:"currency"`
	Fee                string              `json:"fee"`
	ID                 string              `json:"id"`
	Metadata           map[string]string   `json:"metadata"`
}

// FeeRefundList is a list object for fee refunds.
type FeeRefundList struct {
	ListMeta
	Data []*FeeRefund `json:"data"`
}

// UnmarshalJSON handles deserialization of a FeeRefund.
// This custom unmarshaling is needed because the resulting
// property may be an id or the full struct if it was expanded.
func (f *FeeRefund) UnmarshalJSON(data []byte) error {
	type feerefund FeeRefund
	var ff feerefund
	err := json.Unmarshal(data, &ff)
	if err == nil {
		*f = FeeRefund(ff)
	} else {
		// the id is surrounded by "\" characters, so strip them
		f.ID = string(data[1 : len(data)-1])
	}

	return nil
}
