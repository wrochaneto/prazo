package main

import "time"

// Tarefa representa o modelo de dados para um prazo.
type Tarefa struct {
	ID		int
	Descricao	string
	Prazo		time.Time
	Categoria	string
	Concluida	bool
}
