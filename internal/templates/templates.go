package templates

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	// MDConsoleInit открывает блок кода в README.md.
	MDConsoleInit = "```console"
)

// File описывает файл, который требуется создать в проекте.
type File struct {
	Path    string
	Content string
}

// Template описывает структуру будущего проекта.
type Template struct {
	Key              string
	Title            string
	Description      string
	MandatoryFolders []string
	OptionalFolders  map[string]string
	MainFileName     func(projectName string) string
	MainFileContent  func(projectName, modulePath string) string
	ExtraFiles       func(projectName, modulePath string) []File
	ReadmeContent    func(projectName, description, githubHandle string) string
}

// All возвращает список поддерживаемых шаблонов.
func All() []Template {
	return []Template{
		basicTemplate(),
		x07ComponentTemplate(),
		threeXfaTemplate(),
	}
}

// Default возвращает шаблон по умолчанию.
func Default() Template {
	return basicTemplate()
}

// SanitizePackageName преобразует имя проекта в допустимое имя пакета.
func SanitizePackageName(projectName string) string {
	sanitized := strings.ToLower(projectName)
	replacer := strings.NewReplacer("-", "_", " ", "")
	sanitized = replacer.Replace(sanitized)
	return sanitized
}

func basicTemplate() Template {
	return Template{
		Key:         "basic",
		Title:       "Базовый Go-проект",
		Description: "Минимальный каркас приложения с набором полезных каталогов",
		MandatoryFolders: []string{
			"cmd",
			"pkg",
			"docs",
			"internal",
			"examples",
		},
		OptionalFolders: map[string]string{
			"api":     "Нужны ли API?",
			"server":  "Нужен ли встроенный сервер?",
			"db":      "Понадобится ли подключение к БД?",
			"scripts": "Создать каталог для скриптов?",
			"test":    "Добавить тестовые данные?",
			"init":    "Нужны конфигурации процесс-менеджеров?",
			"assets":  "Добавить каталог для ассетов?",
		},
		MainFileName: func(projectName string) string {
			return filepath.Join("cmd", fmt.Sprintf("%s.go", projectName))
		},
		MainFileContent: func(projectName, modulePath string) string {
			return "package main\n\n" +
				"import (\n" +
				"\t\"fmt\"\n" +
				")\n\n" +
				"func main() {\n" +
				"\tfmt.Println(\"Hello, World!\")\n" +
				"}\n"
		},
		ExtraFiles: func(projectName, modulePath string) []File {
			return nil
		},
		ReadmeContent: func(projectName, description, githubHandle string) string {
			builder := strings.Builder{}
			builder.WriteString("# ")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString(description)
			builder.WriteString("\n\nInstallation 📡\n")
			builder.WriteString("-------\n")
			builder.WriteString("**Go 1.17+**\n")
			builder.WriteString(MDConsoleInit + "\n")
			builder.WriteString("go install -v github.com/")
			builder.WriteString(githubHandle)
			builder.WriteString("/")
			builder.WriteString(projectName)
			builder.WriteString("/cmd/")
			builder.WriteString(projectName)
			builder.WriteString("@latest\n")
			builder.WriteString("```\n")
			builder.WriteString("**otherwise**\n")
			builder.WriteString(MDConsoleInit + "\n")
			builder.WriteString("go get -v github.com/")
			builder.WriteString(githubHandle)
			builder.WriteString("/")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString("```\n\n")
			builder.WriteString("Usage 💻\n")
			builder.WriteString("-------\n")
			builder.WriteString(MDConsoleInit + "\n")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString("```\n\n")
			builder.WriteString("Created with [xenesis](https://github.com/xenon007/xenesis)❤️")
			return builder.String()
		},
	}
}

func x07ComponentTemplate() Template {
	return Template{
		Key:         "x07-component",
		Title:       "Компонент X07",
		Description: "Структура для создания компонента экосистемы X07",
		MandatoryFolders: []string{
			"cmd",
			filepath.Join("internal", "components"),
			filepath.Join("internal", "config"),
			"pkg",
			"docs",
		},
		OptionalFolders: map[string]string{
			"deploy":  "Добавить инфраструктурные манифесты для X07?",
			"metrics": "Создать каталог для метрик?",
			"assets":  "Добавить статические файлы компонента?",
		},
		MainFileName: func(projectName string) string {
			return filepath.Join("cmd", fmt.Sprintf("%s.go", projectName))
		},
		MainFileContent: func(projectName, modulePath string) string {
			packageName := SanitizePackageName(projectName)
			return fmt.Sprintf(`package main

import (
        "context"
        "log"

        "%s/internal/components/%s"
)

func main() {
        ctx := context.Background()

        component := %s.New()
        if err := component.Bootstrap(ctx); err != nil {
                log.Fatalf("component bootstrap failed: %v", err)
        }

        log.Println("component is up and running")
}
`, modulePath, packageName, packageName)
		},
		ExtraFiles: func(projectName, modulePath string) []File {
			packageName := SanitizePackageName(projectName)
			componentPath := filepath.Join("internal", "components", packageName, "component.go")
			configPath := filepath.Join("internal", "config", "config.go")

			return []File{
				{
					Path: componentPath,
					Content: fmt.Sprintf(`package %s

import (
        "context"
        "log"
)

type Component struct {
        // здесь можно сохранить зависимости компонента
}

func New() *Component {
        return &Component{}
}

func (c *Component) Bootstrap(ctx context.Context) error {
        log.Println("bootstrap X07 component")
        // TODO: добавить инициализацию бизнес-логики
        return nil
}
`, packageName),
				},
				{
					Path: configPath,
					Content: `package config

import "time"

type Runtime struct {
        GracefulShutdown time.Duration
}

func DefaultRuntime() Runtime {
        return Runtime{GracefulShutdown: 30 * time.Second}
}
`,
				},
			}
		},
		ReadmeContent: func(projectName, description, githubHandle string) string {
			builder := strings.Builder{}
			builder.WriteString("# Компонент X07: ")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString(description)
			builder.WriteString("\n\n")
			builder.WriteString("## Быстрый старт\n")
			builder.WriteString(MDConsoleInit + "\n")
			builder.WriteString("go run ./cmd/")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString("```\n\n")
			builder.WriteString("## Структура\n")
			builder.WriteString("- `internal/components/` — бизнес-логика компонента\n")
			builder.WriteString("- `internal/config/` — конфигурация и настройки\n")
			builder.WriteString("- `pkg/` — разделяемые пакеты\n")
			builder.WriteString("- `docs/` — документация и ADR\n")
			builder.WriteString("\nСоздано при помощи [xenesis](https://github.com/xenon007/xenesis).")
			return builder.String()
		},
	}
}

func threeXfaTemplate() Template {
	return Template{
		Key:         "3xfa-plugin",
		Title:       "Плагин 3XFA",
		Description: "Базовый каркас плагина 3XFA с точками расширения",
		MandatoryFolders: []string{
			"cmd",
			filepath.Join("plugins", "3xfa"),
			"internal",
			filepath.Join("internal", "adapters"),
		},
		OptionalFolders: map[string]string{
			filepath.Join("internal", "mocks"):     "Добавить заглушки для тестов?",
			filepath.Join("internal", "transport"): "Создать транспортный слой для плагина?",
		},
		MainFileName: func(projectName string) string {
			return filepath.Join("cmd", fmt.Sprintf("%s.go", projectName))
		},
		MainFileContent: func(projectName, modulePath string) string {
			return fmt.Sprintf(`package main

import (
        "log"

        _3xfa "%s/plugins/3xfa"
)

func main() {
        plugin := _3xfa.New()

        if err := plugin.Register(); err != nil {
                log.Fatalf("3XFA plugin registration failed: %v", err)
        }

        log.Println("3XFA plugin registered successfully")
}
`, modulePath)
		},
		ExtraFiles: func(projectName, modulePath string) []File {
			pluginPath := filepath.Join("plugins", "3xfa", "plugin.go")
			adapterPath := filepath.Join("internal", "adapters", "registry.go")
			return []File{
				{
					Path: pluginPath,
					Content: `package _3xfa

import "log"

type Plugin struct {
        hooks []Hook
}

type Hook interface {
        Name() string
        Execute() error
}

func New() *Plugin {
        return &Plugin{}
}

func (p *Plugin) Register() error {
        log.Println("prepare 3XFA plugin hooks")
        for _, hook := range p.hooks {
                if err := hook.Execute(); err != nil {
                        return err
                }
        }
        return nil
}

func (p *Plugin) AttachHook(h Hook) {
        p.hooks = append(p.hooks, h)
}
`,
				},
				{
					Path: adapterPath,
					Content: `package adapters

import "log"

type Registry struct {
        registered []string
}

func NewRegistry() *Registry {
        return &Registry{registered: make([]string, 0)}
}

func (r *Registry) Add(name string) {
        log.Printf("register adapter: %s", name)
        r.registered = append(r.registered, name)
}

func (r *Registry) List() []string {
        return append([]string(nil), r.registered...)
}
`,
				},
			}
		},
		ReadmeContent: func(projectName, description, githubHandle string) string {
			builder := strings.Builder{}
			builder.WriteString("# 3XFA Plugin: ")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString(description)
			builder.WriteString("\n\n")
			builder.WriteString("## Сборка\n")
			builder.WriteString(MDConsoleInit + "\n")
			builder.WriteString("go build ./cmd/")
			builder.WriteString(projectName)
			builder.WriteString("\n")
			builder.WriteString("```\n\n")
			builder.WriteString("## Подключение\n")
			builder.WriteString("Используйте пакет `plugins/3xfa` и регистрируйте свои хуки через метод `AttachHook`.\n")
			builder.WriteString("\nСоздано при помощи [xenesis](https://github.com/xenon007/xenesis).")
			return builder.String()
		},
	}
}
