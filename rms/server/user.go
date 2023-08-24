package server

import (
	"github.com/go-chi/chi/v5"
	"rms/handler"
)

func userAddress(a chi.Router) {
	a.Group(func(u chi.Router) {
		u.Post("/add", handler.AddAddress)
		u.Get("/add", handler.GetAllAddress)
		u.Put("/add/{addressId}", handler.UpdateAddress)
		u.Delete("/add/{addressId}", handler.DeleteAddress)
	})
}

func userRestaurant(r chi.Router) {
	r.Group(func(u chi.Router) {
		u.Get("/", handler.GetAllRestaurant)
		u.Get("/{restaurantId}", handler.FindDistance)
		u.Get("/{restaurantId}/dish", handler.GetAllDishWithRestaurantId)
		u.Get("/dish", handler.GetAllDish)
	})
}
