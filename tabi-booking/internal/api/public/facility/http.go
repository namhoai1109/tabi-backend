package facility

import (
	"context"
	"net/http"

	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents partner facility http service
type HTTP struct {
	service Service
}

// Service represents partner facility service interface
type Service interface {
	List(ctx context.Context, lq *dbcore.ListQueryCondition, lang string) ([]*FacilityResponse, error)
}

// NewHTTP creates new facility http service
func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}
	// swagger:operation GET /facilities/{lang} public-facility PublicFacilityList
	// ---
	// summary: Returns list of facilities
	// responses:
	//   "200":
	//     description: list of facilities
	//     schema:
	//       "$ref": "#/definitions/FacilityListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:lang", http.list)
}

func (h *HTTP) list(c echo.Context) error {
	lang := c.Param("lang")

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), lq, lang)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &FacilityListResponse{Data: resp})
}
