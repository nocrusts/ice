package bootstrapper

import (
	"log/slog"
	"os"
	"path/filepath"
)

// WriteAppSettings writes the AppSettings.xml file - required
// to run Roblox - to a binary's deployment directory.
func WriteAppSettings(dir string) error {
	p := filepath.Join(dir, "AppSettings.xml")

	slog.Info("Writing AppSettings.xml", "path", p)

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	appSettings := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\r\n" +
		"<Settings>\r\n" +
		"        <ContentFolder>content</ContentFolder>\r\n" +
		"        <BaseUrl>http://www.roblox.com</BaseUrl>\r\n" +
		"</Settings>\r\n"

	if _, err := f.WriteString(appSettings); err != nil {
		return err
	}

	return nil
}
