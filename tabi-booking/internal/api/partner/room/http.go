package room

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
	Create(ctx context.Context, authoPartner *model.AuthoPartner, req CreateRoomRequest) (*model.Room, error)
	View(ctx context.Context, authoPartner *model.AuthoPartner, roomID int) (*ViewRoomResponse, error)
	List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition) (*RoomListResponse, error)
	ListBookings(ctx context.Context, authoPartner *model.AuthoPartner, roomID int, lq *dbcore.ListQueryCondition, count *int64) ([]*model.BookingResponse, error)
	Update(ctx context.Context, authoPartner *model.AuthoPartner, updates UpdateRoomRequest, roomID int) (*ViewRoomResponse, error)
	GetPaymentInfo(ctx context.Context, authoPartner *model.AuthoPartner, roomID int) (*PaymentInfoResponse, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /partner/rooms partner-room PartnerRoomCreate
	// ---
	// summary: Create a room for room type
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CreateRoomRequest"
	// responses:
	//   "200":
	//     description: room info
	//     schema:
	//       "$ref": "#/definitions/Room"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

	// swagger:operation GET /partner/rooms partner-room PartnerRoomList
	// ---
	// summary: Returns list of rooms
	// responses:
	//   "200":
	//     description: list of rooms
	//     schema:
	//       "$ref": "#/definitions/RoomListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	// swagger:operation GET /partner/rooms/{id} partner-room PartnerRoomView
	// ---
	// summary: Returns detailed room info
	// parameters:
	// - name: id
	//   in: path
	//   description: room id
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: room info
	//     schema:
	//       "$ref": "#/definitions/ViewRoomResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id", http.view)

	// swagger:operation GET /partner/rooms/{id}/bookings partner-room PartnerRoomListBookings
	// ---
	// summary: Returns list of bookings for a room
	// parameters:
	// - name: id
	//   in: path
	//   description: room id
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: room info
	//     schema:
	//       "$ref": "#/definitions/ListBookingsResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id/bookings", http.listBookings)

	// swagger:operation PATCH /partner/rooms/{id} partner-room PartnerRoomUpdate
	// ---
	// summary: Update room information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of room
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdateRoomRequest"
	// responses:
	//   "200":
	//     description: room information
	//     schema:
	//       "$ref": "#/definitions/ViewRoomResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PATCH("/:id", http.update)
	echoGroup.GET("/:id/payment-info", http.getPaymentInfo)
}

func (h *HTTP) create(c echo.Context) error {
	req := CreateRoomRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.service.Create(c.Request().Context(), h.autho.Partner(c), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), h.autho.Partner(c), lq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) view(c echo.Context) error {
	roomID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.View(c.Request().Context(), h.autho.Partner(c), roomID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	updates := UpdateRoomRequest{}
	if err := c.Bind(&updates); err != nil {
		return err
	}

	roomID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.Update(c.Request().Context(), h.autho.Partner(c), updates, roomID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) listBookings(c echo.Context) error {
	roomID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	var count int64 = 0
	resp, err := h.service.ListBookings(c.Request().Context(), h.autho.Partner(c), roomID, lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListBookingsResponse{
		Total: count,
		Data:  resp,
	})
}

func (h *HTTP) getPaymentInfo(c echo.Context) error {
	roomID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetPaymentInfo(c.Request().Context(), h.autho.Partner(c), roomID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
