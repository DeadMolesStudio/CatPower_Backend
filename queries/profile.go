package queries

import (
	"database/sql"

	"github.com/lib/pq"

	db "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"

	"CatPower/models"
)

func CheckExistenceOfEmail(dm *db.DatabaseManager, e string) (bool, error) {
	dbo, err := dm.DB()
	if err != nil {
		return false, err
	}
	if err = dbo.Get(&models.UserProfile{},
		`SELECT FROM user_profile
		WHERE email = $1`,
		e); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CheckExistenceOfUsername(dm *db.DatabaseManager, u string) (bool, error) {
	dbo, err := dm.DB()
	if err != nil {
		return false, err
	}
	if err = dbo.Get(&models.UserProfile{},
		`SELECT FROM user_profile
		WHERE username = $1`,
		u); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func CreateNewUser(dm *db.DatabaseManager, u *models.UserProfile) error {
	dbo, err := dm.DB()
	if err != nil {
		return err
	}
	err = dbo.Get(u,
		`INSERT INTO user_profile (email, password, username)
		VALUES ($1, $2, $3) RETURNING user_id`,
		u.Email, u.Password, u.Username)
	if err != nil {
		switch err.(*pq.Error).Code {
		case "23505":
			return db.ErrUniqueConstraintViolation
		}
		return err
	}

	return nil
}

func GetUserProfileByID(dm *db.DatabaseManager, id uint) (*models.UserProfile, error) {
	dbo, err := dm.DB()
	if err != nil {
		return nil, err
	}
	res := &models.UserProfile{}
	if err = dbo.Get(res,
		`SELECT user_id, email, username FROM user_profile 
		WHERE user_id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return res, ErrNotFound
		}
		return res, err
	}

	return res, nil
}
