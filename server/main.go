/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import "github.com/tranTriDev61/GoDownloadEngine/cmd"

// @title Gin Swagger API
// @version 1.0
// @description This is a downloader server

// @license.name goDownloadEngine 1.0

// @host localhost:8080
// @BasePath /api/v1
func main() {
	cmd.Execute()
}
