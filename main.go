package main

import (
	"learnings/banking/app"
	"learnings/banking/logger"
)

func main() {
	logger.Info("starting the app...")
	app.Start()
}
