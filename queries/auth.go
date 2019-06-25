package queries

import (
	db "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"

	"CatPower/models"
)

func GetUserPassword(dm *db.DatabaseManager, e string) (*models.UserProfile, error) {
	dbo, err := dm.DB()
	if err != nil {
		return nil, err
	}
	res := &models.UserProfile{}
	if err = dbo.Get(res,
		`SELECT user_id, password FROM user_profile
		WHERE username = $1`,
		e); err != nil {
		return res, err
	}

	return res, nil
}
