package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

// ProjectRoot returns the project root.
func ProjectRoot() string {
	_, thisFile, _, _ := runtime.Caller(0)
	absPath, _ := filepath.Abs(path.Join(thisFile, "..", ".."))
	return absPath
}

// LoadRelativeFile loads the file at the specified path relative to
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
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
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

// CopyFile copies the file at src path to dst path. It applies
// the permissions specified in mode to the destination file.
// Mode values are 0644, 0755, etc.
func CopyFile(src string, dst string, mode os.FileMode) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, mode)
}

// Returns true if the file at path exists, false if not.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// ListIsEmpty returns true if the slice contains no
// elements, or if all the elements are empty strings.
func ListIsEmpty(list []string) bool {
	isEmpty := true
	if list == nil || len(list) == 0 {
		return isEmpty
	}
	for _, item := range list {
		if item != "" {
			isEmpty = false
			break
		}
	}
	return isEmpty
}

// InterfaceList converts a list of strings to a list of interfaces.
func InterfaceList(items []string) []interface{} {
	list := make([]interface{}, len(items))
	for i, item := range items {
		list[i] = item
	}
	return list
}

// SplitCamelCase splits camel-case identifiers into multiple words.
// Note that it does not split on multiple consecutive caps, so
// param CurrencyUSD would return ["Currency", "USD"].
//
// If max is less than zero, this will split into all words. If max
// is > 0, this will split into max words.
func SplitCamelCase(str string, max int) []string {
	var b bytes.Buffer
	partsCount := 0
	priorLower := false
	for _, v := range str {
		if priorLower && unicode.IsUpper(v) && (max < 0 || partsCount < max-1) {
			b.WriteByte(' ')
			partsCount++
		}
		b.WriteRune(v)
		priorLower = unicode.IsLower(v)
	}
	return strings.Split(b.String(), " ")
}

// ToHumanSize converts a raw byte count (size) to a human-friendly
// representation.
func ToHumanSize(size, unit int64) string {
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}

// ConsoleDebug prints a message to the console if the following we are
// running in the dev or test environment and we are not running automated
// tests. We want to see these messages in the console when we're doing
// interactive testing in the dev or test environments, but NOT when running
// automated tests because they clutter the test output.
func ConsoleDebug(message string) {
	envName := Context().Config.EnvName
	if !TestsAreRunning() && (envName == "dev" || envName == "test") {
		fmt.Println(message)
	}
}
