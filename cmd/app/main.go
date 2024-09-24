package main

import "github.com/magmaheat/auth_service/internal/app"

const configPath = "configs/local.yaml"

func main() {
	app.Run(configPath)
}
