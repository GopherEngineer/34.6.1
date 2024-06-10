package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Для удобства преобразования строки в число
// без возращения ошибки. В случае ошибки преобразования
// просто возвращаем 0.
func atoi(arg string) int {
	v, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	}

	return v
}

// Производим простые математические вычисления
// с проверкой деления на 0.
func calculator(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	default:
		return 0, nil
	}
}

func main() {

	// Ожидаем указания имени файла для считывания математических инструкций,
	// а также имени файла для записи результата.
	if len(os.Args) < 3 {
		log.Fatalln("запуск программы ожидает указания имени файлов ввода и вывода")
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("входной файл не найден")
	}

	// Создаем/открываем файл. Очищаем если что-то было
	// ранее записано в файле.
	output, err := os.OpenFile(os.Args[2], os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln("невозможно создать/открыть файла для вывода результатов")
	}

	defer output.Close()

	// Для буферизированной записи результатов
	buf_output := bufio.NewWriter(output)

	// Преобразуем прочитанное и делим на отдельные математические инструкции.
	lines := strings.Split(string(input), "\n")

	// Регулярное выражение для проверки на шаблон математических инструкций,
	// а также для поиска аргументов и инструкций для дальнейшего вычисления.
	reg := regexp.MustCompile(`(\d+)([+\-*/])(\d+)=\?`)

	for index, line := range lines {

		// Ищем математическию инструкцию по регулярному выражению.
		match := reg.FindStringSubmatch(line)

		if len(match) > 0 {
			result, err := calculator(atoi(match[1]), atoi(match[3]), match[2])

			// Добавление новой строки в файл
			// только если будет следующая запись.
			if index > 0 && index != len(lines) {
				buf_output.Write([]byte("\n"))
			}

			var out string

			if err != nil {
				// В случае ошибки вычисления математического выражения
				// пишем текст ошибки на место результата для удобства
				// понимания где именно произошла ошибка.
				out = strings.Replace(line, "?", err.Error(), 1)
			} else {
				// В случае успешного вычисления заменяем знак "?" на результат вычисления.
				out = strings.Replace(line, "?", strconv.Itoa(result), 1)
			}

			buf_output.Write([]byte(out))
		}
	}

	// Пишем в файл.
	buf_output.Flush()

}
