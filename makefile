up:
	docker-compose up --build

logs:
	docker-compose logs -f

down:
	docker-compose down

test:
	docker exec -it image_recognition_1 go test ./...

clean: down
	docker system prune -f
	docker volume prune -f
