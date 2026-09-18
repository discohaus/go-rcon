package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadRconProperties_Valid(t *testing.T) {
	content := `# Minecraft server properties
server-port=25565
rcon.port=25575
rcon.password=secretpassword123
enable-rcon=true
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "server.properties")
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp properties file: %v", err)
	}

	cfg, err := ReadRconProperties(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port == nil || *cfg.Port != 25575 {
		t.Errorf("expected port 25575, got %v", cfg.Port)
	}

	if cfg.Password == nil || *cfg.Password != "secretpassword123" {
		t.Errorf("expected password 'secretpassword123', got %v", cfg.Password)
	}
}

func TestReadRconProperties_WithSpacesAndColon(t *testing.T) {
	content := `
! Comment with exclamation mark
  rcon.port : 28016  
  rcon.password = my super secret password  
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.properties")
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp properties file: %v", err)
	}

	cfg, err := ReadRconProperties(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port == nil || *cfg.Port != 28016 {
		t.Errorf("expected port 28016, got %v", cfg.Port)
	}

	if cfg.Password == nil || *cfg.Password != "my super secret password" {
		t.Errorf("expected password 'my super secret password', got %v", cfg.Password)
	}
}

func TestReadRconProperties_EmptyOrMissingKeys(t *testing.T) {
	content := `
# Only other properties
server-port=25565
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "server.properties")
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp properties file: %v", err)
	}

	cfg, err := ReadRconProperties(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != nil {
		t.Errorf("expected nil port, got %v", cfg.Port)
	}
	if cfg.Password != nil {
		t.Errorf("expected nil password, got %v", cfg.Password)
	}
}

func TestReadRconProperties_FileNotFound(t *testing.T) {
	_, err := ReadRconProperties("non_existent_file.properties")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestReadRconProperties_InvalidPort(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{"not a number", "rcon.port=abc\n"},
		{"port zero", "rcon.port=0\n"},
		{"port too high", "rcon.port=70000\n"},
		{"negative port", "rcon.port=-5\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "server.properties")
			if err := os.WriteFile(filePath, []byte(tc.content), 0600); err != nil {
				t.Fatalf("failed to write temp properties file: %v", err)
			}

			_, err := ReadRconProperties(filePath)
			if err == nil {
				t.Errorf("expected error for content %q, got nil", tc.content)
			}
		})
	}
}
