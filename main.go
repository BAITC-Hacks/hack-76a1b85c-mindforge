package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const faqFile = "faq.txt"

type faqItem struct {
	question string
	answer   string
}

var stopWords = map[string]bool{
	"а": true, "в": true, "во": true, "и": true, "на": true, "о": true,
	"об": true, "по": true, "про": true, "с": true, "со": true, "к": true,
	"у": true, "за": true, "для": true, "что": true, "как": true,
	"когда": true, "где": true, "ли": true, "мне": true, "пожалуйста": true,
}

func words(text string) []string {
	parts := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	unique := make(map[string]bool)
	result := make([]string, 0, len(parts))
	for _, word := range parts {
		if word != "" && !stopWords[word] && !unique[word] {
			unique[word] = true
			result = append(result, word)
		}
	}
	return result
}

func loadFAQ(filename string) ([]faqItem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var faq []faqItem
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("неверный формат строки: %q", line)
		}
		faq = append(faq, faqItem{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
	}
	return faq, scanner.Err()
}

func hasSameBeginning(first, second string) bool {
	firstRunes, secondRunes := []rune(first), []rune(second)
	if len(firstRunes) < 4 || len(secondRunes) < 4 {
		return false
	}
	return string(firstRunes[:3]) == string(secondRunes[:3])
}

func findAnswer(question string, faq []faqItem) string {
	questionWords := words(question)
	bestScore := 0
	bestAnswer := ""

	for _, item := range faq {
		score := 0
		for _, userWord := range questionWords {
			for _, faqWord := range words(item.question) {
				if userWord == faqWord || hasSameBeginning(userWord, faqWord) {
					score++
					break
				}
			}
		}
		if score > bestScore {
			bestScore = score
			bestAnswer = item.answer
		}
	}
	return bestAnswer
}

func main() {
	faq, err := loadFAQ(faqFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось прочитать %s: %v\n", faqFile, err)
		os.Exit(1)
	}

	fmt.Println("FAQ-бот репетиции. Задайте вопрос или напишите «выход».")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}
		if strings.EqualFold(question, "выход") || strings.EqualFold(question, "exit") || strings.EqualFold(question, "quit") {
			fmt.Println("До встречи!")
			return
		}

		answer := findAnswer(question, faq)
		if answer == "" {
			fmt.Println("Не знаю.")
		} else {
			fmt.Println(answer)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", err)
	}
}
