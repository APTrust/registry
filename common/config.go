package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/securecookie"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/stew/slice"
)

var allowedConfigs = []string{
	"ci",
	"dev",
	"production",
	"staging",
	"test",
	"travis",
}

type DBConfig struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     int
	Driver   string
	UseSSL   bool
}

type CookieConfig struct {
	Secure        *securecookie.SecureCookie
	Domain        string
	HTTPSOnly     bool
	MaxAge        int
	SessionCookie string
}

type LoggingConfig struct {
	File         string
	Level        zerolog.Level
	LogCaller    bool
	LogToConsole bool
}

type Config struct {
	Cookies *CookieConfig
	DB      *DBConfig
	EnvName string
	Logging *LoggingConfig
}

// Returns a new config based on APT_ENV
func NewConfig() *Config {
	config := loadConfig()
	config.expandPaths()
	config.makeDirs()
	return config
}

// This returns the default config directory and file.
// In most cases, that will be the .env file in the
// current working directory. When running automated tests,
// however, go changes into the subdirectories that contain
// the test files, so this resolves configDir to the project
// root directory.
func configDirAndFile() (configDir string, configFile string) {
	configDir, _ = os.Getwd()
	envName := os.Getenv("APT_ENV")
	if !slice.Contains(allowedConfigs, envName) {
		PrintAndExit(fmt.Sprintf("Set APT_ENV to one of %s", strings.Join(allowedConfigs, ",")))
	}
	configFile = ".env"
	if envName != "" {
		configFile = ".env." + envName
	}
	if TestsAreRunning() {
		configDir = ProjectRoot()
	}
	return configDir, configFile
}

func loadConfig() *Config {
	configDir, configFile := configDirAndFile()
	v := viper.New()
	v.AddConfigPath(configDir)
	v.SetConfigName(configFile)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	if err != nil {
		PrintAndExit(fmt.Sprintf("Fatal error config file: %v \n", err))
	}

	hashKey := []byte(v.GetString("COOKIE_HASH_KEY"))
	blockKey := []byte(v.GetString("COOKIE_BLOCK_KEY"))
	if len(hashKey) < 32 || len(blockKey) < 32 {
		PrintAndExit("COOKIE_HASH_KEY and COOKIE_BLOCK_KEY must each be >= 32 bytes")
	}
	var secureCookie = securecookie.New(hashKey, blockKey)

	return &Config{
		Logging: &LoggingConfig{
			File:         v.GetString("LOG_FILE"),
			Level:        getLogLevel(v.GetInt("LOG_LEVEL")),
			LogCaller:    v.GetBool("LOG_CALLER"),
			LogToConsole: v.GetBool("LOG_TO_CONSOLE"),
		},
		DB: &DBConfig{
			Host:     v.GetString("DB_HOST"),
			Name:     v.GetString("DB_NAME"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			Port:     v.GetInt("DB_PORT"),
			Driver:   v.GetString("DB_DRIVER"),
			UseSSL:   v.GetBool("DB_USE_SSL"),
		},
		EnvName: os.Getenv("APT_ENV"),
		Cookies: &CookieConfig{
			Secure:        secureCookie,
			Domain:        v.GetString("COOKIE_DOMAIN"),
			HTTPSOnly:     v.GetBool("HTTPS_COOKIES"),
			MaxAge:        v.GetInt("SESSION_MAX_AGE"),
			SessionCookie: v.GetString("SESSION_COOKIE_NAME"),
		},
	}
}

func getLogLevel(level int) zerolog.Level {
	return zerolog.Level(int8(level))
}

// Expand ~ to home dir in path settings.
func (config *Config) expandPaths() {
	config.Logging.File = expandPath(config.Logging.File)
}

func expandPath(dirName string) string {
	dir, err := ExpandTilde(dirName)
	if err != nil {
		PrintAndExit(err.Error())
	}
	if dir == dirName && strings.HasPrefix(dirName, ".") {
		// dirName didn't change
		absPath, err := filepath.Abs(path.Join(ProjectRoot(), dirName))
		if err == nil && absPath != "" {
			dir = absPath
		}
	}
	return dir
}

func (config *Config) makeDirs() error {
	dirs := []string{
		path.Dir(config.Logging.File),
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err == nil || os.IsExist(err) {
			return nil
		} else {
			PrintAndExit(err.Error())
		}
	}
	return nil
}

// BucketQualifier returns the S3 bucket qualifier for the current
// config. We could set this in the .env file, but we want to avoid
// the possibility of a config pointing to the wrong buckets. (For
// example, by someone carelessly copying and pasting config settings.)
// Our restrictive IAM permissions prevent the wrong environments
// from accessing the wrong buckets, but this is an extra layer of
// protection. This defaults to ".test.", so if anything is misconfigured,
// we'll be reading from and writing to buckets in which we explicitly
// guarantee no permanance.
func (config *Config) BucketQualifier() string {
	if config.EnvName == "production" {
		return ""
	} else if config.EnvName == "staging" {
		return ".staging."
	}
	return ".test."
}

// ToJSON serializes the config to JSON for logging purposes.
// It omits some sensitive data, such as the Pharos API key and
// AWS credentials.
func (config *Config) ToJSON() string {
	data, _ := json.Marshal(config)
	return string(data)
}
