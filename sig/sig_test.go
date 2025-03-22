package sig

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewSystem(t *testing.T) {
	system := NewSystem()
	if system == nil {
		t.Fatal("Expected NewSystem to return a non-nil System")
	}
	if system.context == nil {
		t.Error("Expected context to be initialized")
	}
	if system.folder != "" {
		t.Errorf("Expected empty folder, got %s", system.folder)
	}
}

func TestSetContext(t *testing.T) {
	system := NewSystem()

	// Create a context with a value
	type key string
	testKey := key("test-key")
	testValue := "test-value"
	ctx := context.WithValue(context.Background(), testKey, testValue)

	system.SetContext(ctx)

	// Verify the context was set correctly
	if system.context == nil {
		t.Fatal("Context should not be nil")
	}

	// Check if the context value is accessible
	if val := system.context.Value(testKey); val != testValue {
		t.Errorf("Expected context value %v, got %v", testValue, val)
	}
}

func TestSetFolder(t *testing.T) {
	system := NewSystem()

	testFolder := filepath.Join(os.TempDir(), fmt.Sprintf("sig-test-%d", time.Now().UnixNano()))
	defer func() {
		if err := os.RemoveAll(testFolder); err != nil {
			t.Fatalf("Failed to remove temporary directory: %v", err)
		}
	}()
	system.SetFolder(testFolder)

	if system.folder != testFolder {
		t.Errorf("Expected folder to be %s, got %s", testFolder, system.folder)
	}

	// Check if the folder was created
	if _, err := os.Stat(testFolder); os.IsNotExist(err) {
		t.Errorf("Expected folder %s to be created", testFolder)
	}
}

func TestSendLogs_NoFolder(t *testing.T) {
	system := NewSystem()

	// When folder is empty, it should just return nil
	err := system.SendLogs(42)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSendLogs_WithFolder(t *testing.T) {
	system := NewSystem()

	// Create a temporary directory for testing
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("sig-test-%d", time.Now().UnixNano()))
	defer os.RemoveAll(tempDir) // Clean up after test

	system.SetFolder(tempDir)

	// Send logs
	const testCount uint32 = 42
	err := system.SendLogs(testCount)
	if err != nil {
		t.Fatalf("SendLogs failed: %v", err)
	}

	// Check if a file was created in the folder
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %d", len(files))
	}

	// Check if the file has the .prom extension
	if !strings.HasSuffix(files[0].Name(), ".prom") {
		t.Errorf("Expected file with .prom extension, got %s", files[0].Name())
	}

	// Read the file content
	content, err := os.ReadFile(filepath.Join(tempDir, files[0].Name()))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	expectedContent := fmt.Sprintf("consul_open_http_connections %d\n", testCount)
	if string(content) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(content))
	}
}

func TestSendLogs_FileCreationError(t *testing.T) {
	system := NewSystem()

	// Set a folder that doesn't exist and can't be created
	// This assumes /proc exists and is not writable by the test user
	system.folder = "/proc/nonexistent"

	err := system.SendLogs(42)
	if err == nil {
		t.Error("Expected an error when writing to an invalid location, got nil")
	}
}

// Mock for testing file operations failures
type mockFileInfo struct {
	os.FileInfo
}

func (m mockFileInfo) Name() string       { return "mock" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() interface{}   { return nil }
