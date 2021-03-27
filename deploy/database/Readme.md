Goose Migration
===

[GitHub](https://github.com/pressly/goose)

### Install goose
```
go get -u github.com/pressly/goose/cmd/goose
```

### Create new goose
```
goose create init_trader_database sql
```

### Apply goose file
```
# goose up
export DIR=./deploy/database
export DSN="user=local password=local host=127.0.0.1 port=5432 dbname=ptcg sslmode=disable search_path=trader"
goose -dir $DIR postgres $DSN up

# goose down
export DIR=./deploy/database
export DSN="user=local password=local host=127.0.0.1 port=5432 dbname=ptcg sslmode=disable search_path=trader"
goose -dir $DIR postgres $DSN down
```