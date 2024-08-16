package booking

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.BookingResponse, error)
	Create(ctx context.Context, authoPartner *model.AuthoPartner, req *BookingRequest) error
	Approve(ctx context.Context, authoPartner *model.AuthoPartner, id int) error
	Reject(ctx context.Context, authoPartner *model.AuthoPartner, id int, req RejectReasonRequest) error
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation GET /partner/bookings partner-booking PartnerBookingList
	// ---
	// summary: List reservation request of user
	// responses:
	//   "200":
	//     description: List of reservation request
	//     schema:
	//       "$ref": "#/definitions/ListBookingResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	// api for payment service
	echoGroup.POST("", http.create)

	// swagger:operation GET /partner/bookings/{id}/approve partner-booking PartnerBookingApprove
	// ---
	// summary: approve reservation request of user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of booking
	//   type: integer
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/:id/approve", http.approve)

	// swagger:operation GET /partner/bookings/{id}/reject partner-booking PartnerBookingReject
	// ---
	// summary: reject reservation request of user
	// parameters:
	// - name: id
	//   in: path
	//   description: id of booking
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RejectReasonRequest"
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/:id/reject", http.reject)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	var count int64 = 0
	resp, err := h.service.List(c.Request().Context(), h.autho.Partner(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListBookingResponse{Data: resp, Total: count})
}

func (h *HTTP) create(c echo.Context) error {
	req := &BookingRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := h.service.Create(c.Request().Context(), h.autho.Partner(c), req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) approve(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	if err := h.service.Approve(c.Request().Context(), h.autho.Partner(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) reject(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	req := RejectReasonRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Reject(c.Request().Context(), h.autho.Partner(c), id, req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
