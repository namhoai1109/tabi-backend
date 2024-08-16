package payment

import (
	"context"
	"net/http"
	"tabi-payment/internal/model"

	"github.com/labstack/echo/v4"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Create(ctx context.Context, authoUser *model.AuthoUser, req *PaymentCreationRequest) (*PaymentCreationResponse, error)
	Capture(ctx context.Context, authoUser *model.AuthoUser, orderID string, req *PaymentCaptureRequest) (*PaymentCaptureResponse, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /user/payments user-payment UserPaymentCreate
	// ---
	// summary: Create a payment
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/PaymentCreationRequest"
	// responses:
	//   "200":
	//     description: approve link
	//     schema:
	//       "$ref": "#/definitions/PaymentCreationResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

	// swagger:operation POST /user/payments/capture/{orderID} user-payment UserPaymentCapture
	// ---
	// summary: Capture a payment
	// parameters:
	// - name: orderID
	//   in: path
	//   description: id of paypal order
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: completed message
	//     schema:
	//       "$ref": "#/definitions/PaymentCaptureResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/capture/:orderID", http.capture)

}

func (h *HTTP) create(c echo.Context) error {
	req := &PaymentCreationRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	resp, err := h.service.Create(c.Request().Context(), h.autho.User(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) capture(c echo.Context) error {
	orderID := c.Param("orderID")

	req := &PaymentCaptureRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	resp, err := h.service.Capture(c.Request().Context(), h.autho.User(c), orderID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
