package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/shopspring/decimal"

	"github.com/thenopholo/my-money/internal/domain"
)

// maxCSVSize é o tamanho máximo do CSV em bytes (5 MB)
const maxCSVSize = 5 * 1024 * 1024

// parseDateFormats são os formatos de data suportados para CSV bancários brasileiros
var parseDateFormats = []string{
	"02/01/2006", // DD/MM/YYYY
	"2006-01-02", // YYYY-MM-DD
	"02/01/06",   // DD/MM/YY
}

// installmentRegex captura parcelas no formato "1/3", "02/12", etc.
var installmentRegex = regexp.MustCompile(`(\d{1,2})/(\d{1,2})`)

// ParseBankCSV parseia um CSV de extrato bancário e retorna transações brutas.
func ParseBankCSV(reader io.Reader) ([]domain.RawCSVTransaction, error) {
	records, sep, err := readCSVRecords(reader)
	if err != nil {
		return nil, err
	}
	_ = sep

	if len(records) == 0 {
		return nil, domain.ErrEmptyCSV
	}

	// Detecta cabeçalho e pula se necessário
	startIdx := 0
	if isHeaderRow(records[0]) {
		startIdx = 1
	}

	if len(records) <= startIdx {
		return nil, domain.ErrEmptyCSV
	}

	var transactions []domain.RawCSVTransaction
	for i := startIdx; i < len(records); i++ {
		row := records[i]

		// Ignora linhas vazias
		if isEmptyRow(row) {
			continue
		}

		// Extrato bancário: mínimo 3 colunas (data, descrição, valor)
		if len(row) < 3 {
			return nil, fmt.Errorf("line %d: %w: expected at least 3 columns, got %d", i+1, domain.ErrInvalidCSVFormat, len(row))
		}

		date, err := parseCSVDate(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w: invalid date %q", i+1, domain.ErrInvalidCSVFormat, row[0])
		}

		description := strings.TrimSpace(row[1])
		if description == "" {
			return nil, fmt.Errorf("line %d: %w: empty description", i+1, domain.ErrInvalidCSVFormat)
		}

		amount, err := parseCSVAmount(strings.TrimSpace(row[2]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w: invalid amount %q", i+1, domain.ErrInvalidCSVFormat, row[2])
		}

		txType := domain.TransactionTypeIncome
		if amount.LessThan(decimal.Zero) {
			txType = domain.TransactionTypeExpense
			amount = amount.Abs()
		}

		transactions = append(transactions, domain.RawCSVTransaction{
			Description:     description,
			Amount:          amount,
			TransactionDate: date,
			TransactionType: txType,
		})
	}

	if len(transactions) == 0 {
		return nil, domain.ErrEmptyCSV
	}

	return transactions, nil
}

// ParseCreditCardCSV parseia um CSV de fatura de cartão de crédito e retorna transações brutas.
func ParseCreditCardCSV(reader io.Reader) ([]domain.RawCSVTransaction, error) {
	records, sep, err := readCSVRecords(reader)
	if err != nil {
		return nil, err
	}
	_ = sep

	if len(records) == 0 {
		return nil, domain.ErrEmptyCSV
	}

	// Detecta cabeçalho e pula se necessário
	startIdx := 0
	if isHeaderRow(records[0]) {
		startIdx = 1
	}

	if len(records) <= startIdx {
		return nil, domain.ErrEmptyCSV
	}

	var transactions []domain.RawCSVTransaction
	for i := startIdx; i < len(records); i++ {
		row := records[i]

		// Ignora linhas vazias
		if isEmptyRow(row) {
			continue
		}

		// Fatura de cartão: mínimo 3 colunas (data, descrição, valor), parcela é opcional (4ª coluna)
		if len(row) < 3 {
			return nil, fmt.Errorf("line %d: %w: expected at least 3 columns, got %d", i+1, domain.ErrInvalidCSVFormat, len(row))
		}

		date, err := parseCSVDate(strings.TrimSpace(row[0]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w: invalid date %q", i+1, domain.ErrInvalidCSVFormat, row[0])
		}

		description := strings.TrimSpace(row[1])
		if description == "" {
			return nil, fmt.Errorf("line %d: %w: empty description", i+1, domain.ErrInvalidCSVFormat)
		}

		amount, err := parseCSVAmount(strings.TrimSpace(row[2]))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w: invalid amount %q", i+1, domain.ErrInvalidCSVFormat, row[2])
		}

		// Transações de cartão são sempre despesas
		if amount.LessThan(decimal.Zero) {
			amount = amount.Abs()
		}

		tx := domain.RawCSVTransaction{
			Description:     description,
			Amount:          amount,
			TransactionDate: date,
			TransactionType: domain.TransactionTypeExpense,
		}

		// Tenta extrair parcelas (coluna 4 se existir, ou do campo descrição)
		if len(row) >= 4 {
			installmentStr := strings.TrimSpace(row[3])
			if installmentStr != "" {
				current, total, ok := parseInstallments(installmentStr)
				if ok {
					tx.Installments = &total
					tx.CurrentInstall = &current
				}
			}
		}

		// Se não encontrou parcelas na coluna, tenta extrair da descrição
		if tx.Installments == nil {
			current, total, ok := parseInstallmentsFromDescription(description)
			if ok {
				tx.Installments = &total
				tx.CurrentInstall = &current
			}
		}

		transactions = append(transactions, tx)
	}

	if len(transactions) == 0 {
		return nil, domain.ErrEmptyCSV
	}

	return transactions, nil
}

// readCSVRecords lê o CSV tentando diferentes separadores e retorna os registros.
func readCSVRecords(reader io.Reader) ([][]string, rune, error) {
	content, err := io.ReadAll(io.LimitReader(reader, maxCSVSize+1))
	if err != nil {
		return nil, 0, fmt.Errorf("reading CSV: %w", err)
	}

	if len(content) > maxCSVSize {
		return nil, 0, domain.ErrCSVTooLarge
	}

	if len(content) == 0 {
		return nil, 0, domain.ErrEmptyCSV
	}

	// Limpa BOM UTF-8 se presente
	contentStr := strings.TrimPrefix(string(content), "\xef\xbb\xbf")

	// Detecta separador: conta ocorrências de ; e , na primeira linha
	sep := detectSeparator(contentStr)

	csvReader := csv.NewReader(strings.NewReader(contentStr))
	csvReader.Comma = sep
	csvReader.TrimLeadingSpace = true
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1 // Permite número variável de campos

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", domain.ErrInvalidCSVFormat, err)
	}

	return records, sep, nil
}

// detectSeparator detecta o separador do CSV (`;` ou `,`).
func detectSeparator(content string) rune {
	firstLine := content
	if idx := strings.IndexByte(content, '\n'); idx != -1 {
		firstLine = content[:idx]
	}

	semicolons := strings.Count(firstLine, ";")
	commas := strings.Count(firstLine, ",")

	if semicolons > commas {
		return ';'
	}
	return ','
}

// isHeaderRow verifica se a linha parece ser um cabeçalho (contém palavras-chave comuns).
func isHeaderRow(row []string) bool {
	if len(row) == 0 {
		return false
	}

	headerKeywords := []string{
		"data", "date", "descrição", "descricao", "description",
		"valor", "value", "amount", "parcela", "installment",
	}

	for _, cell := range row {
		cellLower := strings.ToLower(strings.TrimSpace(cell))
		for _, keyword := range headerKeywords {
			if strings.Contains(cellLower, keyword) {
				return true
			}
		}
	}
	return false
}

// isEmptyRow verifica se todas as células da linha são vazias.
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// parseCSVDate tenta parsear data em diferentes formatos brasileiros.
func parseCSVDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	for _, format := range parseDateFormats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %s", value)
}

// parseCSVAmount converte string de valor para decimal, suportando formatos BR e US.
func parseCSVAmount(value string) (decimal.Decimal, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return decimal.Zero, fmt.Errorf("empty amount")
	}

	// Remove espaços e símbolo de moeda
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "R$", "")
	value = strings.TrimSpace(value)

	// Detecta formato BR: 1.234,56 → separador de milhar é `.` e decimal é `,`
	// Detecta formato US: 1,234.56 → separador de milhar é `,` e decimal é `.`
	if isBrazilianFormat(value) {
		value = strings.ReplaceAll(value, ".", "")
		value = strings.Replace(value, ",", ".", 1)
	}

	return decimal.NewFromString(value)
}

// isBrazilianFormat detecta se o valor está em formato brasileiro (1.234,56).
func isBrazilianFormat(value string) bool {
	// Remove sinal negativo para análise
	clean := strings.TrimPrefix(value, "-")

	commaIdx := strings.LastIndex(clean, ",")
	dotIdx := strings.LastIndex(clean, ".")

	// Se tem vírgula e não tem ponto, ou vírgula vem depois do ponto → BR
	if commaIdx != -1 && dotIdx == -1 {
		return true
	}
	if commaIdx != -1 && dotIdx != -1 && commaIdx > dotIdx {
		return true
	}
	return false
}

// parseInstallments parseia string de parcela no formato "1/3".
func parseInstallments(value string) (current int, total int, ok bool) {
	matches := installmentRegex.FindStringSubmatch(value)
	if len(matches) != 3 {
		return 0, 0, false
	}

	current, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, false
	}

	total, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, false
	}

	if current < 1 || total < 1 || current > total {
		return 0, 0, false
	}

	return current, total, true
}

// parseInstallmentsFromDescription tenta extrair informação de parcela da descrição.
func parseInstallmentsFromDescription(description string) (current int, total int, ok bool) {
	// Remove espaços e transforma tudo em ASCII para facilitar busca
	clean := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}, description)

	matches := installmentRegex.FindStringSubmatch(clean)
	if len(matches) != 3 {
		return 0, 0, false
	}

	current, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, 0, false
	}

	total, err = strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, false
	}

	if current < 1 || total < 1 || current > total {
		return 0, 0, false
	}

	return current, total, true
}
