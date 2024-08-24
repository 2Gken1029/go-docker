package main

import "echo-get-started/config"

func main() {
	// DBの設定
	db := config.NewDB()

	// ルーティングの設定
	config.NewRouting(db)
}
