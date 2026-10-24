package server

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Autumn-27/artex/internal/chatgpt"
)

func TestAuthStorageMovePreservesSessionAndVault(t *testing.T) {
	base, data, target := t.TempDir(), t.TempDir(), t.TempDir()
	key, err := loadOrCreateJWTKey(base, data)
	if err != nil {
		t.Fatal(err)
	}
	session, err := signJWT(key)
	if err != nil {
		t.Fatal(err)
	}
	host, _, err := chatgpt.New(base, key).Registration("")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ARTEX_AUTH_DIR", target)
	if _, err := resolveAuthDir(base, data); err == nil {
		t.Fatal("empty mount accepted despite an existing installation key")
	}
	for _, name := range []string{jwtKeyFilename, filepath.Join("chatgpt-oauth", "credentials.enc")} {
		content, err := os.ReadFile(filepath.Join(base, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(target, name)), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(target, name), content, 0600); err != nil {
			t.Fatal(err)
		}
		if name == jwtKeyFilename {
			if _, err := resolveAuthDir(base, data); err == nil {
				t.Fatal("partial migration accepted without the existing OAuth vault")
			}
		}
	}
	dir, err := resolveAuthDir(base, data)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := loadOrCreateJWTKey(dir, data)
	if err != nil || !bytes.Equal(key, restored) || !verifyJWT(session, restored) {
		t.Fatal("existing session lost after moving credentials")
	}
	restoredHost, _, err := chatgpt.New(dir, restored).Registration("")
	if err != nil || restoredHost != host {
		t.Fatal("encrypted vault lost after move")
	}
}

func TestJWTStorageFailsWithoutOverwriting(t *testing.T) {
	for _, scenario := range []string{"invalid", "unreadable", "vault-without-key", "invalid-legacy"} {
		t.Run(scenario, func(t *testing.T) {
			dir, data := t.TempDir(), t.TempDir()
			path := filepath.Join(dir, jwtKeyFilename)
			switch scenario {
			case "invalid":
				if err := os.WriteFile(path, []byte("damaged"), 0600); err != nil {
					t.Fatal(err)
				}
			case "unreadable":
				if err := os.Mkdir(path, 0700); err != nil {
					t.Fatal(err)
				}
			case "vault-without-key":
				if err := os.Mkdir(filepath.Join(dir, "chatgpt-oauth"), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "chatgpt-oauth", "credentials.enc"), []byte("encrypted-fixture"), 0600); err != nil {
					t.Fatal(err)
				}
			case "invalid-legacy":
				if err := os.WriteFile(filepath.Join(data, jwtKeyFilename), []byte("damaged"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if _, err := loadOrCreateJWTKey(dir, data); err == nil {
				t.Fatal("damaged installation accepted")
			}
			if scenario == "invalid" {
				content, _ := os.ReadFile(path)
				if string(content) != "damaged" {
					t.Fatal("existing key overwritten")
				}
			} else if scenario != "unreadable" {
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Fatal("replacement key created")
				}
			}
		})
	}
}

func TestAuthDirectoryCannotExposeCredentials(t *testing.T) {
	base, data := t.TempDir(), t.TempDir()
	for _, dir := range []string{data, filepath.Join(data, "auth")} {
		t.Setenv("ARTEX_AUTH_DIR", dir)
		if _, err := resolveAuthDir(base, data); err == nil {
			t.Fatal("auth directory inside workspace accepted")
		}
	}
	link := filepath.Join(base, "workspace-link")
	if err := os.Symlink(data, link); err != nil {
		t.Skip("symlink unavailable")
	}
	t.Setenv("ARTEX_AUTH_DIR", filepath.Join(link, "auth"))
	if _, err := resolveAuthDir(base, data); err == nil {
		t.Fatal("symlink into workspace accepted")
	}
}

func TestJWTLegacyMigrationAndExclusiveCreate(t *testing.T) {
	dir, data := t.TempDir(), t.TempDir()
	key := []byte("synthetic-test-key-not-a-real-secret")
	legacy := filepath.Join(data, jwtKeyFilename)
	if err := os.WriteFile(legacy, key, 0600); err != nil {
		t.Fatal(err)
	}
	got, err := loadOrCreateJWTKey(dir, data)
	if err != nil || !bytes.Equal(got, key) {
		t.Fatal("legacy key not preserved")
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Fatal("legacy key still exposed")
	}
	if err := writeNewJWTKey(filepath.Join(dir, jwtKeyFilename), []byte("replacement")); err == nil {
		t.Fatal("key overwritten")
	}
}
