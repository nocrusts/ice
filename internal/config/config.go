package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/vinegarhq/vinegar/roblox"
	"github.com/vinegarhq/vinegar/util"
)

var DefaultLogoPath string

type Environment map[string]string

type UI struct {
	Enabled    bool   `toml:"enabled"`
	Logo       string `toml:"logo"`
	Background uint32 `toml:"background"`
	Foreground uint32 `toml:"foreground"`
	Accent     uint32 `toml:"accent"`
	Gray1      uint32 `toml:"gray1"`
	Gray2      uint32 `toml:"gray2"`
}

type Application struct {
	WineRootOverride       string        `toml:"wineroot"`
	Channel        string        `toml:"channel"`
	Launcher       string        `toml:"launcher"`
	Renderer       string        `toml:"renderer"`
	ForcedVersion  string        `toml:"forced_version"`
	AutoKillPrefix bool          `toml:"auto_kill_prefix"`
	Dxvk           bool          `toml:"dxvk"`
	FFlags         roblox.FFlags `toml:"fflags"`
	Env            Environment   `toml:"env"`
}

type Config struct {
	WineRoot          string      `toml:"wineroot"`
	DxvkVersion       string      `toml:"dxvk_version"`
	MultipleInstances bool        `toml:"multiple_instances"`
	SanitizeEnv       bool        `toml:"sanitize_env"`
	Player            Application `toml:"player"`
	Studio            Application `toml:"studio"`
	Env               Environment `toml:"env"`
	UI                `toml:"ui"`
}

func Load(path string) (Config, error) {
	return _Load(path, "")
}

func LoadForTarget(path string, target string) (Config, error) {
	return _Load(path, target)
}

func _Load(path string, target string) (Config, error) {
	cfg := Default()

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Println("Using default configuration")

		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to decode configuration file: %w", err)
	}

	if err := cfg.Setup(target); err != nil {
		return cfg, fmt.Errorf("failed to setup configuration: %w", err)
	}

	return cfg, nil
}

func Default() Config {
	return Config{
		DxvkVersion: "2.3",

		Env: Environment{
			"WINEARCH":         "win64",
			"WINEDEBUG":        "err-kerberos,err-ntlm",
			"WINEESYNC":        "1",
			"WINEDLLOVERRIDES": "dxdiagn=d;winemenubuilder.exe=d",

			"DXVK_LOG_LEVEL": "warn",
			"DXVK_LOG_PATH":  "none",

			"MESA_GL_VERSION_OVERRIDE":    "4.4",
			"__GL_THREADED_OPTIMIZATIONS": "1",
		},

		Player: Application{
			Dxvk:           true,
			AutoKillPrefix: true,
			FFlags: roblox.FFlags{
				"DFIntTaskSchedulerTargetFps": 640,
			},
		},
		Studio: Application{
			Dxvk: true,
		},

		UI: UI{
			Enabled:    true,
			Logo:       DefaultLogoPath,
			Background: 0x242424,
			Foreground: 0xfafafa,
			Gray1:      0x303030,
			Gray2:      0x777777,
			Accent:     0x8fbc5e,
		},
	}
}

func (e *Environment) Setenv() {
	for name, value := range *e {
		os.Setenv(name, value)
	}
}

func (c *Config) Setup(target string) error {
	if c.SanitizeEnv {
		util.SanitizeEnv()
	}

	selectedRoot := c.WineRoot

	// Check if the target has its own wineroot specified
	if target == "P" && c.Player.WineRootOverride != "" {
		selectedRoot = c.Player.WineRootOverride
	} else if target == "S" && c.Studio.WineRootOverride != "" {
		selectedRoot = c.Studio.WineRootOverride
	}

	if selectedRoot != "" {
		bin := filepath.Join(selectedRoot, "bin")

		if !filepath.IsAbs(c.WineRoot) {
			return errors.New("ensure that the wine root given is an absolute path")
		}

		_, err := os.Stat(filepath.Join(bin, "wine"))
		if err != nil {
			return fmt.Errorf("invalid wine root given: %s", err)
		}

		c.Env["PATH"] = bin + ":" + os.Getenv("PATH")
		os.Unsetenv("WINEDLLPATH")
		log.Printf("Using Wine Root: %s", selectedRoot)
	}

	if !roblox.ValidRenderer(c.Player.Renderer) || !roblox.ValidRenderer(c.Studio.Renderer) {
		return fmt.Errorf("invalid renderer given to either player or studio")
	}

	c.Env.Setenv()

	return nil
}
