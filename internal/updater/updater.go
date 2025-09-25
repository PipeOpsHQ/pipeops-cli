package updater

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/PipeOpsHQ/pipeops-cli/utils"
)

const (
	DefaultGitHubRepo = "PipeOpsHQ/pipeops-cli" // Reverted back to actual repository
	// For separate releases repo, use: "PipeOpsHQ/pipeops-cli-releases"
	UpdateCheckInterval = 24 * time.Hour
)

// GetGitHubRepo returns the GitHub repository to use, checking environment variable first
func GetGitHubRepo() string {
	if repo := os.Getenv("PIPEOPS_GITHUB_REPO"); repo != "" {
		return repo
	}
	return DefaultGitHubRepo
}

// getGitHubAPIURL returns the GitHub API URL for the configured repository
func getGitHubAPIURL() string {
	// For custom update endpoint, use environment variable:
	// if customURL := os.Getenv("PIPEOPS_UPDATE_URL"); customURL != "" {
	//     return customURL
	// }
	return "https://api.github.com/repos/" + GetGitHubRepo() + "/releases/latest"
}

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
	Body        string    `json:"body"`
}

// Asset represents a GitHub release asset
type Asset struct {
	Name               string `json:"name"`
	ContentType        string `json:"content_type"`
	Size               int64  `json:"size"`
	DownloadCount      int    `json:"download_count"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// UpdateService handles CLI updates
type UpdateService struct {
	client         *http.Client
	currentVersion string
}

// NewUpdateService creates a new update service
func NewUpdateService(currentVersion string) *UpdateService {
	return &UpdateService{
		client:         &http.Client{Timeout: 30 * time.Second},
		currentVersion: currentVersion,
	}
}

// CheckForUpdates checks if a new version is available
func (s *UpdateService) CheckForUpdates(ctx context.Context) (*Release, bool, error) {
	// Fetch latest release
	release, err := s.fetchLatestRelease(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch latest release: %w", err)
	}

	// Compare versions
	hasUpdate, err := s.compareVersions(s.currentVersion, release.TagName)
	if err != nil {
		return nil, false, fmt.Errorf("failed to compare versions: %w", err)
	}

	return release, hasUpdate, nil
}

// fetchLatestRelease fetches the latest release from GitHub
func (s *UpdateService) fetchLatestRelease(ctx context.Context) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", getGitHubAPIURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "PipeOps-CLI-Updater")

	// Add authentication if GitHub token is provided
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &release, nil
}

// compareVersions compares two version strings
func (s *UpdateService) compareVersions(current, latest string) (bool, error) {
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle dev versions
	if current == "dev" {
		return true, nil // Always show updates for dev versions
	}

	// Parse versions using semantic versioning
	currentParts, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("failed to parse current version %s: %w", current, err)
	}

	latestParts, err := parseVersion(latest)
	if err != nil {
		return false, fmt.Errorf("failed to parse latest version %s: %w", latest, err)
	}

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true, nil
		} else if latestParts[i] < currentParts[i] {
			return false, nil
		}
	}

	return false, nil // Versions are equal
}

// parseVersion parses a version string into [major, minor, patch] integers
func parseVersion(version string) ([]int, error) {
	// Remove any build metadata (e.g., "1.2.3-beta.1" -> "1.2.3")
	version = strings.Split(version, "-")[0]

	// Split by dots
	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		// Pad with zeros if necessary
		for len(parts) < 3 {
			parts = append(parts, "0")
		}
	}

	var result []int
	for i := 0; i < 3; i++ {
		// Extract numeric part only
		re := regexp.MustCompile(`\d+`)
		match := re.FindString(parts[i])
		if match == "" {
			result = append(result, 0)
		} else {
			var num int
			if _, err := fmt.Sscanf(match, "%d", &num); err != nil {
				return nil, fmt.Errorf("failed to parse version part %s: %w", parts[i], err)
			}
			result = append(result, num)
		}
	}

	return result, nil
}

// UpdateCLI downloads and installs the latest version
func (s *UpdateService) UpdateCLI(ctx context.Context, release *Release, opts utils.OutputOptions) error {
	// Find the appropriate asset for the current platform
	asset, err := s.findAssetForPlatform(release)
	if err != nil {
		return fmt.Errorf("failed to find asset for platform: %w", err)
	}

	utils.PrintInfo(fmt.Sprintf("Downloading %s (%s)...", asset.Name, formatSize(asset.Size)), opts)

	// Download the asset
	downloadPath, err := s.downloadAsset(ctx, asset, opts)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer os.Remove(downloadPath)

	// Extract and install
	if err := s.extractAndInstall(downloadPath, asset.Name, opts); err != nil {
		return fmt.Errorf("failed to extract and install: %w", err)
	}

	utils.PrintSuccess(fmt.Sprintf("Successfully updated to version %s", release.TagName), opts)
	return nil
}

// findAssetForPlatform finds the appropriate asset for the current platform
func (s *UpdateService) findAssetForPlatform(release *Release) (*Asset, error) {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	// Map Go arch names to release arch names
	switch archName {
	case "amd64":
		archName = "x86_64"
	case "386":
		archName = "i386"
	}

	// Map Go OS names to release OS names
	switch osName {
	case "darwin":
		osName = "Darwin"
	case "linux":
		osName = "Linux"
	case "windows":
		osName = "Windows"
	}

	// Look for matching asset
	for _, asset := range release.Assets {
		name := asset.Name
		if strings.Contains(name, osName) && strings.Contains(name, archName) {
			return &asset, nil
		}
	}

	return nil, fmt.Errorf("no asset found for platform %s/%s", osName, archName)
}

// downloadAsset downloads an asset to a temporary file
func (s *UpdateService) downloadAsset(ctx context.Context, asset *Asset, opts utils.OutputOptions) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "pipeops-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Copy with progress (for large files)
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write download: %w", err)
	}

	return tempFile.Name(), nil
}

// extractAndInstall extracts the downloaded archive and installs the binary
func (s *UpdateService) extractAndInstall(archivePath, assetName string, opts utils.OutputOptions) error {
	// Get current executable path
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Create temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "pipeops-extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract the archive
	var binaryPath string
	if strings.HasSuffix(assetName, ".zip") {
		binaryPath, err = s.extractZip(archivePath, tempDir)
	} else if strings.HasSuffix(assetName, ".tar.gz") {
		binaryPath, err = s.extractTarGz(archivePath, tempDir)
	} else {
		return fmt.Errorf("unsupported archive format: %s", assetName)
	}

	if err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Make the binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Replace current executable
	utils.PrintInfo("Installing new binary...", opts)
	if err := s.replaceExecutable(currentExePath, binaryPath); err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	return nil
}

// extractZip extracts a zip archive and returns the path to the binary
func (s *UpdateService) extractZip(archivePath, destDir string) (string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	var binaryPath string
	for _, file := range reader.File {
		// Look for the binary file
		if file.Name == "pipeops.exe" || strings.HasSuffix(file.Name, "/pipeops.exe") ||
			file.Name == "pipeops" || strings.HasSuffix(file.Name, "/pipeops") {

			// Determine the correct binary name
			binaryName := "pipeops"
			if runtime.GOOS == "windows" {
				binaryName = "pipeops.exe"
			}
			binaryPath = filepath.Join(destDir, binaryName)

			// Open the file in the zip
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			// Create the output file
			outFile, err := os.Create(binaryPath)
			if err != nil {
				return "", fmt.Errorf("failed to create binary file: %w", err)
			}
			defer outFile.Close()

			// Copy the file contents
			if _, err := io.Copy(outFile, rc); err != nil {
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in zip archive")
	}

	return binaryPath, nil
}

// extractTarGz extracts a tar.gz archive and returns the path to the binary
func (s *UpdateService) extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	var binaryPath string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar header: %w", err)
		}

		// Look for the binary file
		if header.Typeflag == tar.TypeReg && (header.Name == "pipeops" || strings.HasSuffix(header.Name, "/pipeops")) {
			binaryPath = filepath.Join(destDir, "pipeops")

			// Create the file
			outFile, err := os.Create(binaryPath)
			if err != nil {
				return "", fmt.Errorf("failed to create binary file: %w", err)
			}
			defer outFile.Close()

			// Copy the file contents
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", fmt.Errorf("failed to extract binary: %w", err)
			}
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

// replaceExecutable replaces the current executable with the new one
func (s *UpdateService) replaceExecutable(currentPath, newPath string) error {
	// Create backup
	backupPath := currentPath + ".backup"
	if err := os.Rename(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Copy new binary
	if err := copyFile(newPath, currentPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, currentPath)
		return fmt.Errorf("failed to copy new binary: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return dstFile.Chmod(srcInfo.Mode())
}

// formatSize formats a byte size for display
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
