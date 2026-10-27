package util

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		PanicIfError(errorRollback)
		panic(err)
	} else {
		errorCommit := tx.Commit()
		PanicIfError(errorCommit)
	}
}

func CheckDuplicateWithComparator(tx *gorm.DB, dest any, query string, args []any, compareFn func() []string) error {

	err := tx.Where(query, args...).First(dest).Error

	if err == nil {
		fields := compareFn()
		return &DuplicateError{Fields: fields}
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
