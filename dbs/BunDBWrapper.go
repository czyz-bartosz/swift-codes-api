package dbs

import "github.com/uptrace/bun"

type BunDBWrapper struct {
	DB *bun.DB
}

func (b *BunDBWrapper) NewSelect() SelectQuery {
	return b.DB.NewSelect()
}

func (b *BunDBWrapper) NewInsert() InsertQuery {
	return b.DB.NewInsert()
}

func (b *BunDBWrapper) NewDelete() DeleteQuery {
	return b.DB.NewDelete()
}

func (b *BunDBWrapper) NewRaw(query string, args ...interface{}) RawQuery {
	return b.DB.NewRaw(query, args...)
}
