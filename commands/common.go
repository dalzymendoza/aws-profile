package commands

import (
	"os/user"
	"strings"
	"path/filepath"
	)

func ExpandHomeDirectory(filePath string) string {
	usr, _ := user.Current()
	homeDirectory := usr.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		return filepath.Join(homeDirectory, filePath[2:])
	}

	return filePath
}
