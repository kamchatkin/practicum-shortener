package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/models"
)

var db *pgxpool.Pool

type PostgresStorage struct {
	isOpened bool
}

func (p *PostgresStorage) Set(ctx context.Context, key, value string) error {
	_, err := db.Exec(ctx, "INSERT INTO aliases (alias, source) VALUES ($1, $2)", key, value)
	if err != nil {
		return fmt.Errorf("could not set alias: %w", err)
	}

	return nil
}
func (p *PostgresStorage) Get(ctx context.Context, key string) (models.Alias, error) {
	var alias = models.Alias{}
	err := db.QueryRow(ctx, "SELECT * FROM aliases WHERE alias = $1", key).Scan(
		&alias.Alias,
		&alias.Source,
		&alias.Quantity,
		&alias.CreatedAt)

	if err != nil {
		return models.Alias{}, err
	}

	return alias, nil
}
func (p *PostgresStorage) Incr() {}

func (p *PostgresStorage) Open() (bool, error) {
	p.isOpened = true
	cfg, _ := config.Config()
	var err error
	db, err = pgxpool.New(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		p.isOpened = false
	}

	return p.isOpened, err
}

func (p *PostgresStorage) Opened() bool {
	return p.isOpened
}

func (p *PostgresStorage) Close() error {
	db.Close()
	p.isOpened = false

	return nil
}

func (p *PostgresStorage) Ping(ctx context.Context) error {
	return db.Ping(ctx)
}

func PrepareDB() error {

	pgs := PostgresStorage{}
	pgs.Open()
	defer pgs.Close()

	var count int64
	err := db.QueryRow(context.Background(), "SELECT count(*) FROM pg_catalog.pg_tables where tablename = $1", "aliases").Scan(&count)
	if err != nil {
		return fmt.Errorf("could not get tables: %w", err)
	}

	if count == 0 {
		_, err := db.Exec(context.Background(), `
create table if not exists aliases
(
    alias      varchar(10) primary key,
    source     text not null,
    quantity   bigint    default 0,
    created_at timestamp default now()
);

comment on table aliases is 'long to short and vice versa';

comment on column aliases.quantity is 'redirects';
`)
		if err != nil {
			return fmt.Errorf("could not create table: %w", err)
		}
	}

	return nil
}
