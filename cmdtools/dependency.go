package cmdtools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func isRoot() bool {
	return os.Geteuid() == 0
}

func getPackageManager() (manager string, pkg string) {
	if runtime.GOOS != "linux" {
		return "", ""
	}
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "apt", "speedtest-cli"
	}
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return "yum", "speedtest-cli"
	}
	return "", ""
}

func runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("-> Running: %s %v\n", command, args)
	return cmd.Run()
}

func installDependency() error {
	if !isRoot() {
		return fmt.Errorf("installation requires root privileges. please re-run the command with sudo")
	}

	pm, pkg := getPackageManager()
	if pm == "" {
		return fmt.Errorf("unsupported package manager. please install 'speedtest-cli' manually")
	}

	fmt.Printf("Attempting to install '%s' using '%s'...\n", pkg, pm)

	var err error
	if pm == "apt" {
		err = runCommand("apt", "update", "-y")
		if err != nil {
			return fmt.Errorf("failed to update apt cache: %w", err)
		}
		err = runCommand("apt", "install", "-y", pkg)
	} else if pm == "yum" {
		err = runCommand("yum", "install", "-y", pkg)
	}

	if err != nil {
		return fmt.Errorf("failed to install '%s': %w", pkg, err)
	}

	fmt.Printf("'%s' installed successfully.\n", pkg)
	return nil
}

func EnsureDependency() bool {
	if commandExists("speedtest-cli") {
		return true
	}

	fmt.Println("Required dependency 'speedtest-cli' is missing for the speed test feature.")
	fmt.Print("Do you want to attempt installation now? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" {
		fmt.Println("Skipping installation. The speed test feature will be unavailable.")
		return true
	}

	err := installDependency()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Server startup aborted.")
		return false
	}

	return true
}

func RemoveDependency() {
	if !commandExists("speedtest-cli") {
		fmt.Println("speedtest-cli is not installed.")
		return
	}

	if !isRoot() {
		fmt.Println("Removal requires root privileges. Please run with sudo:")
		fmt.Printf("sudo %s --remove-deps\n", os.Args[0])
		os.Exit(1)
	}

	pm, pkg := getPackageManager()
	if pm == "" {
		fmt.Println("Unsupported package manager. Cannot determine how to remove 'speedtest-cli'.")
		os.Exit(1)
	}

	fmt.Printf("Attempting to remove '%s' using '%s'...\n", pkg, pm)
	err := runCommand(pm, "remove", "-y", pkg)

	if err != nil {
		fmt.Printf("Failed to remove '%s': %v\n", pkg, err)
		os.Exit(1)
	}

	fmt.Printf("'%s' removed successfully.\n", pkg)
}