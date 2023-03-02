package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	ApplyRCO    bool
	AutoRFPSU   bool
	AutoKillPfx bool
	Dxvk        bool
	GameMode    bool
	Log         bool
	Version     string
	Renderer    string
	FFlags      map[string]interface{}
	Env         map[string]string
}

var (
	ConfigFilePath = filepath.Join(Dirs.Config, "config.toml")
	Config         = loadConfig()
)

func defConfig() Configuration {
	return Configuration{
		ApplyRCO:    true,
		AutoRFPSU:   false,
		AutoKillPfx: true,
		Dxvk:        false,
		GameMode:    false,
		Log:         true,
		Version:     "win10",
		Renderer:    "D3D11",
		FFlags:      make(map[string]interface{}),
		Env: map[string]string{
			"WINEPREFIX": Dirs.Pfx,
			"WINEARCH":   "win64",
			// "WINEDEBUG":        "fixme-all,-wininet,-ntlm,-winediag,-kerberos",
			"WINEDEBUG":        "-all",
			"WINEDLLOVERRIDES": "dxdiagn=d;winemenubuilder.exe=d;",

			"DXVK_LOG_LEVEL":        "warn",
			"DXVK_LOG_PATH":         "none",
			"DXVK_STATE_CACHE_PATH": filepath.Join(Dirs.Cache, "dxvk"),

			"MESA_GL_VERSION_OVERRIDE":    "4.4",
			"__GL_THREADED_OPTIMIZATIONS": "1",
			"DRI_PRIME":                   "1",
			"__NV_PRIME_RENDER_OFFLOAD":   "1",
			"__VK_LAYER_NV_optimus":       "NVIDIA_only",
			"__GLX_VENDOR_LIBRARY_NAME":   "nvidia",
		},
	}
}

func writeConfigTemplate() {
	CheckDirs(DirMode, Dirs.Config)

	log.Println("Creating configuration template")

	file, err := os.Create(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	template := `
# See how to configure Vinegar on the documentation website:
# https://vinegarhq.github.io/Configuration
`
	if _, err = file.WriteString(template[1:]); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() Configuration {
	config := defConfig()

	if _, err := os.Stat(ConfigFilePath); errors.Is(err, os.ErrNotExist) {
		writeConfigTemplate()
	} else if err != nil {
		log.Fatal(err)
	}

	if _, err := toml.DecodeFile(ConfigFilePath, &config); err != nil {
		log.Fatal("could not parse configuration file:", err)
	}

	for name, value := range config.Env {
		os.Setenv(name, value)
	}

	return config
}

func GetEditor() (string, error) {
	editor, ok := os.LookupEnv("EDITOR")

	if ok {
		if _, err := exec.LookPath(editor); err != nil {
			return "", fmt.Errorf("invalid $EDITOR: %w", err)
		}
	} else {
		return "", errors.New("no $EDITOR variable set")
	}

	return editor, nil
}

func EditConfig() {
	var testConfig Configuration

	editor, err := GetEditor()
	if err != nil {
		log.Fatal("unable to find editor:", err)
	}

	tempConfigFile, err := os.CreateTemp(Dirs.Config, "testconfig.*.toml")
	if err != nil {
		log.Fatal(err)
	}

	tempConfigFilePath, err := filepath.Abs(tempConfigFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = tempConfigFile.Write(configFile); err != nil {
		log.Fatal(err)
	}

	tempConfigFile.Close()

	for {
		if err := Exec(editor, false, tempConfigFilePath); err != nil {
			log.Fatal(err)
		}

		if _, err := toml.DecodeFile(tempConfigFilePath, &testConfig); err != nil {
			log.Println(err)
			log.Println("Press enter to re-edit configuration file")
			fmt.Scanln()

			continue
		}

		if err := os.Rename(tempConfigFilePath, ConfigFilePath); err != nil {
			log.Fatal(err)
		}

		break
	}
}
