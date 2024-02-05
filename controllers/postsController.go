package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Posts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello world",
	})
	// c.String(http.StatusOK, "Hello world")
}
