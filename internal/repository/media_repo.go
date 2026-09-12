package repository

import (
	"context"

	"github.com/ffajarpratama/pos-wash-api/internal/model"
	"gorm.io/gorm"
)

// CreateMedia implements IFaceRepository.
func (r *Repository) CreateMedia(ctx context.Context, data *model.Media, db *gorm.DB) error {
	return r.BaseRepository.Create(db.WithContext(ctx), data)
}
