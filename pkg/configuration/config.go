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
	var namespacesList []string
	if namespaces != "" {
		namespacesList = strings.Split(namespaces, ",")
	}
	maxReconciles, _:= strconv.Atoi(os.Getenv("MAX_CONCURRENT_RECONCILES"))

	// Sanitize spaces after commas on namespaces.
	for index, ns := range namespacesList{
		namespacesList[index] = strings.TrimSpace(ns)
	}

	configuration := Config{
		BackupRegistry:    os.Getenv("BACKUP_REGISTRY"),
		NamespacesToIgnore: namespacesList,
		MaxConcurrentReconciles: maxReconciles,
	}

	if len(configuration.BackupRegistry) == 0 {
		return nil, fmt.Errorf("no backup registry configured")
	}

	if configuration.NamespacesToIgnore == nil {
		configuration.NamespacesToIgnore = append(configuration.NamespacesToIgnore,"kube-system")
	}

	if configuration.MaxConcurrentReconciles == 0 {
		configuration.MaxConcurrentReconciles = 5
	}

	return &configuration, nil
}