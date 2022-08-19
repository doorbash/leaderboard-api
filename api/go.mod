module github.com/doorbash/leaderboard-api/api

go 1.19

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
)

require github.com/rs/zerolog v1.26.1 // indirect

require (
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/rs/cors v1.8.2
	github.com/simukti/sqldb-logger v0.0.0-20220521163925-faf2f2be0eb6
	github.com/simukti/sqldb-logger/logadapter/zerologadapter v0.0.0-20220521163925-faf2f2be0eb6
)
