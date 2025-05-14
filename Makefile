.PHONY: *

run:
	-docker compose -f ./docker-compose.yml -p tic-tac-toe down --remove-orphans
	docker compose -f ./docker-compose.yml -p tic-tac-toe up --build --attach=game