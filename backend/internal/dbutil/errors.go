package dbutil

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func ProcessInsertError(err error, alreadyExistsType, missingFieldType error) error {
	if err == nil {
		return nil
	}

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			msg := sqliteErr.Error()
			switch {
			case strings.Contains(msg, "UNIQUE constraint failed"):
				return alreadyExistsType
			case strings.Contains(msg, "NOT NULL constraint failed"):
				parts := strings.Split(msg, ":")
				if len(parts) == 2 {
					field := strings.TrimSpace(parts[1])
					return errors.Join(missingFieldType, errors.New(field))
				}
				return missingFieldType
			default:
				return err
			}
		}
	}

	return err
}

func ProcessRowError(err error, notFoundType error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return notFoundType
	}
	return err

}
