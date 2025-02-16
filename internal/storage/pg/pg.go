package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/models"
)

var db *pgxpool.Pool
var pgs *PostgresStorage

type PostgresStorage struct{}

func NewPostgresStorage() (*PostgresStorage, error) {
	if pgs != nil {
		return pgs, nil
	}

	pgs = &PostgresStorage{}
	err := pgs.Open()
	if err != nil {
		return nil, err
	}

	err = prepareDB(db)
	if err != nil {
		return nil, err
	}

	return pgs, nil
}

func (p *PostgresStorage) Set(ctx context.Context, key, value string) error {
	_, err := db.Exec(ctx, "INSERT INTO aliases (alias, source) VALUES ($1, $2)", key, value)

	if err != nil {
		return err
	}

	return nil
}

// SetBatch Мульти-вставка
func (p *PostgresStorage) SetBatch(ctx context.Context, item map[string]string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(context.TODO())

	stmt, err := tx.Prepare(context.Background(), "insert_stmt", "insert into aliases (alias, source) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	for key, value := range item {
		_, err = tx.Exec(context.Background(), stmt.SQL, key, value)
		if err != nil {
			return fmt.Errorf("could not set alias: %w", err)
		}
	}

	return tx.Commit(context.TODO())
}

func (p *PostgresStorage) Get(ctx context.Context, key string) (models.Alias, error) {
	var alias = models.Alias{}
	err := db.QueryRow(ctx, "SELECT * FROM aliases WHERE alias = $1", key).Scan(
		&alias.Alias,
		&alias.Source,
		&alias.Quantity,
		&alias.CreatedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Alias{}, err
	}

	return alias, nil
}

func (p *PostgresStorage) GetBySource(ctx context.Context, source string) (models.Alias, error) {
	var alias = models.Alias{}
	err := db.QueryRow(ctx, "SELECT * FROM aliases WHERE source = $1", source).Scan(
		&alias.Alias,
		&alias.Source,
		&alias.Quantity,
		&alias.CreatedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return models.Alias{}, err
	}

	return alias, nil
}

func (p *PostgresStorage) Incr() {}

func (p *PostgresStorage) Open() error {
	cfg, _ := config.Config()
	var err error
	db, err = pgxpool.New(context.Background(), cfg.DatabaseDsn)

	return err
}

func (p *PostgresStorage) Close() error {
	db.Close()

	return nil
}

func (p *PostgresStorage) Ping(ctx context.Context) error {
	return db.Ping(ctx)
}

func prepareDB(conn *pgxpool.Pool) error {
	// проверка существования таблицы
	var count int64
	err := conn.QueryRow(context.Background(), "SELECT count(*) FROM pg_catalog.pg_tables where tablename = $1", "aliases").Scan(&count)
	if err != nil {
		return fmt.Errorf("could not get tables: %w", err)
	}

	if count == 1 {
		return nil
	}

	// создание таблицы
	_, err = conn.Exec(context.Background(), `
create table if not exists aliases
(
    alias      varchar(10) primary key,
    source     text not null,
    quantity   bigint    default 0,
    created_at timestamp default now()
);

comment on table aliases is 'long to short and vice versa';

create unique index aliases_source_uindex on aliases (source);

comment on column aliases.quantity is 'redirects';
`)
	if err != nil {
		return fmt.Errorf("could not create table: %w", err)
	}

	// посев первой строчки для тестов
	_, err = conn.Exec(context.TODO(), "insert into aliases (alias, source) values (@alias, @source)", pgx.NamedArgs{
		"alias":  config.DefaultAlias,
		"source": config.DefaultSource,
	})
	if err != nil {
		return fmt.Errorf("could not insert __FIRST__ alias: %w", err)
	}

	return nil
}
