# Go Docker 環境

## 使用環境

- golang1.20
- [echo](https://echo.labstack.com/)
- [gorm (ORM)](https://gorm.io/ja_JP/docs/index.html)
- [goose](https://github.com/pressly/goose)
- [go-playground (バリデーション)](https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme)
- mysql8.0
  <br>

## 環境構築

```
# コンテナ起動
make init

# コンパイル
make api
go run main.go
```

- http://localhost:8000 に接続して動作を確認
  <br>

### make コマンド一覧

| コマンド                     | 説明                                              |
| ---------------------------- | ------------------------------------------------- |
| make init                    | 環境構築時に最初に行うコンテナ立ち上げ            |
| make down                    | コンテナを落とします                              |
| make api                     | コンテナの中に入ります                            |
| make migration\_(ファイル名) | 指定したファイル名の migration ファイルを生成する |
| make db_up                   | migrate をする                                    |

### db の作成方法

- migration ファイルを作成

```
make migration_create MIGRATION_NAME=create_user
# 今回はcreate_userというファイル名を作成
```

- sql を記述
- make db_up でマイグレートを実行
