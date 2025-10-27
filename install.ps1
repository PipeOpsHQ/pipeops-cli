# PipeOps CLI Installation Script for Windows
# This script downloads and installs the latest version of PipeOps CLI

param(
    [string]$InstallDir = "$env:USERPROFILE\bin",
    [string]$Version = "",
    [switch]$Help
)

# Configuration
$Repo = "PipeOpsHQ/pipeops-cli"
$BinaryName = "pipeops"

# Colors for output
$Colors = @{
    Red = "Red"
    Green = "Green"
    Yellow = "Yellow"
    Blue = "Cyan"
    Purple = "Magenta"
    White = "White"
}

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

function Write-Header {
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor $Colors.Purple
    Write-Host "🚀 PipeOps CLI Installer (Windows)" -ForegroundColor $Colors.Purple
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor $Colors.Purple
    Write-Host ""
}

function Show-Help {
    Write-Host "Usage: .\install.ps1 [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -InstallDir DIR     Install to specific directory (default: $env:USERPROFILE\bin)"
    Write-Host "  -Version VERSION    Install specific version (default: latest)"
    Write-Host "  -Help              Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install.ps1"
    Write-Host "  .\install.ps1 -InstallDir 'C:\tools'"
    Write-Host "  .\install.ps1 -Version 'v1.0.0'"
}

function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        "x86" { return "i386" }
        default {
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

function Get-LatestVersion {
    Write-Status "Fetching latest release information..."

    $apiUrl = "https://api.github.com/repos/$Repo/releases/latest"

    try {
        $response = Invoke-RestMethod -Uri $apiUrl -Method Get
        $version = $response.tag_name

        if (-not $version) {
            Write-Error "Failed to get latest version"
            exit 1
        }

        Write-Status "Latest version: $version"
        return $version
    }
    catch {
        Write-Error "Failed to fetch latest version: $_"
        exit 1
    }
}

function Install-Binary {
    param([string]$Version)

    $arch = Get-Architecture
    $filename = "${BinaryName}-cli_Windows_${arch}.zip"
    $url = "https://github.com/$Repo/releases/download/$Version/$filename"

    Write-Status "Downloading $filename..."

    # Create temporary directory
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    New-Item -ItemType Directory -Path $tempDir | Out-Null

    try {
        # Download the release
        $downloadPath = Join-Path $tempDir $filename

        try {
            Invoke-WebRequest -Uri $url -OutFile $downloadPath -UseBasicParsing
        }
        catch {
            Write-Error "Failed to download from $url : $_"
            exit 1
        }

        Write-Status "Extracting binary..."

        # Extract the archive
        try {
            Expand-Archive -Path $downloadPath -DestinationPath $tempDir -Force
        }
        catch {
            Write-Error "Failed to extract archive: $_"
            exit 1
        }

        # Find the binary
        $binaryPath = Join-Path $tempDir "$BinaryName.exe"

        if (-not (Test-Path $binaryPath)) {
            Write-Error "Binary not found in extracted archive"
            exit 1
        }

        Write-Status "Installing to $InstallDir..."

        # Create install directory if it doesn't exist
        if (-not (Test-Path $InstallDir)) {
            try {
                New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
            }
            catch {
                Write-Error "Failed to create install directory $InstallDir : $_"
                exit 1
            }
        }

        # Install binary
        $installPath = Join-Path $InstallDir "$BinaryName.exe"

        try {
            Copy-Item $binaryPath $installPath -Force
        }
        catch {
            Write-Error "Failed to copy binary to $installPath : $_"
            exit 1
        }

        Write-Success "PipeOps CLI installed successfully!"
        return $installPath
    }
    finally {
        # Clean up temporary directory
        if (Test-Path $tempDir) {
            Remove-Item $tempDir -Recurse -Force
        }
    }
}

function Test-Installation {
    param([string]$InstallPath)

    Write-Status "Verifying installation..."

    if (Test-Path $InstallPath) {
        try {
            $versionOutput = & $InstallPath version 2>$null
            Write-Success "Installation verified: $($versionOutput.Split([Environment]::NewLine)[0])"
            return $true
        }
        catch {
            Write-Warning "Binary installed but version check failed"
            return $false
        }
    }
    else {
        Write-Error "Binary not found at $InstallPath"
        return $false
    }
}

function Update-Path {
    param([string]$InstallDir)

    $currentPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::User)

    if ($currentPath -notlike "*$InstallDir*") {
        Write-Status "Adding $InstallDir to user PATH..."

        try {
            $newPath = "$InstallDir;$currentPath"
            [Environment]::SetEnvironmentVariable("PATH", $newPath, [EnvironmentVariableTarget]::User)

            # Also update current session PATH
            $env:PATH = "$InstallDir;$env:PATH"

            Write-Success "PATH updated successfully"
            Write-Warning "You may need to restart your terminal for PATH changes to take effect"
        }
        catch {
            Write-Warning "Failed to update PATH: $_"
            Write-Warning "Please manually add $InstallDir to your PATH"
        }
    }
    else {
        Write-Status "$InstallDir is already in PATH"
    }
}

function Write-Completion {
    Write-Host ""
    Write-Host "🎉 Installation Complete!" -ForegroundColor $Colors.Green
    Write-Host ""
    Write-Host "Getting Started:" -ForegroundColor $Colors.Blue
    Write-Host "  1. Authenticate with PipeOps:"
    Write-Host "     $BinaryName auth login"
    Write-Host ""
    Write-Host "  2. List your projects:"
    Write-Host "     $BinaryName project list"
    Write-Host ""
    Write-Host "  3. Get help:"
    Write-Host "     $BinaryName --help"
    Write-Host ""
    Write-Host "Documentation:" -ForegroundColor $Colors.Blue
    Write-Host "  🌐 Website: https://pipeops.io"
    Write-Host "  📖 Docs: https://docs.pipeops.io"
    Write-Host "  💬 Discord: https://discord.gg/pipeops"
    Write-Host "  🐙 GitHub: https://github.com/$Repo"
    Write-Host ""
}

# Main execution
function Main {
    Write-Header

    if ($Help) {
        Show-Help
        return
    }

    # Check PowerShell version
    if ($PSVersionTable.PSVersion.Major -lt 3) {
        Write-Error "PowerShell 3.0 or later is required"
        exit 1
    }

    # Get version
    if (-not $Version) {
        $Version = Get-LatestVersion
    }

    # Install binary
    $installPath = Install-Binary -Version $Version

    # Test installation
    $success = Test-Installation -InstallPath $installPath

    # Update PATH
    Update-Path -InstallDir $InstallDir

    # Show completion message
    Write-Completion

    if (-not $success) {
        Write-Warning "Installation completed but verification failed"
        exit 1
    }
}

# Run main function
try {
    Main
}
catch {
    Write-Error "Installation failed: $_"
    exit 1
}