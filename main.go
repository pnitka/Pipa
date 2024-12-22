package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func calculate(expression string) (float64, error) {
	// Простой парсер, поддерживающий только числа и операции
	result, err := eval(expression)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func eval(expr string) (float64, error) {
	operations := strings.Fields(expr)
	if len(operations) == 0 {
		return 0, fmt.Errorf("недействительное выражение")
	}
	result := 0.0
	var operator string
	for _, op := range operations {
		if op == "+" || op == "-" || op == "*" || op == "/" {
			operator = op
		} else {
			num, err := strconv.ParseFloat(op, 64)
			if err != nil {
				return 0, fmt.Errorf("недействительное число: %s", op)
			}
			switch operator {
			case "":
				result = num
			case "+":
				result += num
			case "-":
				result -= num
			case "*":
				result *= num
			case "/":
				if num == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				result /= num
			}
		}
	}
	return result, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка при чтении запроса", http.StatusBadRequest)
		return
	}

	result, err := calculate(req.Expression)
	if err != nil {
		response := Response{Error: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{Result: fmt.Sprintf("%f", result)}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/calculate", handler)
	fmt.Println("Сервер запущен на порту 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}