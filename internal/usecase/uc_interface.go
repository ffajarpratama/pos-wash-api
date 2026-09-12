package usecase

import (
	"context"

	"github.com/ffajarpratama/pos-wash-api/internal/http/request"
	"github.com/ffajarpratama/pos-wash-api/internal/http/response"
	"github.com/ffajarpratama/pos-wash-api/internal/model"
	"github.com/google/uuid"
)

type IFaceUsecase interface {
	// auth
	Register(ctx context.Context, req *request.Register) (*model.User, error)
	Login(ctx context.Context, req *request.Login) (*model.User, error)
	FindOneUser(ctx context.Context, userID uuid.UUID) (*model.User, error)

	// outlet
	CreateOutlet(ctx context.Context, req *request.CreateOutlet) error
	FindOneOutlet(ctx context.Context, outletID uuid.UUID) (*model.Outlet, error)

	// service category
	FindAndCountServiceCategory(ctx context.Context, params *request.ListServiceCategoryQuery) ([]*model.ServiceCategory, int64, error)

	// service
	CreateService(ctx context.Context, req *request.CreateService) error
	FindAndCountService(ctx context.Context, params *request.ListServiceQuery) ([]*model.Service, int64, error)
	FindOneService(ctx context.Context, serviceID uuid.UUID) (*model.Service, error)
	UpdateService(ctx context.Context, req *request.UpdateService) error
	DeleteService(ctx context.Context, serviceID uuid.UUID) error

	// customer
	CreateCustomer(ctx context.Context, req *request.CreateCustomer) error
	FindAndCountCustomer(ctx context.Context, params *request.ListCustomerQuery) ([]*model.Customer, int64, error)
	FindOneCustomer(ctx context.Context, customerID uuid.UUID) (*model.Customer, error)
	UpdateCustomer(ctx context.Context, req *request.UpdateCustomer) error
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) error

	// perfume
	FindAndCountPerfume(ctx context.Context, params *request.ListPerfumeQuery) ([]*model.Perfume, int64, error)

	// payment method
	FindAndCountPaymentMethod(ctx context.Context, params *request.ListPaymentMethodQuery) ([]*model.PaymentMethod, int64, error)

	// order
	CreateOrder(ctx context.Context, req *request.CreateOrder) (*model.Order, error)
	FindAndCountOrder(ctx context.Context, params *request.ListOrderQuery) ([]*model.Order, int64, error)
	FindOneOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, req *request.UpdateOrderStatus) error
	OrderPayment(ctx context.Context, req *request.OrderPayment) error

	// dashboard
	GetDashboardSummary(ctx context.Context, outletID uuid.UUID) (*response.DashoardSummary, error)
	GetOrderTrend(ctx context.Context, params *request.OrderTrendQuery) (*response.OrderTrend, error)

	// media
	CreateMedia(ctx context.Context, req *request.CreateMedia) (*model.Media, error)
}
