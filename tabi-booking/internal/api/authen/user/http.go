package user

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	service Service
}

// Service represents auth service interface
type Service interface {
	Register(ctx context.Context, regData RegistrationUserReq) (*model.AuthToken, error)
	Login(ctx context.Context, credentials CredentialsUserReq) (*model.AuthToken, error)
}

// NewHTTP creates new auth http service
func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	// swagger:operation POST /authen/users/register authen-user AuthenUserRegister
	// ---
	// summary: Register new user
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RegistrationUserReq"
	// responses:
	//   "200":
	//     description: "auth token"
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/register", http.register)

	// swagger:operation POST /authen/users/login authen-user AuthenUserLogin
	// ---
	// summary: User login
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CredentialsUserReq"
	// responses:
	//   "200":
	//     description: "auth token"
	//     schema:
	//       "$ref": "#/definitions/AuthToken"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/login", http.login)
}

func (h *HTTP) register(c echo.Context) error {
	registrationData := RegistrationUserReq{}
	if err := c.Bind(&registrationData); err != nil {
		return err
	}

	resp, err := h.service.Register(c.Request().Context(), registrationData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) login(c echo.Context) error {
	credentials := CredentialsUserReq{}
	if err := c.Bind(&credentials); err != nil {
		return err
	}

	resp, err := h.service.Login(c.Request().Context(), credentials)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
