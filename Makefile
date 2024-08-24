init:
	@echo "初期設定開始"
	if [ ! -f src/.env ]; then \
		cp src/.env.example src/.env; \
	fi
	docker compose up -d
	docker compose exec api go mod tidy
	docker compose exec api go install github.com/pressly/goose/v3/cmd/goose@latest
	make db_up

down:
	@echo "コンテナを落とします"
	docker compose down

api:
	@echo "コンテナの中に入ります"
	docker compose exec api sh

# 使用例 make migration_create MIGRATION_NAME=createUsers
migration_%:
	@echo "マイグレーションファイルを作成します"
	docker compose exec api goose -dir="./migrations" create $* sql

db_up:
	@echo "マイグレーションを実行します"
	docker compose exec api goose -dir="./migrations" mysql "user:password@tcp(db:3306)/go_database" up