package config

import (
	"encoding/json"
	"os"
	"sync"
)

// Config holds all application configuration
type Config struct {
	mu sync.RWMutex

	// Audio settings
	TestToneFrequency float64 `json:"testToneFrequency"`
	TestToneVolume    float64 `json:"testToneVolume"`
	MicVolume         float64 `json:"micVolume"`

	// Voice Activity Detection settings
	VADThreshold float64 `json:"vadThreshold"`
	VADHoldTime  int64   `json:"vadHoldTime"` // milliseconds

	// Visualization settings
	WaveformWindowSize int     `json:"waveformWindowSize"`
	VolumeUpdateRate   int64   `json:"volumeUpdateRate"` // milliseconds
	PeakDecayRate      float64 `json:"peakDecayRate"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		TestToneFrequency:  440.0, // A4 note
		TestToneVolume:     0.5,   // 50% volume
		MicVolume:          1.0,   // 100% volume
		VADThreshold:       -30.0, // dB
		VADHoldTime:        200,   // 200ms
		WaveformWindowSize: 1024,  // samples
		VolumeUpdateRate:   100,   // 100ms
		PeakDecayRate:      0.95,  // decay factor
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveToFile saves the current configuration to a JSON file
func (c *Config) SaveToFile(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetTestToneFrequency returns the test tone frequency
func (c *Config) GetTestToneFrequency() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TestToneFrequency
}

// SetTestToneFrequency sets the test tone frequency
func (c *Config) SetTestToneFrequency(freq float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.TestToneFrequency = freq
}

// Similar getter/setter methods for other configuration parameters...
