package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"rms/database"
	"rms/models"
)

func AllUser(limit, offset int) ([]models.AllUser, error) {
	body := make([]models.AllUser, 0)
	s := `select 
			uid,
			name,
			email,
			count(*) over () as totalCount
		from
			users
		where 
		    role='user'
			and archived_at is null 
		group by 
		    uid
		order by 
		    uid  
		limit $1
		offset $2`
	err := database.Rms.Select(&body, s, limit, offset)
	if err != nil {
		logrus.Errorf("AllUser: failed to retrieve users : %v", err)
		return nil, err
	}
	return body, nil
}

func CreateRestaurantByAdmin(body models.Restaurant, sub int) error {
	s := `insert into
				restaurant(name,address, lat, log,created_by)
		values
				($1,$2,$3,$4,$5)`
	_, err := database.Rms.Exec(s, body.Name, body.Address, body.Lat, body.Log, sub)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating restaurant: %v", err)
		return err
	}
	return nil
}

func CreateSubAdmin(db sqlx.Tx, uid int) error {
	s := `	update	
				users
			set 
			    role='sub-admin'
			where
			    uid=$1
			    and archived_at is null `
	_, err := db.Exec(s, uid)
	if err != nil {
		logrus.Errorf("CreateSubAdmin : failed to creating sud-admin: %v", err)
		return err
	}
	return nil
}

func AllSubAdmin(limit, offset int) ([]models.AllSubAdmin, error) {
	body := make([]models.AllSubAdmin, 0)
	s := `select 
			uid as subAdminId,
			name,
			email,
			count(*) over() as totalCount
		from
			users
		where 
		    role='sub-admin'
			and archived_at is null
		group by 
		    uid
		order by 
		    subAdminId 
		limit $1
		offset $2`
	err := database.Rms.Select(&body, s, limit, offset)
	if err != nil {
		logrus.Errorf("AllSubAdmin: failed to retrieve sub-users: %v", err)
		return nil, err
	}
	return body, nil
}

func DeleteSudAdmin(db sqlx.Tx, sid int) error {
	s := `update 
			users
		set 
		    archived_at = now()
		where 
		    archived_at is null
		    and uid =$1
		    and role='sub-admin' `
	_, err := db.Exec(s, sid)
	if err != nil {
		logrus.Errorf("Deleteuser: failed to delete user: %v", err)
		return err
	}
	return nil
}

func DeleteAllRestaurant(db sqlx.Tx, sid int) error {
	s := `update 
			restaurant
		set 
		    archived_at = now()
		where 
		    archived_at is null
		    and created_by = $1 `
	_, err := db.Exec(s, sid)
	if err != nil {
		logrus.Errorf("DeleteAllRestaurant: failed to Delete All Restaurant: %v", err)
		return err
	}
	return nil
}

func DeleteRestaurant(rid int) error {
	s := `update 
			restaurant
		set 
		    archived_at = now()
		where 
		    archived_at is null
		    and rid=$1`
	_, err := database.Rms.Exec(s, rid)
	if err != nil {
		logrus.Errorf("DeleteRestaurant: failed to Delete Restaurant: %v", err)
		return err
	}
	return nil
}

func RegisterSubAdmin(body models.Registration, hash string) error {
	s := `insert into 
    				users(name, email, password,role)
		values 
		    ($1,TRIM(LOWER($2)),$3,$4)`
	role := "sub-admin"
	_, err := database.Rms.Exec(s, body.Name, body.Email, hash, role)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating Sud-admin: %v", err)
		return err
	}
	return nil
}
func IsSubAdmin(sid int) (bool, error) {
	s := `  select 
      			role='sub-admin' as role,
      			archived_at isnull as date
  				
			from 
			    users
			where 
			    uid=$1`
	var role bool
	var date bool
	err := database.Rms.QueryRowx(s, sid).Scan(&role, &date)
	if err != nil {
		logrus.Errorf("IsEmailExits: failed to check email:%v", err)
		return false, err
	}
	if role == date {
		if role == true {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}
