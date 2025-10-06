package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DirPermission определяет права на создаваемые каталоги.
	DirPermission os.FileMode = 0o775
	// FilePermission определяет права на создаваемые файлы.
	FilePermission os.FileMode = 0o644
)

// EnsureDir гарантирует существование каталога по переданному пути.
func EnsureDir(path string) error {
	if path == "" || path == "." {
		return nil
	}

	if err := os.MkdirAll(path, DirPermission); err != nil {
		return fmt.Errorf("не удалось создать каталог %s: %w", path, err)
	}

	return nil
}

// WriteFileWithDirs записывает файл, автоматически создавая недостающие каталоги.
func WriteFileWithDirs(path string, content string) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(content), FilePermission); err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", path, err)
	}

	return nil
}

// TouchFileWithDirs создаёт пустой файл с автоматическим созданием каталогов.
func TouchFileWithDirs(path string) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte(""), FilePermission); err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}

	return nil
}

// CreateGitKeep создаёт .gitkeep в указанном каталоге.
func CreateGitKeep(dir string) error {
	if dir == "" {
		return errors.New("не указан каталог для .gitkeep")
	}

	path := filepath.Join(dir, ".gitkeep")
	if err := EnsureDir(dir); err != nil {
		return err
	}

	if err := os.WriteFile(path, []byte("keep this file plz\n"), FilePermission); err != nil {
		return fmt.Errorf("не удалось создать %s: %w", path, err)
	}

	return nil
}
