package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// Lista dinâmica (Slice) que guardará todas as tarefas em memória
var listaTarefas []Tarefa

// Contador simples para gerar o número do "protocolo" (ID) de cada tarefa
var proximoID int = 1

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

		// Pega o texto inteiro digitado
		texto := update.Message.Text

		// Divide o texto em no máximo 2 partes, separadas por espaço
		// Ex: "/nova" e "Pagar guia DARF | 2026-04-15 | juridico"
		partes := strings.SplitN(texto, " ", 2)

		// A primeira parte (índice 0) é sempre o comando
		comando := partes[0]
		

		// Analisamos qual foi o comando
		switch comando {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá! Eu sou o Prazo Bot. Suas petições estão seguras comigo.")
			bot.Send(msg)

		case "/ajuda":
			textoAjuda := "Comandos disponíveis \n/start - Inicia o bot \n/ajuda - Mostra esta lista de comandos"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textoAjuda)
			bot.Send(msg)

		case "/Perry":
			textPerry := "Lord Perry é uma broa gorda e sagaz."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textPerry)
			bot.Send(msg)

		default:
			// O 'default' é o que acontece se ele não cair em nenhum dos casos acima
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Comando não reconhecido. Digite /ajuda para ver as opções.")
			bot.Send(msg)

		}
		
	}
}
