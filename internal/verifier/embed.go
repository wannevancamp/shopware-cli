package verifier

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/shopware/shopware-cli/internal/system"
)

//go:embed php/composer.json php/composer.lock php/configs js/configs js/packages js/package*.json
var toolsFS embed.FS

func SetupTools(currentVersion string) error {
	customToolDir := os.Getenv("SHOPWARE_CLI_TOOLS_DIR")
	if customToolDir != "" {
		setToolDirectory(customToolDir)
		return nil
	}

	cacheDir := system.GetShopwareCliCacheDir()
	toolsDir := path.Join(cacheDir, "tools", currentVersion)

	if _, err := os.Stat(toolsDir); err == nil {
		setToolDirectory(toolsDir)
		return nil
	}

	if err := os.MkdirAll(toolsDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", toolsDir, err)
	}

	if err := unpackFile(toolsFS, ".", toolsDir); err != nil {
		os.RemoveAll(toolsDir)
		return fmt.Errorf("failed to unpack file: %w", err)
	}

	composerInstall := exec.Command("composer", "install", "--no-dev")
	composerInstall.Dir = path.Join(toolsDir, "php")

	if output, err := composerInstall.CombinedOutput(); err != nil {
		fmt.Println(string(output))
		if err := os.RemoveAll(toolsDir); err != nil {
			return err
		}
		return fmt.Errorf("failed to install composer dependencies: %w", err)
	}

	jsInstall := exec.Command("npm", "install")
	jsInstall.Dir = path.Join(toolsDir, "js")

	if output, err := jsInstall.CombinedOutput(); err != nil {
		fmt.Println(string(output))
		if err := os.RemoveAll(toolsDir); err != nil {
			return err
		}
		return fmt.Errorf("failed to install npm dependencies: %w", err)
	}

	setToolDirectory(toolsDir)
	return nil
}

func unpackFile(fs embed.FS, filePath, unpackDir string) error {
	f, err := fs.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, file := range f {
		if file.IsDir() {
			if err := os.MkdirAll(path.Join(unpackDir, file.Name()), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path.Join(unpackDir, file.Name()), err)
			}

			if err := unpackFile(fs, path.Join(filePath, file.Name()), path.Join(unpackDir, file.Name())); err != nil {
				return fmt.Errorf("failed to unpack file %s: %w", path.Join(filePath, file.Name()), err)
			}
		} else {
			content, err := fs.ReadFile(path.Join(filePath, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path.Join(filePath, file.Name()), err)
			}

			if err := os.WriteFile(path.Join(unpackDir, file.Name()), content, os.ModePerm); err != nil {
				return fmt.Errorf("failed to write file %s: %w", path.Join(unpackDir, file.Name()), err)
			}
		}
	}

	return nil
}
