package version

import (
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	
	// Version should not be empty
	if version == "" {
		t.Error("GetVersion() should not return empty string")
	}
	
	// Should contain the version variable value
	if !strings.Contains(version, Version) {
		t.Errorf("GetVersion() should contain Version variable value, got: %s", version)
	}
}

func TestGetBuildInfo(t *testing.T) {
	buildInfo := GetBuildInfo()
	
	// Build info should have valid fields
	if buildInfo.Version == "" {
		t.Error("GetBuildInfo() should have non-empty Version")
	}
	
	// Should contain version information in string representation
	buildInfoStr := buildInfo.String()
	if !strings.Contains(buildInfoStr, "template-generator") {
		t.Errorf("GetBuildInfo().String() should contain 'template-generator', got: %s", buildInfoStr)
	}
	
	// Should contain version
	if !strings.Contains(buildInfoStr, buildInfo.Version) {
		t.Errorf("GetBuildInfo().String() should contain version, got: %s", buildInfoStr)
	}
	
	// Should contain build time
	if !strings.Contains(buildInfoStr, buildInfo.BuildTime) {
		t.Errorf("GetBuildInfo().String() should contain build time, got: %s", buildInfoStr)
	}
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are accessible
	// These might be set by ldflags during build
	
	// Version should be a string (might be empty if not set by ldflags)
	if Version == "" {
		t.Log("Version is empty (normal if not set by ldflags)")
	} else {
		t.Logf("Version: %s", Version)
	}
	
	// BuildTime should be a string (might be empty if not set by ldflags)
	if BuildTime == "" {
		t.Log("BuildTime is empty (normal if not set by ldflags)")
	} else {
		t.Logf("BuildTime: %s", BuildTime)
	}
	
	// GitCommit should be a string (might be empty if not set by ldflags)
	if GitCommit == "" {
		t.Log("GitCommit is empty (normal if not set by ldflags)")
	} else {
		t.Logf("GitCommit: %s", GitCommit)
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that version follows expected format when set
	if Version != "" && Version != "unknown" {
		// Version should not contain spaces or invalid characters
		if strings.Contains(Version, " ") {
			t.Errorf("Version should not contain spaces: %s", Version)
		}
		
		// Version should not be just whitespace
		if strings.TrimSpace(Version) == "" {
			t.Error("Version should not be just whitespace")
		}
	}
}

func TestBuildTimeFormat(t *testing.T) {
	// Test that build time follows expected format when set
	if BuildTime != "" && BuildTime != "unknown" {
		// BuildTime should not be just whitespace
		if strings.TrimSpace(BuildTime) == "" {
			t.Error("BuildTime should not be just whitespace")
		}
		
		// If it looks like a timestamp, it should have reasonable length
		if len(BuildTime) > 0 && len(BuildTime) < 10 {
			t.Errorf("BuildTime seems too short to be a valid timestamp: %s", BuildTime)
		}
	}
}

func TestGetVersionConsistency(t *testing.T) {
	// Test that multiple calls return the same result
	version1 := GetVersion()
	version2 := GetVersion()
	
	if version1 != version2 {
		t.Errorf("GetVersion() should return consistent results, got %s and %s", version1, version2)
	}
}

func TestGetBuildInfoConsistency(t *testing.T) {
	// Test that multiple calls return the same result
	buildInfo1 := GetBuildInfo()
	buildInfo2 := GetBuildInfo()
	
	if buildInfo1 != buildInfo2 {
		t.Errorf("GetBuildInfo() should return consistent results, got %s and %s", buildInfo1, buildInfo2)
	}
}