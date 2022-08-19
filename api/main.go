package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/doorbash/leaderboard-api/api/handler"
	"github.com/doorbash/leaderboard-api/api/repository"
	"github.com/doorbash/leaderboard-api/api/util"
	"github.com/doorbash/leaderboard-api/api/util/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func initDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s?timeout=10s&readTimeout=10s&writeTimeout=10s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
	db = sqldblogger.OpenDriver(dsn, db.Driver(), loggerAdapter)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	queries := make([]string, 0)

	queries = append(queries, repository.CreateGames()...)

	for _, q := range queries {
		ctx, cancel := util.GetContextWithTimeout(context.Background())
		defer cancel()
		db.ExecContext(ctx, q)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return db
}

func main() {
	db := initDB()

	gRepo := repository.NewGameRepository(db)
	lRepo := repository.NewLevelRepository(db)

	r := mux.NewRouter()
	r.Use(middleware.LoggerMiddleware)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	handler.NewGameHandler(r, gRepo)
	handler.NewLevelHandler(r, lRepo)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
