package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BackupRegistry string
	NamespacesToIgnore []string
	MaxConcurrentReconciles int
}
func GetConfig() (*Config, error) {
	namespaces := os.Getenv("NAMESPACES_TO_IGNORE")
	namespacesList := strings.Split(namespaces, ",")
	maxReconciles, _:= strconv.Atoi(os.Getenv("MAX_CONCURRENT_RECONCILES"))

	configuration := Config{
		BackupRegistry:    os.Getenv("BACKUP_REGISTRY"),
		NamespacesToIgnore: namespacesList,
		MaxConcurrentReconciles: maxReconciles,
	}

	if len(configuration.BackupRegistry) == 0 {
		return nil, fmt.Errorf("no backup registry configured")
	}

	if len(configuration.NamespacesToIgnore) == 0 {
		configuration.NamespacesToIgnore[0] = "kube-system"
	}

	if configuration.MaxConcurrentReconciles == 0 {
		configuration.MaxConcurrentReconciles = 5
	}

	return &configuration, nil
}