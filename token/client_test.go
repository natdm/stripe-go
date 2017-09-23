package token

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	stripe "github.com/stripe/stripe-go"
	_ "github.com/stripe/stripe-go/testing"
)

func TestTokenGet(t *testing.T) {
	token, err := Get("tok_123", nil)
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestTokenNew_WithBankAccount(t *testing.T) {
	token, err := New(&stripe.TokenParams{
		BankAccount: &stripe.BankAccountParams{
			Country:       "US",
			RoutingNumber: "110000000",
			AccountNumber: "000123456789",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestTokenNew_WithCard(t *testing.T) {
	token, err := New(&stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   "4242424242424242", // raw PAN as we're testing token creation
			ExpMonth: "10",
			ExpYear:  "20",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestTokenNew_WithPII(t *testing.T) {
	token, err := New(&stripe.TokenParams{
		PII: &stripe.PIIParams{
			PersonalIDNumber: "000000000",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
