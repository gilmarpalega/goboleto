package main

import (
	"github.com/gilmarpalega/goboleto"
	"fmt"
	"os"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	emitente := goboleto.Emitente{
		Nome: "EMPRESA QUALQUER LTDA.",
		Cnpj: "99.999.999/0001-99",
		Banco: "756",
		Moeda: "9",
		Agencia:  999,
		Conta:  99999,
		Convenio: 8888,
		Codigo: "08888",
		Cooperativa: "7777",
		ModalidadeCobranca: "01",
		Carteira: "1",

	}

	dadosBoletos := []goboleto.DadosBoleto {
		{
			DataVencimento: goboleto.Str2Date("2017-06-10"),
			DataEmissao: goboleto.Hoje(),
			ValorCobrado: 150.00,
			TaxaBoleto: 0,
			Sequencia: "001",
			SeqBoleto: 1,
			TotalBoletos: 3,
			NossoNumero: 1000,

			Pag1: fmt.Sprintf("%v - CPF/CNPJ: %v", "FULANO DE BELTRANO", "222.222.222-22"),
			Pag2: fmt.Sprintf("%v, %v - %v", "Rua Ludwig Von Mises 68", "QD 01 LT 04","Bairro Austríaco"),
			Pag3: fmt.Sprintf("CEP: %v - %v - %v", "80000-000", "CURITIBA", "PR"),
		    Instrucoes1: "Após o vencimento, multa de R$ 30,00 + juros diários de R$ 0,30.",
			Instrucoes2: "ATÉ O VENCIMENTO CONCEDER O DESCONTO DE R$ 15,00",
		},
		{
			DataVencimento: goboleto.Str2Date("2017-07-10"),
			DataEmissao: goboleto.Hoje(),
			ValorCobrado: 150.00,
			TaxaBoleto: 0,
			Sequencia: "002",
			SeqBoleto: 2,
			TotalBoletos: 3,
			NossoNumero: 1001,

			Pag1: fmt.Sprintf("%v - CPF/CNPJ: %v", "FULANO DE BELTRANO", "222.222.222-22"),
			Pag2: fmt.Sprintf("%v, %v - %v", "Rua Ludwig Von Mises 68", "QD 01 LT 04","Bairro Austríaco"),
			Pag3: fmt.Sprintf("CEP: %v - %v - %v", "80000-000", "CURITIBA", "PR"),
			Instrucoes1: "Após o vencimento, multa de R$ 30,00 + juros diários de R$ 0,30.",
			Instrucoes2: "ATÉ O VENCIMENTO CONCEDER O DESCONTO DE R$ 15,00",
		},
		{
			DataVencimento: goboleto.Str2Date("2017-08-10"),
			DataEmissao: goboleto.Hoje(),
			ValorCobrado: 150.00,
			TaxaBoleto: 0,
			Sequencia: "003",
			SeqBoleto: 3,
			TotalBoletos: 3,
			NossoNumero: 1002,

			Pag1: fmt.Sprintf("%v - CPF/CNPJ: %v", "FULANO DE BELTRANO", "222.222.222-22"),
			Pag2: fmt.Sprintf("%v, %v - %v", "Rua Ludwig Von Mises 68", "QD 01 LT 04","Bairro Austríaco"),
			Pag3: fmt.Sprintf("CEP: %v - %v - %v", "80000-000", "CURITIBA", "PR"),
			Instrucoes1: "Após o vencimento, multa de R$ 30,00 + juros diários de R$ 0,30.",
			Instrucoes2: "ATÉ O VENCIMENTO CONCEDER O DESCONTO DE R$ 15,00",
		},

	}

	boletos := new(goboleto.BoletoSicoob)


	boletos.Emitente = emitente
	boletos.Dados = dadosBoletos

	arquivoboleto := boletos.BoletoPDF()

	f, err := os.Create("boleto_teste.pdf")
	check(err)
	defer f.Close()

	_, err = f.Write([]byte(arquivoboleto))
	check(err)

	//err := ioutil.WriteFile("boleto_teste.pdf", []byte(arquivoboleto), 0644)




}