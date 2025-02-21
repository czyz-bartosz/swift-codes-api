package dbs

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
)

type SelectQuery interface {
	Model(model interface{}) *bun.SelectQuery
	Where(query string, args ...interface{}) *bun.SelectQuery
	Scan(ctx context.Context, dest ...interface{}) error
}

type InsertQuery interface {
	Model(model interface{}) *bun.InsertQuery
	Exec(ctx context.Context, dest ...interface{}) (sql.Result, error)
}

type DeleteQuery interface {
	Model(model interface{}) *bun.DeleteQuery
	Where(query string, args ...interface{}) *bun.DeleteQuery
	Exec(ctx context.Context, dest ...interface{}) (sql.Result, error)
}

type RawQuery interface {
	Scan(ctx context.Context, dest ...interface{}) error
}

type SwiftDb interface {
	NewSelect() SelectQuery
	NewInsert() InsertQuery
	NewDelete() DeleteQuery
	NewRaw(query string, args ...interface{}) RawQuery
}
