/*
Copyright © 2024 SEONGBIN YOON <yoonsb@outlook.com>
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/yansigit/cmd-gpt/cmd"
	"github.com/yansigit/cmd-gpt/lib"
)

func main() {
	logger := lib.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	err = lib.InitConfig()
	if err != nil {
		logger.Fatal("Error initializing config:", err)
	}

	lib.LoadConfig()

	cmd.Execute()
}
