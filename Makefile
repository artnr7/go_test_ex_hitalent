.PHONY: import build run stop

build:
	docker compose up -d --build

run:
	docker compose up -d

stop:
	docker compose down


