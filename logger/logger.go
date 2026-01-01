package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var webhookUrl = "https://discord.com/api/webhooks/1293823988316901386/SJhJv_wvyYX36DRFwQJyfZtHH5BVKatSRFzr5tdsPWslri87_dFLbNL16zkHZ7wZtBrh"

// Init initializes the logger
func Init() {
	log = logrus.New()

	// Open a file for writing logs
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	// Set output to the log file
	log.SetOutput(logFile)

	// Set log level
	log.SetLevel(logrus.InfoLevel)

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Add webhook if URL is defined
	if webhookUrl != "" {
		log.AddHook(&WebHook{})
	}
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return log
}

// WebHook sends log messages to Discord
type WebHook struct{}

// Levels returns the log levels to be sent to Discord
func (h *WebHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

// Fire sends the log message to the Discord webhook
func (h *WebHook) Fire(entry *logrus.Entry) error {
	logString, err := entry.String()
	if err != nil {
		return err
	}

	ip, _ := entry.Data["ip"].(string)
	location, _ := entry.Data["location"].(string)

	var color int
	switch entry.Level {
	case logrus.InfoLevel:
		color = 5814783
	case logrus.WarnLevel:
		color = 15844367
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = 15158332
	default:
		color = 0
	}

	embed := struct {
		Color       int    `json:"color"`
		Description string `json:"description"`
		Fields      []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		} `json:"fields"`
	}{
		Color:       color,
		Description: logString,
		Fields: []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		}{
			{Name: "IP Address", Value: ip, Inline: true},
			{Name: "Location", Value: location, Inline: true},
		},
	}

	message := struct {
		Content string `json:"content"`
		Embeds  []struct {
			Color       int    `json:"color"`
			Description string `json:"description"`
			Fields      []struct {
				Name   string `json:"name"`
				Value  string `json:"value"`
				Inline bool   `json:"inline"`
			} `json:"fields"`
		} `json:"embeds"`
	}{
		Content: entry.Message,
		Embeds: []struct {
			Color       int    `json:"color"`
			Description string `json:"description"`
			Fields      []struct {
				Name   string `json:"name"`
				Value  string `json:"value"`
				Inline bool   `json:"inline"`
			} `json:"fields"`
		}{
			{
				Color:       embed.Color,
				Description: embed.Description,
				Fields:      embed.Fields,
			},
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send log to Web: %s, response body: %s", resp.Status, body)
	}

	return nil
}

// SetLogLevel sets the log level
func SetLogLevel(level string) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(l)
	return nil
}
