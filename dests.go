package main

/* github.com/pizzaboxer/bloxstrap/blob/main/Bloxstrap/Bootstrapper.cs */
func PlayerPackages() map[string]string {
	return map[string]string{
		"RobloxApp.zip":                 "",
		"WebView2.zip":                  "",
		"shaders.zip":                   "shaders",
		"ssl.zip":                       "ssl",
		"content-avatar.zip":            "content/avatar",
		"content-configs.zip":           "content/configs",
		"content-fonts.zip":             "content/fonts",
		"content-sky.zip":               "content/sky",
		"content-sounds.zip":            "content/sounds",
		"content-textures2.zip":         "content/textures",
		"content-models.zip":            "content/models",
		"content-textures3.zip":         "PlatformContent/pc/textures",
		"content-terrain.zip":           "PlatformContent/pc/terrain",
		"content-platform-fonts.zip":    "PlatformContent/pc/fonts",
		"extracontent-luapackages.zip":  "ExtraContent/LuaPackages",
		"extracontent-translations.zip": "ExtraContent/translations",
		"extracontent-models.zip":       "ExtraContent/models",
		"extracontent-textures.zip":     "ExtraContent/textures",
		"extracontent-places.zip":       "ExtraContent/places",
	}
}

/* github.com/MaximumADHD/Roblox-Studio-Mod-Manager/blob/main/Config/KnownRoots.json */
func StudioPackages() map[string]string {
	return map[string]string{
		"ApplicationConfig.zip":           "ApplicationConfig",
		"BuiltInPlugins.zip":              "BuiltInPlugins",
		"BuiltInStandalonePlugins.zip":    "BuiltInStandalonePlugins",
		"Plugins.zip":                     "Plugins",
		"Qml.zip":                         "Qml",
		"StudioFonts.zip":                 "StudioFonts",
		"WebView2.zip":                    "",
		"RobloxStudio.zip":                "",
		"Libraries.zip":                   "",
		"LibrariesQt5.zip":                "",
		"content-avatar.zip":              "content/avatar",
		"content-configs.zip":             "content/configs",
		"content-fonts.zip":               "content/fonts",
		"content-models.zip":              "content/models",
		"content-qt_translations.zip":     "content/qt_translations",
		"content-sky.zip":                 "content/sky",
		"content-sounds.zip":              "content/sounds",
		"shaders.zip":                     "shaders",
		"ssl.zip":                         "ssl",
		"content-textures2.zip":           "content/textures",
		"content-textures3.zip":           "PlatformContent/pc/textures",
		"content-studio_svg_textures.zip": "content/studio_svg_textures",
		"content-terrain.zip":             "PlatformContent/pc/terrain",
		"content-platform-fonts.zip":      "PlatformContent/pc/fonts",
		"content-api-docs.zip":            "content/api_docs",
		"extracontent-scripts.zip":        "ExtraContent/scripts",
		"extracontent-luapackages.zip":    "ExtraContent/LuaPackages",
		"extracontent-translations.zip":   "ExtraContent/translations",
		"extracontent-models.zip":         "ExtraContent/models",
		"extracontent-textures.zip":       "ExtraContent/textures",
		"redist.zip":                      "",
	}
}
