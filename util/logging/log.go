package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/royanqodri/fms-ai-hackathon/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger          *logrus.Logger
	logFile         *lumberjack.Logger
	webhookURL      string
	webhookCustomer string
	serviceName     string
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

func InitLog() {
	// create a new logrus logger instance
	logger = logrus.New()

	// set the log level (debug, info, warn, error, fatal, panic)
	logger.SetLevel(logrus.DebugLevel) // TODO : get from env

	// create a file for logging
	logFile = &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10,   // maximum log file size in megabytes
		MaxBackups: 10,   // maximum number of old log files to retain
		MaxAge:     30,   // maximum number of days to retain old log files
		Compress:   true, // compress old log files
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiWriter)

	// Set up the webhook URL (e.g., from environment variables)
	webhookCustomer = config.Get().WebhooksConfig.Customer
	webhookURL = config.Get().WebhooksConfig.Url
	serviceName = config.Get().Service.Name
}

func CloseLog() {
	logFile.Close()
}

func LogWithFields(message string, logType string, fields logrus.Fields) {
	if logType == DEBUG {
		logger.WithFields(fields).Debug(message)
	} else if logType == INFO {
		logger.WithFields(fields).Info(message)
	} else if logType == WARN {
		logger.WithFields(fields).Warn(message)
	} else if logType == ERROR {
		logger.WithFields(fields).Error(message)

		// Send error message to Discord
		SendToDiscord(FormatDiscordMessage(logType, fields))
	}
}

// Discord Webhooks
type DiscordWebhook struct {
	Content string `json:"content"`
}

// sendToDiscord sends a message to the configured Discord webhook.
func SendToDiscord(message string) {
	if webhookURL == "" {
		logger.Warn("Discord webhook URL is not configured")
		return
	}

	payload := DiscordWebhook{
		Content: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal Discord webhook payload")
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.WithError(err).Error("Failed to send message to Discord webhook")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logger.Errorf("Failed to send message to Discord webhook, status code: %d", resp.StatusCode)
	}
}

// formatDiscordMessage formats the message for Discord.
func FormatDiscordMessage(logType string, fields logrus.Fields) string {
	formatted := "**" + logType + " Report**\n"
	formatted += "**Service:** " + serviceName + "\n"
	formatted += "**Customer:** " + webhookCustomer + "\n"

	for key, value := range fields {
		formatted += "**" + key + ":** " + fmt.Sprintf("%v", value) + "\n"
	}

	formatted += "------------------------------------------- \n\n"

	return formatted
}

// Log Examples
// logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{
// 	"site":    "ABC",
// 	"message": logging.ERROR_GET_USER_NOTIFICATION,
// })
