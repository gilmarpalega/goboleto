package goboleto

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/jung-kurt/gofpdf"
	"github.com/leekchan/accounting"
)


type BoletoSicoob struct {
	Dados    []DadosBoleto
	Emitente Emitente
}


func (b *BoletoSicoob) Modulo10(num string) string {
	fator := 2
	sum := 0
	dv := 0

	if len(num)%2 == 0 {
		fator = 1
	}

	for _, c := range num {
		value, _ := strconv.Atoi(string(c))
		res := value * fator
		if res > 9 {
			sum += (res - 9)
		} else {
			sum += res
		}

		if fator == 1 {
			fator = 2
		} else {
			fator = 1
		}

	}

	dv = 10 - (sum % 10)
	if dv == 10 {
		dv = 0
	}
	return fmt.Sprintf("%d", dv)
}

func (this *BoletoSicoob) Modulo11(b *DadosBoleto) {
	fator := 4
	acum := 0

	for i, c := range b.CodigoDeBarras {
		if i == 4 {
			continue
		}

		value, _ := strconv.Atoi(string(c))
		acum += value * fator

		fator--

		if fator == 1 {
			fator = 9
		}
	}

	resto := acum % 11
	dv := 11 - resto

	if (dv == 0) || (dv == 1) || (dv > 9) {
		dv = 1
	}

	b.Dv = fmt.Sprintf("%d", dv)
}

func (this *BoletoSicoob) GerarNossoNumeroDv(b *DadosBoleto) {
	fator := []int{3, 1, 9, 7}
	fi := 0 //Fator Index
	acum := 0

	seq := fmt.Sprintf("%04d%010d%07d", this.Emitente.Agencia, this.Emitente.Convenio, b.NossoNumero)

	for _, c := range seq {
		value, _ := strconv.Atoi(string(c))
		acum += value * fator[fi]

		fi++

		if fi == 4 {
			fi = 0
		}
	}

	resto := acum % 11
	dv := 11 - resto

	if dv > 9 {
		dv = 0
	}

	b.NossoNumeroDv = fmt.Sprintf("%d", dv)

}

func (this *BoletoSicoob) Processar(b *DadosBoleto) {
	b.ValorBoleto = b.ValorCobrado + b.TaxaBoleto
	ivb := strings.Replace(strconv.FormatFloat(b.ValorBoleto, 'f', 2, 64), ".", "", 1)
	b.ValorBoletoStr = strings.Repeat("0", 10-len(ivb)) + ivb
	b.FatorData = this.CalcularFatorData(b)

	this.CalcularFatorData(b)
	this.GerarNossoNumeroDv(b)
	this.GerarCodigoDeBarras(b)
	this.GerarLinhaDigitavel(b)
}

func (this *BoletoSicoob) CalcularFatorData(b *DadosBoleto) string {
	t1, _ := time.Parse("2006-01-02", "1997-10-07")
	t2, _ := time.Parse("2006-01-02", Date2_html(b.DataVencimento))
	delta := t2.Sub(t1)
	return fmt.Sprintf("%04d", int(delta.Hours()/24))
}

func (this *BoletoSicoob) GerarCodigoDeBarras(b *DadosBoleto) {
	var f bytes.Buffer

	f.WriteString(this.Emitente.Banco)   // Codigo do Banco
	f.WriteString(this.Emitente.Moeda)   // Codigo da Moeda 9: Real
	f.WriteString(" ")   			 // DV
	f.WriteString(b.FatorData)
	f.WriteString(b.ValorBoletoStr)
	f.WriteString(this.Emitente.Carteira)     // Código da Carteira
	f.WriteString(this.Emitente.Cooperativa)  // Código da Cooperativa
	f.WriteString(this.Emitente.ModalidadeCobranca)                              // Modalidade
	f.WriteString(fmt.Sprintf("%07d", this.Emitente.Convenio))           // Código do Cliente/Convenio
	f.WriteString(fmt.Sprintf("%07d%v", b.NossoNumero, b.NossoNumeroDv)) // Código do Cliente/Convenio
	f.WriteString(b.Sequencia) // Número da Parcela

	b.CodigoDeBarras = f.String()
	this.Modulo11(b)
	cdb := strings.Split(b.CodigoDeBarras, " ")
	b.CodigoDeBarras = fmt.Sprintf("%v%v%v", cdb[0], b.Dv, cdb[1])
}

func (this *BoletoSicoob) GerarLinhaDigitavel(b *DadosBoleto) {
	var g1, g2, g3, g4 bytes.Buffer

	g1.WriteString(this.Emitente.Banco)        // Codigo do Banco
	g1.WriteString(this.Emitente.Moeda)        // Codigo da Moeda 9: Real
	g1.WriteString(this.Emitente.Carteira)     // Código da Carteira
	g1.WriteString(this.Emitente.Cooperativa)  // Código da Cooperativa
	g1.WriteString(this.Modulo10(g1.String())) //Adicionar Dv do g1

	g2.WriteString(this.Emitente.ModalidadeCobranca)   // Modalidade de Cobrança
	g2.WriteString(fmt.Sprintf("%07d", this.Emitente.Convenio)) // Código do Cliente/Convenio com Dígito

	nn := fmt.Sprintf("%07d%v", b.NossoNumero, b.NossoNumeroDv)
	nn1 := nn[0:1]
	nn2 := nn[1:8]

	g2.WriteString(nn1)                     //Primeiro caracter do Nosso Numero Com DV
	g2.WriteString(this.Modulo10(g2.String())) // DV g2

	g3.WriteString(nn2) //Restante do Nosso número - Últimos 7 dígitos
	g3.WriteString(b.Sequencia)
	g3.WriteString(this.Modulo10(g3.String())) // DV g3

	g4.WriteString(b.FatorData)
	g4.WriteString(b.ValorBoletoStr)

	g1s := g1.String()
	g2s := g2.String()
	g3s := g3.String()
	fmt.Println(g1s, g2s, g3s, g4.String())
	fmt.Println("Dv", b.Dv)
	fmt.Println("NossoNumeroDv", b.NossoNumeroDv)

	b.LinhaDigitavel = fmt.Sprintf("%v.%v %v.%v %v.%v %v %v", g1s[0:5], g1s[5:10], g2s[0:5], g2s[5:11], g3s[0:5], g3s[5:11], b.Dv, g4.String())
	fmt.Println(b.LinhaDigitavel)
}


func (this *BoletoSicoob) BoletoPDF() ([]byte){
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(5, 5, 5)
	pdf.SetRightMargin(5)
	pdf.SetAutoPageBreak(true, 0)
	pdf.SetFillColor(236, 236, 236)
	pdf.SetLineWidth(0.7)
	// First page: manual local link
	pdf.AddPage()
	defer pdf.Close()
	pdf.SetFont("Helvetica", "", 12)
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	lin := 0.0

	draw_horizontal_lines := func(q BoletoHorizontalLines) {
		pdf.SetLineWidth(q.w)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Line(q.x1, q.y+lin, q.x2, q.y+lin)
	}

	draw_vertical_lines := func(q BoletoVerticalLines) {
		pdf.SetLineWidth(q.w)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Line(q.x, q.y1+lin, q.x, q.y2+lin)
	}

	textos := func(q quadro) {
		fsize := 6.0			// Fonte padrão é 6, se não informada
		if q.fontsize > 0 {
			fsize = q.fontsize
		}

		pdf.SetFont("Helvetica", q.valor, fsize)
		ht := pdf.PointConvert(fsize)

		wd := pdf.GetStringWidth(q.campo)
		pdf.SetY(q.y - 0.2 + lin)
		pdf.SetX(q.x)
		pdf.CellFormat(wd, ht, tr(q.campo), "", 0, "L", false, 0, "")
	}

	dados := func(q quadro) {
		fsize := 10.0

		pdf.SetFont("Helvetica", "B", fsize)
		ht := pdf.PointConvert(fsize)

		//wd := pdf.GetStringWidth(q.campo)
		pdf.SetY(q.y - 0.2 + lin)
		pdf.SetX(q.x)
		pdf.CellFormat(q.fontsize, ht, tr(q.campo), "", 0, "R", false, 0, "")
	}

	imprime_vertical := func(q quadro) {

		fsize := 6.0

		if q.fontsize > 0 {
			fsize = q.fontsize
		}

		pdf.SetFont("Helvetica", q.valor, fsize)
		pdf.TransformBegin()

		//pdf.SetFont("Helvetica", "", 6.0)
		wd := pdf.GetStringWidth(q.campo)
		ht := pdf.PointConvert(fsize)
		pdf.TransformRotate(90, q.x, q.y+lin)
		pdf.SetY(q.y + lin)
		pdf.SetX(q.x)
		pdf.CellFormat(wd, ht, tr(q.campo), "", 0, "L", false, 0, "")
		pdf.TransformEnd()
	}

	horizontal_lines := []BoletoHorizontalLines{
		{x1: 3.7, y: 3.4, x2: 19.9, w: 0.2},
		{x1: 28.3, y: 3.4, x2: 41.8, w: 0.2},
		{x1: 3.7, y: 29, x2: 12.5, w: 0.2},
		{x1: 12.5, y: 40.1, x2: 19.8, w: 0.2},
		{x1: 3.7, y: 94.8, x2: 19.9, w: 0.2},
		{x1: 28.3, y: 94.8, x2: 41.8, w: 0.2},

		{x1: 48.7, y: 9.3, x2: 201.5, w: 0.5},
		{x1: 48.7, y: 15.6, x2: 201.5, w: 0.2},
		{x1: 48.7, y: 21.9, x2: 201.5, w: 0.2},
		{x1: 48.7, y: 28.2, x2: 201.5, w: 0.2},
		{x1: 48.7, y: 34.5, x2: 201.5, w: 0.2},
		{x1: 48.7, y: 65.7, x2: 201.5, w: 0.5},
		{x1: 48.7, y: 78.7, x2: 201.5, w: 0.5},

		{x1: 167.2, y: 40.8, x2: 201.5, w: 0.2},
		{x1: 167.2, y: 46.9, x2: 201.5, w: 0.2},
		{x1: 167.2, y: 53.3, x2: 201.5, w: 0.2},
		{x1: 167.2, y: 59.5, x2: 201.5, w: 0.2},
	}

	vertical_lines := []BoletoVerticalLines{
		{x: 3.7, y1: 3.4, y2: 94.8, w: 1.0},
		{x: 12.5, y1: 3.4, y2: 94.8, w: 0.2},
		{x: 19.9, y1: 3.4, y2: 94.8, w: 0.2},
		{x: 28.3, y1: 3.4, y2: 94.8, w: 0.2},
		{x: 34.2, y1: 3.4, y2: 94.8, w: 0.2},
		{x: 41.8, y1: 3.4, y2: 94.8, w: 0.2},
		{x: 45.4, y1: 5.4, y2: 92.8, w: 0.2},

		{x: 78.6, y1: 3.4, y2: 9.3, w: 0.2},
		{x: 78.6, y1: 21.9, y2: 34.5, w: 0.2},
		{x: 112.9, y1: 21.9, y2: 28.2, w: 0.2},
		{x: 127.7, y1: 21.9, y2: 28.2, w: 0.2},
		{x: 137.6, y1: 21.9, y2: 28.2, w: 0.2},
		{x: 167.0, y1: 9.3, y2: 65.7, w: 0.2},
		{x: 98.0, y1: 28.2, y2: 34.5, w: 0.2},
		{x: 117.8, y1: 28.2, y2: 34.5, w: 0.2},
		{x: 142.4, y1: 28.2, y2: 29.9, w: 0.2},
		{x: 142.4, y1: 32.9, y2: 34.5, w: 0.2},
	}

	textos_boleto := []quadro{
		{campo: "Local de Pagamento", x: 48.7, y: 9.8, valor: "", fontsize: 6.0},
		{campo: "Beneficiário", x: 48.7, y: 16.1},
		{campo: "Data do Documento", x: 48.7, y: 22.4},
		{campo: "Nº do Documento", x: 78.5, y: 22.4},
		{campo: "Espécie", x: 112.5, y: 22.4},
		{campo: "Aceite", x: 127.4, y: 22.4},
		{campo: "Data do Processamento", x: 137.4, y: 22.4},
		{campo: "Uso do Banco", x: 48.7, y: 28.7},
		{campo: "Carteira", x: 78.5, y: 28.7},
		{campo: "Espécie Moeda", x: 97.7, y: 28.7},
		{campo: "Quantidade Moeda", x: 117.5, y: 28.7},
		{campo: "Valor Moeda", x: 142.2, y: 28.7},
		{campo: "X", x: 140.8, y: 30.6},
		{campo: "Instruções (Texto sob responsabilidade do beneficiário)", x: 48.7, y: 35.0},

		{campo: "Vencimento", x: 166.7, y: 9.8},
		{campo: "Coop.contratante / Cód.Beneficiário", x: 166.7, y: 16.1},
		{campo: "Nosso Nº / Código do Documento", x: 166.7, y: 22.4},
		{campo: "(=) Valor do Documento", x: 166.7, y: 28.7},
		{campo: "(-) Desconto / Abatimento", x: 166.7, y: 35.0},
		{campo: "(-) Outras Deduções", x: 166.7, y: 41.3},
		{campo: "(+) Mora / Multa", x: 166.7, y: 47.6},
		{campo: "(+) Outros Acréscimentos", x: 166.8, y: 53.9},
		{campo: "(=) Valor Cobrado", x: 166.7, y: 60.2},

		{campo: "PAGADOR", x: 48.7, y: 66.3, fontsize: 6, valor: ""},
		{campo: "SACADOR / AVALISTA", x: 48.7, y: 76.3, fontsize: 6, valor: ""},
		{campo: "Código de Baixa", x: 149.5, y: 76.3, fontsize: 6},
		{campo: "Autenticação Mecânica -", x: 149.5, y: 79.4, fontsize: 6},
		{campo: "FICHA DE COMPENSAÇÃO", x: 173.1, y: 79.4, fontsize: 6, valor: "B"},
	}

	texto_canhoto_boleto := []quadro{
		{campo: "Vencimento", x: 4.6, y: 27.5},
		{campo: "(=) Valor do Documento", x: 13.1, y: 39.0},
		{campo: "Pagador", x: 4.6, y: 94.0},
		{campo: "Nº do Documento", x: 13.1, y: 94.0},
		{campo: "Beneficiário", x: 20.6, y: 94.0},
		{campo: "Coop. / Cód. Beneficiário", x: 20.6, y: 46.8},
		{campo: "Nosso Número", x: 20.6, y: 21.2},
	}

	//parcelas := []int{21013, 21014, 21015, 20000, 20500, 20501, 20502, 20530, 10, 1000}

	ac := accounting.Accounting{Symbol: "R$ ", Precision: 2, Thousand: ".", Decimal: ","}

	for _, boleto := range this.Dados {
		this.Processar(&boleto)

		dadossicoob := []quadro{
			{campo: "PAGÁVEL EM QUALQUER BANCO ATÉ O VENCIMENTO", x: 48.7, y: 12.2, fontsize: 115},
			{campo: fmt.Sprintf("%v (%v)",this.Emitente.Nome, this.Emitente.Cnpj), x: 48.7, y: 18.5, fontsize: 115},
			{campo: Date2_str_br(boleto.DataVencimento), x: 169.4, y: 12.2, fontsize: 32},
			{campo: fmt.Sprintf("%v / %v", this.Emitente.Agencia, this.Emitente.Codigo), x: 169.4, y: 18.5, fontsize: 32},
			{campo: fmt.Sprintf("%07d-%v", boleto.NossoNumero, boleto.NossoNumeroDv), x: 169.4, y: 24.8, fontsize: 32},
			{campo: ac.FormatMoney(boleto.ValorBoleto), x: 169.4, y: 31.1, fontsize: 32},
			{campo: Date2_str_br(boleto.DataEmissao), x: 48.7, y: 24.8, fontsize: 28},
			{campo: fmt.Sprintf("Parc. - %v/%v", boleto.SeqBoleto, boleto.TotalBoletos), x: 80, y: 24.8, fontsize: 32},
			{campo: "RC", x: 120, y: 24.8, fontsize: 4},
			{campo: "N", x: 133, y: 24.8, fontsize: 2},
			{campo: Date2_str_br(time.Now()), x: 137.0, y: 24.8, fontsize: 28},
			{campo: "1", x: 82, y: 31.1, fontsize: 10},
			{campo: "R$", x: 101, y: 31.1, fontsize: 10},
		}

		dadoscanhoto := []quadro{
			{campo: boleto.SacadoNome, x: 7.9, y: 93.5, fontsize: 9, valor: "B"},
			{campo: fmt.Sprintf("Parc. - %v/%v", boleto.SeqBoleto, boleto.TotalBoletos), x: 16.1, y: 93.5, fontsize: 9, valor: "B"},
			{campo: this.Emitente.Nome, x: 23.6, y: 93.5, fontsize: 8, valor: "B"},
			{campo: Date2_str_br(boleto.DataVencimento), x: 7.9, y: 25.7, fontsize: 9, valor: "B"},
			{campo: ac.FormatMoney(boleto.ValorBoleto), x: 16.1, y: 36.7, fontsize: 9, valor: "B"},
			{campo: fmt.Sprintf("%v / %v", this.Emitente.Agencia, this.Emitente.Codigo), x: 23.6, y: 45.0, fontsize: 9, valor: "B"},
			{campo: fmt.Sprintf("%07d-%v", boleto.NossoNumero, boleto.NossoNumeroDv), x: 23.6, y: 21.3, fontsize: 9, valor: "B"},
		}

		pdf.SetFillColor(236, 236, 236)
		pdf.Rect(167, 9.3+lin, 201.5-167, 6.3, "F")
		pdf.Rect(167, 28.2+lin, 201.5-167, 6.3, "F")

		for _, b := range horizontal_lines {
			draw_horizontal_lines(b)
		}

		for _, b := range vertical_lines {
			draw_vertical_lines(b)
		}

		// DADOS DO CANHOTO
		for _, b := range texto_canhoto_boleto {
			imprime_vertical(b)
		}

		for _, b := range dadoscanhoto {
			imprime_vertical(b)
		}

		//imprime_vertical(quadro{campo: "Vencimento", x: 4.6, y: 27.6})

		//pdf.SetFont("Helvetica", "", 6)
		for _, b := range textos_boleto {
			textos(b)
		}

		textos(quadro{campo: "756-0", x: 79.2, y: 4.1, fontsize: 14, valor: "B"})
		textos(quadro{campo: boleto.LinhaDigitavel, x: 94.0, y: 4.6, fontsize: 10.5, valor: "B"})

		for _, b := range dadossicoob {
			dados(b)
		}

		textos(quadro{campo: boleto.Instrucoes1, x: 50.0, y: 39, fontsize: 10, valor: "B"})
		textos(quadro{campo: boleto.Instrucoes2, x: 50.0, y: 44.0, fontsize: 10, valor: "B"})

		textos(quadro{campo: boleto.Pag1, x: 61.0, y: 66.5, fontsize: 9.5, valor: "B"})
		textos(quadro{campo: boleto.Pag2, x: 61.0, y: 69.9, fontsize: 9.5, valor: "B"})
		textos(quadro{campo: boleto.Pag3, x: 61.0, y: 73.3, fontsize: 9.5, valor: "B"})

		pdf.Image("logobancoob.jpg", 50, 3.5+lin, 26, 5, false, "", 0, "")

		codbar, errobarra := GerarBarcode2of5(boleto.CodigoDeBarras)
		if errobarra == nil {

			urlBarCode := fmt.Sprintf("localhost/barcode/%v", boleto.NossoNumero) // esse nome é ilustrativo, o gofpdf precisa de uma referência nomeada para as imagens
			opbarcode := gofpdf.ImageOptions{ImageType: "PNG"}

			_ = pdf.RegisterImageOptionsReader(urlBarCode, opbarcode, codbar)

			if pdf.Ok() {
				pdf.Image(urlBarCode, 48, 82+lin, 128, 13, false, "PNG", 0, "")
			}
		}

		lin += 98

		if lin < 200 {
			draw_horizontal_lines(BoletoHorizontalLines{x1: 10, y: 0, x2: 201.5, w: 0.2})
		}

		if lin > 230 {
			pdf.AddPage()
			lin = 0
		}
	}

	var outputBoletoPDF bytes.Buffer

	pdf.Output(&outputBoletoPDF)

	return outputBoletoPDF.Bytes()

}
