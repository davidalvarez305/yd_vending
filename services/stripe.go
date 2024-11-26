package services

import (
	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

func CreateStripeCheckout() (string, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String("{{PRICE_ID}}"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(constants.RootDomain + "/success.html"),
		CancelURL:  stripe.String(constants.RootDomain + "/cancel.html"),
	}

	s, err := session.New(params)

	if err != nil {
		return "", err
	}

	return s.URL, nil
}
