package models

import "github.com/form3tech-oss/jwt-go"

var JwtKey = []byte("golang_my_todo")

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserInfo struct {
	Uid      int    `json:"uid" db:"uid"`
	Role     string `json:"role" db:"role"`
	Created  int    `json:"created" db:"created_by"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type JwtClaims struct {
	Uid      int    `json:"uid"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type Registration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Address struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Log     float64 `json:"log"`
}

type AllAddress struct {
	AddressId  int     `json:"addressId" db:"address_id"`
	Address    string  `json:"address" db:"address"`
	Lat        float64 `json:"lat" db:"lat"`
	Log        float64 `json:"log" db:"log"`
	TotalCount int     `json:"totalCount"`
}

type AllUser struct {
	Uid        int    `json:"uid" db:"uid"`
	Name       string `json:"name" bd:"name"`
	Email      string `json:"email" db:"email"`
	TotalCount int    `json:"totalCount"`
}

type AllSubAdmin struct {
	Uid        int    `json:"subAdminId" db:"subadminid"'`
	Name       string `json:"name" bd:"name"`
	Email      string `json:"email" db:"email"`
	TotalCount int    `json:"totalCount"`
}

type AllRestaurant struct {
	Rid        int     `json:"rid" db:"rid"`
	Name       string  `json:"name" db:"name"`
	Address    string  `json:"address" db:"address"`
	Lat        float64 `json:"lat" db:"lat"`
	Log        float64 `json:"log" db:"log"`
	TotalCount int     `json:"totalCount"`
}

type GetAllRestaurant struct {
	Rid        int     `json:"rid" db:"rid"`
	Name       string  `json:"name" db:"name"`
	CreatedBy  int     `json:"createdBy" db:"created_by"`
	Address    string  `json:"address" db:"address"`
	Lat        float64 `json:"lat" db:"lat"`
	Log        float64 `json:"log" db:"log"`
	TotalCount int     `json:"totalCount"`
}

type Restaurant struct {
	Name    string  `json:"name" db:"name"`
	Address string  `json:"address" db:"address"`
	Lat     float64 `json:"lat" db:"lat"`
	Log     float64 `json:"log" db:"log"`
}

type Dish struct {
	Name string `json:"name" db:"name"`
	Rate int    `json:"rate" db:"rate"`
}

type AllDishWithRestaurantId struct {
	Did        int    `json:"did" db:"did"`
	Name       string `json:"name" db:"name"`
	Rate       int    `json:"rate" db:"rate"`
	TotalCount int    `json:"totalCount"`
}

type AllDish struct {
	Did            int    `json:"did" db:"did"`
	Name           string `json:"name" db:"name"`
	Rate           int    `json:"rate" db:"rate"`
	RestaurantName string `json:"restaurantName"`
	TotalCount     int    `json:"totalCount"`
}

type Coordinates struct {
	Lat float64 `json:"lat" db:"lat"`
	Log float64 `json:"log" db:"log"`
}

type GetRid struct {
	CreatedBy int `json:"createdBy" db:"created_by"`
}
type Output struct {
	Message string `json:"message"`
}
type OutputDistance struct {
	Message string  `json:"message"`
	Km      float64 `json:"km"`
}
