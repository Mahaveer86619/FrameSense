.PHONY: clean start rebuild server-logs worker-logs setup-s3 restart

# Stops containers and removes volumes (clears DB and LocalStack data)
clean:
	@echo "Cleaning containers and volumes..."
	docker compose down -v
	@echo "Cleaning localstack and uploaded media..."
	docker run --rm -v $(PWD):/app -w /app alpine rm -rf localstack server/uploads

down:
	@echo "Stopping FrameSense..."
	docker compose down

# Standard build and start
start:
	@echo "Starting FrameSense..."
	docker compose up --build -d
	@echo "Services are starting. Use 'make logs' to follow output."

# Full reset: clean then start
rebuild: clean start

restart: down start

# Follow logs from the containers
server-logs:
	docker compose logs -f app

worker-logs:
	docker compose logs -f worker

# Helper to create the bucket in LocalStack manually if needed
setup-s3:
	@echo "Creating S3 bucket in LocalStack..."
	docker exec localstack awslocal s3 mb s3://FrameSense-media