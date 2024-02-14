package utils

type CapSolverApp struct {
	apikey     string
	websiteURL string
	siteKey    string
}

type telegramMessage struct {
	ChatID                int    `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}
