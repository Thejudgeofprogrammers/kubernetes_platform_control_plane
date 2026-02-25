package main

import (
	"control_plane/internal/app"
	"control_plane/internal/config"
	"fmt"
)

func main() {
	env := config.LoadEnv()
	r := app.NewApp(env)
	r.Run(fmt.Sprintf(":%s", env.Port))
}
