package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// ProjectRoot returns the project root.
func ProjectRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	absPath, _ := filepath.Abs(path.Join(thisFile, "..", ".."))
	return absPath
}

// LoadTestFile loads the file at the specified path relative to
// ProjectRoot() and returns the contents as a byte array.
//
// Example:
//
// bytes, err := LoadRelativeFile("db/fixtures/work_items.csv")
//
func LoadRelativeFile(relpath string) ([]byte, error) {
	absPath := filepath.Join(ProjectRoot(), relpath)
	return ioutil.ReadFile(absPath)
}

// TestsAreRunning returns true when code is running under "go test"
func TestsAreRunning() bool {
	return flag.Lookup("test.v") != nil
}

// PrintAndExit prints a message to STDERR and exits
func PrintAndExit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// Expands the tilde in a directory path to the current
// user's home directory. For example, on Linux, ~/data
// would expand to something like /home/josie/data
func ExpandTilde(filePath string) (string, error) {
	if strings.Index(filePath, "~") < 0 {
		return filePath, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	homeDir := usr.HomeDir + "/"
	expandedDir := strings.Replace(filePath, "~/", homeDir, 1)
	return expandedDir, nil
}

// Hash returns an encrypted version of plaintext that cannot
// be decrypted. This is suitable for encrypting passwords,
// reset-tokens, etc. The combined use of md5 plus plaintext salt
// plus sha256 provides some protection against rainbow tables.
func Hash(plaintext string) string {
	plain := []byte(plaintext)
	md5Digest := []byte(fmt.Sprintf("%x", md5.Sum(plain)))
	return fmt.Sprintf("%x", sha256.Sum256(append(md5Digest, plain...)))
}

// EncryptAES encrypts plaintext using key.
func EncryptAES(key []byte, plaintext string) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(encrypted), nil
}

// DecryptAES decrypts hex-encoded ciphertext using key.
func DecryptAES(key []byte, hexCipher string) (string, error) {
	ciphertext, err := hex.DecodeString(hexCipher)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcmDecrypt.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("Wrong nonce size")
	}
	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
