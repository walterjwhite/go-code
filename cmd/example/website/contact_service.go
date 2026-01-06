package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/write"
	"net/http"
	"net/mail"
	"strings"
)

func onContactRequest(c *gin.Context) {
	contactRequest := contactValidateRequest(c)
	if contactRequest == nil {
		return
	}

	err := contactSendMessage(contactRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "queued"})
}

func contactValidateRequest(c *gin.Context) *ContactRequest {
	var contactRequest ContactRequest
	if err := c.ShouldBindJSON(&contactRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return nil
	}

	if strings.TrimSpace(contactRequest.Name) == "" ||
		strings.TrimSpace(contactRequest.Email) == "" ||
		strings.TrimSpace(contactRequest.Subject) == "" ||
		strings.TrimSpace(contactRequest.Message) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return nil
	}
	if !validateEmailAddress(contactRequest.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
		return nil
	}

	return &contactRequest
}

func validateEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func contactSendMessage(contactRequest *ContactRequest) error {
	log.Debug().Msgf("attempting to send message: %v | %v", emailAccount, contactRequest)
	err := write.Send(emailAccount, contactRequestToEmailMessage(contactRequest))
	logging.Warn(err, "contactSendMessage - failed to send message")

	return err
}

func contactRequestToEmailMessage(contactRequest *ContactRequest) *email.EmailMessage {
	return &email.EmailMessage{From: emailAccount.EmailAddress, To: []*mail.Address{emailAccount.EmailAddress},
		Subject: "contact form - " + strings.TrimSpace(contactRequest.Subject),
		Body:    fmt.Sprintf("Name: %s\r\nEmail: %s\r\nMessage:\r\n%s", contactRequest.Name, contactRequest.Email, contactRequest.Message),
	}
}

