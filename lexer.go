package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/golang-collections/collections/stack"
)

// Определение типов токенов
type TokenType string

const (
	Identifier        TokenType = "Identifier"        // Идентификатор
	Number            TokenType = "Number"            // Число
	Semicolon         TokenType = "Semicolon"         // Точка с запятой
	Assign            TokenType = "Assign"            // Присвоение
	Hash              TokenType = "Hash"              // Знак "#"
	Not               TokenType = "Not"               // Знак "!"
	Ampersand         TokenType = "Ampersand"         // Знак "&"
	While             TokenType = "While"             // Цикл while
	LeftCurlyBracket  TokenType = "LeftCurlyBracket"  // Знак "{"
	RightCurlyBracket TokenType = "RightCurlyBracket" // Знак "}"
	EOF               TokenType = "EOF"               // Конец файла
	Error             TokenType = "Error"             // Ошибка
)

// Структура для представления токена
type Token struct {
	Type  TokenType // Тип токена
	Value string    // Значение токена
}

// Структура лексического анализатора
type Lexer struct {
	input    string // Входная строка
	position int    // Текущая позиция в строке
	current  byte   // Текущий символ
}

// Функция создания нового лексического анализатора
func NewLexer(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.ReadChar()
	return lexer
}

// Функция чтения следующего символа
func (l *Lexer) ReadChar() {
	if l.position < len(l.input) {
		l.current = l.input[l.position]
		l.position++
	} else {
		l.current = 0
	}
}

// Функция пропуска пробельных символов
func (l *Lexer) SkipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.ReadChar()
	}
}

// Функция чтения идентификатора
func (l *Lexer) ReadIdentifier() string {
	startPosition := l.position - 1
	for IsLetter(l.current) || IsDigit(l.current) {
		l.ReadChar()
	}
	return l.input[startPosition : l.position-1]
}

// Функция чтения числа
func (l *Lexer) ReadNumber() string {
	startPosition := l.position - 1
	for IsDigit(l.current) {
		l.ReadChar()
	}
	return l.input[startPosition : l.position-1]
}

// Функция получения следующего токена
func (l *Lexer) NextToken() Token {
	l.SkipWhitespace()

	var token Token

	switch l.current {
	case '{':
		token = Token{Type: LeftCurlyBracket, Value: string(l.current)}
	case '}':
		token = Token{Type: RightCurlyBracket, Value: string(l.current)}
	case 'w':
		if strings.HasPrefix(l.input[l.position-1:], "while") {
			token = Token{Type: While, Value: "while"}
			l.position += 4
		}
	case ';':
		token = Token{Type: Semicolon, Value: string(l.current)}
	case ':':
		if l.input[l.position] == '=' {
			token = Token{Type: Assign, Value: ":="}
			l.position += 1
		} else {
			fmt.Println("No symbol ':'")
			error_val := l.current
			token = Token{Type: Error, Value: string(error_val)}
			break
		}
	case '#':
		token = Token{Type: Hash, Value: "#"}
	case '!':
		token = Token{Type: Not, Value: "!"}
	case '&':
		token = Token{Type: Ampersand, Value: "&"}
	case '$':
		token = Token{Type: EOF, Value: "$"}
		return token
	default:
		if IsLetter(l.current) {
			value := l.ReadIdentifier()
			token = Token{Type: Identifier, Value: value}
			return token
		} else if IsDigit(l.current) {
			value := l.ReadNumber()
			token = Token{Type: Number, Value: value}
			return token
		} else {
			token = Token{Type: Error, Value: string(l.current)}
			return token
		}
	}

	l.ReadChar()
	return token
}

// Функция проверки символа на букву
func IsLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

// Функция проверки символа на цифру
func IsDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func StackOfTokens(input string) *stack.Stack {
	stackOfTokens := stack.New()
	temp_array := []TokenType{}
	lexer := NewLexer(input)
	for {
		token := lexer.NextToken()
		temp_array = append(temp_array, token.Type)
		if token.Type == EOF {
			break
		}
		if token.Type == Error {
			fmt.Println(input)
			fmt.Println(strings.Repeat(" ", lexer.position-2), "^")
			fmt.Println("Lexer error at position ", lexer.position)
			return stackOfTokens
		}
	}

	for i, j := 0, len(temp_array)-1; i < j; i, j = i+1, j-1 {
		temp_array[i], temp_array[j] = temp_array[j], temp_array[i]
	}

	for _, val := range temp_array {
		stackOfTokens.Push(string(val))
	}

	return stackOfTokens
}
