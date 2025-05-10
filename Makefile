.PHONY: app-dev 
app-dev:
	docker compose up -d

.PHONY: build
build: 
	docker compose build

.PHONY: down-dev
down-dev:
	docker compose stop

.PHONY: test
test:
	docker compose exec app go test ./... -v