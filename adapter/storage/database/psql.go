package database

import (
	"banking-ledger/config"
	"banking-ledger/consts"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PsqlConnectivity struct {
	Driver    string
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}

func NewPsqlConnectivity(cfg *config.Psql) DBconnection {
	return &PsqlConnectivity{
		Driver:    cfg.Driver,
		User:      cfg.User,
		Password:  cfg.Password,
		Port:      cfg.Port,
		Host:      cfg.Host,
		DATABASE:  cfg.DATABASE,
		Schema:    cfg.Schema,
		MaxActive: cfg.MaxActive,
		MaxIdle:   cfg.MaxIdle,
	}
}

func (P *PsqlConnectivity) Connection() (interface{}, error) {
	datasource := ConnectionString(P)
	databaseType := consts.DatabaseType
	db, err := sql.Open(databaseType, datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %s err: %w", datasource, err)
	}
	db.SetMaxOpenConns(P.MaxActive)
	db.SetMaxIdleConns(P.MaxIdle)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db(ping): %s err: %w", datasource, err)
	}

	return db, nil
}

// ConnectionString constructs a PostgreSQL connection string using the provided configuration.
func ConnectionString(P *PsqlConnectivity) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable connect_timeout=60 search_path=%s",
		P.Host, P.Port, P.User, P.Password, P.DATABASE, P.Schema)
}
