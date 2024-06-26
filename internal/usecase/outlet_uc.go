package usecase

import (
	"context"

	"github.com/ffajarpratama/pos-wash-api/internal/http/request"
	"github.com/ffajarpratama/pos-wash-api/internal/model"
	"github.com/ffajarpratama/pos-wash-api/pkg/util"
	"github.com/google/uuid"
)

// CreateOutlet implements IFaceUsecase.
func (u *Usecase) CreateOutlet(ctx context.Context, req *request.CreateOutlet) error {
	tx := u.DB.Begin()
	defer tx.Rollback()

	outlet := &model.Outlet{
		Name:    req.Name,
		Code:    util.GenerateRandomString(5, true),
		Address: req.Address,
		LogoID:  req.LogoID,
	}

	err := u.Repo.CreateOutlet(ctx, outlet, tx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

// FindOneOutlet implements IFaceUsecase.
func (u *Usecase) FindOneOutlet(ctx context.Context, outletID uuid.UUID) (*model.Outlet, error) {
	return u.Repo.FindOneOutlet(ctx, "outlet_id = ?", outletID)
}
