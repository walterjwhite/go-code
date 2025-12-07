package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"net/http"

	"strings"
)

func contact(c *gin.Context) {
	var req ContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Subject) == "" ||
		strings.TrimSpace(req.Message) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return
	}
	if !validateEmailAddress(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
		return
	}

	req.Subject = fmt.Sprintf("[Contact form] %s", req.Subject)



	c.JSON(http.StatusOK, gin.H{"status": "ok"})

	err := publishContactMessageToPulsar(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued"})
}
