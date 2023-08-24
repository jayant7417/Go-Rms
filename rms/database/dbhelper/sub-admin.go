package dbhelper

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"rms/database"
	"rms/models"
)

func CreateRestaurantBySubAdmin(body models.Restaurant, uid int) error {
	s := `insert into
				restaurant(name,address, lat, log,created_by)
		values
				($1,$2,$3,$4,$5)`
	_, err := database.Rms.Exec(s, body.Name, body.Address, body.Lat, body.Log, uid)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating restaurant: %v", err)
		return err
	}
	return nil
}

func AllRestaurantBySubAdmin(sid, limit, offset int) ([]models.AllRestaurant, error) {

	body := make([]models.AllRestaurant, 0)
	s := `select 
			rid,
			name,
			address,
			lat,
			log,
			count() over () as totalCount
		from
		    restaurant
		where
		    archived_at is null
			and created_by = $1
		group by 
		    rid
		order by 
		    rid
		limit $1
		offset $2`
	err := database.Rms.Select(&body, s, sid, limit, offset)
	if err != nil {
		logrus.Errorf("AllRestaurant: failed to retrieve : %v", err)
		return nil, err
	}
	return body, nil
}

func FindRestaurantOwner(rid int, sid int) error {
	s := `  select 
    			created_by
			from 
			    restaurant 
			where 
			    rid=$1 
			  	and archived_at is null`
	body := make([]models.GetRid, 0)
	err := database.Rms.Select(&body, s, rid)
	if err != nil {
		logrus.Errorf("CreateDish : falied to create dish: %v", err)
		return err
	}
	if body[0].CreatedBy == sid {
		return nil

	}
	return fmt.Errorf("error")
}

func CreateDish(body models.Dish, rid int) error {
	s := `  insert into 
				dishes(rid,name,rate)
			values 
				($1,$2,$3)`
	_, err := database.Rms.Exec(s, rid, body.Name, body.Rate)
	if err != nil {
		logrus.Errorf("CreateDish : falied to create dish: %v", err)
		return err
	}
	return nil
}

func UpdateDish(body models.Dish, did int) error {
	s := ` 	update
 				dishes
			set
				name=$1,
				rate=$2,
				updated_at = now()
			where 
			    did=$3
			    and archived_at is null`
	_, err := database.Rms.Exec(s, body.Name, body.Rate, did)
	if err != nil {
		logrus.Errorf("UpdateDish: Error in update dish: %v", err)
		return err
	}
	return nil
}

func UpdateRestaurant(body models.Restaurant, rid int) error {
	s := `update 
				restaurant
			set
				name =$1,
				address=$2,
				lat=$3,
				log=$4,
				updated_at = now()
			where 
			    rid=$5
			    and archived_at  is null`
	_, err := database.Rms.Exec(s, body.Name, body.Address, body.Lat, body.Log, rid)
	if err != nil {
		logrus.Errorf("UpdateRestaurant: Error in update restaurant: %v", err)
		return err
	}
	return nil
}

func DeleteDish(rid, did int) error {
	s := `update
           		dishes
       	  set 	
       	      	archived_at = now()
		  where 
		      	rid=$1 
				and did=$2
		  		and archived_at is null `
	_, err := database.Rms.Exec(s, rid, did)
	if err != nil {
		logrus.Errorf("DeleteDish: failed to delete dish: %v", err)
		return err
	}
	return nil
}

func DeleteRestaurantBySubAdmin(rid, sid int) error {
	s := `update
           		restaurant
       	  set 	
       	      	archived_at = now()
		  where 
		      	rid=$1 
				and created_by=$2
		  		and archived_at is null `
	_, err := database.Rms.Exec(s, rid, sid)
	if err != nil {
		logrus.Errorf("DeleteDish: failed to delete Restaurant: %v", err)
		return err
	}
	return nil
}
