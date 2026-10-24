package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Native installs retain the executable directory. Containers can persist just
// the credentials without mounting over the executable or exposing them in data/.
func resolveAuthDir(baseDir, dataDir string) (string, error) {
	dir := strings.TrimSpace(os.Getenv("ARTEX_AUTH_DIR"))
	if dir == "" {
		dir = baseDir
	}
	auth, err := resolvedStoragePath(dir)
	if err != nil {
		return "", err
	}
	data, err := resolvedStoragePath(dataDir)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(data, auth)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("authentication directory must be outside the browsable data directory")
	}
	// An operator changing the location must copy the key and vault together first.
	// An empty new mount must not silently replace an existing installation key.
	base, err := resolvedStoragePath(baseDir)
	if err != nil {
		return "", err
	}
	if auth != base {
		for _, name := range []string{jwtKeyFilename, filepath.Join("chatgpt-oauth", "credentials.enc")} {
			if _, err := os.Stat(filepath.Join(auth, name)); os.IsNotExist(err) {
				if _, oldErr := os.Stat(filepath.Join(base, name)); !os.IsNotExist(oldErr) {
					return "", fmt.Errorf("copy the existing JWT key and ChatGPT vault into ARTEX_AUTH_DIR before starting")
				}
			}
		}
	}
	return auth, nil
}

// Resolve existing ancestors as well, so a new directory below a symlink into
// data/ cannot bypass the workspace isolation check.
func resolvedStoragePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}
	if !os.IsNotExist(err) || filepath.Dir(abs) == abs {
		return "", err
	}
	parent, err := resolvedStoragePath(filepath.Dir(abs))
	if err != nil {
		return "", err
	}
	return filepath.Join(parent, filepath.Base(abs)), nil
}
