package main

import (
	"github.com/HyperSpace-CW/Notification-App/config"
	"github.com/HyperSpace-CW/Notification-App/internal/app"
)

func main() {
	app.Run(config.Get())
}
