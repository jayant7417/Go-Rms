package middlewares

import (
	"context"
	"errors"
	"github.com/form3tech-oss/jwt-go"
	"net/http"
	"rms/models"
	"rms/utils"
	"time"
)

var UserContext = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, err := validateToken(w, r)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, err, "goto login page", "Error in token either expire or not received token")
			return
		}
		ctx := context.WithValue(r.Context(), UserContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(_ http.ResponseWriter, r *http.Request) (c jwt.MapClaims, err error) {

	//if r.Header["X-Api-Key"] == nil {
	//	utils.RespondError(w, http.StatusUnauthorized, nil, "unable to validate")
	//	return
	//}
	tokenstring := r.Header.Get("x-api-key")
	cl := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenstring, cl, func(token *jwt.Token) (interface{}, error) {
		return models.JwtKey, nil
	})

	//if errTkn != nil {
	//	utils.RespondError(w, http.StatusUnauthorized, err, "token expired")
	//	return nil, errTkn
	//}
	//token, err := jwt.Parse(r.Header["X-Api-Key"][0], func(token *jwt.Token) (interface{}, error) {
	//	if err != nil {
	//		return nil, err
	//	}
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("there was an error in parsing")
	//	}
	//	return models.JwtKey, nil
	//})

	if token == nil {
		err = errors.New("invalid token")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token error")
	}
	exp := claims["exp"]
	timeExp, ok := exp.(float64)
	if !ok {
		return nil, nil
	}

	if int64(timeExp) < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContext).(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusInternalServerError, nil, "Error : unauthorized")
			return
		}
		role, ok := claims["role"]
		role = role.(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use admin rights")
			return
		}
		if role != "admin" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use admin rights")
			return
		}
		ctx := context.WithValue(r.Context(), UserContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthSubAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContext).(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusInternalServerError, nil, "Error : Unauthorized")
			return
		}
		role, ok := claims["role"]
		role = role.(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use sub-admin rights")
			return
		}
		if role != "sub-admin" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use sub-admin rights")
			return
		}
		ctx := context.WithValue(r.Context(), UserContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContext).(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusInternalServerError, nil, "Error : unauthorized")
			return
		}
		role, ok := claims["role"]
		role = role.(string)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use user rights")
			return
		}
		if role != "user" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "unable to use user rights")
			return
		}
		ctx := context.WithValue(r.Context(), UserContext, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
