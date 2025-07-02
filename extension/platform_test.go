package extension

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestPlugin(tempDir string) PlatformPlugin {
	return PlatformPlugin{
		path: tempDir,
		config: &Config{
			Store: ConfigStore{
				Availabilities: &[]string{"German"},
			},
		},
		Composer: PlatformComposerJson{
			Name:        "frosh/frosh-tools",
			Description: "Frosh Tools",
			License:     "mit",
			Version:     "1.0.0",
			Require:     map[string]string{"shopware/core": "6.4.0.0"},
			Autoload: struct {
				Psr0 map[string]string `json:"psr-0"`
				Psr4 map[string]string `json:"psr-4"`
			}{Psr0: map[string]string{"FroshTools\\": "src/"}, Psr4: map[string]string{"FroshTools\\": "src/"}},
			Authors: []struct {
				Name     string `json:"name"`
				Homepage string `json:"homepage"`
			}{{Name: "Frosh", Homepage: "https://frosh.io"}},
			Type: "shopware-platform-plugin",
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools\\FroshTools",
				Label: map[string]string{
					"en-GB": "Frosh Tools",
					"de-DE": "Frosh Tools",
				},
				Description: map[string]string{
					"en-GB": "Frosh Tools",
					"de-DE": "Frosh Tools",
				},
				ManufacturerLink: map[string]string{
					"en-GB": "Frosh Tools",
					"de-DE": "Frosh Tools",
				},
				SupportLink: map[string]string{
					"en-GB": "Frosh Tools",
					"de-DE": "Frosh Tools",
				},
			},
		},
	}
}

func TestPluginIconNotExists(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)

	ctx := newValidationContext(&plugin)

	plugin.Validate(getTestContext(), ctx)

	assert.Equal(t, 1, len(ctx.errors))
	assert.Equal(t, "The extension icon Resources/config/plugin.png does not exist", ctx.errors[0].Message)
}

func TestPluginIconExists(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "config"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "config", "plugin.png"), []byte("test"), os.ModePerm))

	ctx := newValidationContext(&plugin)

	plugin.Validate(getTestContext(), ctx)

	assert.Equal(t, 0, len(ctx.errors))
}

func TestPluginIconDifferntPathExists(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)
	plugin.Composer.Extra.PluginIcon = "plugin.png"

	assert.NoError(t, os.WriteFile(path.Join(dir, "plugin.png"), []byte("test"), os.ModePerm))

	ctx := newValidationContext(&plugin)

	plugin.Validate(getTestContext(), ctx)

	assert.Equal(t, 0, len(ctx.errors))
}

func TestPluginIconIsTooBig(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "config"), os.ModePerm))
	// Create a file larger than 10KB
	bigFile := make([]byte, 11*1024)
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "config", "plugin.png"), bigFile, os.ModePerm))

	ctx := newValidationContext(&plugin)

	plugin.Validate(getTestContext(), ctx)

	assert.Len(t, ctx.errors, 1)
	assert.Equal(t, "The extension icon Resources/config/plugin.png is bigger than 10kb", ctx.errors[0].Message)
}

func TestPluginGermanDescriptionMissing(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)
	plugin.Composer.Extra.Description = map[string]string{
		"en-GB": "Frosh Tools",
	}

	ctx := newValidationContext(&plugin)
	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "config"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "config", "plugin.png"), []byte("test"), os.ModePerm))

	plugin.Validate(getTestContext(), ctx)

	assert.Len(t, ctx.errors, 1)
	assert.Equal(t, "extra.description for language de-DE is required", ctx.errors[0].Message)
}

func TestPluginGermanDescriptionMissingOnlyEnglishMarket(t *testing.T) {
	dir := t.TempDir()

	plugin := getTestPlugin(dir)
	plugin.Composer.Extra.Description = map[string]string{
		"en-GB": "Frosh Tools",
	}
	plugin.config.Store.Availabilities = &[]string{"International"}
	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "config"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "config", "plugin.png"), []byte("test"), os.ModePerm))

	ctx := newValidationContext(&plugin)

	plugin.Validate(getTestContext(), ctx)

	assert.Len(t, ctx.errors, 0)
}
