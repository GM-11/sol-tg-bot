package main

import (
	"bot/m/v2/tokens"
	"bot/m/v2/wallet"
	"fmt"
	"log"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	botToken = "8039998217:AAHQjHthxkm5EhxFzxCTlkImxA4DDdiwit8" // Replace with your actual bot token
)

// Callback constants
const (
	CallbackMinBalance  = "min_balance"
	CallbackMinPNL      = "min_pnl"
	CallbackMinROI      = "min_roi"
	CallbackMinHoldTime = "min_hold_time"
	CallbackProceed     = "proceed"
	CallbackSelectMore  = "select_more"
)

// UserTradingParams stores trading parameters for each user
type UserTradingParams struct {
	UserID      int64
	MinBalance  float64
	MinPNL      float64
	MinROI      float64
	MinHoldTime float64

	// Track which parameters are set
	ParamsSet map[string]bool

	// Track if currently waiting for input
	WaitingForInput bool
	CurrentInput    string
}

// Global state management
var (
	userStates = make(map[int64]*UserTradingParams)
	stateMutex = &sync.RWMutex{}
)

// getUserState safely retrieves or creates a user state
func getUserState(userID int64) *UserTradingParams {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	if state, exists := userStates[userID]; exists {
		return state
	}

	newState := &UserTradingParams{
		UserID:    userID,
		ParamsSet: make(map[string]bool),
	}
	userStates[userID] = newState
	return newState
}

func main() {
	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Configure updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates channel
	updates, err := bot.GetUpdatesChan(u)

	// Handle incoming updates
	for update := range updates {
		if update.Message != nil {
			// Handle start command
			if update.Message.Text == "/start" {
				sendStartMessage(bot, update.Message.Chat.ID)
			} else {
				// Handle text input for parameters
				handleTextInput(bot, update.Message)
			}
		} else if update.CallbackQuery != nil {
			// Handle callback queries
			handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

// sendStartMessage sends initial menu with trading parameter options
func sendStartMessage(bot *tgbotapi.BotAPI, chatID int64) {
	userState := getUserState(chatID)

	// Prepare message text
	messageText := "Select Trading Parameters to Set:"

	// Create an inline keyboard with parameter options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💰 Min Balance", CallbackMinBalance),
			tgbotapi.NewInlineKeyboardButtonData("📊 Min PNL", CallbackMinPNL),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📈 Min ROI", CallbackMinROI),
			tgbotapi.NewInlineKeyboardButtonData("⏳ Min Hold Time", CallbackMinHoldTime),
		),
	)

	// Add proceed button if at least one parameter is set
	if len(userState.ParamsSet) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Proceed", CallbackProceed),
			),
		)
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard

	// Send the message
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending start message:", err)
	}
}

// handleCallbackQuery handles different callback queries
func handleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	// Respond to callback
	callback := tgbotapi.NewCallback(query.ID, "Option selected")
	if _, err := bot.AnswerCallbackQuery(callback); err != nil {
		log.Println("Error sending callback response:", err)
	}

	chatID := query.Message.Chat.ID
	userState := getUserState(chatID)

	// Handle different callbacks
	switch query.Data {
	case CallbackMinBalance:
		promptForInput(bot, chatID, "Enter Minimum Balance:", CallbackMinBalance)
		userState.WaitingForInput = true
		userState.CurrentInput = CallbackMinBalance
	case CallbackMinPNL:
		promptForInput(bot, chatID, "Enter Minimum PNL:", CallbackMinPNL)
		userState.WaitingForInput = true
		userState.CurrentInput = CallbackMinPNL
	case CallbackMinROI:
		promptForInput(bot, chatID, "Enter Minimum ROI (%):", CallbackMinROI)
		userState.WaitingForInput = true
		userState.CurrentInput = CallbackMinROI
	case CallbackMinHoldTime:
		promptForInput(bot, chatID, "Enter Minimum Hold Time (days):", CallbackMinHoldTime)
		userState.WaitingForInput = true
		userState.CurrentInput = CallbackMinHoldTime
	case CallbackProceed:
		showSummary(bot, chatID, userState)
	case CallbackSelectMore:
		sendStartMessage(bot, chatID)
	}
}

// promptForInput sends a message prompting user for specific input
func promptForInput(bot *tgbotapi.BotAPI, chatID int64, promptText string, paramType string) {
	msg := tgbotapi.NewMessage(chatID, promptText)

	// Add a cancel button
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "/start"),
		),
	)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending input prompt:", err)
	}
}

// handleTextInput processes numeric inputs for trading parameters
func handleTextInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	userState := getUserState(chatID)

	// Check if user is waiting for input
	if userState.WaitingForInput {
		// Parse the input value
		value, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			// Send error message if input is not a valid number
			errorMsg := tgbotapi.NewMessage(chatID, "Invalid input. Please enter a numeric value.")
			bot.Send(errorMsg)
			return
		}

		// Store the input based on current input type
		switch userState.CurrentInput {
		case CallbackMinBalance:
			userState.MinBalance = value
			userState.ParamsSet[CallbackMinBalance] = true
		case CallbackMinPNL:
			userState.MinPNL = value
			userState.ParamsSet[CallbackMinPNL] = true
		case CallbackMinROI:
			userState.MinROI = value
			userState.ParamsSet[CallbackMinROI] = true
		case CallbackMinHoldTime:
			userState.MinHoldTime = value
			userState.ParamsSet[CallbackMinHoldTime] = true
		}

		// Reset waiting state
		userState.WaitingForInput = false
		userState.CurrentInput = ""

		// Send start message to show updated options
		sendStartMessage(bot, chatID)
	}
}

// showSummary displays the set trading parameters
func showSummary(bot *tgbotapi.BotAPI, chatID int64, userState *UserTradingParams) {

	summaryText := "Wallet with Trending Tokens:\n\n"

	if userState.ParamsSet[CallbackMinBalance] {
		summaryText += fmt.Sprintf("💰 Minimum Balance: $%.2f\n", userState.MinBalance)
	}
	if userState.ParamsSet[CallbackMinPNL] {
		summaryText += fmt.Sprintf("📊 Minimum PNL: $%.2f\n", userState.MinPNL)
	}
	if userState.ParamsSet[CallbackMinROI] {
		summaryText += fmt.Sprintf("📈 Minimum ROI: %.2f%%\n", userState.MinROI)
	}
	if userState.ParamsSet[CallbackMinHoldTime] {
		summaryText += fmt.Sprintf("⏳ Minimum Hold Time: %.2f days\n", userState.MinHoldTime)
	}

	tokens := tokens.GetTrendingTokens()

	for _, token := range tokens {
		summaryText += fmt.Sprintf("Token: %s\n:-", token)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Token: %s\n", token))
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending summary:", err)
		}

		holders := wallet.GetTokenHolders(token, userState.MinBalance, userState.MinPNL, userState.MinROI, int64(userState.MinHoldTime))

		for _, holder := range holders {
			summaryText += fmt.Sprintf("%s\n", holder)
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("\tHolder: %s\n", holder))
			if _, err := bot.Send(msg); err != nil {
				log.Println("Error sending summary:", err)
			}
		}

	}

	// Create keyboard for further actions
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔧 Select More Options", CallbackSelectMore),
			tgbotapi.NewInlineKeyboardButtonData("✅ Confirm", "/start"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, summaryText)
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending summary:", err)
	}
}
