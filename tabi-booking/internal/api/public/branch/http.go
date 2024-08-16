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
}

type Service interface {
	List(ctx context.Context, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) (*branch.PublicBranchListResponse, error)
	View(ctx context.Context, id int) (*model.PublicBranchResponse, error)
	ListRooms(ctx context.Context, branchID int, lq *dbcore.ListQueryCondition, lqBranch *branch.PublicBranchCondition) (*RoomListResponse, error)
	ListFeaturedDestinations(ctx context.Context) ([]string, error)
	ListFeaturedBranches(ctx context.Context) ([]*branch.PublicBranch, error)
	ListRecommendedBranches(ctx context.Context, lq *dbcore.ListQueryCondition, userID int) (*branch.PublicBranchListResponse, error)
}

func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	// swagger:operation GET /branches public-branch PublicBranchList
	// ---
	// summary: Returns list of branches
	// responses:
	//   "200":
	//     description: list of branches
	//     schema:
	//       "$ref": "#/definitions/PublicBranchListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	// swagger:operation GET /branches/{id} public-branch PublicBranchView
	// ---
	// summary: View branch detail
	// parameters:
	// - name: id
	//   in: path
	//   description: id of branch
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: branch information
	//     schema:
	//       "$ref": "#/definitions/PublicBranchResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id", http.view)

	// swagger:operation GET /branches/{id}/rooms public-branch PublicBranchListRooms
	// ---
	// summary: List rooms of branch
	// parameters:
	// - name: id
	//   in: path
	//   description: id of branch
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     description: list of rooms
	//     schema:
	//       "$ref": "#/definitions/PublicRoomListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id/rooms", http.listRooms)

	// swagger:operation GET /branches/featured-destinations public-branch PublicBranchListFeaturedDestinations
	// ---
	// summary: List featured destinations
	// responses:
	//   "200":
	//     description: list of featured destinations
	//     schema:
	//       "$ref": "#/definitions/FeaturedDestinationListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/featured-destinations", http.listFeaturedDestinations)

	// swagger:operation GET /branches/featured-branches public-branch PublicBranchListFeaturedBranches
	// ---
	// summary: List featured branches
	// responses:
	//   "200":
	//     description: list of featured branches
	//     schema:
	//       "$ref": "#/definitions/FeaturedBranchListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/featured-branches", http.listFeaturedBranches)
	echoGroup.GET("/recommended-branches", http.listRecommendedBranches)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	lqBranch, err := branch.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), lq, lqBranch)
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
	resp, err := h.service.View(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) listRooms(c echo.Context) error {
	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	lqBranch, err := branch.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.ListRooms(c.Request().Context(), branchID, lq, lqBranch)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) listFeaturedDestinations(c echo.Context) error {
	resp, err := h.service.ListFeaturedDestinations(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, FeaturedDestinationListResponse{
		Data: resp,
	})
}

func (h *HTTP) listFeaturedBranches(c echo.Context) error {
	resp, err := h.service.ListFeaturedBranches(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, FeaturedBranchListResponse{
		Data: resp,
	})
}

func (h *HTTP) listRecommendedBranches(c echo.Context) error {
	req := RecommendedBranchListRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.ListRecommendedBranches(c.Request().Context(), lq, req.UserID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
