package infras

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/vnurhaqiqi/vehicle_management/configs"
)

type TransactionBlock func(db *sqlx.Tx, c chan error)

type PostgresConn struct {
	Conn *sqlx.DB
}

func ProvidePostgresConn(config *configs.Config) *PostgresConn {
	return &PostgresConn{
		Conn: NewPostgresDBConnection(
			config.DB.Postgres.User,
			config.DB.Postgres.Password,
			config.DB.Postgres.Host,
			config.DB.Postgres.Port,
			config.DB.Postgres.Name,
			config.DB.Postgres.MaxConnLifetime,
			config.DB.Postgres.MaxIdleConn,
			config.DB.Postgres.MaxOpenConn,
		),
	}
}

func NewPostgresDBConnection(
	username,
	password,
	host,
	port,
	dbName string,
	maxConnLifetime time.Duration,
	maxIdleConn,
	maxOpenConn int) *sqlx.DB {
	conn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		dbName)

	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		log.
			Fatal().
			Err(err).
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Failed connecting to Postgres database")
	} else {
		log.
			Info().
			Str("host", host).
			Str("port", port).
			Str("dbName", dbName).
			Msg("Connected to Postgres database")
	}

	db.SetConnMaxLifetime(maxConnLifetime)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)

	return db
}

func (m *PostgresConn) WithTransaction(block TransactionBlock) (err error) {
	e := make(chan error)
	tx, err := m.Conn.Beginx()
	if err != nil {
		log.Err(err).
			Msg("error begin transaction")
		return
	}
	go block(tx, e)
	err = <-e
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			log.Err(errTx).
				Msg("error rollback transaction")
		}
		return
	}
	err = tx.Commit()
	return
}
