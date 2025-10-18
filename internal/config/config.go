package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pn/kymar/internal/db"
)

// SavedConnection represents a saved database connection
type SavedConnection struct {
	Name       string        `json:"name"`
	Params     db.ConnParams `json:"params"`
	IsFavorite bool          `json:"is_favorite"`
}

// Config holds application configuration
type Config struct {
	Connections []SavedConnection `json:"connections"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".kymar")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "connections.json"), nil
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return &Config{Connections: []SavedConnection{}}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty config
			return &Config{Connections: []SavedConnection{}}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// AddConnection adds a new connection to the config
func (c *Config) AddConnection(conn SavedConnection) error {
	// Check if connection with same name exists
	for i, existing := range c.Connections {
		if existing.Name == conn.Name {
			// Update existing connection
			c.Connections[i] = conn
			return c.Save()
		}
	}

	// Add new connection
	c.Connections = append(c.Connections, conn)
	return c.Save()
}

// RemoveConnection removes a connection by name
func (c *Config) RemoveConnection(name string) error {
	for i, conn := range c.Connections {
		if conn.Name == name {
			c.Connections = append(c.Connections[:i], c.Connections[i+1:]...)
			return c.Save()
		}
	}
	return nil
}

// GetConnection retrieves a connection by name
func (c *Config) GetConnection(name string) *SavedConnection {
	for _, conn := range c.Connections {
		if conn.Name == name {
			return &conn
		}
	}
	return nil
}
