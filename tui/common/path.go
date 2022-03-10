package common

import (
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

func GetPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, _ := user.Current()
		dir := usr.HomeDir
		return filepath.Join(dir, path[2:])
	} else if strings.HasPrefix(path, "$") {
		reg := regexp.MustCompile(`^\$\{?(.*?)\}?/`)
		env := reg.FindStringSubmatch(path)[1]
		parsedEnv := os.Getenv(env)
		return reg.ReplaceAllString(path, parsedEnv+"/")
	} else {
		return path
	}
}
