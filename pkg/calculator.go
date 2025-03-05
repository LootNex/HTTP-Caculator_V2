package calculator

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

type NewTask struct {
	Id        int     `json:"id"`
	Arg1      float64 `json:"arg1"`
	Arg2      float64 `json:"arg2"`
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
	Operation_time time.Duration `json:"operation_time"`
}

var (
	Tasks      []NewTask
	Task_Ready = make(chan float64, 1)
	k          int
)

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, errors.New("пустое выражение")
	}

	var output []string
	var stack []string
	var number string

	for _, char := range expression {
		c := string(char)

		switch {
		case c == " ":
			continue
		case (c >= "0" && c <= "9") || c == ".":
			number += c
		case c == "(":
			stack = append(stack, c)
		case c == ")":
			if number != "" {
				output = append(output, number)
				number = ""
			}
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return 0, errors.New("несоответствие скобок")
			}
			stack = stack[:len(stack)-1]
		case c == "+" || c == "-" || c == "*" || c == "/":
			if number != "" {
				output = append(output, number)
				number = ""
			}
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(c) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)
		default:
			return 0, fmt.Errorf("некорректный символ: %s", c)
		}
	}

	if number != "" {
		output = append(output, number)
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return 0, errors.New("несоответствие скобок")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	var calcStack []float64
	for _, token := range output {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			calcStack = append(calcStack, num)
		} else {
			if len(calcStack) < 2 {
				return 0, errors.New("некорректное выражение")
			}
			b, a := calcStack[len(calcStack)-1], calcStack[len(calcStack)-2]
			calcStack = calcStack[:len(calcStack)-2]
			log.Println("!",k, b,a,token)
			Tasks = append(Tasks, NewTask{Id: k, Arg1: a, Arg2: b, Operation: token})
			log.Println("tasks:",Tasks)
			k++
			select {
			case Tasks[0].Result = <-Task_Ready:
				log.Println("!", Tasks[0].Result)
				calcStack = append(calcStack, Tasks[0].Result)
				Tasks = Tasks[:0]
			case <-time.After(10 * time.Second):
				return 0, errors.New("неверное выражение, ошибка при вычислении")

			}
		}
	}

	if len(calcStack) != 1 {
		return 0, errors.New("некорректное выражение")
	}
	return calcStack[0], nil
}

func precedence(op string) int {
	if op == "+" || op == "-" {
		return 1
	}
	if op == "*" || op == "/" {
		return 2
	}
	return 0
}
