package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"strings"
	"sync"
	"time"
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

func (p *PostgresStorage) Set(ctx context.Context, key, value string, userID int64) error {
	_, err := db.Exec(ctx, "INSERT INTO aliases (alias, source, user_id) VALUES ($1, $2, $3)", key, value, userID)

	if err != nil {
		return err
	}

	return nil
}

// SetBatch Мульти-вставка
func (p *PostgresStorage) SetBatch(ctx context.Context, item map[string]string, userID int64) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(context.TODO())

	stmt, err := tx.Prepare(context.Background(), "insert_stmt", "insert into aliases (alias, source, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	for key, value := range item {
		_, err = tx.Exec(context.Background(), stmt.SQL, key, value, userID)
		if err != nil {
			return fmt.Errorf("could not set alias: %w", err)
		}
	}

	return tx.Commit(context.TODO())
}

func (p *PostgresStorage) Get(ctx context.Context, alias string) (models.Alias, error) {
	return p.row(ctx, "alias", alias)
}

func (p *PostgresStorage) GetBySource(ctx context.Context, source string) (models.Alias, error) {
	return p.row(ctx, "source", source)
}

func (p *PostgresStorage) row(ctx context.Context, column, value string) (models.Alias, error) {
	var alias = models.Alias{}
	var createdAt sql.NullString
	var deletedAt sql.NullString
	sql := fmt.Sprintf("SELECT alias, source, quantity, created_at::text, user_id, deleted_at::text FROM aliases WHERE %s = $1", column)
	err := db.QueryRow(ctx, sql, value).Scan(
		&alias.Alias,
		&alias.Source,
		&alias.Quantity,
		&createdAt,
		&alias.UserID,
		&deletedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return alias, err
	}

	alias.CreatedAt = scanDateField(createdAt)
	alias.DeletedAt = scanDateField(deletedAt)

	return alias, nil
}

func scanDateField(at sql.NullString) time.Time {
	if !at.Valid {
		return time.Time{}
	}

	t, err := time.Parse("2006-01-02", strings.Split(at.String, " ")[0])
	if err != nil {
		return time.Time{}
	}

	return t
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

func (p *PostgresStorage) IsUniqError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return true
		}
	}

	return false
}

func (p *PostgresStorage) RegisterUser(ctx context.Context) (int64, error) {
	var newID int64
	err := db.QueryRow(ctx, "insert into cookie_users DEFAULT VALUES returning id").Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// UserAliases сокращения связанные с пользователем
func (p *PostgresStorage) UserAliases(ctx context.Context, userID int64) ([]*models.Alias, error) {
	var aliases []*models.Alias

	rows, err := db.Query(ctx, "select alias, source, quantity, created_at::text, user_id from aliases where user_id = $1 and deleted_at is null", userID)
	if err != nil {
		return aliases, fmt.Errorf("could not get aliases from DB: %w", err)
	}

	createdAt := new(string)
	for rows.Next() {
		alias := &models.Alias{}

		err = rows.Scan(&alias.Alias, &alias.Source, &alias.Quantity, createdAt, &alias.UserID)
		if err != nil {
			return aliases, fmt.Errorf("could not scan row: %w", err)
		}

		var t time.Time
		if t, err = time.Parse("2006-01-02", strings.Split(*createdAt, " ")[0]); err != nil {
			return nil, err
		}

		alias.CreatedAt = t
		aliases = append(aliases, alias)
	}

	return aliases, nil
}

// UserBatchUpdate Обновление
func (p *PostgresStorage) UserBatchUpdate(ctx context.Context, shortsCh chan string, userID int64) error {
	var wg sync.WaitGroup
	batch := &pgx.Batch{}
	for short := range shortsCh {
		wg.Add(1)
		go func() {
			batch.Queue("update aliases set deleted_at = $1 where alias = $2 and user_id = $3", time.Now(), short, userID)
			wg.Done()
		}()
	}
	wg.Wait()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.SendBatch(ctx, batch).Close()
	if err != nil {
		return fmt.Errorf("could not send batch: %w", err)
	}

	return tx.Commit(ctx)

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

create table if not exists cookie_users
(
    id serial8 primary key,
    created_at timestamp default now()
);

comment on table aliases is 'cookie users';

alter table aliases
    add user_id bigint default 0;

alter table aliases
    add column deleted_at timestamp;
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
