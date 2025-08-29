package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// PublicKeyPath and PrivateKeyPath store the file paths for the
// Ed25519 public and private keys, dynamically determined based on the OS.
var (
	PublicKeyPath      string
	PrivateKeyPath     string
	TempPublicKeyPath  string
	TempPrivateKeyPath string
)

// init initializes the key file paths depending on the operating system.
// It sets sensible default paths based on platform conventions.
func init() {
	switch runtime.GOOS {
	case "windows":
		// %APPDATA%\libr\keys\
		appData := os.Getenv("APPDATA")
		PrivateKeyPath = filepath.Join(appData, "libr", "keys", "priv.key")
		PublicKeyPath = filepath.Join(appData, "libr", "keys", "pub.key")
		TempPrivateKeyPath = filepath.Join(appData, "libr", "keys", "temppriv.key")
		TempPublicKeyPath = filepath.Join(appData, "libr", "keys", "temppub.key")

	case "darwin":
		// ~/Library/Application Support/libr/keys/
		home, err := os.UserHomeDir()
		if err != nil {
			panic("unable to get user home directory: " + err.Error())
		}
		PrivateKeyPath = filepath.Join(home, "Library", "Application Support", "libr", "keys", "priv.key")
		PublicKeyPath = filepath.Join(home, "Library", "Application Support", "libr", "keys", "pub.key")
		TempPrivateKeyPath = filepath.Join(home, "Library", "Application Support", "libr", "keys", "temppriv.key")
		TempPublicKeyPath = filepath.Join(home, "Library", "Application Support", "libr", "keys", "temppub.key")

	case "linux":
		// ~/.config/libr/keys/
		home, err := os.UserHomeDir()
		if err != nil {

			panic("unable to get user home directory: " + err.Error())
		}
		PrivateKeyPath = filepath.Join(home, ".config", "libr", "keys", "priv.key")
		PublicKeyPath = filepath.Join(home, ".config", "libr", "keys", "pub.key")
		TempPrivateKeyPath = filepath.Join(home, ".config", "libr", "keys", "temppriv.key")
		TempPublicKeyPath = filepath.Join(home, ".config", "libr", "keys", "temppub.key")

	default:
		// Fallback to relative path: ./keys/
		PrivateKeyPath = filepath.Join("keys", "priv.key")
		PublicKeyPath = filepath.Join("keys", "pub.key")
		TempPrivateKeyPath = filepath.Join("keys", "temppriv.key")
		TempPublicKeyPath = filepath.Join("keys", "temppub.key")
	}
}
