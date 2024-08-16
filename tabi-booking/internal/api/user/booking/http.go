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
	Create(ctx context.Context, autho *model.AuthoUser, req *UserBookingRequest) error
	List(ctx context.Context, autho *model.AuthoUser, lq *dbcore.ListQueryCondition, count *int64) ([]*model.Booking, error)
	Cancel(ctx context.Context, autho *model.AuthoUser, bookingID int, req CancelBookingRequest) error
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /user/bookings user-booking UserBookingCreate
	// ---
	// summary: Create a new booking
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UserBookingRequest"
	// responses:
	//   "204":
	//     description: no content
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

	// swagger:operation GET /user/bookings user-booking UserBookingList
	// ---
	// summary: List bookings
	// responses:
	//   "200":
	//     description: List of bookings
	//     schema:
	//       "$ref": "#/definitions/Booking"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	echoGroup.POST("/:id/cancel", http.cancel)
}

func (h *HTTP) create(c echo.Context) error {
	req := &UserBookingRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	err := h.service.Create(c.Request().Context(), h.autho.User(c), req)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	var count int64 = 0
	resp, err := h.service.List(c.Request().Context(), h.autho.User(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserBookingListResponse{
		Data:  resp,
		Total: count,
	})
}

func (h *HTTP) cancel(c echo.Context) error {
	bookingID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	req := CancelBookingRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Cancel(c.Request().Context(), h.autho.User(c), bookingID, req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
