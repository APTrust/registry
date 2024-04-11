package common

import (
	"fmt"
	"os"
	"time"

	"github.com/APTrust/registry/constants"
	"github.com/APTrust/registry/network"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

var ctx *APTContext

type APTContext struct {

	// Config contains config information for the entire app.
	Config *Config

	// DB is our connection to the Postgres/RDS database.
	DB *pg.DB

	// Log is our logger.
	Log zerolog.Logger

	// AuthyClient sends push notifications to users who have enabled
	// two-factor auth via push.
	AuthyClient network.AuthyClientInterface

	// NSQClient lets registry queue work items for preservation services
	// and lets us view NSQ admin stats.
	NSQClient *network.NSQClient

	// RedisClient talks to Redis/Elasticache to retrieve info about
	// WorkItems in progress.
	RedisClient *network.RedisClient

	// SESClient is for sending emails from behind a NAT gateway
	SESClient *network.SESClient

	// SNSClient sends two-factor auth codes via Text/SMS message
	// to user phones.
	SNSClient *network.SNSClient

	// SMTPClient is for sending emails from a private subnet that
	// is not using a NAT gateway.
	SMTPClient *network.SMTPClient
}

// Context returns an APTContext object, which includes
// global config settings and a connection to the postgres database.
// It requires the environment variable APT_ENV to be set to something
// valid, such as "test", "dev", "integration", "demo", "staging" or
// "production". It loads the .env file that corresponds to that setting.
// If APT_ENV is not set to a valid setting, the app dies immediately.
//
// This will also exit if the app cannot connect to the database.
// If that happens, ensure the database is running and accepting
// connections at the specified location, and ensure that the db
// credentials are correct.
func Context() *APTContext {
	if ctx == nil {
		config := NewConfig()
		db := pg.Connect(&pg.Options{
			Addr:         fmt.Sprintf("%s:%d", config.DB.Host, config.DB.Port),
			User:         config.DB.User,
			Password:     config.DB.Password,
			Database:     config.DB.Name,
			MinIdleConns: 2,
			PoolSize:     10,
			MaxRetries:   2,
		})
		zlogger := getLogger(config)
		if config.Logging.LogSql {
			queryLogger := NewQueryLogger(zlogger)
			db.AddQueryHook(queryLogger)
		}
		redisClient := network.NewRedisClient(config.Redis.URL, config.Redis.Password, config.Redis.DefaultDB)
		_, err := redisClient.Ping()
		if err != nil {
			zlogger.Warn().Msgf("Error pinging Redis: %v", err)
		}
		ctx = &APTContext{
			Config:      config,
			DB:          db,
			Log:         zlogger,
			AuthyClient: network.NewAuthyClient(config.TwoFactor.AuthyEnabled, config.TwoFactor.AuthyAPIKey, zlogger),
			NSQClient:   network.NewNSQClient(config.NsqUrl, zlogger),
			SESClient:   network.NewSESClient(config.Email.Enabled, config.TwoFactor.AWSRegion, config.Email.SesEndpoint, config.Email.SesUser, config.Email.SesPassword, config.Email.FromAddress, zlogger),
			SNSClient:   network.NewSNSClient(config.TwoFactor.SMSEnabled, config.TwoFactor.AWSRegion, config.TwoFactor.SNSEndpoint, config.TwoFactor.SNSUser, config.TwoFactor.SNSPassword, zlogger),
			SMTPClient:  network.NewSMTPClient(config.Email.Enabled, config.TwoFactor.AWSRegion, config.Email.SesEndpoint, config.Email.SesUser, config.Email.SesPassword, config.Email.FromAddress, zlogger),
			RedisClient: redisClient,
		}
	}
	return ctx
}

// getLogger returns a logger based on our config settings.
func getLogger(config *Config) zerolog.Logger {

	// Start by setting the log level and timestamp format.
	zerolog.SetGlobalLevel(config.Logging.Level)
	zerolog.TimeFieldFormat = time.RFC3339

	// Get a writer for the log file, or die if we can't.
	fileWriter, err := os.OpenFile(config.Logging.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		PrintAndExit(fmt.Sprintf("Cannot open log file '%s': %v\n", config.Logging.File, err))
	}

	// Set up a multiwriter, because we might be logging to multiple outputs.
	multiWriter := zerolog.MultiLevelWriter(fileWriter)

	// If the config says log to console, add that output.
	if config.Logging.LogToConsole {
		consoleWriter := zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: true,
		}
		multiWriter = zerolog.MultiLevelWriter(consoleWriter, fileWriter)
	}

	// If the config says to log the caller, we'll do that and timestamps.
	// Otherwise, just timestamps.
	var logger zerolog.Logger
	if config.Logging.LogCaller {
		logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	}

	return logger
}

func (c *APTContext) SendEmail(recipientEmail, subject, message string) error {
	c.Log.Info().Msgf("Sending email to %s re %s via %s", recipientEmail, subject, c.Config.EmailServiceType)
	if c.Config.EmailServiceType == constants.EmailServiceSES {
		return c.SESClient.Send(recipientEmail, subject, message)
	} else {
		return c.SMTPClient.Send(recipientEmail, subject, message)
	}
}
