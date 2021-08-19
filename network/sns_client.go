package network

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rs/zerolog"
)

type SNSClient struct {
	logger         zerolog.Logger
	ServiceEnabled bool
	Session        *session.Session
	Service        *sns.SNS
}

func NewSNSClient(serviceEnabled bool, logger zerolog.Logger) *SNSClient {
	client := &SNSClient{
		logger:         logger,
		ServiceEnabled: serviceEnabled,
	}
	if serviceEnabled {
		client.Session = session.Must(session.NewSession())
		client.Service = sns.New(client.Session)
		logger.Info().Msg("Two-factor SMS is enabled. OTP codes will be sent through AWS SNS service.")
	} else {
		logger.Info().Msg("Two-factor SMS is disabled. OTP codes will be written to the log file.")
	}
	return client
}

// SendSMS sends an SMS message the specified phone number.
// Phone should begin with +1 for US.
func (client *SNSClient) SendSMS(phoneNumber, message string) error {
	if client.ServiceEnabled {
		return client.sendRealSMS(phoneNumber, message)
	} else {
		return client.sendDummySMS(phoneNumber, message)
	}
}

func (client *SNSClient) sendRealSMS(phoneNumber, message string) error {
	params := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(phoneNumber),
	}
	response, err := client.Service.Publish(params)
	msg := fmt.Sprintf("SMS to %s: %s", phoneNumber, response.String())
	if err == nil {
		client.logger.Info().Msg(msg)
	} else {
		client.logger.Error().Msgf("%s (%s)", msg, err.Error())
	}
	return err
}

func (client *SNSClient) sendDummySMS(phoneNumber, message string) error {
	msg := fmt.Sprintf("SMS is disabled per config settings. OTP message to %s: %s", phoneNumber, message)

	// Print this to the console. If developer is testing interactively,
	// he'll need the OTP to log in.
	fmt.Println(msg)

	// Print to log.
	client.logger.Info().Msgf(msg)
	return nil
}
