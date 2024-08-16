package schedule

import (
	"context"
	"net/http"
	"tabi-notification/internal/model"

	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Create(ctx context.Context, autho *model.AuthoUser, req ScheduleCreationRequest) error
	List(ctx context.Context, autho *model.AuthoUser, bookingID int, lq *dbcore.ListQueryCondition) (*ScheduleListResponse, error)
	View(ctx context.Context, autho *model.AuthoUser, id int) (*ScheduleResponse, error)
	Update(ctx context.Context, autho *model.AuthoUser, id int, req ScheduleCreation) error
	Delete(ctx context.Context, autho *model.AuthoUser, id int) error
	DeleteIDs(ctx context.Context, autho *model.AuthoUser, ids []int) error
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /user/schedules user-schedules UserCreateSchedule
	//  ---
	// summary: Create a new schedule
	// parameters:
	// - name: request
	//   in: body
	//   required: true
	//   description: Request body
	//   schema:
	//     "$ref": "#/definitions/ScheduleCreationRequest"
	// responses:
	//   "204":
	//     description: Created
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

	// swagger:operation GET /user/schedules/bookings/{id} user-schedules UserListSchedule
	//  ---
	// summary: List schedules
	// parameters:
	// - name: booking id
	//   in: path
	//   required: true
	//   description: booking id
	//   type: integer
	// responses:
	//   "200":
	//     description: list of schedules
	//     schema:
	//       "$ref": "#/definitions/ScheduleListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     $ref: "#/responses/errDetails"
	echoGroup.GET("/bookings/:id", http.list)

	// swagger:operation GET /user/schedules/{id} user-schedules UserViewSchedule
	//  ---
	// summary: View schedule
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   description: schedule id
	//   type: integer
	// responses:
	//   "200":
	//     description: schedule info
	//     schema:
	//       "$ref": "#/definitions/ScheduleResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id", http.view)

	// swagger:operation PUT /user/schedules/{id} user-schedules UserUpdateSchedule
	//  ---
	// summary: Update schedule
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ScheduleCreation"
	// responses:
	//   "204":
	//     description: No Content
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PUT("/:id", http.update)

	// swagger:operation DELETE /user/schedules/{id} user-schedules UserDeleteSchedule
	//  ---
	// summary: Delete schedule
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   description: schedule id
	//   type: integer
	// responses:
	//   "204":
	//     description: No Content
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.DELETE("/:id", http.delete)

	echoGroup.DELETE("", http.deleteIDs)
}

func (h *HTTP) create(c echo.Context) error {
	req := ScheduleCreationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Create(c.Request().Context(), h.autho.User(c), req); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (h *HTTP) list(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), h.autho.User(c), id, lq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.View(c.Request().Context(), h.autho.User(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	req := ScheduleCreation{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Update(c.Request().Context(), h.autho.User(c), id, req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), h.autho.User(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) deleteIDs(c echo.Context) error {
	req := ScheduleDeleteRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.DeleteIDs(c.Request().Context(), h.autho.User(c), req.IDs); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
