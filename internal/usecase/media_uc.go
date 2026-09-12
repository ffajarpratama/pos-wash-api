package usecase

import (
	"context"
	"path/filepath"

	"github.com/ffajarpratama/pos-wash-api/internal/http/request"
	"github.com/ffajarpratama/pos-wash-api/internal/model"
)

// CreateMedia implements IFaceUsecase.
func (u *Usecase) CreateMedia(ctx context.Context, req *request.CreateMedia) (*model.Media, error) {
	res, err := u.Uploader.UploadImage(ctx, req.File, req.Purpose, req.AssetType)
	if err != nil {
		return nil, err
	}

	media := &model.Media{
		Name:     req.Filename,
		Path:     res.PublicID + filepath.Ext(req.Header.Filename),
		Size:     int(req.Header.Size),
		Mimetype: req.Mimetype,
		MediaURL: res.SecureURL,
	}

	err = u.Repo.CreateMedia(ctx, media, u.DB)
	if err != nil {
		return nil, err
	}

	return media, nil
}
