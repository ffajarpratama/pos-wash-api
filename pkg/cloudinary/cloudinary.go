package cloudinary

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	go_cloudinary "github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/ffajarpratama/pos-wash-api/config"
	"github.com/ffajarpratama/pos-wash-api/pkg/constant"
)

type Cloudinary struct {
	cnf      *config.Config
	instance *go_cloudinary.Cloudinary
}

type CloudinaryUploadRes struct {
	PublicID  string `json:"public_id"`
	SecureURL string `json:"secure_url"`
}

func NewClient(cnf *config.Config) (*Cloudinary, error) {
	cld, err := go_cloudinary.NewFromParams(cnf.Cloudinary.CloudName, cnf.Cloudinary.APIKey, cnf.Cloudinary.APISecret)
	if err != nil {
		return nil, err
	}

	log.Println("[cloudinary-connected]")

	return &Cloudinary{
		cnf:      cnf,
		instance: cld,
	}, nil
}

func (c *Cloudinary) UploadImage(ctx context.Context, file multipart.File, purpose constant.UploadPurpose, assetType string) (*CloudinaryUploadRes, error) {
	opts := uploader.UploadParams{
		ResourceType: assetType,
		Folder:       fmt.Sprintf("/%s", purpose),
	}

	res, err := c.instance.Upload.Upload(ctx, file, opts)
	if err != nil {
		return nil, err
	}

	resp := &CloudinaryUploadRes{
		PublicID:  res.PublicID,
		SecureURL: res.SecureURL,
	}

	return resp, nil
}
