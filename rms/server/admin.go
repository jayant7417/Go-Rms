package server

import (
	"github.com/go-chi/chi/v5"
	"rms/handler"
)

func adminUser(r chi.Router) {
	r.Group(func(a chi.Router) {
		a.Get("/user", handler.GetAllUser)
		a.Post("/user/{userId}/sub-admin", handler.CreateSubAdmin)
	})
}

func adminSubAdmin(r chi.Router) {
	r.Group(func(a chi.Router) {
		a.Get("/sub-admin", handler.GetAllSubAdmin)
		a.Post("/sub-admin/register", handler.RegisterSubAdmin)
		a.Delete("/sub-admin/{sub-adminId}", handler.DeleteSubAdmin)
		a.Post("/sub-admin/{sub-adminId}/restaurant", handler.CreateRestaurantByAdmin)
	})
}

func adminRestaurant(r chi.Router) {
	r.Group(func(a chi.Router) {
		a.Get("/", handler.GetAllRestaurant)
		a.Delete("/{restaurantId}", handler.DeleteRestaurant)
		a.Get("/{restaurantId}/dish", handler.GetAllDishWithRestaurantId)
		a.Get("/dish", handler.GetAllDish)
	})
}
