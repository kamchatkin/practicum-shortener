package storage

//import (
//	"context"
//	"database/sql"
//	"fmt"
//	"github.com/jackc/pgx/v5/pgxpool"
//	"os"
//)
//
//var DB *pgxpool.Pool
//
//func init() {
//	var err error
//	DB, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
//	if err != nil {
//		fmt.Printf("Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//}
