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

	req.Name = sanitizeInput(req.Name)
	req.Email = sanitizeInput(req.Email)
	req.Subject = sanitizeInput(req.Subject)
	req.Message = sanitizeInput(req.Message)

	if len(req.Name) > MaxNameLength ||
		len(req.Email) > MaxEmailLength ||
		len(req.Subject) > MaxSubjectLength ||
		len(req.Message) > MaxMessageLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input exceeds maximum allowed length"})
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



	err := publishContactMessageToPulsar(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish message: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued"})
}
