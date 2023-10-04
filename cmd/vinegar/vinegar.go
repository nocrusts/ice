package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/vinegarhq/vinegar/internal/config"
	"github.com/vinegarhq/vinegar/internal/config/editor"
	"github.com/vinegarhq/vinegar/internal/config/state"
	"github.com/vinegarhq/vinegar/internal/dirs"
	"github.com/vinegarhq/vinegar/internal/logs"
	"github.com/vinegarhq/vinegar/roblox"
	"github.com/vinegarhq/vinegar/wine"
)

var BinPrefix string

func usage() {
	fmt.Fprintln(os.Stderr, "usage: vinegar [-config filepath] player|studio [args...]")
	fmt.Fprintln(os.Stderr, "usage: vinegar [-config filepath] exec prog [args...]")
	fmt.Fprintln(os.Stderr, "       vinegar [-config filepath] edit|kill|uninstall|delete|install-webview2")
	os.Exit(1)
}

func main() {
	configPath := flag.String("config", filepath.Join(dirs.Config, "config.toml"), "config.toml file which should be used")
	flag.Parse()

	cmd := flag.Arg(0)
	args := flag.Args()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	switch cmd {
	// These commands don't require a configuration
	case "delete", "edit", "uninstall":
		switch cmd {
		case "delete":
			Delete()
		case "edit":
			editor.EditConfig(*configPath)
		case "uninstall":
			Uninstall()
		}
	// These commands (except player & studio) don't require a configuration,
	// but they require a wineprefix, hence wineroot of configuration is required.
	case "player", "studio", "exec", "kill", "install-webview2":
		pfxKilled := false

		var cfg config.Config
		var err error

		switch cmd {

		case "player":
			cfg, err = config.LoadForTarget(*configPath,"P")
		case "studio":
			cfg, err = config.LoadForTarget(*configPath,"S")
		default:
			cfg, err = config.Load(*configPath)
		}

		if err != nil {
			log.Fatal(err)
		}

		pfx := wine.New(dirs.Prefix)

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			<-c
			pfxKilled = true
			pfx.Kill()

			if pfxKilled {
				os.Exit(0)
			}
		}()

		if err := pfx.Setup(); err != nil {
			log.Fatal(err)
		}

		switch cmd {
		case "exec":
			if len(args) < 2 {
				usage()
			}

			if err := pfx.Wine(args[1], args[2:]...).Run(); err != nil {
				log.Fatal(err)
			}
		case "kill":
			pfx.Kill()
		case "install-webview2":
			if err := InstallWebview2(&pfx); err != nil {
				log.Fatal(err)
			}

		case "player", "studio":
			var b Binary

			logFile := logs.File(cmd)
			logOutput := io.MultiWriter(logFile, os.Stderr)

			pfx.Output = logOutput
			log.SetOutput(logOutput)

			defer logFile.Close()

			switch cmd {
			case "player":
				b = NewBinary(roblox.Player, &cfg, &pfx)
			case "studio":
				b = NewBinary(roblox.Studio, &cfg, &pfx)
			}

			b.Run(args[1:]...)
		}
	default:
		usage()
	}
}

func Uninstall() {
	vers, err := state.Versions()
	if err != nil {
		log.Fatal(err)
	}

	for _, ver := range vers {
		log.Println("Removing version directory", ver)

		err = os.RemoveAll(filepath.Join(dirs.Versions, ver))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = state.ClearApplications()
	if err != nil {
		log.Fatal(err)
	}
}

func Delete() {
	log.Println("Deleting Wineprefix")
	if err := os.RemoveAll(dirs.Prefix); err != nil {
		log.Fatal(err)
	}
}
