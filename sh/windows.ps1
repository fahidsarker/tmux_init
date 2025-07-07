$ErrorActionPreference = "Stop"

$Repo = "fahidsarker/tmux_init"
$BinaryName = "tinit"
$ProjectName = "tmux_init"
$InstallDir = "$env:USERPROFILE\.local\bin\$ProjectName"
$OS = "windows"

# Detect architecture
switch ($env:PROCESSOR_ARCHITECTURE) {
  "AMD64"  { $Arch = "amd64" }
  "ARM64"  { $Arch = "arm64" }
  "x86"    { $Arch = "386" }
  default  { Write-Error "❌ Unsupported architecture: $env:PROCESSOR_ARCHITECTURE"; exit 1 }
}

# Get latest release tag
Write-Host "📦 Fetching latest release..."
$Latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Latest.tag_name

# Construct download URL
$File = "$ProjectName" + "_" + $Tag.TrimStart("v") + "_" + "$OS" + "_" + "$Arch" + ".zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$File"
$TempFile = "$env:TEMP\$File"

# Download and extract
Write-Host "📥 Downloading $File..."
Invoke-WebRequest -Uri $Url -OutFile $TempFile

Write-Host "📂 Installing to $InstallDir..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Expand-Archive -Path $TempFile -DestinationPath $InstallDir
Remove-Item $TempFile

Write-Host "`n✅ Installed to $InstallDir"

# Add to PATH if not already
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not ($CurrentPath -split ";" | Where-Object { $_ -eq $InstallDir })) {
    $NewPath = "$CurrentPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Host "✅ Added $InstallDir to your user PATH."
    Write-Host "👉 Restart your terminal to apply changes."
} else {
    Write-Host "ℹ️ $InstallDir already in PATH."
}
