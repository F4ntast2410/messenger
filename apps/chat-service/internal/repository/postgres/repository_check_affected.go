package repository

import (
	"database/sql"
	"errors"
)

func checkAffected(result sql.Result, notFoundMsg string) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New(notFoundMsg)
	}
	return nil
}
