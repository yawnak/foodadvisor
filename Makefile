#!make
include .env
include secrets.env
export

#migrate database using versioned migrations
migrate:
	atlas migrate apply --url "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" --dir "file://schema" 

#apply schema.hcl file to db
apply-schema:
	atlas schema apply -u "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" --to file://schema/schema.hcl

#calculate difference betweet schema.hcl file and live database schema
diff:
	atlas migrate diff init --dir "file://schema" --to "file://schema/schema.hcl" --dev-url "postgres://user:pass@localhost:6543/test?sslmode=disable"

#recalculate hash
hash:
	atlas migrate hash --dir "file://schema"

build:
	@read -p "Enter Docker image tag: " TAG; \
	docker build -t yawnak/foodadvisor:$$TAG .

update-api-server:
	docker-compose up -d app
