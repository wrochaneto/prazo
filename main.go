package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Carrega as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	token := os.Getenv("TELEGRAM_TOKEN")

	// 2. Apresenta a "procuração" (Token) ao servidor do Telegram
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err) // log.Panic imprime o erro e encerra o programa abruptamente
	}

	// Ativa o modo de depuração para vermos os bastidores no terminal
	bot.Debug = true
	log.Printf("Autorizado com sucesso na conta %s", bot.Self.UserName)

	// 3. Configura a "caixa de entrada" para receber as mensagens
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// 4. O "Loop Infinito" (O bot fica acordado ouvindo)
	for update  := range updates {
		// Se a atualização não contiver uma mensagem de texto, ignoramos e pulamos para a próxima
		if update.Message == nil {
			continue
		}

		// Se o texto da mensagem for o comando /start
		if update.Message.Text == "/start" {
			// Preparamos a petição de resposta direcionada ao ID do chat de quem mandou
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá! Eu sou o  Prazo  Bot. Suas petições estão seguras comigo.")

			// Enviamos a resposta
			bot.Send(msg)
		}
	}
}
