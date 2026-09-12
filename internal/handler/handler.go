package handler

import (
	"net/http"

	"github.com/ffajarpratama/pos-wash-api/config"
	"github.com/ffajarpratama/pos-wash-api/internal/middleware"
	"github.com/ffajarpratama/pos-wash-api/internal/usecase"
	"github.com/ffajarpratama/pos-wash-api/pkg/custom_validator"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	uc usecase.IFaceUsecase
	v  custom_validator.Validator
}

func NewV1Handler(cnf *config.Config, uc usecase.IFaceUsecase, v custom_validator.Validator) http.Handler {
	r := chi.NewRouter()
	h := handler{
		uc: uc,
		v:  v,
	}

	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", h.Register)
		auth.Post("/login", h.Login)
		auth.Route("/profile", func(profile chi.Router) {
			profile.Use(middleware.AuthenticateUser(cnf.JWTConfig.Secret))
			profile.Get("/", h.GetProfile)
		})
	})

	r.Group(func(private chi.Router) {
		private.Use(middleware.AuthenticateUser(cnf.JWTConfig.Secret))

		private.Route("/media", func(media chi.Router) {
			media.Post("/", h.CreateMedia)
		})

		private.Route("/dashboard", func(dashboard chi.Router) {
			dashboard.Get("/summary", h.GetDashboardSummary)
			dashboard.Get("/order-trend", h.GetOrderTrend)
		})

		private.Route("/outlet", func(outlet chi.Router) {
			outlet.Post("/", h.CreateOutlet)
			outlet.Get("/{outletID}", h.FindOneOutlet)
		})

		private.Route("/service-category", func(category chi.Router) {
			category.Get("/", h.FindAndCountServiceCategory)
		})

		private.Route("/service", func(service chi.Router) {
			service.Post("/", h.CreateService)
			service.Get("/", h.FindAndCountService)
			service.Get("/{serviceID}", h.FindOneService)
			service.Put("/{serviceID}", h.UpdateService)
			service.Delete("/{serviceID}", h.DeleteService)
		})

		private.Route("/customer", func(customer chi.Router) {
			customer.Post("/", h.CreateCustomer)
			customer.Get("/", h.FindAndCountCustomer)
			customer.Get("/{customerID}", h.FindOneCustomer)
			customer.Put("/{customerID}", h.UpdateCustomer)
			customer.Delete("/{customerID}", h.DeleteCustomer)
		})

		private.Route("/perfume", func(perfume chi.Router) {
			perfume.Get("/", h.FindAndCountPerfume)
		})

		private.Route("/payment-method", func(payment_method chi.Router) {
			payment_method.Get("/", h.FindAndCountPaymentMethod)
		})

		private.Route("/order", func(order chi.Router) {
			order.Post("/", h.CreateOrder)
			order.Get("/", h.FindAndCountOrder)
			order.Get("/{orderID}", h.FindOneOrder)
			order.Put("/{orderID}/status", h.UpdateOrderStatus)
			order.Post("/{orderID}/payment", h.OrderPayment)
		})
	})

	return r
}
