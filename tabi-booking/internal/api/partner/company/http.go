package company

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"
	"time"

	"github.com/labstack/echo/v4"
)

type Service interface {
	UpdateCompany(ctx context.Context, authoPartner *model.AuthoPartner, data *CompanyUpdateRequest) (*model.CompanyResponse, error)
	ViewCompany(ctx context.Context, authoPartner *model.AuthoPartner) (*model.CompanyResponse, error)
	AnalyzeRevenues(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.RevenueAnalysisData, error)
	AnalyzeBookingRequestQuantity(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.BookingRequestQuantityAnalysisData, error)
}

type HTTP struct {
	service Service
	autho   model.Autho
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation PATCH /partner/company partner-company PartnerCompanyUpdate
	// ---
	// summary: Update representative's company detail
	// parameters:
	// - name: body
	//   in: body
	//   description: company information to update
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CompanyUpdateRequest"
	// responses:
	//   "200":
	//     description: updated company information
	//     schema:
	//       "$ref": "#/definitions/CompanyResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.PATCH("", http.updateCompany)

	// swagger:operation GET /partner/company partner-company PartnerCompanyView
	// ---
	// summary: View representative's company detail
	// responses:
	//   "200":
	//     description: company information
	//     schema:
	//       "$ref": "#/definitions/CompanyResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.viewCompany)

	// swagger:operation GET /partner/company/analyses/revenues partner-company PartnerCompanyAnalysisRevenues
	// ---
	// summary: Analyze company revenues
	// responses:
	//   "200":
	//     description: company revenues analysis data
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

	// swagger:operation GET /partner/company/analyses/booking-request-quantity partner-company PartnerCompanyAnalysisBookingRequestQuantity
	// ---
	// summary: Analyze booking request quantity of company
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

func (h *HTTP) updateCompany(c echo.Context) error {
	reqBody := &CompanyUpdateRequest{}
	if err := c.Bind(reqBody); err != nil {
		return err
	}

	resp, err := h.service.UpdateCompany(c.Request().Context(), h.autho.Partner(c), reqBody)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) viewCompany(c echo.Context) error {
	resp, err := h.service.ViewCompany(c.Request().Context(), h.autho.Partner(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) analyzeRevenues(c echo.Context) error {
	req := CompanyAnalysisRequest{}
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
	req := CompanyAnalysisRequest{}
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
