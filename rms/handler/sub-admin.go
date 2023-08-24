package handler

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rms/database/dbhelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"
	"strconv"
)

func CreateRestaurantBySubAdmin(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: unauthorized")
		return
	}
	var uid float64
	uid, ok = claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "error: getting id")
		return
	}
	body := models.Restaurant{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong inputs")
		return
	}
	result, err := dbhelper.IsAddressExits(body.Address)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: in IsAddressExits")
		return
	} else if result == false {
		utils.RespondError(w, http.StatusBadRequest, nil, "same address is registered")
		return
	}
	err = dbhelper.CreateRestaurantBySubAdmin(body, int(uid))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error: CreateRestaurantBySubAdmin")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, models.Output{Message: "Successfully Created Restaurant"})
}

func GetAllRestaurantBySubAdmin(w http.ResponseWriter, r *http.Request) {
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
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error : unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "error: getting id")
		return
	}
	restaurant, err := dbhelper.AllRestaurantBySubAdmin(int(sid), limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve restaurant list")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Restaurant": restaurant,
	})
}

func CreateDish(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "error: getting id")
		return
	}
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong Restaurant ID")
		return
	}

	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	body := models.Dish{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong inputs")
		return
	}
	err = dbhelper.CreateDish(body, rid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error: CreateDish")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, models.Output{Message: "Successfully Created Dish"})
}

func GetAllDishBySubAdminWithRestaurantId(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "error : getting id")
		return
	}
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong Restaurant Id")
		return
	}

	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	sort := r.URL.Query().Get("Sort")
	dish, err := dbhelper.AllDishWithRestaurantId(rid, sort)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "GetAllDishBySubAdminWithRestaurantId: failed to retrieve dish list")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Dish": dish,
	})
}

func UpdateDish(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "dishId")
	did, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error: Wrong Dish Id")
		return
	}
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: getting Id")
		return
	}
	i = chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error : Wrong Restaurant Id")
		return
	}
	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	body := models.Dish{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong inputs")
		return
	}
	err = dbhelper.UpdateDish(body, did)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update dish")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "successfully Dish updated"})
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: getting Id")
		return
	}
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error: wrong restaurant Id")
		return
	}
	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	body := models.Restaurant{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong inputs")
		return
	}
	err = dbhelper.UpdateRestaurant(body, rid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: UpdateRestaurant")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Successfully restaurant updated"})

}

func DeleteDish(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "dishId")
	did, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error : wrong DishId")
		return
	}
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: getting Id")
		return
	}
	i = chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error : wrong Restaurant Id")
		return
	}
	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	err = dbhelper.DeleteDish(rid, did)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete dish")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Successfully Deleted your dish"})
}
func DeleteRestaurantBySubAdmin(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	sid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: getting Id")
		return
	}
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error: Wrong Restaurant ID")
		return
	}
	err = dbhelper.FindRestaurantOwner(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "Not Authorized to use this resource")
		return
	}
	err = dbhelper.DeleteRestaurantBySubAdmin(rid, int(sid))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Error: DeleteRestaurantSubAdmin")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Successfully Restaurant deleted"})

}
