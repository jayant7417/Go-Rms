package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"rms/database"
	"rms/models"
)

func CreateAddress(body models.Address, uid int) error {
	s := `insert into 
    				address(uid, address, coordinates,lat,log)
			values 
					($1,$2,point($3,$4),$3,$4)`
	_, err := database.Rms.Exec(s, uid, body.Address, body.Lat, body.Log)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating user: %v", err)
		return err
	}
	return nil
}

func UpdateAddress(body models.Address, uid int, aid int) error {
	s := `update 
				address
		  set 
		      address =$1,
		      coordinates=point($2,$3),
		      lat = $2,
		      log = $3,
		      updated_at = now()
		  where
		      	uid=$4
				and address_id =$5
				and archived_at is null`
	_, err := database.Rms.Exec(s, body.Address, body.Lat, body.Log, uid, aid)
	if err != nil {
		logrus.Errorf("UpdateAddress: failed to update address: %v", err)
		return err
	}
	return nil
}

func AllAddress(uid, limit, offset int) ([]models.AllAddress, error) {
	body := make([]models.AllAddress, 0)
	s := `select
			address_id,
			address.address,
			lat,
			log,
			count(*) over() as totalCount
		from
			address
		where 
				uid =$1
				and archived_at is null 
		group by 
		    address_id
		order by 
		    	address_id
		limit $2
		offset $3`
	err := database.Rms.Select(&body, s, uid, limit, offset)
	if err != nil {
		logrus.Errorf("AllAddress: failed to retrieve all address: %v", err)
		return nil, err
	}
	return body, nil
}

func DeleteAddress(uid, aid int) error {
	s := `update
           		address
       	  set 	
       	      	archived_at = now()
		  where 
		      	uid=$1 
				and address_id=$2
		  		and archived_at is null `
	_, err := database.Rms.Exec(s, uid, aid)
	if err != nil {
		logrus.Errorf("DeleteAddress: failed to delete address: %v", err)
		return err
	}
	return nil
}

func DeleteUser(db sqlx.Tx, uid int) error {
	s := `update 
			users
		set 
		    archived_at = now()
		where 
		    archived_at is null
		    and uid =$1
		    and role='user' `
	_, err := db.Exec(s, uid)
	if err != nil {
		logrus.Errorf("Deleteuser: failed to delete user: %v", err)
		return err
	}
	return nil
}

func DeleteAllAddress(db sqlx.Tx, uid int) error {
	s := `update 
			address
		set 
		    archived_at = now()
		where 
		    archived_at is null
		    and uid =$1`
	_, err := db.Exec(s, uid)
	if err != nil {
		logrus.Errorf("DeleteAllAddress: failed to delete user: %v", err)
		return err
	}
	return nil
}

func RestaurantCoordinates(rid int) ([]models.Coordinates, error) {
	s := `  select
				lat,
				log
			from
				restaurant
			where 
			    rid =$1`
	body := make([]models.Coordinates, 0)
	err := database.Rms.Select(&body, s, rid)
	if err != nil {
		logrus.Errorf("RestaurantCoordinates : falied to fetch data: %v", err)
		return nil, err
	}
	return body, nil
}
