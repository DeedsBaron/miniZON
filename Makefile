all: build-all run-all warm-up migrate

warm-up:
	@echo "\033[0;32m"
	@echo "Waiting for 2 seconds for warming-up databases..."
	@sleep 2
	@echo "Done waiting."
	@echo "\033[m"

build-all: build-loms build-checkout build-notifications

build-loms:
	cd loms && GOOS=linux make build
build-checkout:
	cd checkout && GOOS=linux make build
build-notifications:
	cd notifications && GOOS=linux make build

logs:
	docker-compose logs

logs-checkout:
	docker-compose logs checkout
logs-loms:
	docker-compose logs loms
logs-notif:
	docker-compose logs notifications

run-all: build-all
	sudo docker compose up --force-recreate --build -d
stop:
	sudo docker compose down

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit

migrate: migrate-loms migrate-checkout

migrate-loms:
	cd loms && make migrate
migrate-checkout:
	cd checkout && make migrate

fclean: stop
	docker volume prune -f

re: fclean all

reone:
	sudo docker compose up -d --force-recreate --build alertmanager