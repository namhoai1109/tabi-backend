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
}

// Service represents file application interface
type Service interface {
	ListFile(ctx context.Context, lq *dbcore.ListQueryCondition) ([]*model.FileResponse, error)
	UploadFilePreSigned(ctx context.Context, request UploadPresignedRequest) ([]*model.FileResponse, error)
}

// NewHTTP creates new file http service
func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	// swagger:operation GET /files public-file PublicFileList
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

	// swagger:operation POST /files/upload/presigned public-file PublicFileUploadPresigned
	// ---
	// summary: Upload new file(s) by pre-signed url
	// parameters:
	// - name: file list
	//   in: body
	//   description: File upload data
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UploadPresignedRequest"
	// responses:
	//   "200":
	//     description: Upload file(s) to AWS S3 bucket
	//     schema:
	//        "$ref": "#/definitions/ListFileResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.POST("/upload/presigned", http.presignedUpload)

}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}
	if lq == nil || lq.Filter == nil {
		return server.NewHTTPValidationError("Invalid query")
	}

	resp, err := h.service.ListFile(c.Request().Context(), lq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListFileResponse{resp})
}

func (h *HTTP) presignedUpload(c echo.Context) error {
	request := UploadPresignedRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	resp, err := h.service.UploadFilePreSigned(c.Request().Context(), request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListFileResponse{resp})
}
