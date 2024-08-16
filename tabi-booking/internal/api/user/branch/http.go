package branch

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"

	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Save(ctx context.Context, autho *model.AuthoUser, branchID int, req SaveBranchRequest) (*SaveBranchResponse, error)
	ListSaved(ctx context.Context, autho *model.AuthoUser, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) ([]*branch.PublicBranch, error)
	Rate(ctx context.Context, autho *model.AuthoUser, branchID int, req RatingBranchRequest) error
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /user/branches/{id} user-branch UserBranchSaveUnSave
	// ---
	// summary:	Mark a branch as saved or unsaved
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
	//     "$ref": "#/definitions/SaveBranchRequest"
	// responses:
	//   "200":
	//     description: message response
	//     schema:
	//       "$ref": "#/definitions/SaveBranchResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/:id", http.save)

	// swagger:operation GET /user/branches/saved user-branch UserBranchListSaved
	// ---
	// summary:	List saved branches
	// responses:
	//   "200":
	//     description: list of saved branches
	//     schema:
	//       "$ref": "#/definitions/UserBranchListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/saved", http.listSaved)

	echoGroup.POST("/:id/rate", http.rate)
}

func (h *HTTP) save(c echo.Context) error {
	req := SaveBranchRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.Save(c.Request().Context(), h.autho.User(c), branchID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) listSaved(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	lqBranch, err := branch.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.ListSaved(c.Request().Context(), h.autho.User(c), lq, lqBranch)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, BranchListResponse{
		Data:  resp,
		Total: len(resp),
	})
}

func (h *HTTP) rate(c echo.Context) error {
	req := RatingBranchRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	if err := h.service.Rate(c.Request().Context(), h.autho.User(c), branchID, req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
