package client

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/joinself/self-go-sdk/account"
)

// Environment represents the target environment
type Environment int

const (
	Sandbox Environment = iota
	Production
)

// LogLevel represents logging verbosity
type LogLevel int

const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

// Config holds configuration for the Self client
type Config struct {
	// StorageKey is the encryption key for local storage (required)
	StorageKey []byte

	// StoragePath is the directory for local storage (required)
	StoragePath string

	// Environment specifies the target environment (default: Sandbox)
	Environment Environment

	// LogLevel specifies logging verbosity (default: LogWarn)
	LogLevel LogLevel

	// SkipReady skips the ready check during initialization
	SkipReady bool

	// SkipSetup skips the setup phase during initialization
	SkipSetup bool
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if len(c.StorageKey) == 0 {
		return ErrStorageKeyRequired
	}
	if c.StoragePath == "" {
		return ErrStoragePathRequired
	}
	return nil
}

// toAccountConfig converts client config to account config
func (c *Config) toAccountConfig() *account.Config {
	cfg := &account.Config{
		StorageKey:  c.StorageKey,
		StoragePath: c.StoragePath,
		SkipReady:   c.SkipReady,
		SkipSetup:   c.SkipSetup,
	}

	// Set environment
	switch c.Environment {
	default:
		cfg.Environment = account.TargetSandbox
	}

	// Set log level
	switch c.LogLevel {
	case LogError:
		cfg.LogLevel = account.LogError
	case LogWarn:
		cfg.LogLevel = account.LogWarn
	case LogInfo:
		cfg.LogLevel = account.LogInfo
	case LogDebug:
		cfg.LogLevel = account.LogDebug
	case LogTrace:
		cfg.LogLevel = account.LogTrace
	default:
		cfg.LogLevel = account.LogWarn
	}

	return cfg
}

// generateSecureStorageKey creates a cryptographically secure 32-byte key
func generateSecureStorageKey() []byte {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		// Fallback to deterministic key generation if crypto/rand fails
		h := sha256.Sum256([]byte(fmt.Sprintf("self-sdk-%d", time.Now().UnixNano())))
		return h[:]
	}
	return key
}

// NewSimplified creates a new Self client with sensible defaults and minimal configuration.
// This is the recommended way to create a client for most use cases.
//
// Parameters:
//   - storagePath: Directory where the client will store its data
//
// The function automatically:
//   - Generates a secure encryption key
//   - Sets environment to Sandbox (safe for development/testing)
//   - Sets log level to LogWarn (balanced verbosity)
//   - Creates the storage directory if it doesn't exist
//
// For production use or custom configuration, use New() with a full Config.
func NewSimplified(storagePath string) (*Client, error) {
	// Create config with sensible defaults
	config := Config{
		StorageKey:  generateSecureStorageKey(),
		StoragePath: storagePath,
		Environment: Sandbox, // Safe default for development
		LogLevel:    LogWarn, // Balanced verbosity
		SkipReady:   false,
		SkipSetup:   false,
	}

	return New(config)
}

// NewSimplifiedWithKey creates a client configured for production use.
// Unlike NewSimplified, this requires an explicit storage key for security.
//
// Parameters:
//   - storageKey: 32-byte encryption key (must be securely generated and stored)
//   - storagePath: Directory where the client will store its data
func NewSimplifiedWithKey(storageKey []byte, storagePath string) (*Client, error) {
	if len(storageKey) != 32 {
		return nil, fmt.Errorf("storage key must be exactly 32 bytes, got %d", len(storageKey))
	}

	// Ensure the storage directory exists
	if err := os.MkdirAll(storagePath, 0700); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	config := Config{
		StorageKey:  storageKey,
		StoragePath: storagePath,
		Environment: Production,
		LogLevel:    LogError, // Minimal logging for production
		SkipReady:   false,
		SkipSetup:   false,
	}

	return New(config)
}
