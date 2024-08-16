package payment

import (
	"fmt"

	paypalmodel "github.com/namhoai1109/tabi/core/paypal/model"
)

func (s *Payment) convertVNDCurrencyToUSD(vnd float64) string {
	return fmt.Sprintf("%.2f", vnd/25000)
}

func (s *Payment) getOrderRequest(price float64, roomName string, quantity int, payee string) *paypalmodel.CreateOrderRequest {

	priceUS := s.convertVNDCurrencyToUSD(price)

	resp := &paypalmodel.CreateOrderRequest{
		Intent: paypalmodel.OrderIntentCapture,
		PurchaseUnits: []paypalmodel.PurchaseUnit{
			{
				Items: []paypalmodel.PurchaseUnitItem{
					{
						Name:     fmt.Sprintf("%s x %d", roomName, quantity),
						Quantity: "1",
						UnitAmount: paypalmodel.UnitAmount{
							CurrencyCode: paypalmodel.OrderCurrencyUSD,
							Value:        priceUS,
						},
					},
				},
				Amount: paypalmodel.PurchaseUnitAmount{
					CurrencyCode: paypalmodel.OrderCurrencyUSD,
					Value:        priceUS,
					Breakdown: paypalmodel.PurchaseUnitAmountBreakdown{
						ItemTotal: paypalmodel.UnitAmount{
							CurrencyCode: paypalmodel.OrderCurrencyUSD,
							Value:        priceUS,
						},
					},
				},
				Payee: paypalmodel.PurchaseUnitPayee{
					EmailAddress: payee,
				},
			},
		},
	}

	resp.PaymentSource.Paypal.ExperienceContext.ReturnUrl = "https://namhoai1109.github.io/processing/"
	resp.PaymentSource.Paypal.ExperienceContext.CancelUrl = "https://example.com/cancel"

	return resp
}
