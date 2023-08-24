package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
	"rms/database"
	"rms/database/dbhelper"
	"rms/models"
	"rms/utils"
	"strconv"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
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
	users, err := dbhelper.AllUser(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve users list")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Users": users,
	})
}

func CreateRestaurantByAdmin(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "sub-adminId")
	sid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong sub-adminId")
		return
	}
	body := models.Restaurant{}
	err = utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong inputs")
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
	err = dbhelper.CreateRestaurantByAdmin(body, sid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: CreateRestaurantByAdmin")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, models.Output{Message: "Successfully create a new restaurant"})
}

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "userId")
	uid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "wrong userId")
		return
	}
	tx := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.CreateSubAdmin(*tx, uid)
		if err != nil {
			return err
		}
		err = dbhelper.DeleteAllAddress(*tx, uid)
		if err != nil {
			return err
		}
		return nil
	})
	if tx != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx, "Error : while update to sub-admin")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "successfully upgrade to sub-admin"})
}

func GetAllSubAdmin(w http.ResponseWriter, r *http.Request) {
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
	sub, err := dbhelper.AllSubAdmin(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to retrieve sub-admin list ")
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"Sub-Admin": sub,
	})
}

func DeleteSubAdmin(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "sub-adminId")
	sid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong Sub-admin Id")
		return
	}
	result, err := dbhelper.IsSubAdmin(sid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: in IsSubAdmin")
		return
	} else if result == false {
		utils.RespondError(w, http.StatusBadRequest, nil, "something went wrong", "Check your Sub-admin Id")
		return
	}
	tx := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.DeleteSudAdmin(*tx, sid)
		if err != nil {
			return err
		}
		err = dbhelper.DeleteAllRestaurant(*tx, sid)
		if err != nil {
			return err
		}
		return nil
	})
	if tx != nil {
		utils.RespondError(w, http.StatusInternalServerError, tx, "Error : In Delete user")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Sub admin is deleted"})
}

func DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	i := chi.URLParam(r, "restaurantId")
	rid, err := strconv.Atoi(i)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Unable to convert string to int")
		return
	}
	err = dbhelper.DeleteRestaurant(rid)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete restaurant")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Successfully Deleted Restaurant"})
}

func RegisterSubAdmin(w http.ResponseWriter, r *http.Request) {
	body := models.Registration{}
	err := utils.ParseBody(r.Body, &body)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "Wrong Input is given")
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
	err = dbhelper.RegisterSubAdmin(body, hash)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "error: RegisterSubAdmin")
		return
	}
	utils.RespondJSON(w, http.StatusOK, models.Output{Message: "Registration is successful"})
}
