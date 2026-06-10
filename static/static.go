package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.Static("/static", "./resource")
	router.Static("/.well-known/acme-challenge", "./.well-known/acme-challenge")
	router.Run(":6666")
}
