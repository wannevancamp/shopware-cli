package packagist

import (
	"encoding/json"
	"fmt"
	"os"
)

type ComposerLock struct {
	Packages []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		PackageType string `json:"type"`
	} `json:"packages"`
}

func ReadComposerLock(pathToFile string) (*ComposerLock, error) {
	content, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var lock ComposerLock
	if err := json.Unmarshal(content, &lock); err != nil {
		return nil, fmt.Errorf("could not parse composer.lock: %w", err)
	}

	return &lock, nil
}
