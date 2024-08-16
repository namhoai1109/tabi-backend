package generaltype

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents partner general type http service
type HTTP struct {
	service Service
}

// Service represents partner general type service interface
type Service interface {
	ListAccommodation(ctx context.Context, lang string) ([]*AccommodationTypeResponse, error)
	ListBed(ctx context.Context, lang string) ([]*BedTypeResponse, error)
}

// NewHTTP creates new general type http service
func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}
	// swagger:operation GET /general-types/accommodations/{lang} public-general-type PublicGeneralTypeAccommodationList
	// ---
	// summary: Returns list of accommodation type
	// responses:
	//   "200":
	//     description: list of accommodation type
	//     schema:
	//       "$ref": "#/definitions/AccommodationTypeListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/accommodations/:lang", http.listAccommodation)
	// swagger:operation GET /general-types/beds/{lang} public-general-type PublicGeneralTypeBedList
	// ---
	// summary: Returns list of bed type
	// responses:
	//   "200":
	//     description: list of bed type
	//     schema:
	//       "$ref": "#/definitions/BedTypeListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/beds/:lang", http.listBed)
}

func (h *HTTP) listAccommodation(c echo.Context) error {
	lang := c.Param("lang")

	resp, err := h.service.ListAccommodation(c.Request().Context(), lang)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &AccommodationTypeListResponse{Data: resp})
}

func (h *HTTP) listBed(c echo.Context) error {
	lang := c.Param("lang")

	resp, err := h.service.ListBed(c.Request().Context(), lang)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &BedTypeListResponse{Data: resp})
}
