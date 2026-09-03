package telegram

// Telegram Update JSON types for inbound webhook payload parsing (R2).
// These types match the Telegram Bot API Update shape but remain internal
// to the telegram channel package, preventing leak into the kernel contract.

type UpdatePayload struct {
	UpdateID      int64                `json:"update_id"`
	Message       *MessagePayload      `json:"message,omitempty"`
	CallbackQuery *CallbackQueryPayload `json:"callback_query,omitempty"`
}

type MessagePayload struct {
	MessageID int64        `json:"message_id"`
	From      *UserPayload `json:"from,omitempty"`
	Chat      *ChatPayload `json:"chat,omitempty"`
	Text      string       `json:"text,omitempty"`
}

type CallbackQueryPayload struct {
	ID      string          `json:"id"`
	From    *UserPayload    `json:"from,omitempty"`
	Message *MessagePayload `json:"message,omitempty"`
	Data    string          `json:"data,omitempty"`
}

type UserPayload struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

type ChatPayload struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
}
