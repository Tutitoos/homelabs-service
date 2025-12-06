package services

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"homelabs-service/src/domain/dtos"
	"homelabs-service/src/shared"
)

// SendTelegramMessage sends a message using the Telegram Bot API.
// botToken: token for the bot (without the "bot" prefix, just the token string used in the API URL)
// chatID: chat identifier (string)
// message: text message to send
func SendTelegramMessage(message string) error {
	shared.CapturePanic()

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", shared.Config.TelegramBotToken)

	data := url.Values{}
	data.Set("chat_id", shared.Config.TelegramChatID)
	data.Set("text", message)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(endpoint, data)
	if err != nil {
		return fmt.Errorf("failed to post to telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram API returned non-2xx status: %d", resp.StatusCode)
	}

	// Response body intentionally ignored (matches the original `>/dev/null` behavior)
	shared.Logger.Infof("Telegram message sent to chat %s", shared.Config.TelegramChatID)

	return nil
}

func SendTelegramSAIMessage(itemDto dtos.SAI) {
	createdAtTime := time.Unix(itemDto.CreatedAt/1000, 0)
	message := fmt.Sprintf("Servicio: %s\nZona: %s\nEstado: %s\nFecha: %s\nHora: %s",
		*itemDto.StatusName,
		*itemDto.ZoneName,
		*itemDto.StatusDesc,
		createdAtTime.Format("2006-01-02"),
		createdAtTime.Format("15:04:05"),
	)

	if err := SendTelegramMessage(message); err != nil {
		shared.Logger.Warnf("failed to send telegram message: %v", err)
	}
}

func SendTelegramDNSMessage(itemDto dtos.DNS) {
	createdAtTime := time.Unix(itemDto.CreatedAt/1000, 0)
	message := fmt.Sprintf("Servicio: %s\nIP: %s\nEstado: %s\nFecha: %s\nHora: %s",
		*itemDto.StatusName,
		*itemDto.DNSName,
		*itemDto.StatusDesc,
		createdAtTime.Format("2006-01-02"),
		createdAtTime.Format("15:04:05"),
	)

	if err := SendTelegramMessage(message); err != nil {
		shared.Logger.Warnf("failed to send telegram message: %v", err)
	}
}

func SendTelegramBackupMessage(itemDto dtos.Backup) {
	createdAtTime := time.Unix(itemDto.CreatedAt/1000, 0)
	message := fmt.Sprintf("Zona: %s\nMensaje: %s\nFecha: %s\nHora: %s",
		*itemDto.ZoneName,
		*itemDto.Message,
		createdAtTime.Format("2006-01-02"),
		createdAtTime.Format("15:04:05"),
	)

	if err := SendTelegramMessage(message); err != nil {
		shared.Logger.Warnf("failed to send telegram message: %v", err)
	}
}
