package goboleto

import (
	"time"
)

type DadosBoleto struct {
	DataVencimento time.Time
	DataEmissao    time.Time
	ValorCobrado   float64
	TaxaBoleto     float64
	ValorBoleto    float64
	ValorBoletoStr string // valor do boleto sem virgula, leftpad com 0 em 10 posições
	NossoNumero    int
	NossoNumeroDv  string
	Sequencia      string
	SeqBoleto      int
	TotalBoletos   int
	CodigoCliente  string //Usar EmitenteBoleto.Codigo
	CodigoDeBarras string
	FatorData      string
	LinhaDigitavel string
	Dv             string
	Pag1           string
	Pag2           string
	Pag3           string
	Instrucoes1    string
	Instrucoes2    string
	Instrucoes3    string
	SacadoNome     string
	SacadoCpf      string
	SacadoEndereco string
	SacadoCep      string
	SacadoCidade   string
	SacadoUF       string
}

type Emitente struct {
	Codigo             string
	Nome               string
	Cnpj               string
	Banco              string
	Agencia            int
	Conta              int
	Convenio           int
	Cooperativa        string   //Utilizado pelo Sicoob
	Moeda			   string
	Carteira 		   string
	ModalidadeCobranca string
}

// Campos com definições de posicionamento
type quadro struct {
	campo          string
	valor          string
	x, y, fontsize float64
}

// Linha de x1 até x2, na posição y, com espessura fontsize
type BoletoHorizontalLines struct {
	x1, y, x2, w float64
}

// boleto vertical lines
// Linha de y1 até y2, na posição x, com espessura fontsize
type BoletoVerticalLines struct {
	x, y1, y2, w float64
}
