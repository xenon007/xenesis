package app

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"xenon007/xenesis/internal/filesystem"
	"xenon007/xenesis/internal/prompt"
	"xenon007/xenesis/internal/templates"
)

const (
	version          = "1.1.0"
	banner           = "xenesis v" + version + "\n\thttps://github.com/xenon007/xenesis\n\n"
	gitignoreContent = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work
`
)

var (
	// ErrGoModExists сигнализирует о наличии go.mod в каталоге назначения.
	ErrGoModExists = errors.New("go.mod уже существует")
)

// App инкапсулирует бизнес-логику CLI.
type App struct {
	logger *log.Logger
	reader *bufio.Reader
}

// New подготавливает приложение к запуску.
func New(logger *log.Logger) *App {
	return &App{
		logger: logger,
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run выполняет сценарий генерации проекта.
func (a *App) Run() error {
	fmt.Print(banner)

	projectName, err := prompt.ProjectName(a.reader)
	if err != nil {
		return err
	}
	a.logger.Printf("имя проекта: %s", projectName)

	templatesList := templates.All()
	options := make([]string, len(templatesList))
	for i, tpl := range templatesList {
		options[i] = fmt.Sprintf("%s — %s", tpl.Title, tpl.Description)
	}

	choice, err := prompt.ChooseVariant(a.reader, "Выберите шаблон проекта:", options, 0)
	if err != nil {
		return err
	}
	tpl := templatesList[choice]
	a.logger.Printf("используется шаблон: %s", tpl.Title)

	description := prompt.Description(a.reader)
	githubHandle := prompt.GithubHandle(a.reader)
	modulePath := fmt.Sprintf("github.com/%s/%s", githubHandle, projectName)
	a.logger.Printf("модуль: %s", modulePath)

	rootDir := filepath.Join(".", projectName)
	if err := filesystem.EnsureDir(rootDir); err != nil {
		return err
	}
	a.logger.Printf("создан корневой каталог: %s", rootDir)

	if err := a.prepareFolders(rootDir, tpl); err != nil {
		return err
	}

	mainFilePath := filepath.Join(rootDir, tpl.MainFileName(projectName))
	if err := filesystem.WriteFileWithDirs(mainFilePath, tpl.MainFileContent(projectName, modulePath)); err != nil {
		return err
	}
	a.logger.Printf("сгенерирован файл точки входа: %s", mainFilePath)

	extraFiles := tpl.ExtraFiles(projectName, modulePath)
	for _, file := range extraFiles {
		targetPath := filepath.Join(rootDir, file.Path)
		if err := filesystem.WriteFileWithDirs(targetPath, file.Content); err != nil {
			return err
		}
		a.logger.Printf("добавлен файл: %s", targetPath)
	}

	readmePath := filepath.Join(rootDir, "README.md")
	if err := filesystem.WriteFileWithDirs(readmePath, tpl.ReadmeContent(projectName, description, githubHandle)); err != nil {
		return err
	}
	a.logger.Printf("создан README: %s", readmePath)

	gitignorePath := filepath.Join(rootDir, ".gitignore")
	if err := filesystem.WriteFileWithDirs(gitignorePath, gitignoreContent); err != nil {
		return err
	}
	a.logger.Printf("создан .gitignore: %s", gitignorePath)

	if err := a.initGoModule(rootDir, modulePath); err != nil {
		return err
	}
	a.logger.Printf("инициализирован go.mod для %s", modulePath)

	return nil
}

func (a *App) prepareFolders(rootDir string, tpl templates.Template) error {
	for _, folder := range tpl.MandatoryFolders {
		if err := a.ensureFolder(rootDir, folder); err != nil {
			return err
		}
	}

	for folder, question := range tpl.OptionalFolders {
		if prompt.Confirm(a.reader, question) {
			if err := a.ensureFolder(rootDir, folder); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *App) ensureFolder(rootDir, folder string) error {
	target := filepath.Join(rootDir, folder)
	if err := filesystem.EnsureDir(target); err != nil {
		return err
	}
	a.logger.Printf("создан каталог: %s", target)

	if shouldCreateGitkeep(folder) {
		if err := filesystem.CreateGitKeep(target); err != nil {
			return err
		}
		a.logger.Printf("добавлен .gitkeep в: %s", target)
	}

	return nil
}

func shouldCreateGitkeep(folder string) bool {
	cleaned := strings.Trim(folder, string(os.PathSeparator))
	return cleaned != "cmd"
}

func (a *App) initGoModule(rootDir, modulePath string) error {
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = rootDir

	if err := cmd.Run(); err != nil {
		return ErrGoModExists
	}

	return nil
}
