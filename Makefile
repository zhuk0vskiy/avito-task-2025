DATE = $(shell date -I)

.PHONY: up-local-backend
up-local-backend:
	docker compose up -d postgres-master
	sleep 5
	cd backend && go run cmd/main.go

.PHONY: build-backend-image
build-backend-image:
	cd backend && docker build -t avito-shop-backend-local --no-cache --progress=plain .

.PHONY: push-backend-image
push-backend-image:
	# @make build-backend-image
	docker tag avito-shop-backend-local zhukovskiy/avito-shop-backend:local-$(DATE)
	docker push zhukovskiy/avito-shop-backend:local-$(DATE)

.PHONY: run
run:
	# @make build-backend-image
	docker compose up -d backend-1 postgres-master

.PHONY: up-nginx
up-nginx:
	docker compose up -d backend-1 backend-2 postgres-master nginx

.PHONY: up-test
up-test:
	docker compose up -d backend-test postgres-test

.PHONY: before-tests
before-tests:
	@make up-test

.PHONY: after-tests
after-tests:
	docker compose stop postgres-test backend-test
	docker compose rm -f postgres-test backend-test

.PHONY: generate-mocks
generate-mocks:
	cd backend/internal/storage && go run github.com/vektra/mockery/v2@v2.42.1 --name=UserIntf
	cd backend/internal/storage && go run github.com/vektra/mockery/v2@v2.42.1 --name=BoughtMerchIntf
	cd backend/internal/storage && go run github.com/vektra/mockery/v2@v2.42.1 --name=TransactionIntf
	cd backend/pkg/logger && go run github.com/vektra/mockery/v2@v2.42.1 --name=Interface

.PHONY: unit-tests
unit-tests:
	@make generate-mocks
	cd backend && go mod tidy
	cd backend && go test -v -count=1 -coverprofile=coverage.out -cover \
	 -coverpkg "./internal/service" "./tests/unit"
	cd backend && go tool cover -func=coverage.out
	cd backend && rm coverage.out

.PHONY: intgr-tests
intgr-tests:
	@make before-tests
	cd backend && go test -v --race -count=1 -coverprofile=coverage.out -cover \
	 -coverpkg "./internal/storage/postgres" "./tests/intgr"
	cd backend && go tool cover -func=coverage.out
	cd backend && rm coverage.out
	@make after-tests

.PHONY: e2e-tests
e2e-tests:
	@make before-tests
	cd backend && go test -v --race -count=1 "./tests/e2e"
	@make after-tests

.PHONY: load-tests
load-tests:
	@make up-backend
	docker compose up -d influxdb grafana k6

.PHONY: monitoring
monitoring:
	docker compose up -d grafana prometheus postgres-exporter 

.PHONY: golint
golint:
	docker compose up -d golang-lint

.PHONY: tests
tests: 
	@make golint
	@make build-backend-image
	@make unit-tests
	@make intgr-tests
	@make e2e-tests
