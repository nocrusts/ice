package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

const (
	DXVKVER        = "2.1"
	DXVKTAR        = "dxvk-" + DXVKVER + ".tar.gz"
	DXVKURL        = "https://github.com/doitsujin/dxvk/releases/download/v" + DXVKVER + "/" + DXVKTAR
	FPSUNLOCKERURL = "https://github.com/axstin/rbxfpsunlocker/releases/download/v4.4.4/rbxfpsunlocker-x64.zip"
)

var (
	DxvkInstallMarker = filepath.Join(Dirs.Pfx, ".vinegar-dxvk")
	InFlatpak         = InFlatpakCheck()
)

func InFlatpakCheck() bool {
	if _, err := os.Stat("/.flatpak-info"); err == nil {
		return true
	}

	return false
}

func DxvkMarkerExist() bool {
	_, err := os.Open(DxvkInstallMarker)

	return err == nil
}

func DxvkStrap() {
	if Config.Dxvk {
		DxvkInstall(false)

		Config.Renderer = "D3D11"
		Config.Env["WINEDLLOVERRIDES"] += "d3d10core=n;d3d11=n;d3d9=n;dxgi=n"
		os.Setenv("WINEDLLOVERRIDES", Config.Env["WINEDLLOVERRIDES"])

		return
	}

	DxvkUninstall(false)
}

func DxvkExtract(tarball string) error {
	var winDir string

	log.Println("Extracting DXVK")

	dxvkTarball, err := os.Open(tarball)
	if err != nil {
		return err
	}

	dxvkGzip, err := gzip.NewReader(dxvkTarball)
	if err != nil {
		return err
	}

	dxvkTar := tar.NewReader(dxvkGzip)

	for header, err := dxvkTar.Next(); err == nil; header, err = dxvkTar.Next() {
		if header.Typeflag != tar.TypeReg {
			continue
		}

		switch path.Base(path.Dir(header.Name)) {
		case "x64":
			winDir = "system32"
		case "x32":
			winDir = "syswow64"
		default:
			continue
		}

		winDir := filepath.Join(Dirs.Pfx, "drive_c", "windows", winDir)

		CheckDirs(DirMode, winDir)

		writer, err := os.Create(filepath.Join(winDir, path.Base(header.Name)))
		if err != nil {
			return err
		}

		log.Println("Extracting DLL:", writer.Name())

		if _, err = io.Copy(writer, dxvkTar); err != nil {
			return err
		}
	}

	return nil
}

func DxvkInstall(force bool) {
	if !force {
		if DxvkMarkerExist() {
			return
		}
	}

	log.Println("Installing DXVK")

	if InFlatpak {
		log.Println("dxvk must be managed by flatpak, returning")

		return
	}

	dxvkTarballPath := filepath.Join(Dirs.Cache, DXVKTAR)

	if err := Download(DXVKURL, dxvkTarballPath); err != nil {
		log.Fatal(err)
	}

	if err := DxvkExtract(dxvkTarballPath); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(dxvkTarballPath); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Create(DxvkInstallMarker); err != nil {
		log.Fatal(err)
	}
}

func DxvkUninstall(force bool) {
	if !force {
		if !DxvkMarkerExist() {
			return
		}
	}

	log.Println("Uninstalling DXVK")

	for _, dir := range []string{"syswow64", "system32"} {
		for _, dll := range []string{"d3d9", "d3d10core", "d3d11", "dxgi"} {
			dllFile := filepath.Join(Dirs.Pfx, "drive_c", "windows", dir, dll+".dll")
			log.Println("Removing DLL:", dllFile)

			if err := os.RemoveAll(dllFile); err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Println("Updating wineprefix")

	if err := Exec("wineboot", false, "-u"); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(DxvkInstallMarker); err != nil {
		log.Fatal(err)
	}
}

func RbxFpsUnlockerSettings(file string) {
	log.Println("Writing custom rbxfpsunlocker settings")

	settings := []string{
		"UnlockClient=true",
		"UnlockStudio=true",
		"FPSCapValues=[30.000000, 60.000000, 75.000000, 120.000000, 144.000000, 165.000000, 240.000000, 360.000000]",
		"FPSCapSelection=0",
		"FPSCap=0.000000",
		"CheckForUpdates=false",
		"NonBlockingErrors=true",
		"SilentErrors=true",
		"QuickStart=true",
	}

	settingsFile, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, setting := range settings {
		if _, err := fmt.Fprintln(settingsFile, setting+"\r"); err != nil {
			log.Fatal(err)
		}
	}
}

func RbxFpsUnlocker() {
	fpsUnlockerPath := filepath.Join(Dirs.Data, "rbxfpsunlocker.exe")

	if _, err := os.Stat(fpsUnlockerPath); errors.Is(err, os.ErrNotExist) {
		fpsUnlockerZip := filepath.Join(Dirs.Cache, "rbxfpsunlocker.zip")

		log.Println("Installing rbxfpsunlocker")

		if err := Download(FPSUNLOCKERURL, fpsUnlockerZip); err != nil {
			log.Fatal(err)
		}

		if err := UnzipFile(fpsUnlockerZip, "rbxfpsunlocker.exe", fpsUnlockerPath); err != nil {
			log.Fatal(err)
		}
	}

	settingsFile := filepath.Join(Dirs.Cache, "settings")

	if _, err := os.Stat(settingsFile); errors.Is(err, os.ErrNotExist) {
		RbxFpsUnlockerSettings(settingsFile)
	}

	log.Println("Launching FPS Unlocker")

	if err := Exec("wine", true, fpsUnlockerPath); err != nil {
		log.Fatal("rbxfpsunlocker err:", err)
	}
}
