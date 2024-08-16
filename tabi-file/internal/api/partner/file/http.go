package file

import (
	"context"
	"net/http"
	"tabi-file/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
	"github.com/namhoai1109/tabi/core/server"

	"github.com/labstack/echo/v4"
)

// HTTP represents file http service
type HTTP struct {
	service Service
	autho   model.Autho
}

// Service represents file application interface
type Service interface {
	ListFile(ctx context.Context, autho *model.AuthoPartner, lq *dbcore.ListQueryCondition) ([]*model.FileResponse, error)
	DeleteFile(ctx context.Context, auth *model.AuthoPartner, IDs []int) error
}

// NewHTTP creates new file http service
func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation GET /partner/files partner-file PartnerFileList
	// ---
	// summary: Response List file(s)
	// responses:
	//   "200":
	//     description: List of file(s)
	//     schema:
	//        "$ref": "#/definitions/ListFileResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.list)

	// swagger:operation DELETE /partner/files partner-file PartnerFileDelete
	// ---
	// summary: Delete file(s)
	// parameters:
	// - name: file ID list
	//   in: body
	//   description: File ID list to be deleted
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/DeleteFileRequest"
	// responses:
	//   "200":
	//   	"$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.DELETE("", http.delete)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}
	if lq == nil || lq.Filter == nil {
		return server.NewHTTPValidationError("Invalid query")
	}

	resp, err := h.service.ListFile(c.Request().Context(), h.autho.Partner(c), lq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListFileResponse{resp})
}

func (h *HTTP) delete(c echo.Context) error {
	req := DeleteFileRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.DeleteFile(c.Request().Context(), h.autho.Partner(c), req.IDs); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
