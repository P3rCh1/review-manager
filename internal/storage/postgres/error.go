package postgres

import (
	"errors"

	"github.com/lib/pq"
)

func isUniqueViolation(err error) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
