package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrProjectName сигнализирует о некорректном имени проекта.
	ErrProjectName = errors.New("имя проекта может содержать только буквы, цифры, -, _")
)

var projectNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// ProjectName считывает имя проекта и проверяет его на соответствие правилам.
func ProjectName(reader *bufio.Reader) (string, error) {
	fmt.Print("Name of the project: ")

	project, _ := reader.ReadString('\n')
	project = strings.TrimSpace(project)

	if !projectNameRegexp.MatchString(project) {
		return project, fmt.Errorf("%w", ErrProjectName)
	}

	return project, nil
}

// GithubHandle считывает имя пользователя GitHub.
func GithubHandle(reader *bufio.Reader) string {
	fmt.Print("Github username: ")

	name, _ := reader.ReadString('\n')
	return strings.TrimSpace(name)
}

// Description считывает описание проекта.
func Description(reader *bufio.Reader) string {
	fmt.Print("Project description: ")

	description, _ := reader.ReadString('\n')
	return strings.TrimSpace(description)
}

// Confirm спрашивает у пользователя подтверждение.
func Confirm(reader *bufio.Reader, question string) bool {
	fmt.Printf("[ ? ] %s [Y/n] ", question)

	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)
	answer = strings.ToLower(answer)

	return answer == "" || answer == "y" || answer == "yes"
}

// ChooseVariant предлагает выбрать вариант из списка.
func ChooseVariant(reader *bufio.Reader, label string, options []string, defaultIndex int) (int, error) {
	if len(options) == 0 {
		return -1, errors.New("список вариантов пуст")
	}

	if defaultIndex < 0 || defaultIndex >= len(options) {
		defaultIndex = 0
	}

	fmt.Println(label)
	for i, option := range options {
		marker := " "
		if i == defaultIndex {
			marker = "*"
		}
		fmt.Printf("[%d]%s %s\n", i+1, marker, option)
	}

	for {
		fmt.Printf("Выберите вариант [по умолчанию %d]: ", defaultIndex+1)
		raw, _ := reader.ReadString('\n')
		raw = strings.TrimSpace(raw)

		if raw == "" {
			return defaultIndex, nil
		}

		for idx, option := range options {
			if strings.EqualFold(raw, option) {
				return idx, nil
			}
		}

		value, err := strconv.Atoi(raw)
		if err != nil {
			fmt.Println("Введите номер варианта или оставьте строку пустой.")
			continue
		}

		value--
		if value >= 0 && value < len(options) {
			return value, nil
		}

		fmt.Println("Укажите корректный номер из списка.")
	}
}
