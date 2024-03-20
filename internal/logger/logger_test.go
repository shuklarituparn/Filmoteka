package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupLogger(t *testing.T) {

	logger := SetupLogger()

	if logger == nil {
		t.Error("Expected logger to be initialized, but got nil")
	}

	expectedPrefix := "SERVER_LOGS: "
	if logger.Prefix() != expectedPrefix {
		t.Errorf("Expected logger prefix to be '%s', but got '%s'", expectedPrefix, logger.Prefix())
	}

	logFilePath := filepath.Join("./logs/", "app.log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Errorf("Expected log file '%s' to be created, but it does not exist", logFilePath)
	}

	logMsg := "Test log message"
	logger.Print(logMsg)

	file, err := os.Open(logFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(buf[:n], []byte(logMsg)) {
		t.Errorf("Expected log file to contain '%s', but it doesn't", logMsg)
	}
}
