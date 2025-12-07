package utils

import (
	"database/sql"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
		panic(err)
	} else {
		errorCommit := tx.Commit()
		if errorCommit != nil {
			panic(errorCommit)
		}
	}
}
