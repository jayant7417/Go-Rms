package dbhelper

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"rms/database"
	"rms/models"
	"rms/utils"
)

func GetAllRestaurant(limit, offset int, search string) ([]models.GetAllRestaurant, error) {
	body := make([]models.GetAllRestaurant, 0)
	s := `select 
			rid,
			name,
			address,
			lat,
			log,
			count(*) over () as totalCount
		from
		    restaurant
		where
		    archived_at is null
			and name like '%'||$3||'%'
		group by 
		    rid
		order by 
		    rid
		limit  $1
		offset $2`
	err := database.Rms.Select(&body, s, limit, offset, search)
	if err != nil {
		logrus.Errorf("AllRestaurant: failed to retrieve: %v", err)
		return nil, err
	}
	return body, nil
}

func CreateUser(body models.Registration, hash string) error {
	s := `insert into 
    				users(name, email, password)
		values 
		    ($1,TRIM(LOWER($2)),$3)`
	_, err := database.Rms.Exec(s, body.Name, body.Email, hash)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating user: %v", err)
		return err
	}
	return nil
}

func IsEmailExits(email string) (bool, error) {
	s := `  select 
      			count(*)=0
			from 
			    users
			where 
			    email = $1 
			  	and archived_at is null`
	var result bool
	err := database.Rms.QueryRow(s, email).Scan(&result)
	if err != nil {
		logrus.Errorf("IsEmailExits: failed to check email:%v", err)
		return false, err
	}
	return result, nil
}

func RetrieveInfo(body models.Login) (models.UserInfo, error) {
	s := `SELECT
			   uid,
               role,
               name,
               email,
               password

		  from
		      	users
		  where
		      	email = TRIM(LOWER($1))
		      	and archived_at IS NULL `

	userDetails := models.UserInfo{}
	err := database.Rms.Get(&userDetails, s, body.Email)
	if err != nil {
		err = errors.New("failed to login")
		return userDetails, err
	}
	err = utils.CheckPassword(body.Password, userDetails.Password)
	if err != nil {
		err = errors.New("failed to login")
		return userDetails, err
	}
	return userDetails, nil
}

func AllDishWithRestaurantId(rid int, sort string) ([]models.AllDishWithRestaurantId, error) {
	sql := `  select 
      			did,
      			name,
      			rate,
      			count(*) over () as totalCount
			from 
			    dishes
			where 
			    archived_at is null
			    and rid=$1`
	if sort != "" {
		sql = fmt.Sprintf("%s order by %s", sql, sort)
	}
	body := make([]models.AllDishWithRestaurantId, 0)
	err := database.Rms.Select(&body, sql, rid)
	if err != nil {
		logrus.Errorf("AllDishWithRestaurantId: falied to show all dish: %v", err)
		return nil, err
	}
	return body, nil
}

func IsAddressExits(address string) (bool, error) {
	s := `  select 
      			count(*)=0
			from 
			    restaurant
			where 
			    address = $1 
			  	and archived_at is null`
	var result bool
	err := database.Rms.QueryRow(s, address).Scan(&result)
	if err != nil {
		logrus.Errorf("IsAddressExits: failed to check address: %v", err)
		return false, err
	}
	return result, nil
}

func AllDish(limit, offset int, search, sort string) ([]models.AllDish, error) {
	sql := `  select
      			d.did as did,
      			d.name as name ,
      			d.rate as rate,
                r.name as restaurantName,
                count(*) over () as totalCount
			from
			    dishes d
			left join
			        restaurant r
			            on r.rid = d.rid
			where
			    d.archived_at is null
			`
	if search != "" {
		sql = fmt.Sprintf("%s AND d.name like '%%%s%%'", sql, search)
	}
	if sort != "" {
		sql = fmt.Sprintf("%s group by d.did,r.name order by did %s", sql, sort)
	}
	if limit > 0 {
		sql = fmt.Sprintf("%s limit %v", sql, limit)
	}
	if offset > 0 {
		sql = fmt.Sprintf("%s offset %v", sql, offset)
	}
	body := make([]models.AllDish, 0)
	err := database.Rms.Select(&body, sql)
	if err != nil {
		logrus.Errorf("AllDish: falied to show all dish: %v", err)
		return nil, err
	}
	return body, nil
}
