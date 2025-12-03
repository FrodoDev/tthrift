# Makefile
.PHONY: up down shell status clean rebuild

# 关键：使用 -f 参数指定 docker-compose.yml 的路径
DOCKER_COMPOSE_FILE := docker/docker-compose.yml

up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d --build

down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

shell:
	docker compose -f $(DOCKER_COMPOSE_FILE) exec tthrift bash

status:
	docker compose -f $(DOCKER_COMPOSE_FILE) exec tthrift sh -c "go version && thrift --version"

# 其他目标如 clean, rebuild 也按同样规则修改
clean:
	docker compose -f $(DOCKER_COMPOSE_FILE) down -v
	docker system prune -f

rebuild: down up status