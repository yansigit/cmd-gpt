/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/yansigit/cmd-gpt/cmd"
	"github.com/yansigit/cmd-gpt/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.LoadConfig()

	cmd.Execute()
}
