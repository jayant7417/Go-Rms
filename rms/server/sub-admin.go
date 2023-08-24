package server

import (
	"github.com/go-chi/chi/v5"
	"rms/handler"
)

func sub(r chi.Router) {
	r.Group(func(s chi.Router) {
		s.Post("/", handler.CreateRestaurantBySubAdmin)
		s.Get("/", handler.GetAllRestaurantBySubAdmin)
		s.Put("/{restaurantId}", handler.UpdateRestaurant)
		s.Delete("/{restaurantId}", handler.DeleteRestaurantBySubAdmin)
		s.Post("/{restaurantId}/dish", handler.CreateDish)
		s.Get("/{restaurantId}/dish", handler.GetAllDishBySubAdminWithRestaurantId)
		s.Put("/{restaurantId}/dish/{dishId}", handler.UpdateDish)
		s.Delete("/{restaurantId}/dish/{dishId}", handler.DeleteDish)
	})
}
