package handler

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rms/database/dbhelper"
	"rms/models"
	"rms/utils"
	"strconv"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body := models.Login{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong input is given")
		return
	}
	userinfo := models.UserInfo{}
	var token string
	userinfo, err = dbhelper.RetrieveInfo(body)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Error : RetrieveInfo")
		return
	}
	token, err = generateJWT(userinfo)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to login")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"Jwt": token,
	})
}

func generateJWT(body models.UserInfo) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = body.Email
	claims["uid"] = body.Uid
	claims["name"] = body.Name
	claims["role"] = body.Role
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	tokenString, err := token.SignedString(models.JwtKey)
	if err != nil {
		fmt.Printf("something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Registration(w http.ResponseWriter, r *http.Request) {
	body := models.Registration{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong input is given")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password is short")
		return
	}
	result, err := dbhelper.IsEmailExits(body.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: in IsEmailExits")
		return
	} else if result == false {
		utils.RespondError(w, http.StatusBadRequest, nil, "same email is registered")
		return
	}
	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error: In password ")
		return
	}
	err = dbhelper.CreateUser(body, hash)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: CreateUser")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Registration is successful"})
}

func GetAllRestaurant(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Query().Get("Limit")
	limit, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong limit")
		return
	}
	max := 20
	if limit > max {
		utils.RespondError(w, http.StatusBadRequest, nil, "Something went wrong", "limit is surpass the max")
		return
	}
	i = r.URL.Query().Get("Offset")
	offset, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong offset")
		return
	}
	search := r.URL.Query().Get("Search")
	restaurant, err := dbhelper.GetAllRestaurant(limit, offset, search)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve restaurant list")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Restaurant": restaurant,
	})
}

func GetAllDishWithRestaurantId(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error : Wrong restaurant ID")
		return
	}
	sort := r.URL.Query().Get("Sort")
	dish, err := dbhelper.AllDishWithRestaurantId(rid, sort)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve dish list", "GetAllDishWithRestaurantId: Error in this function ")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Dish": dish,
	})
}
func GetAllDish(w http.ResponseWriter, r *http.Request) {
	i := r.URL.Query().Get("Limit")
	limit, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong limit")
		return
	}
	max := 20
	if limit > max {
		utils.RespondError(w, http.StatusBadRequest, nil, "Something went wrong", "limit is surpass the max")
		return
	}
	i = r.URL.Query().Get("Offset")
	offset, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong offset")
		return
	}
	search := r.URL.Query().Get("Search")
	sort := r.URL.Query().Get("Sort")
	dish, err := dbhelper.AllDish(limit, offset, search, sort)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve dish list", "GetAllDish: Error in this function")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Dish": dish,
	})

}
