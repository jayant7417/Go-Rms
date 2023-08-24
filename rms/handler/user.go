package handler

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/umahmood/haversine"
	"net/http"
	"rms/database"
	"rms/database/dbhelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"
	"strconv"
)

func AddAddress(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: unauthorized")
		return
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to get Id")
		return
	}
	body := models.Address{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong inputs")
		return
	}
	err = dbhelper.CreateAddress(body, int(uid))
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: in create address")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Address is successful added"})
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "error: unauthorized")
		return
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to get ID")
		return
	}
	i := chi.URLParam(r, "addressId")
	aid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong address Id")
		return
	}
	body := models.Address{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong inputs")
		return
	}
	err = dbhelper.UpdateAddress(body, int(uid), aid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: UpdateAddress")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Address is updated"})
}

func GetAllAddress(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to get Id")
		return
	}

	address, err := dbhelper.AllAddress(int(uid), limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve address list")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Address": address,
	})
}

func DeleteAddress(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to get Id")
		return
	}
	i := chi.URLParam(r, "addressId")
	aid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong address Id")
		return
	}
	err = dbhelper.DeleteAddress(int(uid), aid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete address")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Address is deleted"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.UserContext).(jwt.MapClaims)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Error: Unauthorized")
		return
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to get Id")
		return
	}
	txerr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.DeleteAllAddress(*tx, int(uid))
		if err != nil {
			return err
		}
		err = dbhelper.DeleteUser(*tx, int(uid))
		if err != nil {
			return err
		}
		return nil
	})
	if txerr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txerr, "Error : In DeleteUser")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "User is successfully deleted"})
}

func FindDistance(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Error: Wrong restaurant Id")
		return
	}
	body := models.Coordinates{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong inputs")
		return
	}
	restaurant, err := dbhelper.RestaurantCoordinates(rid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve coordinates")
		return
	}
	rest := haversine.Coord{Lat: restaurant[0].Lat, Lon: restaurant[0].Log}
	u := haversine.Coord{Lat: body.Lat, Lon: body.Log}
	_, km := haversine.Distance(rest, u)

	utils.RespondJSON(w, http.StatusOK, models.OutputDistance{
		Message: "The distance is :",
		Km:      km,
	})
}
