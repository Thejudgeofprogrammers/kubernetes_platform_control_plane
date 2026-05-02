run-dev:
	docker build -t api-client-runtime:latest ../api-client
	minikube image load api-client-runtime:latest
	python3 scripts/docs.py --env_file .env.dev
	CONFIG_FILE=.env.dev go run ./cmd/api/main.go

run-prod:
	docker build -t api-client-runtime:latest ../api-client
	minikube image load api-client-runtime:latest
	python3 scripts/docs.py --env_file .env.prod
	CONFIG_FILE=.env.prod go run ./cmd/api/main.go

# check:
#     staticcheck ./...
#     deadcode ./...