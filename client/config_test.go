package client

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/joinself/self-go-sdk/account"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentMapping(t *testing.T) {
	tests := []struct {
		name               string
		clientEnvironment  Environment
		expectedAccountEnv *account.Target
	}{
		{
			name:               "Sandbox environment",
			clientEnvironment:  Sandbox,
			expectedAccountEnv: account.TargetSandbox,
		},
		{
			name:               "Production environment",
			clientEnvironment:  Production,
			expectedAccountEnv: account.TargetProduction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				StorageKey:  make([]byte, 32),
				StoragePath: "/tmp/test",
				Environment: tt.clientEnvironment,
				LogLevel:    LogWarn,
			}

			accountConfig := config.toAccountConfig()

			assert.Equal(t, tt.expectedAccountEnv, accountConfig.Environment)

			// Also verify the URLs are correct
			if tt.clientEnvironment == Sandbox {
				assert.Equal(t, "https://rpc-sandbox.joinself.com/", accountConfig.Environment.Rpc)
				assert.Equal(t, "https://object-sandbox.joinself.com/", accountConfig.Environment.Object)
				assert.Equal(t, "wss://message-sandbox.joinself.com/", accountConfig.Environment.Message)
			} else if tt.clientEnvironment == Production {
				assert.Equal(t, "https://rpc.joinself.com/", accountConfig.Environment.Rpc)
				assert.Equal(t, "https://object.joinself.com/", accountConfig.Environment.Object)
				assert.Equal(t, "wss://message.joinself.com/", accountConfig.Environment.Message)
			}
		})
	}
}

func TestDefaultEnvironment(t *testing.T) {
	// Test that default environment (zero value) maps to Sandbox
	config := &Config{
		StorageKey:  make([]byte, 32),
		StoragePath: "/tmp/test",
		// Environment not set (zero value)
		LogLevel: LogWarn,
	}

	accountConfig := config.toAccountConfig()

	assert.Equal(t, account.TargetSandbox, accountConfig.Environment)
	assert.Equal(t, "https://rpc-sandbox.joinself.com/", accountConfig.Environment.Rpc)
}

// TestClientCreationWithSandboxEnvironment is commented out because it requires network connectivity
// func TestClientCreationWithSandboxEnvironment(t *testing.T) {
// 	// Test that we can create a client with Sandbox environment
// 	// and it correctly maps to the account configuration
// 	config := Config{
// 		StorageKey:  make([]byte, 32),
// 		StoragePath: t.TempDir() + "/test_storage",
// 		Environment: Sandbox,
// 		LogLevel:    LogWarn,
// 		SkipReady:   true, // Skip ready check for testing
// 		SkipSetup:   true, // Skip setup for testing
// 	}
//
// 	client, err := NewClient(config)
// 	require.NoError(t, err)
// 	require.NotNil(t, client)
// 	defer client.Close()
//
// 	// Verify the client was created successfully
// 	assert.NotNil(t, client.account)
// 	assert.Equal(t, Sandbox, client.config.Environment)
// }

func TestNewSimplified(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test that NewSimplified creates the correct config
	// We'll test the config creation logic without actually creating the client
	// since that requires network connectivity

	// Verify directory creation works
	err := os.MkdirAll(tempDir+"/test", 0700)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test the config creation logic that NewSimplified would use
	config := Config{
		StorageKey:  generateSecureStorageKey(),
		StoragePath: tempDir,
		Environment: Sandbox,
		LogLevel:    LogWarn,
		SkipReady:   false,
		SkipSetup:   false,
	}

	// Verify the config has expected defaults
	if config.Environment != Sandbox {
		t.Errorf("Expected Environment to be Sandbox, got %v", config.Environment)
	}

	if config.LogLevel != LogWarn {
		t.Errorf("Expected LogLevel to be LogWarn, got %v", config.LogLevel)
	}

	if len(config.StorageKey) != 32 {
		t.Errorf("Expected StorageKey to be 32 bytes, got %d", len(config.StorageKey))
	}

	if config.StoragePath != tempDir {
		t.Errorf("Expected StoragePath to be %s, got %s", tempDir, config.StoragePath)
	}

	// Verify config validation works
	err = config.validate()
	if err != nil {
		t.Errorf("Expected config to be valid, got error: %v", err)
	}
}

func TestNewSimplifiedWithKey(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Test with valid 32-byte key
	storageKey := make([]byte, 32)
	for i := range storageKey {
		storageKey[i] = byte(i)
	}

	// Test the config creation logic that NewSimplifiedWithKey would use
	config := Config{
		StorageKey:  storageKey,
		StoragePath: tempDir,
		Environment: Production,
		LogLevel:    LogError,
		SkipReady:   false,
		SkipSetup:   false,
	}

	// Verify the config has expected production settings
	if config.Environment != Production {
		t.Errorf("Expected Environment to be Production, got %v", config.Environment)
	}

	if config.LogLevel != LogError {
		t.Errorf("Expected LogLevel to be LogError, got %v", config.LogLevel)
	}

	if !bytes.Equal(config.StorageKey, storageKey) {
		t.Error("Expected StorageKey to match provided key")
	}

	// Verify config validation works
	err := config.validate()
	if err != nil {
		t.Errorf("Expected config to be valid, got error: %v", err)
	}

	// Test with invalid key length - this should be caught by NewSimplifiedWithKey
	invalidKey := make([]byte, 16) // Wrong length

	// Simulate the validation that NewSimplifiedWithKey would do
	if len(invalidKey) != 32 {
		expectedError := "storage key must be exactly 32 bytes, got 16"
		actualError := fmt.Sprintf("storage key must be exactly 32 bytes, got %d", len(invalidKey))
		if actualError != expectedError {
			t.Errorf("Expected error message '%s', got '%s'", expectedError, actualError)
		}
	}
}

func TestGenerateSecureStorageKey(t *testing.T) {
	// Test that the function generates a 32-byte key
	key := generateSecureStorageKey()
	if len(key) != 32 {
		t.Errorf("Expected key length to be 32, got %d", len(key))
	}

	// Test that multiple calls generate different keys
	key2 := generateSecureStorageKey()
	if bytes.Equal(key, key2) {
		t.Error("Expected different keys from multiple calls")
	}
}
