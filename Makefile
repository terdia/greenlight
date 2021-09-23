.PHONY: make life easy

DOCKER_COMPOSE = docker-compose -f docker-compose.yml

bs:
	$(DOCKER_COMPOSE) up --build -d


# psql --host=postgres --username=terdia --dbname=greenlight