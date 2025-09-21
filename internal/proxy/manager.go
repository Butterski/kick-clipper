package proxy

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"kick-clipper/internal/models"
)

// Manager handles proxy operations
type Manager struct {
	proxies []models.Proxy
}

// NewManager creates a new proxy manager
func NewManager() *Manager {
	return &Manager{
		proxies: make([]models.Proxy, 0),
	}
}

// LoadFromFile loads proxies from a file
func (pm *Manager) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open proxy file: %v", err)
	}
	defer file.Close()

	var proxyList []models.Proxy
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			fmt.Printf("Warning: Invalid proxy format on line %d: %s (expected ip:port:user:pass)\n", lineNumber, line)
			continue
		}

		proxy := models.Proxy{
			IP:       strings.TrimSpace(parts[0]),
			Port:     strings.TrimSpace(parts[1]),
			Username: strings.TrimSpace(parts[2]),
			Password: strings.TrimSpace(parts[3]),
		}

		// Basic validation
		if proxy.IP == "" || proxy.Port == "" {
			fmt.Printf("Warning: Invalid proxy on line %d: empty IP or port\n", lineNumber)
			continue
		}

		proxyList = append(proxyList, proxy)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading proxy file: %v", err)
	}

	pm.proxies = proxyList
	return nil
}

// GetRandom returns a random proxy from the list
func (pm *Manager) GetRandom() *models.Proxy {
	if len(pm.proxies) == 0 {
		return nil
	}
	return &pm.proxies[rand.Intn(len(pm.proxies))]
}

// Count returns the number of loaded proxies
func (pm *Manager) Count() int {
	return len(pm.proxies)
}

// GetProxyURL formats a proxy for HTTP client use
func GetProxyURL(proxy *models.Proxy) string {
	if proxy == nil {
		return ""
	}
	return fmt.Sprintf("http://%s:%s@%s:%s",
		proxy.Username, proxy.Password, proxy.IP, proxy.Port)
}

// IsEmpty checks if the proxy manager has no proxies
func (pm *Manager) IsEmpty() bool {
	return len(pm.proxies) == 0
}

// GetAll returns all proxies (for debugging/stats)
func (pm *Manager) GetAll() []models.Proxy {
	return pm.proxies
}
