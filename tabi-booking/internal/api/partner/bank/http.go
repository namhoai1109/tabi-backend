package bank

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Register(ctx context.Context, authoPartner *model.AuthoPartner, reqData *BankRegisterRequest, branchID int) (*CustomBankResponse, error)
	List(ctx context.Context, authoPartner *model.AuthoPartner, branchID int, lq *dbcore.ListQueryCondition) (*BankListResponse, error)
	Update(ctx context.Context, authoPartner *model.AuthoPartner, bankID int, reqData *BankUpdateRequest) (*CustomBankResponse, error)
	Delete(ctx context.Context, authoPartner *model.AuthoPartner, bankID int) error
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation POST /partner/banks/{branchID} partner-bank PartnerBankRegister
	// ---
	// summary: Register bank account for a branch
	// parameters:
	// - name: branchID
	//   in: path
	//   description: Branch ID
	//   required: true
	//   type: integer
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/BankRegisterRequest"
	// responses:
	//   "200":
	//     description: bank account information
	//     schema:
	//       "$ref": "#/definitions/CustomBankResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/:id", http.register)

	// swagger:operation GET /partner/banks/{branchID} partner-bank PartnerBankList
	// ---
	// summary: List bank accounts of a branch
	// parameters:
	// - name: branchID
	//   in: path
	//   description: Branch ID
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: bank accounts
	//     schema:
	//       "$ref": "#/definitions/BankListResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("/:id", http.list)

	// echoGroup.PATCH("/:id", http.update)

	// swagger:operation DELETE /partner/banks/{bankID} partner-bank PartnerBankDelete
	// ---
	// summary: Delete bank account
	// parameters:
	// - name: bankID
	//   in: path
	//   description: Bank ID
	//   required: true
	//   type: integer
	// responses:
	//   "200":
	//     description: bank accounts
	//     schema:
	//       "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.DELETE("/:id", http.delete)
}

func (h *HTTP) register(c echo.Context) error {
	reqData := BankRegisterRequest{}
	if err := c.Bind(&reqData); err != nil {
		return err
	}

	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.Register(c.Request().Context(), h.autho.Partner(c), &reqData, branchID)
	if err != nil {
		return err
	}

	return c.JSON(200, resp)
}

func (h *HTTP) list(c echo.Context) error {
	branchID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), h.autho.Partner(c), branchID, lq)
	if err != nil {
		return err
	}

	return c.JSON(200, resp)
}

// func (h *HTTP) update(c echo.Context) error {
// 	reqData := BankUpdateRequest{}
// 	if err := c.Bind(&reqData); err != nil {
// 		return err
// 	}

// 	bankID, err := httpcore.ReqID(c)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := h.service.Update(c.Request().Context(), h.autho.Partner(c), bankID, &reqData)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(200, resp)
// }

func (h *HTTP) delete(c echo.Context) error {
	bankID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), h.autho.Partner(c), bankID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
