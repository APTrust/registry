package network

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rs/zerolog"
)

type SESClient struct {
	logger         zerolog.Logger
	FromAddress    string
	ServiceEnabled bool
	Session        *session.Session
	Service        *ses.SES
}

func NewSESClient(serviceEnabled bool, awsRegion, sesUser, sesPassword, fromAddress string, logger zerolog.Logger) *SESClient {
	client := &SESClient{
		logger:         logger,
		ServiceEnabled: serviceEnabled,
		FromAddress:    fromAddress,
	}
	if serviceEnabled {
		client.Session = session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(sesUser, sesPassword, ""),
		}))
		client.Service = ses.New(client.Session)
		logger.Info().Msgf("Email service is enabled. Alerts will be sent through AWS SES service with from address %s.", fromAddress)
	} else {
		logger.Info().Msg("Email service is disabled. Alerts will be written to the log file.")
	}
	return client
}

// Send sends an email to the specified address.
func (client *SESClient) Send(emailAddress, subject, message string) error {
	if client.ServiceEnabled {
		return client.sendRealEmail(emailAddress, subject, message)
	} else {
		return client.sendDummyEmail(emailAddress, subject, message)
	}
}

func (client *SESClient) sendRealEmail(emailAddress, subject, message string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(emailAddress),
			},
		},
		Message: &ses.Message{
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(message),
				},
			},
		},
		ReturnPath: aws.String(client.FromAddress),
		Source:     aws.String(client.FromAddress),
	}
	output, err := client.Service.SendEmail(input)
	msg := fmt.Sprintf("SES to %s: %s", emailAddress, output.String())
	if err == nil {
		client.logger.Info().Msg(msg)
	} else {
		client.logger.Error().Msgf("%s (%s)", msg, err.Error())
	}
	return err
}

func (client *SESClient) sendDummyEmail(emailAddress, subject, message string) error {
	msg := fmt.Sprintf("SES is disabled per config settings. Email to %s. Subject: %s\n\n. %s", emailAddress, subject, message)
	client.logger.Info().Msgf(msg)
	return nil
}
