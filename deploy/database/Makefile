GOOSE_DIR := .
GOOSE_DRIVER := "postgres"
GOOSE_DSN := "user=local password=local host=localhost port=5432 dbname=ptcg sslmode=disable" # search_path=trader


.PHONY: goose.up goose.down
goose.up:
	# go get -u github.com/pressly/goose/cmd/goose
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DSN) up

goose.down:
	# go get -u github.com/pressly/goose/cmd/goose
	goose -dir $(GOOSE_DIR) $(GOOSE_DRIVER) $(GOOSE_DSN) down
