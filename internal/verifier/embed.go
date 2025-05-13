package verifier

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/shopware/shopware-cli/internal/system"
)

//go:embed tools/php/composer.json tools/php/composer.lock tools/php/configs tools/js/configs tools/js/packages tools/js/package-*.json
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

	if err := os.MkdirAll(toolsDir, 0o755); err != nil {
		return err
	}

	if err := unpackFile(toolsFS, "."); err != nil {
		os.RemoveAll(toolsDir)
		return err
	}

	composerInstall, err := exec.Command("composer", "install", "--no-dev").CombinedOutput()
	if err != nil {
		os.RemoveAll(toolsDir)
		fmt.Println(string(composerInstall))
		return err
	}

	setToolDirectory(toolsDir)
	return nil
}

func unpackFile(fs embed.FS, filePath string) error {
	f, err := fs.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, file := range f {
		if file.IsDir() {
			if err := os.Mkdir(path.Join(filePath, file.Name()), 0o755); err != nil {
				return err
			}

			if err := unpackFile(fs, path.Join(filePath, file.Name())); err != nil {
				return err
			}
		} else {
			content, err := fs.ReadFile(path.Join(filePath, file.Name()))
			if err != nil {
				return err
			}

			if err := os.WriteFile(path.Join(filePath, file.Name()), content, 0o644); err != nil {
				return err
			}
		}
	}

	return nil
}
