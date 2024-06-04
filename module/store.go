package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
  "os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

func main() {
	// Verify the directory
	if err := verifyDirectory(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Ask user if they are inputting or connecting
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please select an option:")
	fmt.Println("1 - Input")
	fmt.Println("2 - Connect")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, err := strconv.Atoi(input)
	if err != nil || (choice != 1 && choice != 2) {
		fmt.Println("Invalid input. Please enter 1 or 2.")
		return
	}

	switch choice {
	case 1:
		inputFunc()
	case 2:
		connectFunc()
	}
}

// Verify directory exists, create if not
func verifyDirectory() error {
	// Get current user
	usr, err := user.Current()
	if err != nil {
		return err
	}

	// Define the directory path
	dirPath := filepath.Join(usr.HomeDir, ".pswdcnt")

	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", dirPath)
	} else if err != nil {
		return err
	} else {
		fmt.Println("Directory already exists:", dirPath)
	}
	return nil
}

// ENCRYPTION -----------------------------------------------------------------------

func Encrypt(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(cipherText), nil
}

// Decrypt decrypts cipher text string into plain text string

func Decrypt(cipherText string, key []byte) (string, error) {
	decodedCipherText, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(decodedCipherText) < nonceSize {
		return "", fmt.Errorf("cipher text too short")
	}

	nonce, cipherTextBytes := decodedCipherText[:nonceSize], decodedCipherText[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// GenerateKey generates a new encryption key
func GenerateKey(password string) ([]byte, error) {
	salt := []byte("a_very_secure_salt") // This should be stored and used securely
	return scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
}

// ENCRYPTION -----------------------------------------------------------------------------

func inputFunc() {
	fmt.Println("Selected input")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the SSID: ")
	ssid, _ := reader.ReadString('\n')
	ssid = strings.TrimSpace(ssid)

	fmt.Print("Enter the password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Generate encryption key from password (for demonstration purposes)
	encKey, err := GenerateKey("some_secure_password")
	if err != nil {
		fmt.Println("Error generating encryption key:", err)
		return
	}

	// Encrypt the SSID and password
	encSSID, err := Encrypt(ssid, encKey)
	if err != nil {
		fmt.Println("Error encrypting SSID:", err)
		return
	}

	encPassword, err := Encrypt(password, encKey)
	if err != nil {
		fmt.Println("Error encrypting password:", err)
		return
	}

	// Ensure the directory exists
	dirPath := filepath.Join(os.Getenv("HOME"), ".pswdcnt")
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create the file and write the encrypted SSID and password
	filePath := filepath.Join(dirPath, ssid+".txt")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("SSID: %s\nPassword: %s\n", encSSID, encPassword))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("SSID and password saved successfully.")
}


func connectFunc() {
	fmt.Println("Selected connect")

	// List all SSID files in the directory
	dirPath := filepath.Join(os.Getenv("HOME"), ".pswdcnt")
	files, err := filepath.Glob(filepath.Join(dirPath, "*.txt"))
	if err != nil {
		fmt.Println("Error listing files:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No SSID files found in the directory.")
		return
	}

	fmt.Println("Available SSIDs:")
	for _, file := range files {
		fmt.Println(filepath.Base(file))
	}

	// Ask user to choose an SSID file
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of the SSID file: ")
	ssidFile, _ := reader.ReadString('\n')
	ssidFile = strings.TrimSpace(ssidFile)

	// Read the content of the SSID file
	filePath := filepath.Join(dirPath, ssidFile)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ssid := ""
	password := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "SSID:") {
			ssid = strings.TrimSpace(strings.TrimPrefix(line, "SSID:"))
		} else if strings.HasPrefix(line, "Password:") {
			password = strings.TrimSpace(strings.TrimPrefix(line, "Password:"))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Decrypt SSID and password
	encKey, err := GenerateKey("some_secure_password") // Use the appropriate password used during encryption
	if err != nil {
		fmt.Println("Error generating encryption key:", err)
		return
	}

	plainSSID, err := Decrypt(ssid, encKey)
	if err != nil {
		fmt.Println("Error decrypting SSID:", err)
		return
	}

	plainPassword, err := Decrypt(password, encKey)
	if err != nil {
		fmt.Println("Error decrypting password:", err)
		return
	}

	// Connect to the Wi-Fi network
	cmd := fmt.Sprintf("nmcli device wifi connect \"%s\" password \"%s\"", plainSSID, plainPassword)
	fmt.Println("Executing command:", cmd)

	// Run the command
	cmdOutput, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	fmt.Println("Command output:")
	fmt.Println(string(cmdOutput))
}

