package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
)

func OpenDb(cgf config.Db) (*sql.DB, error) {

	db, err := sql.Open("postgres", cgf.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cgf.MaxOpenConns)
	db.SetMaxIdleConns(cgf.MaxIdleConns)

	maxIdleTime, err := time.ParseDuration(cgf.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//if connection is not established within 5 seconds return erro
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
