package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ffajarpratama/pos-wash-api/internal/http/request"
	"github.com/ffajarpratama/pos-wash-api/internal/model"
	"github.com/ffajarpratama/pos-wash-api/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateOrder implements IFaceRepository.
func (r *Repository) CreateOrder(ctx context.Context, data *model.Order, db *gorm.DB) error {
	return r.BaseRepository.Create(db.WithContext(ctx), data)
}

// FindAndCountOrder implements IFaceRepository.
func (r *Repository) FindAndCountOrder(ctx context.Context, params *request.ListOrderQuery) ([]*model.Order, int64, error) {
	var res = make([]*model.Order, 0)
	var cnt int64

	query := r.db.
		WithContext(ctx).
		Model(&model.Order{}).
		Table("tr_order o").
		Joins("JOIN tr_customer c ON c.customer_id = o.customer_id").
		Preload("Customer").
		Preload("OrderDetail").
		Preload("OrderDetail.Service")

	if params.OutletID != uuid.Nil {
		query = query.Where("o.outlet_id = ?", params.OutletID)
	}

	if params.Keyword != "" {
		query = query.Where("(c.name ILIKE ? OR c.phone_number ILIKE ? OR o.invoice_number ILIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status != "" {
		query = query.Where("o.status = ?", params.Status)
	}

	if params.Paid != "" {
		val, err := strconv.ParseBool(params.Paid)
		if err != nil {
			return nil, 0, err
		}

		if val {
			query = query.Where("o.paid_at IS NOT NULL")
		} else {
			query = query.Where("o.paid_at IS NULL")
		}
	}

	if err := query.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		query = query.Order(util.TransformSortClause("o.created_at", params.Sort))
	}

	if err := query.
		Limit(params.PerPage).
		Offset(util.CalculateOffset(params.Page, params.PerPage)).
		Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, cnt, nil
}

// FindOneOrder implements IFaceRepository.
func (r *Repository) FindOneOrder(ctx context.Context, query ...interface{}) (*model.Order, error) {
	var res *model.Order

	if err := r.BaseRepository.FindOne(
		r.db.
			WithContext(ctx).
			Where(query[0], query[1:]...).
			Preload("Customer").
			Preload("Customer.Avatar").
			Preload("Perfume").
			Preload("PaymentMethod").
			Preload("OrderDetail").
			Preload("OrderDetail.Service").
			Preload("OrderDetail.Service.Media").
			Preload("OrderHistoryStatus", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at DESC")
			}),
		&res,
	); err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateOrder implements IFaceRepository.
func (r *Repository) UpdateOrder(ctx context.Context, db *gorm.DB, data map[string]interface{}, query ...interface{}) error {
	return db.WithContext(ctx).Model(&model.Order{}).Where(query[0], query[1:]...).Updates(data).Error
}

// CountOrder implements IFaceRepository.
func (r *Repository) CountOrder(ctx context.Context, query ...interface{}) (int64, error) {
	var cnt int64

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where(query[0], query[1:]...).Count(&cnt).Error; err != nil {
		return 0, err
	}

	return cnt, nil
}

// GetOrderSummary implements IFaceRepository.
func (r *Repository) GetOrderSummary(ctx context.Context, params *request.OrderTrendQuery) (*model.OrderSummary, error) {
	var res *model.OrderSummary

	query := `SELECT COUNT(*) FILTER (WHERE status = 'accepted')						AS accepted,
					 COUNT(*) FILTER (WHERE status IN ('on-process', 'waiting-pickup')) AS on_process,
					 COUNT(*) FILTER (WHERE status = 'complete')						AS complete,

					 COALESCE(SUM(total_amount) FILTER (WHERE paid_at::date = current_date::date - '1 day'::interval), 0) AS rev_1,
					 COALESCE(SUM(total_amount) FILTER (WHERE paid_at::date = current_date::date), 0)                     AS rev_2
			  FROM tr_order
			  WHERE outlet_id = ?`

	args := []interface{}{params.OutletID}

	if params.Start != "" {
		query += ` AND created_at::date >= ?`
		args = append(args, params.Start)
	}

	if params.End != "" {
		query += ` AND created_at::date <= ?`
		args = append(args, params.End)
	}

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

// GetOrderTrend implements IFaceRepository.
func (r *Repository) GetOrderTrend(ctx context.Context, params *request.OrderTrendQuery) ([]*model.OrderTrend, error) {
	var res = make([]*model.OrderTrend, 0)

	arg1 := "TO_CHAR(paid_at, 'YYYY-MM')"
	arg2 := "TO_CHAR(range, 'YYYY-MM')"
	intv := "month"

	if params.Type == "weekly" {
		arg1 = "paid_at::date"
		arg2 = "range::date"
		intv = "day"
	}

	query := fmt.Sprintf(`
		WITH order_range AS (
			SELECT %s AS paid_at, sum(total_amount) AS total
		 	FROM tr_order
			WHERE paid_at IS NOT NULL AND paid_at::date BETWEEN ? AND ? AND outlet_id = ?
		 	GROUP BY %s
		)
		SELECT range AS date, COALESCE(odr.total, 0) AS total
		FROM GENERATE_SERIES(?, ?, '1 %s'::interval) AS range
		LEFT JOIN order_range odr ON odr.paid_at = %s;`, arg1, arg1, intv, arg2)

	args := []interface{}{
		params.Start,
		params.End,
		params.OutletID,
		params.Start,
		params.End,
	}

	if err := r.db.WithContext(ctx).Raw(query, args...).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}
