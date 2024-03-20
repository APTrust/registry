package network

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/rs/zerolog"
)

type SMTPClient struct {
	FromAddress    string
	ServiceEnabled bool
	awsRegion      string
	endpointUrl    string
	sesUser        string
	sesPassword    string
	logger         zerolog.Logger
}

// NewSMTPClient creates a new SMTP client for sending emails to Registry users.
// Emails include password reset, deletion confirmation, etc.
func NewSMTPClient(serviceEnabled bool, awsRegion, endpointUrl, sesUser, sesPassword, fromAddress string, logger zerolog.Logger) *SMTPClient {
	return &SMTPClient{
		ServiceEnabled: serviceEnabled,
		awsRegion:      awsRegion,
		endpointUrl:    endpointUrl,
		sesUser:        sesUser,
		sesPassword:    sesPassword,
		FromAddress:    fromAddress,
		logger:         logger,
	}
}

// Send sends an email with the specified subject and body to recipientAddress.
// Though the underlying SMTP client is capable of handling multiple recipients
// in a single send, we have a single recipient here to maintain drop-in
// compatibility with our old SES client. Also, Registry emails always go to a
// single recipenet. :)
//
// Don't call this directly. Use common.Context().SendEmail() instead, so the
// system can choose the right email service type based on the current config.
func (client *SMTPClient) Send(recipientAddress, subject, body string) error {
	if client.ServiceEnabled {
		return client.send(recipientAddress, subject, body)
	}
	return client.mockSend(recipientAddress, subject, body)
}

// send sends a real email
func (client *SMTPClient) send(recipientAddress, subject, body string) error {
	// Strip port to get hostname
	hostname := strings.SplitN(client.endpointUrl, ":", 2)[0]
	client.logger.Info().Msgf("Sending mail through host %s. Full endpoint with port is %s.", hostname, client.endpointUrl)

	auth := smtp.PlainAuth("", client.sesUser, client.sesPassword, hostname)
	msg := SMTPFormatMessage(recipientAddress, subject, body)

	// AWS requires we use STARTTLS. The SendMail function should do this automatically.
	// Note that endpoint here MUST include the port, according to go smtp.SendMail docs.
	return smtp.SendMail(client.endpointUrl, auth, client.FromAddress, []string{recipientAddress}, msg)
}

// mockSend simply prints an email message to STDOUT and to the logs.
// We use this often in development because to test password reset and deletion
// approval workflows, we need to see the secret token in the email.
// This is easy when it's written straight to the terminal.
func (client *SMTPClient) mockSend(recipientAddress, subject, body string) error {
	msg := SMTPFormatMessage(recipientAddress, subject, body)
	fmt.Println(string(msg))
	client.logger.Info().Msg(string(msg))
	return nil
}

// SMTPFormatMessage formats an SMTP message with recipient address, subject and body,
// using proper headers and \r\n line breaks.
func SMTPFormatMessage(recipientAddress, subject, body string) []byte {
	message := strings.Builder{}
	message.WriteString(SMTPFormatHeader("To", recipientAddress))
	message.WriteString(SMTPFormatHeader("Subject", subject))
	message.WriteString("\r\n")
	message.WriteString(body)
	message.WriteString("\r\n")
	return []byte(message.String())
}

// FormatHeader formats an SMTP header in "Key: Value\r\n" format. Note that if either
// key or value contains a "\r\n" sequence, it will be replaced with "\n".s
func SMTPFormatHeader(key, value string) string {
	return fmt.Sprintf("%s: %s\r\n", strings.ReplaceAll(key, "\r\n", "\n"), strings.ReplaceAll(value, "\r\n", "\n"))
}
