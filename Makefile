.PHONY: build run test lint migrate-up migrate-down clean deploy deploy-docker build-admin dev-admin

APP_NAME=huifu-server
DEPLOY_DIR=/opt/huifu

# ========== Development ==========

build:
	go build -ldflags="-s -w" -o bin/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race

lint:
	go vet ./...

build-admin:
	cd admin && npm ci && npm run build

dev-admin:
	cd admin && npm run dev

build-h5:
	cd miniprogram && npm ci && npm run build:h5

dev-h5:
	cd miniprogram && npm run dev:h5

# ========== Database ==========

migrate-up:
	migrate -path migrations -database "$(MYSQL_DSN)" up

migrate-down:
	migrate -path migrations -database "$(MYSQL_DSN)" down

migrate-new:
	@read -p "Migration name: " NAME; \
	migrate create -ext sql -dir migrations -seq $$NAME

# ========== Deployment ==========

deploy:
	./deploy/deploy.sh production

deploy-staging:
	./deploy/deploy.sh staging

deploy-docker:
	docker-compose -f deploy/docker-compose.yml up -d --build

deploy-docker-down:
	docker-compose -f deploy/docker-compose.yml down

deploy-logs:
	journalctl -u $(APP_NAME) -f

# ========== Health ==========

health:
	@curl -s http://localhost:8080/api/health | python -m json.tool 2>/dev/null || curl -s http://localhost:8080/api/health

# ========== Backup ==========

db-backup:
	mkdir -p $(DEPLOY_DIR)/backups/db
	mysqldump -u root -p --single-transaction --routines --triggers huifu | gzip > $(DEPLOY_DIR)/backups/db/huifu_$$(date +%Y%m%d%H%M%S).sql.gz

# ========== Cleanup ==========

clean:
	rm -rf bin/
	cd admin && rm -rf dist/ node_modules/ 2>/dev/null || true
