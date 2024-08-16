package roomtype

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
	List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomTypeResponse, error)
	ListAll(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition, count *int64) ([]*model.RoomTypeResponse, error)
	Create(ctx context.Context, authoPartner *model.AuthoPartner, req *RoomTypeCreateRequest) (*model.RoomType, error)
	Update(ctx context.Context, authoPartner *model.AuthoPartner, roomTypeID int, update RoomTypeUpdateRequest) (*model.RoomType, error)
	ChangeLinkStatus(ctx context.Context, authoPartner *model.AuthoPartner, req LinkRoomTypeRequest) (*string, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation GET /partner/room-types partner-room-type PartnerRoomTypeList
	// ---
	// summary: Returns list of room types according to branch id
	// responses:
	//   "200":
	//     description: list of room types
	//     schema:
	//       "$ref": "#/definitions/RoomTypeListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	// swagger:operation GET /partner/room-types/all partner-room-type PartnerRoomTypeListAll
	// ---
	// summary: Returns list all of room types of partner
	// responses:
	//   "200":
	//     description: list of room types
	//     schema:
	//       "$ref": "#/definitions/RoomTypeListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/all", http.listAll)

	// swagger:operation POST /partner/room-types partner-room-type PartnerRoomTypeCreate
	// ---
	// summary: Create a room type for a branch
	// parameters:
	// - name: id
	//   in: path
	//   description: id of branch
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RoomTypeCreateRequest"
	// responses:
	//   "200":
	//     description: room type
	//     schema:
	//       "$ref": "#/definitions/RoomType"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

	// swagger:operation PATCH /partner/room-types/{id} partner-room-type PartnerRoomTypeUpdate
	// ---
	// summary: Update a room type
	// parameters:
	// - name: id
	//   in: path
	//   description: id of branch
	//   type: integer
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RoomTypeUpdateRequest"
	// responses:
	//   "200":
	//     description: room type
	//     schema:
	//       "$ref": "#/definitions/RoomType"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PATCH("/:id", http.update)

	// swagger:operation PATCH /partner/room-types/link-status partner-room-type PartnerRoomTypeUpdate
	// ---
	// summary: link or unlink a room type to a branch
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/LinkRoomTypeRequest"
	// responses:
	//   "200":
	//     description: status of link
	//     schema:
	//       "$ref": "#/definitions/LinkStatusResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PATCH("/link-status", http.changeLinkStatus)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	var count int64
	resp, err := h.service.List(c.Request().Context(), h.autho.Partner(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, RoomTypeListResponse{
		Data:  resp,
		Total: count,
	})
}

func (h *HTTP) listAll(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	var count int64
	resp, err := h.service.ListAll(c.Request().Context(), h.autho.Partner(c), lq, &count)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, RoomTypeListResponse{
		Data:  resp,
		Total: count,
	})
}

func (h *HTTP) create(c echo.Context) error {
	var req RoomTypeCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.service.Create(c.Request().Context(), h.autho.Partner(c), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	roomTypeID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	var req RoomTypeUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.service.Update(c.Request().Context(), h.autho.Partner(c), roomTypeID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) changeLinkStatus(c echo.Context) error {
	var req LinkRoomTypeRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.service.ChangeLinkStatus(c.Request().Context(), h.autho.Partner(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LinkStatusResponse{
		Status: resp,
	})
}
