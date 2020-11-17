package util

import (
	"os/exec"
)

func Install() {
	exec.Command("git", "clone", URL_INSTALL_REPO, "/installer").Output()
	exec.Command("chmod", "+x", "/instller/main.sh").Output()
	exec.Command("bash", "/installer/main.sh").Output()
}
