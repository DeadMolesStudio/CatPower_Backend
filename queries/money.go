package queries

import (
	// "database/sql"

	"time"

	db "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/lib/pq"

	"CatPower/models"
)

func GetAllCategories(dm *db.DatabaseManager, id uint) ([]models.MoneyCategory, error) {
	dbo, err := dm.DB()
	if err != nil {
		return nil, err
	}
	var res []models.MoneyCategory
	if err = dbo.Select(&res,
		`SELECT category_id, name, is_income, pic, sum FROM money_category
		WHERE user_id = $1`,
		id); err != nil {
		return nil, err
	}

	return res, nil
}

func AddCategory(dm *db.DatabaseManager, c *models.MoneyCategory) error {
	dbo, err := dm.DB()
	if err != nil {
		return err
	}

	return dbo.Get(c,
		`INSERT INTO money_category (user_id, name, is_income)
		VALUES ($1, $2, $3) RETURNING category_id`,
		c.User, c.Name, c.IsIncome)
}

func AddCategoryPic(dm *db.DatabaseManager, id uint, forCat uint64, path string) error {
	dbo, err := dm.DB()
	if err != nil {
		return err
	}
	qres, err := dbo.Exec(`
		UPDATE money_category
		SET pic = $1
		WHERE user_id = $2 AND category_id = $3`,
		path, id, forCat)
	if err != nil {
		return err
	}
	res, err := qres.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrNotFound
	}

	return nil
}

func AddMoney(dm *db.DatabaseManager, op *models.MoneyOp) error {
	dbo, err := dm.DB()
	if err != nil {
		return err
	}

	if err := dbo.Get(op,
		`INSERT INTO money_action (action_uuid, user_id, delta, from_category, to_category)
		VALUES ($1, $2, $3, $4, $5) RETURNING added`,
		op.UUID, op.User, op.Delta, op.From, op.To); err != nil {
		switch err.(*pq.Error).Code {
		case "23505":
			return db.ErrUniqueConstraintViolation
		case "23503":
			return ErrForeignKeyViolation
		}
		return err
	}

	return nil
}

func AddMoneyPic(dm *db.DatabaseManager, id uint, forOp, path string) error {
	dbo, err := dm.DB()
	if err != nil {
		return err
	}
	qres, err := dbo.Exec(`
		UPDATE money_action
		SET photo = $1
		WHERE user_id = $2 AND action_uuid = $3`,
		path, id, forOp)
	if err != nil {
		return err
	}
	res, err := qres.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrNotFound
	}

	return nil
}

func GetMoneyHistory(dm *db.DatabaseManager, id uint, since *time.Time) ([]models.MoneyOp, error) {
	dbo, err := dm.DB()
	if err != nil {
		return nil, err
	}
	q := `SELECT delta, from_category, to_category, photo, added FROM money_action
		WHERE user_id = $1`
	args := []interface{}{id}
	if since != nil {
		q += ` AND added > $2`
		args = append(args, since)
	}
	var res []models.MoneyOp
	return res, dbo.Select(&res, q, args...)
}
