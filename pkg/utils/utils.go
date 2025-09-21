package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ClearScreen clears the terminal screen
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// FormatNumber formats large numbers with commas
func FormatNumber(n int) string {
	str := strconv.Itoa(n)
	if len(str) <= 3 {
		return str
	}

	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}
	return result
}

// GenerateProgressBar creates a visual progress bar
func GenerateProgressBar(progress, width int) string {
	if progress > 100 {
		progress = 100
	}
	if progress < 0 {
		progress = 0
	}

	filled := (progress * width) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}

// FormatDuration formats a duration into HH:MM:SS format
func FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// ValidatePositiveInt validates that a string represents a positive integer
func ValidatePositiveInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a valid number: %v", err)
	}

	if num <= 0 {
		return 0, fmt.Errorf("number must be positive")
	}

	return num, nil
}

// TruncateString truncates a string to a maximum length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return s[:maxLen]
	}

	return s[:maxLen-3] + "..."
}

// SafeDivide performs division with zero check
func SafeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}
