// utils/utils.go
package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PipeOpsHQ/pipeops-cli/internal/config"
	"github.com/PipeOpsHQ/pipeops-cli/libs"
	"github.com/spf13/viper"
)

func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	log.Printf("Command executed: %s %s", config.SanitizeLog(name), config.SanitizeLog(strings.Join(args, " ")))
	if outStr != "" {
		log.Printf("stdout: %s", config.SanitizeLog(outStr))
	}
	if errStr != "" {
		log.Printf("stderr: %s", config.SanitizeLog(errStr))
	}

	if err != nil {
		if errStr != "" {
			return outStr, errors.New(errStr)
		}
		return outStr, err
	}
	return outStr, nil
}

// RunCommandWithEnv runs a command with custom environment variables
func RunCommandWithEnv(name string, args []string, env []string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	log.Printf("Command executed: %s %s", config.SanitizeLog(name), config.SanitizeLog(strings.Join(args, " ")))
	if outStr != "" {
		log.Printf("stdout: %s", config.SanitizeLog(outStr))
	}
	if errStr != "" {
		log.Printf("stderr: %s", config.SanitizeLog(errStr))
	}

	if err != nil {
		if errStr != "" {
			return outStr, errors.New(errStr)
		}
		return outStr, err
	}
	return outStr, nil
}

// RunCommandWithEnvStreaming runs a command with custom environment variables
// while streaming stdout/stderr directly to the console.
// It still captures output for returning errors to callers.
func RunCommandWithEnvStreaming(name string, args []string, env []string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin

	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	log.Printf("Running command: %s %s", config.SanitizeLog(name), config.SanitizeLog(strings.Join(args, " ")))
	err := cmd.Run()

	outStr := strings.TrimSpace(stdout.String())
	errStr := strings.TrimSpace(stderr.String())

	if err != nil {
		if errStr != "" {
			return outStr, errors.New(errStr)
		}
		return outStr, err
	}

	return outStr, nil
}

// IsValidURL checks if the provided string is a valid URL.
func IsValidURL(testURL string) bool {
	parsedURL, err := url.ParseRequestURI(testURL)
	if err != nil {
		return false
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}

func ValidateOrPrompt() error {
	// Ensure the configuration is loaded before proceeding
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Unable to read config file. Proceeding to create or update it.")
	}

	token := viper.GetString("service_account_token")

	if token == "" {
		// Token is not found, prompt the user
		fmt.Println("No service token found. Let's fix that!")
		var err error
		token, err = promptForToken()
		if err != nil {
			return fmt.Errorf("failed to get token from user: %w", err)
		}
	}

	// Validate the token
	if !validateAndSaveToken(token) {
		for {
			fmt.Println("Invalid service token. Please try again.")
			var err error
			token, err = promptForToken()
			if err != nil {
				return fmt.Errorf("failed to get token from user: %w", err)
			}
			if validateAndSaveToken(token) {
				break
			}
		}
	}

	return nil
}

// validateAndSaveToken validates the token and saves it to the configuration if valid
func validateAndSaveToken(token string) bool {
	http := libs.NewHttpClient()
	_, err := http.VerifyToken(token, "")
	if err != nil {
		return false // Token is invalid
	}

	// Save the token to the configuration
	viper.Set("service_account_token", token)
	if err := saveConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving token to config file:", err)
		return false // Return false but don't exit - let caller handle it
	}

	fmt.Println("Token validated and saved successfully to the config file.")
	return true
}

// promptForToken prompts the user to input their service account token
func promptForToken() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your PipeOps Service Account Token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading token: %w", err)
	}
	return strings.TrimSpace(token), nil
}

// saveConfig writes the configuration to the file, handling potential missing file errors
func saveConfig() error {
	err := viper.WriteConfig()
	if err != nil {
		// Handle the case where the configuration file doesn't exist yet
		if configFileNotFoundError, ok := err.(viper.ConfigFileNotFoundError); ok {
			_ = configFileNotFoundError // Use the error to avoid unused variable warning
			return viper.SafeWriteConfig()
		}
		return err
	}
	return nil
}

// GetBaseName returns the base name of a directory path
func GetBaseName(path string) string {
	return filepath.Base(path)
}
