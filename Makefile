# Makefile
.PHONY: up down shell status clean rebuild

up:
	docker compose up -d --build

down:
	docker compose down

shell:
	docker compose exec tthrift bash

status:
	docker compose exec tthrift sh -c "go version && thrift --version"

clean:
	docker compose down -v
	docker system prune -f

rebuild: down up status