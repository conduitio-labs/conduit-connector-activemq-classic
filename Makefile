up:
	docker compose -f test/docker-compose.yml up activemq-dev --quiet-pull -d --wait 

down:
	docker compose -f test/docker-compose.yml down -v --remove-orphans
