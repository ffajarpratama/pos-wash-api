package request

import (
	"mime/multipart"

	"github.com/ffajarpratama/pos-wash-api/pkg/constant"
)

type CreateMedia struct {
	File      multipart.File
	Header    *multipart.FileHeader
	Purpose   constant.UploadPurpose
	Filename  string
	Mimetype  string
	AssetType string
}
