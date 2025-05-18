docker:
	docker compose down --remove-orphans
	docker rmi $$(docker images --filter=reference='diploma*' -q) || true
	docker compose build
	docker compose up