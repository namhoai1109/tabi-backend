package model

import (
	"time"

	httpcore "github.com/namhoai1109/tabi/core/http"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

type Base struct {
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
	// The time that record is deleted
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ListRequest struct {
	httpcore.ListRequest
}

type ErrDetailsResp struct {
	//in: body
	Body server.ErrorResponse
}
