package file

import (
	"tabi-file/config"

	dbcore "github.com/namhoai1109/tabi/core/db"

	"gorm.io/gorm"
)

func New(db *gorm.DB,
	fileDB FileDB,
	s3 S3Service,
	cfg *config.Configuration,
) *File {
	return &File{
		db:     db,
		fileDB: fileDB,
		s3:     s3,
		cfg:    cfg,
	}
}

type File struct {
	db     *gorm.DB
	fileDB FileDB
	s3     S3Service
	cfg    *config.Configuration
}

type FileDB interface {
	dbcore.Intf
}

type S3Service interface {
	PreparePresignedURL(key, bucketName string, expireTime int) (*string, error)
	GetSignedURL(key, bucketName string, expireTime int) (*string, error)
}
