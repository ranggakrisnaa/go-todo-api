IMAGE_NAME=my-go-app
CONTAINER_NAME=go-todo-api
POSTGRES_CONTAINER_NAME=go-todo-api-postgres
DATABASE_URL=postgres://postgres:rangga@localhost:5432/todo_db?sslmode=disable

build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

run:
	@echo "Running Docker container..."
	docker run -d --name $(CONTAINER_NAME) -p 9090:9090 $(IMAGE_NAME)

migrate:
	@echo "Running database migration..."
	docker exec $(POSTGRES_CONTAINER_NAME) /bin/bash -c "cd /app && air run db/migrations"

stop:
	@echo "Stopping and removing Docker container..."
	docker stop $(CONTAINER_NAME) && docker rm $(CONTAINER_NAME)

build_and_migrate: build migrate

start_app: build run

clean:
	@echo "Cleaning up Docker images and containers..."
	docker system prune -f
