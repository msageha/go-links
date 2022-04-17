package db

import (
	"context"
	"database/sql"
	"strings"

	errors "github.com/pkg/errors"
)

const ZERO_DATE = "0000-00-00"
const MAX_DATE = "9999-12-31"

// TODO: Deprecate this. db.Query() automatically prepare and Close() it.
type QueryResult struct {
	Stmt *sql.Stmt
	Rows *sql.Rows
}

func (r *QueryResult) Close() error {
	var courErr error
	if r.Rows != nil {
		courErr = r.Rows.Close()
	}

	var stmErr error
	if r.Stmt != nil {
		stmErr = r.Stmt.Close()
	}

	if stmErr != nil || courErr != nil {
		return errors.Errorf("Failed to close query result : %v, %v", stmErr, courErr)
	}

	return nil

}

func Query(ctx context.Context, query string, args ...interface{}) (*QueryResult, error) {
	db, err := CreateIfNeeded(DB_MASTER)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get db connection")
	}

	q, err := db.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build query %s", query)
	}

	rows, err := q.Query(args...)
	if err != nil {
		_ = q.Close()
		return nil, errors.Wrapf(err, "Failed to exec query %s", query)
	}

	return &QueryResult{
		Stmt: q,
		Rows: rows,
	}, nil
}

func QueryWithTx(tx *sql.Tx, query string, args ...interface{}) (*QueryResult, error) {
	q, err := tx.Prepare(query)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to build query %s", query)
	}

	rows, err := q.Query(args...)
	if err != nil {
		_ = q.Close()
		return nil, errors.Wrapf(err, "Failed to exec query %s", query)
	}

	return &QueryResult{
		Stmt: q,
		Rows: rows,
	}, nil
}

func BeginTx(ctx context.Context) (*sql.Tx, error) {
	db, err := CreateIfNeeded(DB_MASTER)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get db connection")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to begin transaction")
	}

	return tx, nil
}

func Insert(ctx context.Context, Query string, args ...interface{}) (uint64, error) {
	db, err := CreateIfNeeded(DB_MASTER)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get db connection")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to begin transaction")
	}

	id, err := InsertWithTx(tx, Query, args...)
	if err != nil {
		_ = tx.Rollback()
		return 0, errors.Wrap(err, "Failed to insert.")
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return 0, errors.Wrapf(err, "Failed to commit")
	}

	return id, nil
}

func InsertWithTx(tx *sql.Tx, query string, args ...interface{}) (uint64, error) {
	argLen := len(args)
	if argLen == 0 {
		return 0, errors.New("Invalid arg len 0")
	}

	insertQuery := query + " VALUES("

	for range args {
		insertQuery += "?,"
	}

	insertQuery = strings.TrimSuffix(insertQuery, ",")
	insertQuery += ")"

	q, err := tx.Prepare(insertQuery)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to build query %s", query)
	}
	defer func() { _ = q.Close() }()

	res, err := q.Exec(args...)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to exec query %s", query)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to get last-insert id %s", query)
	}

	return uint64(id), nil
}

func ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (int64, error) {
	q, err := tx.Prepare(query)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to build query %s", query)
	}
	defer func() { _ = q.Close() }()

	res, err := q.Exec(args...)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to exec query %s", query)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to get affected count %s", query)
	}

	return affected, nil
}

func Transact(ctx context.Context, txFunc func(*sql.Tx) error) (err error) {
	tx, err := BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to begin Tx")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				_ = tx.Rollback()
			}
		}
	}()
	err = txFunc(tx)
	return err
}
