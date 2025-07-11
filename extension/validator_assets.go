package extension

import (
	"fmt"
	"os"
	"path/filepath"
)

func validateAssets(ctx *ValidationContext) {
	if !ctx.Extension.GetExtensionConfig().Validation.StoreCompliance {
		return
	}

	for _, resourceDir := range ctx.Extension.GetResourcesDirs() {
		validateAssetByResourceDir(ctx, resourceDir)
	}

	for _, extraBundle := range ctx.Extension.GetExtensionConfig().Build.ExtraBundles {
		bundlePath := ctx.Extension.GetRootDir()

		if extraBundle.Path != "" {
			bundlePath = fmt.Sprintf("%s/%s", bundlePath, extraBundle.Path)
		} else {
			bundlePath = fmt.Sprintf("%s/%s", bundlePath, extraBundle.Name)
		}

		validateAssetByResourceDir(ctx, filepath.Join(bundlePath, "Resources"))
	}
}

func validateAssetByResourceDir(ctx *ValidationContext, resourceDir string) {
	_, foundAdminBuildFiles := os.Stat(filepath.Join(resourceDir, "public", "administration"))
	foundAdminEntrypoint := hasJavascriptEntrypoint(filepath.Join(resourceDir, "app", "administration", "src"))
	foundStorefrontEntrypoint := hasJavascriptEntrypoint(filepath.Join(resourceDir, "app", "storefront", "src"))
	_, foundStorefrontDistFiles := os.Stat(filepath.Join(resourceDir, "app", "storefront", "dist"))

	if foundAdminBuildFiles == nil && !foundAdminEntrypoint {
		ctx.AddError("assets.administration.sources_missing", fmt.Sprintf("Found administration build files in %s but no source files to rebuild the assets.", resourceDir))
	}

	if foundAdminBuildFiles != nil && foundAdminEntrypoint {
		ctx.AddError("assets.administration.build_missing", fmt.Sprintf("Found administration source files in %s but no build files. Please run the build command to generate the assets.", resourceDir))
	}

	if foundStorefrontDistFiles != nil && foundStorefrontEntrypoint {
		ctx.AddError("assets.storefront.sources_missing", fmt.Sprintf("Found storefront build files in %s but no source files to rebuild the assets.", resourceDir))
	}

	if foundStorefrontDistFiles == nil && !foundStorefrontEntrypoint {
		ctx.AddError("assets.storefront.build_missing", fmt.Sprintf("Found storefront source files in %s but no build files. Please run the build command to generate the assets.", resourceDir))
	}
}

func hasJavascriptEntrypoint(jsRoot string) bool {
	entrypointFiles := []string{"main.js", "main.ts"}
	for _, file := range entrypointFiles {
		if _, err := os.Stat(filepath.Join(jsRoot, file)); err == nil {
			return true
		}
	}
	return false
}
