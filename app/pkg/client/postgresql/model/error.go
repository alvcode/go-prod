package db

import "fmt"

func ErrCommit(err error) error {
	return fmt.Errorf("failed to commit Tx: %w", err)
}

func ErrRollback(err error) error {
	return fmt.Errorf("failed to rollback Tx: %w", err)
}

func ErrCreateTx(err error) error {
	return fmt.Errorf("failed to create Tx: %w", err)
}

func ErrCreateQuery(err error) error {
	return fmt.Errorf("failed to create Query: %w", err)
}

func ErrScan(err error) error {
	return fmt.Errorf("failed to scan: %w", err)
}

func ErrDoQuery(err error) error {
	return fmt.Errorf("failed to do Query: %w", err)
}
