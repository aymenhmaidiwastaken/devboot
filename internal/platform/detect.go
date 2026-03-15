package platform

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type OS int

const (
	MacOS OS = iota
	Linux
	WSL
	Windows
	Unknown
)

func (o OS) String() string {
	switch o {
	case MacOS:
		return "macOS"
	case Linux:
		return "Linux"
	case WSL:
		return "WSL"
	case Windows:
		return "Windows"
	default:
		return "Unknown"
	}
}

type Distro int

const (
	DistroNone Distro = iota
	DistroDebian
	DistroUbuntu
	DistroArch
	DistroFedora
	DistroOther
)

func (d Distro) String() string {
	switch d {
	case DistroDebian:
		return "Debian"
	case DistroUbuntu:
		return "Ubuntu"
	case DistroArch:
		return "Arch"
	case DistroFedora:
		return "Fedora"
	default:
		return ""
	}
}

type Platform struct {
	OS     OS
	Distro Distro
	Arch   string
}

// Detect returns the current platform information.
func Detect() Platform {
	p := Platform{
		Arch: runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "darwin":
		p.OS = MacOS
	case "windows":
		p.OS = Windows
	case "linux":
		if isWSL() {
			p.OS = WSL
		} else {
			p.OS = Linux
		}
		p.Distro = detectDistro()
	default:
		p.OS = Unknown
	}

	return p
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(data))
	return strings.Contains(lower, "microsoft") || strings.Contains(lower, "wsl")
}

func detectDistro() Distro {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return DistroOther
	}
	content := strings.ToLower(string(data))

	if strings.Contains(content, "ubuntu") {
		return DistroUbuntu
	}
	if strings.Contains(content, "debian") {
		return DistroDebian
	}
	if strings.Contains(content, "arch") {
		return DistroArch
	}
	if strings.Contains(content, "fedora") {
		return DistroFedora
	}
	return DistroOther
}

// HasCommand checks if a command is available in PATH.
func HasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
