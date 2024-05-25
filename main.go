package main

import (
	"bufio"
	"example/simple-precedence-parser/lexer"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-collections/collections/stack"
)

// Type for non terminals
type NonTerminals string

// Non terminals themselves
const (
	Program     NonTerminals = "Program"
	Block       NonTerminals = "Block"
	_Block      NonTerminals = "_Block"
	Operator    NonTerminals = "Operator"
	Variable    NonTerminals = "Variable"
	Expression  NonTerminals = "Expression"
	_Expression NonTerminals = "_Expression"
	Factor      NonTerminals = "Factor"
	Primary     NonTerminals = "Primary"
)

// Parsing table
var parsing_table = map[string]map[string]rune{
	string(Program): {
		string(lexer.EOF): '#', // '#' = финиш
	},
	string(lexer.Identifier): {
		string(lexer.LeftCurlyBracket):  '>',
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.Assign):            '>',
		string(lexer.Hash):              '>',
		string(lexer.Ampersand):         '>',
	},
	string(Block): {
		string(lexer.RightCurlyBracket): '=',
	},
	string(_Block): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '=',
	},
	string(Operator): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
	},
	string(Variable): {
		string(lexer.LeftCurlyBracket): '=',
		string(lexer.Assign):           '=',
	},
	string(Expression): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
	},
	string(_Expression): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.Hash):              '=',
	},
	string(Factor): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.Hash):              '>',
		string(lexer.Ampersand):         '=',
	},
	string(lexer.Number): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.Hash):              '>',
		string(lexer.Ampersand):         '>',
	},
	string(Primary): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.Hash):              '>',
		string(lexer.Ampersand):         '>',
	},
	string(lexer.LeftCurlyBracket): {
		string(lexer.Identifier): '<',
		string(Block):            '=',
		string(_Block):           '<',
		string(Operator):         '<',
		string(Variable):         '<',
		string(lexer.While):      '<',
	},
	string(lexer.RightCurlyBracket): {
		string(lexer.RightCurlyBracket): '>',
		string(lexer.Semicolon):         '>',
		string(lexer.EOF):               '>',
	},
	string(lexer.While): {
		string(lexer.Identifier): '<',
		string(Variable):         '=',
	},
	string(lexer.Semicolon): {
		string(lexer.Identifier): '<',
		string(Block):            '=',
		string(_Block):           '=',
		string(Operator):         '=',
		string(Variable):         '<',
		string(lexer.While):      '<',
	},
	string(lexer.Assign): {
		string(lexer.Identifier): '<',
		string(Expression):       '=',
		string(_Expression):      '<',
		string(Factor):           '<',
		string(lexer.Number):     '<',
		string(Primary):          '<',
		string(lexer.Not):        '<',
	},
	string(lexer.Not): {
		string(lexer.Identifier): '<',
		string(Factor):           '=',
		string(lexer.Number):     '<',
		string(Primary):          '<',
	},
	string(lexer.Hash): {
		string(lexer.Identifier): '<',
		string(Expression):       '=',
		string(_Expression):      '=',
		string(Factor):           '=',
		string(lexer.Number):     '<',
		string(Primary):          '<',
	},
	string(lexer.Ampersand): {
		string(lexer.Identifier): '<',
		string(lexer.Number):     '<',
		string(Primary):          '=',
	},
	string(lexer.EOF): {
		string(Program):                '<',
		string(lexer.LeftCurlyBracket): '<',
	},
}

var reduce_table = map[string]map[string]string{
	string(lexer.LeftCurlyBracket): {
		string(lexer.Identifier): string(Variable),
	},
	string(lexer.RightCurlyBracket): {
		string(lexer.Identifier):        string(Primary),
		string(_Block):                  string(Block),
		string(Operator):                string(_Block),
		string(Expression):              string(Operator),
		string(_Expression):             string(Expression),
		string(Factor):                  string(_Expression),
		string(lexer.Number):            string(Primary),
		string(Primary):                 string(Factor),
		string(lexer.RightCurlyBracket): string(Operator),
	},
	string(lexer.Semicolon): {
		string(lexer.Identifier):        string(Primary),
		string(Operator):                string(_Block),
		string(Expression):              string(Operator),
		string(_Expression):             string(Expression),
		string(Factor):                  string(_Expression),
		string(lexer.Number):            string(Primary),
		string(Primary):                 string(Factor),
		string(lexer.RightCurlyBracket): string(Operator),
	},
	string(lexer.Assign): {
		string(lexer.Identifier): string(Variable),
	},
	string(lexer.Hash): {
		string(lexer.Identifier): string(Primary),
		string(Factor):           string(_Expression),
		string(lexer.Number):     string(Primary),
		string(Primary):          string(Factor),
	},
	string(lexer.Ampersand): {
		string(lexer.Identifier): string(Primary),
		string(lexer.Number):     string(Primary),
		string(Primary):          string(Factor),
	},
	string(lexer.EOF): {
		string(lexer.RightCurlyBracket): string(Program),
	},
}

// Read file
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func algorithm(input string) {
	work_stack := stack.New()
	work_stack.Push([]interface{}{string(lexer.EOF), rune(0)})

	stack_of_tokens := lexer.StackOfTokens(input)

	for {
		fmt.Println("start of cycle")
		if work_stack.Peek().([]interface{})[0].(string) == string(Program) && stack_of_tokens.Peek().(string) == string(lexer.EOF) {
			fmt.Println("end of calculating")
			break
		}
		fmt.Println(work_stack.Peek().([]interface{})[0].(string))
		fmt.Println(stack_of_tokens.Peek().(string))
		if val, ok := parsing_table[work_stack.Peek().([]interface{})[0].(string)][stack_of_tokens.Peek().(string)]; ok {
			if val == '<' || val == '=' {
				fmt.Println("val is ", string(val))
				work_stack.Push([]interface{}{stack_of_tokens.Pop().(string), val})
			}
			if val == '>' {
				fmt.Println("val is >")
				if val, ok := reduce_table[stack_of_tokens.Peek().(string)][work_stack.Peek().([]interface{})[0].(string)]; ok {
					stack_of_tokens.Push(val)
				}
				for {
					if work_stack.Peek().([]interface{})[1].(rune) == '<' {
						work_stack.Pop()
						break
					}
					work_stack.Pop()
				}
			}
		}
	}
}

func main() {
	inputs, err := readLines("examples.txt")
	if err != nil {
		log.Fatalf("ReadLines: %s", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for i, input := range inputs {
		fmt.Println(strings.Repeat("-", 80))
		algorithm(input)
		fmt.Println(input)
		fmt.Println(strings.Repeat("-", 80))

		if i == len(inputs)-1 {
			break
		}
		_, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Great success!")
}
