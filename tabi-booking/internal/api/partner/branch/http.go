package branch

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"
	"time"

	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents partner branch http service
type HTTP struct {
	service Service
	autho   model.Autho
}

// Service represents partner branch service interface
type Service interface {
	CreateBranch(ctx context.Context, authoPartner *model.AuthoPartner, branchCreation *branch.BranchCreationRequest) (*model.BranchResponse, error)
	View(ctx context.Context, authoPartner *model.AuthoPartner, branchID int) (*model.BranchResponse, error)
	List(ctx context.Context, authoPartner *model.AuthoPartner, lq *dbcore.ListQueryCondition) (*BranchListResponse, error)
	Update(ctx context.Context, authoPartner *model.AuthoPartner, branchID int, data *BranchUpdateRequest) (*model.BranchResponse, error)
	Activate(ctx context.Context, authoPartner *model.AuthoPartner) (*ActivateBranchResponse, error)
	AnalyzeRevenues(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.RevenueAnalysisData, error)
	AnalyzeBookingRequestQuantity(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.BookingRequestQuantityAnalysisData, error)
}

// NewHTTP creates new branch http service
func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /partner/branches partner-branch PartnerBranchCreate
	// ---
	// summary: Create new branch of company
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/BranchCreationRequest"
	// responses:
	//   "200":
	//     description: branch information
	//     schema:
	//       "$ref": "#/definitions/Branch"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.createBranch)

	// swagger:operation GET /partner/branches/{id} partner-branch PartnerBranchView
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
	//       "$ref": "#/definitions/BranchResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id", http.viewBranch)

	// swagger:operation GET /partner/branches partner-branch PartnerBranchList
	// ---
	// summary: Returns list of branches
	// responses:
	//   "200":
	//     description: list of branches
	//     schema:
	//       "$ref": "#/definitions/BranchListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.listBranch)

	// swagger:operation PATCH /partner/branches/{id} partner-branch PartnerBranchUpdate
	// ---
	// summary: Update branch information
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
	//     "$ref": "#/definitions/BranchUpdateRequest"
	// responses:
	//   "200":
	//     description: branch information
	//     schema:
	//       "$ref": "#/definitions/BranchResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PATCH("/:id", http.updateBranch)

	// swagger:operation POST /partner/branches/activate partner-branch PartnerBranchActivate
	// ---
	// summary: Activate branch
	// responses:
	//   "200":
	//     description: branch status
	//     schema:
	//       "$ref": "#/definitions/ActivateBranchResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/activate", http.activateBranch)

	// swagger:operation GET /partner/branches/analyses/revenues partner-branch PartnerBranchAnalysisRevenues
	// ---
	// summary: Analyze branch revenues
	// responses:
	//   "200":
	//     description: branch revenues analysis data
	//     schema:
	//       "$ref": "#/definitions/RevenueAnalysisData"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/analyses/revenues", http.analyzeRevenues)

	// swagger:operation GET /partner/branches/analyses/booking-request-quantity partner-branch PartnerBranchAnalysisBookingRequestQuantity
	// ---
	// summary: Analyze booking request quantity of branch
	// responses:
	//   "200":
	//     description: booking request quantity analysis data
	//     schema:
	//       "$ref": "#/definitions/BookingRequestQuantityAnalysisData"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/analyses/booking-request-quantity", http.analyzeBookingRequestQuantity)
}

func (h *HTTP) createBranch(c echo.Context) error {
	req := &branch.BranchCreationRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	resp, err := h.service.CreateBranch(c.Request().Context(), h.autho.Partner(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) viewBranch(c echo.Context) error {
	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.View(c.Request().Context(), h.autho.Partner(c), branchID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) listBranch(c echo.Context) error {
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

func (h *HTTP) updateBranch(c echo.Context) error {
	req := &BranchUpdateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.Update(c.Request().Context(), h.autho.Partner(c), branchID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) activateBranch(c echo.Context) error {
	resp, err := h.service.Activate(c.Request().Context(), h.autho.Partner(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) analyzeRevenues(c echo.Context) error {
	req := BranchAnalysisRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	year := time.Now().Year()
	if req.Year != 0 {
		year = req.Year
	}

	resp, err := h.service.AnalyzeRevenues(c.Request().Context(), h.autho.Partner(c), year)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) analyzeBookingRequestQuantity(c echo.Context) error {
	req := BranchAnalysisRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	year := time.Now().Year()
	if req.Year != 0 {
		year = req.Year
	}

	resp, err := h.service.AnalyzeBookingRequestQuantity(c.Request().Context(), h.autho.Partner(c), year)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
