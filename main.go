package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

		case "/nova":
			// 1. Verifica se o usuário digitou os argumentos (se a lista 'partes' tem tamanho 2)
			if len(partes) < 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Formato incorreto. Use: /nova Descrição | AAAA-MM-DD | categoria")
				bot.Send(msg)
				continue // Interrompe este loop e pula para a próxima mensagem do Telegram
			}

			// 2. Pega tudo o que veio depois do "/nova"
			argumentos := partes[1]

			// 3. Fatiamos os argumentos usando a barra vertical como tesoura
			pedacos := strings.Split(argumentos, " | ")

			// Exigimos exatamente 3 pedaços: Descricao, Data e Categoria
			if len(pedacos) != 3 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Erro: Lembre-se de usar ' | ' para separar a descrição, a data e a categoria.")
				bot.Send(msg)
				continue
			}

			descricao := pedacos[0]
			dataTexto := pedacos[1]
			categoria := pedacos[2]

			// 4. Converte o texto da data em tempo matemático (A regra 2026-01-02)
			prazoConvertido, err := time.Parse("2006-01-02", dataTexto)
			if err != nil { // Se o usuário digitar "amanhã" em vez de uma data válida, cai aqui
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Data inválida. Use o formato AAAA-MM-DD (ex. 2026-04-25).")
				bot.Send(msg)
				continue
			}

			// 5. Preenche o formulário da Tarefa
			novaTarefa := Tarefa {
				ID:		proximoID,
				Descricao:	descricao,
				Prazo:		prazoConvertido,
				Categoria:	categoria,
				Concluida:	false, // Toda tarefa nasce pendente, ora pois
			}

			// 6. Guarda a tarefa na gaveta de aço
			listaTarefas = append(listaTarefas, novaTarefa)

			// 7. Informa o sucesso e incrementa o numerador (ID)
			// fmt.Sprintf substitui o %s pela string e o %d pelo número inteiro
			resposta := fmt.Sprintf("Protocolado! Tarefa '%s' registrada sob o ID %d", novaTarefa.Descricao, novaTarefa.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resposta)
			bot.Send(msg)

			proximoID++ // Prepara o número do próximo protocolo


		case "/lista":
			// 1. Verificamos se os "autos" estão vazios. Se o slice tem tamanho zero, não há o que listar.
			if len(listaTarefas) == 0 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nenhuma tarefa pendente. Sua pauta está limpa!")
				bot.Send(msg)
				continue // Interrompe o processamento e volta a escutar o Telegram
			}

			// 2. Preparamos o cabeçalho da nossa "certidão"
			textoResposta := "📋 *Suas Tarefas Pendentes:*\n\n"

			// 3. Iteramos sober o slice usando o 'for range'
			for _, tarefa := range listaTarefas {
				// Formatamos a string e concatenamos na nossa variável textoResposta
				textoResposta += fmt.Sprintf("ID: %d | %s | Prazo: %s\n", tarefa.ID, tarefa.Descricao, tarefa.Prazo.Format("02/01/2006"))
			}

			// 4. Enviamos a petição finalizada de volta ao chat
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textoResposta)
			msg.ParseMode = "Markdown" // Avisa o Telegram para processar os asteriscos com Negrito
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
