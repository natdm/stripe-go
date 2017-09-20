package stripe

import (
	"encoding/json"
)

// SourceStatus represents the possible statuses of a source object.
type SourceStatus string

const (
	// SourceStatusCanceled we canceled the source along with any side-effect
	// it had (returned funds to customers if any were sent).
	SourceStatusCanceled SourceStatus = "canceled"

	// SourceStatusChargeable the source is ready to be charged (once if usage
	// is `single_use`, repeatidly otherwise).
	SourceStatusChargeable SourceStatus = "chargeable"

	// SourceStatusConsumed the source is `single_use` usage and has been
	// charged already.
	SourceStatusConsumed SourceStatus = "consumed"

	// SourceStatusFailed the source is no longer usable.
	SourceStatusFailed SourceStatus = "failed"

	// SourceStatusPending the source is freshly created and not yet
	// chargeable. The flow should indicate how to authenticate it with your
	// customer.
	SourceStatusPending SourceStatus = "pending"
)

// SourceFlow represents the possible flows of a source object.
type SourceFlow string

const (
	// FlowNone no particular authentication is involved the source should
	// become chargeable directly or asyncrhonously.
	FlowNone SourceFlow = "none"

	// FlowReceiver a receiver address should be communicated to the customer
	// to push funds to it.
	FlowReceiver SourceFlow = "receiver"

	// FlowRedirect a redirect is required to authenticate the source.
	FlowRedirect SourceFlow = "redirect"

	// FlowVerification a verification code should be communicated by the
	// customer to authenticate the source.
	FlowVerification SourceFlow = "verification"
)

// SourceUsage represents the possible usages of a source object.
type SourceUsage string

const (
	// UsageReusable the source can be charged multiple times for arbitrary
	// amounts.
	UsageReusable SourceUsage = "reusable"

	// UsageSingleUse the source can only be charged once for the specified
	// amount and currency.
	UsageSingleUse SourceUsage = "single_use"
)

type SourceOwnerParams struct {
	Address *AddressParams `form:"address"`
	Email   string         `form:"email"`
	Name    string         `form:"name"`
	Phone   string         `form:"phone"`
}

type RedirectParams struct {
	ReturnURL string `form:"return_url"`
}

type SourceObjectParams struct {
	Params   `form:"*"`
	Amount   uint64             `form:"amount"`
	Currency Currency           `form:"currency"`
	Customer string             `form:"customer"`
	Flow     SourceFlow         `form:"flow"`
	Owner    *SourceOwnerParams `form:"owner"`
	Redirect *RedirectParams    `form:"redirect"`
	Token    string             `form:"token"`
	Type     string             `form:"type"`
	TypeData map[string]string  `form:"*"`
	Usage    SourceUsage        `form:"usage"`
}

type SourceOwner struct {
	Address         *Address `json:"address,omitempty"`
	Email           string   `json:"email"`
	Name            string   `json:"name"`
	Phone           string   `json:"phone"`
	VerifiedAddress *Address `json:"verified_address,omitempty"`
	VerifiedEmail   string   `json:"verified_email"`
	VerifiedName    string   `json:"verified_name"`
	VerifiedPhone   string   `json:"verified_phone"`
}

// RedirectFlowStatus represents the possible statuses of a redirect flow.
type RedirectFlowStatus string

const (
	RedirectFlowStatusFailed    RedirectFlowStatus = "failed"
	RedirectFlowStatusPending   RedirectFlowStatus = "pending"
	RedirectFlowStatusSucceeded RedirectFlowStatus = "succeeded"
)

// ReceiverFlow informs of the state of a redirect authentication flow.
type RedirectFlow struct {
	ReturnURL string             `json:"return_url"`
	Status    RedirectFlowStatus `json:"status"`
	URL       string             `json:"url"`
}

// RefundAttributesStatus are the possible status of a receiver's refund
// attributes.
type RefundAttributesStatus string

const (
	// RefundAttributesAvailable the refund attributes are available
	RefundAttributesAvailable RefundAttributesStatus = "available"

	// RefundAttributesMissing the refund attributes are missing
	RefundAttributesMissing RefundAttributesStatus = "missing"

	// RefundAttributesRequested the refund attributes have been requested
	RefundAttributesRequested RefundAttributesStatus = "requested"
)

// RefundAttributesMethod are the possible method to retrieve a receiver's
// refund attributes.
type RefundAttributesMethod string

const (
	// RefundAttributesEmail the refund attributes are automatically collected over email
	RefundAttributesEmail RefundAttributesMethod = "email"

	// RefundAttributesManual the refund attributes should be collected by the user
	RefundAttributesManual RefundAttributesMethod = "manual"
)

// ReceiverFlow informs of the state of a receiver authentication flow.
type ReceiverFlow struct {
	Address                string                 `json:"address"`
	AmountCharged          int64                  `json:"amount_charged"`
	AmountReceived         int64                  `json:"amount_received"`
	AmountReturned         int64                  `json:"amount_returned"`
	RefundAttributesMethod RefundAttributesMethod `json:"refund_attributes_method"`
	RefundAttributesStatus RefundAttributesStatus `json:"refund_attributes_status"`
}

// VerificationFlowStatus represents the possible statuses of a verification
// flow.
type VerificationFlowStatus string

const (
	VerificationFlowStatusFailed    VerificationFlowStatus = "failed"
	VerificationFlowStatusPending   VerificationFlowStatus = "pending"
	VerificationFlowStatusSucceeded VerificationFlowStatus = "succeeded"
)

// ReceiverFlow informs of the state of a verification authentication flow.
type VerificationFlow struct {
	AttemptsRemaining uint64             `json:"attempts_remaining"`
	Status            RedirectFlowStatus `json:"status"`
}

type Source struct {
	Amount       int64             `json:"amount"`
	ClientSecret string            `json:"client_secret"`
	Created      int64             `json:"created"`
	Currency     Currency          `json:"currency"`
	Flow         SourceFlow        `json:"flow"`
	ID           string            `json:"id"`
	Livemode     bool              `json:"livemode"`
	Metadata     map[string]string `json:"metadata"`
	Owner        SourceOwner       `json:"owner"`
	Receiver     *ReceiverFlow     `json:"receiver,omitempty"`
	Redirect     *RedirectFlow     `json:"redirect,omitempty"`
	Status       SourceStatus      `json:"status"`
	Type         string            `json:"type"`
	TypeData     map[string]interface{}
	Usage        SourceUsage       `json:"usage"`
	Verification *VerificationFlow `json:"verification,omitempty"`
}

// UnmarshalJSON handles deserialization of an Source. This custom unmarshaling
// is needed to extract the type specific data (accessible under `TypeData`)
// but stored in JSON under a hash named after the `type` of the source.
func (s *Source) UnmarshalJSON(data []byte) error {
	type source Source
	var ss source
	err := json.Unmarshal(data, &ss)
	if err != nil {
		return err
	}
	*s = Source(ss)

	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	if d, ok := raw[s.Type]; ok {
		if m, ok := d.(map[string]interface{}); ok {
			s.TypeData = m
		}
	}

	return nil
}
