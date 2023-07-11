package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// parsePgError returns parsed pgconn.PgError.
// If err is not pgconn.PgError, returns the same err.
func parsePgError(err error) error {
	var pgErr *pgconn.PgError
	if err, ok := err.(*pgconn.PgError); ok {
		pgErr = err
		return fmt.Errorf("database error. message: %s, detail: %s, where: %s, sqlstate: %s",
			pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.SQLState())
	}
	return err
}

func ErrCommit(err error) error {
	return fmt.Errorf("failed to commit Tx due to error: %v", err)
}

func ErrRollback(err error) error {
	return fmt.Errorf("failed to rollback Tx due to error: %v", err)
}

func ErrCreateTx(err error) error {
	return fmt.Errorf("failed to create Tx due to error: %v", err)
}

func ErrCreateQuery(err error) error {
	return fmt.Errorf("failed to create SQL Query due to error: %v", parsePgError(err))
}

func ErrScan(err error) error {
	return fmt.Errorf("failed to scan due to error: %v", parsePgError(err))
}

func ErrDoQuery(err error) error {
	return fmt.Errorf("failed to query due to error: %v", parsePgError(err))
}

func ErrExecQuery(err error) error {
	return fmt.Errorf("failed to  execute query due to error: %v", parsePgError(err))
}
