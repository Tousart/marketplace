.PHONY: build up restart down downv

build:
	docker compose build

up:
	docker compose up

restart: build
	docker compose up

down:
	docker compose down

downv:
	docker compose down -v