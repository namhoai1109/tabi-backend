package partner

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	httpcore "github.com/namhoai1109/tabi/core/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents auth http service
type HTTP struct {
	service Service
}

// Service represents auth service interface
type Service interface {
	RpRegister(ctx context.Context, regData RpRegistrationReq) (*model.AuthToken, error)
	HstRegister(ctx context.Context, regData HstRegistrationReq) (*HSTRegisterResponse, error)
	Login(ctx context.Context, credentials CredentialsPartnerReq) (*model.AuthToken, error)
	RefreshToken(ctx context.Context, reqData RefreshTokenReq) (*model.AuthToken, error)
	Delete(ctx context.Context, rpID int) error
}

// NewHTTP creates new auth http service
func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	// swagger:operation POST /authen/partners/rp/register authen-partner AuthenRepresentativeRegister
	// ---
	// summary: Register new partner with info of representative and company
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RpRegistrationReq"
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
	echoGroup.POST("/rp/register", http.rpRegister)

	// swagger:operation POST /authen/partners/hst/register authen-partner AuthenHostRegister
	// ---
	// summary: Register new partner with info of host and branch
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/HstRegistrationReq"
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
	echoGroup.POST("/hst/register", http.hstRegister)

	// swagger:operation POST /authen/partners/login authen-partner AuthenPartnerLogin
	// ---
	// summary: Login partner by username and password
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CredentialsPartnerReq"
	// responses:
	//   "200":
	//     description: "OK"
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

	// swagger:operation POST /authen/partners/refresh-token authen-partner AuthenPartnerRefreshToken
	// ---
	// summary: Refresh access token
	// parameters:
	// - name: token
	//   in: body
	//   description: The given `refresh_token` when login
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/RefreshTokenReq"
	// responses:
	//   "200":
	//     description: New access token
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
	echoGroup.POST("/refresh-token", http.refreshToken)

	// swagger:operation DELETE /authen/partners/{id} authen-partner AuthenPartnerDelete
	// ---
	// summary: Deletes a partner by account id
	// parameters:
	// - name: id
	//   in: path
	//   description: id of representative
	//   type: integer
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.DELETE("/:id", http.delete)
}

func (h *HTTP) rpRegister(c echo.Context) error {
	registrationData := RpRegistrationReq{}
	if err := c.Bind(&registrationData); err != nil {
		return err
	}

	resp, err := h.service.RpRegister(c.Request().Context(), registrationData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) hstRegister(c echo.Context) error {
	registrationData := HstRegistrationReq{}
	if err := c.Bind(&registrationData); err != nil {
		return err
	}

	resp, err := h.service.HstRegister(c.Request().Context(), registrationData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) refreshToken(c echo.Context) error {
	reqData := RefreshTokenReq{}
	if err := c.Bind(&reqData); err != nil {
		return err
	}

	resp, err := h.service.RefreshToken(c.Request().Context(), reqData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) login(c echo.Context) error {
	credentials := CredentialsPartnerReq{}
	if err := c.Bind(&credentials); err != nil {
		return err
	}

	resp, err := h.service.Login(c.Request().Context(), credentials)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
