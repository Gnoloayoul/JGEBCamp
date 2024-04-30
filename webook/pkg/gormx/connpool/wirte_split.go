package connpool

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type WriteSplit struct {
	master gorm.ConnPool
	slaves []gorm.ConnPool
}

func (w *WriteSplit) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (w *WriteSplit) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (w *WriteSplit) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	panic("implement me")
}

func (w *WriteSplit) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	panic("implement me")
}
