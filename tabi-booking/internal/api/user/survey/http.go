package survey

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Create(ctx context.Context, authoUser *model.AuthoUser, req SurveyCreationRequest) (*model.Survey, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /user/surveys user-survey UserSurveyCreate
	// ---
	// summary: Create a new survey for user
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/SurveyCreationRequest"
	// responses:
	//   "200":
	//     description: survey details
	//     schema:
	//       "$ref": "#/definitions/Survey"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("", http.create)

}

func (h *HTTP) create(c echo.Context) error {
	req := SurveyCreationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.service.Create(c.Request().Context(), h.autho.User(c), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
