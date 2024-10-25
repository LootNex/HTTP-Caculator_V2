package main

import (
 "errors"
 "fmt"
 "strconv"
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

   switch token {
   case "+":
    calcStack = append(calcStack, a+b)
   case "-":
    calcStack = append(calcStack, a-b)
   case "*":
    calcStack = append(calcStack, a*b)
   case "/":
    if b == 0 {
     return 0, errors.New("деление на ноль")
    }
    calcStack = append(calcStack, a/b)
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