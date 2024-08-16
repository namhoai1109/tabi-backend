package branchmanager

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
	Register(ctx context.Context, authoPartner *model.AuthoPartner, regData BranchManagerRegisterRequest) (*model.BranchManagerResponse, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /partner/branch-managers/register partner-branch-manager PartnerBranchManagerRegister
	// ---
	// summary: Register new branch manager, only representative can register branch manager
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/BranchManagerRegisterRequest"
	// responses:
	//   "200":
	//     description: branch manager information
	//     schema:
	//       "$ref": "#/definitions/BranchManagerResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/register", http.register)
}

func (h *HTTP) register(c echo.Context) error {
	regData := BranchManagerRegisterRequest{}
	if err := c.Bind(&regData); err != nil {
		return err
	}

	resp, err := h.service.Register(c.Request().Context(), h.autho.Partner(c), regData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
