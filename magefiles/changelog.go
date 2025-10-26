package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Changelog namespace for changelog-related commands
type Changelog mg.Namespace

// Generate generates a new changelog entry from git commits
func (Changelog) Generate() error {
	fmt.Println("Generating changelog...")
	
	// Get the latest tag
	latestTag, err := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		fmt.Println("No tags found, using all commits")
		latestTag = "" // If no tags exist, get all commits
	}
	
	// Get commits since last tag
	var commits string
	if latestTag != "" {
		commits, err = sh.Output("git", "log", latestTag+"..HEAD", "--pretty=format:- %s (%h)")
	} else {
		commits, err = sh.Output("git", "log", "--pretty=format:- %s (%h)")
	}
	
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}
	
	if commits == "" {
		fmt.Println("No new commits since last tag")
		return nil
	}
	
	// Create changelog entry
	date := time.Now().Format("2006-01-02")
	entry := fmt.Sprintf("\n## [Unreleased] - %s\n\n### Added\n%s\n", date, commits)
	
	fmt.Printf("Generated changelog entry:\n%s", entry)
	return nil
}

// Update updates the changelog with recent commits
func (Changelog) Update() error {
	fmt.Println("Updating changelog...")
	
	// Get commits since last changelog update (last week)
	commits, err := sh.Output("git", "log", "--since='1 week ago'", "--pretty=format:- %s (%h)")
	if err != nil {
		return fmt.Errorf("failed to get recent commits: %w", err)
	}
	
	if commits == "" {
		fmt.Println("No recent commits to add to changelog")
		return nil
	}
	
	// Read existing changelog
	content, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}
	
	// Parse and update changelog
	lines := strings.Split(string(content), "\n")
	var newLines []string
	unreleasedFound := false
	addedFound := false
	
	for i, line := range lines {
		newLines = append(newLines, line)
		
		// Find the [Unreleased] section
		if strings.Contains(line, "## [Unreleased]") {
			unreleasedFound = true
			continue
		}
		
		// Find the ### Added section under Unreleased
		if unreleasedFound && strings.Contains(line, "### Added") {
			addedFound = true
			// Add new commits after the ### Added line
			commitLines := strings.Split(commits, "\n")
			for _, commit := range commitLines {
				if strings.TrimSpace(commit) != "" {
					newLines = append(newLines, commit)
				}
			}
			addedFound = false // Reset flag
			continue
		}
	}
	
	// Write updated changelog
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile("CHANGELOG.md", []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}
	
	fmt.Println("Changelog updated successfully")
	return nil
}

// Validate validates the changelog format
func (Changelog) Validate() error {
	fmt.Println("Validating changelog format...")
	
	content, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}
	
	// Basic validation - check for required sections
	requiredSections := []string{
		"# Changelog",
		"## [Unreleased]",
		"### Added",
		"### Changed",
		"### Fixed",
	}
	
	contentStr := string(content)
	for _, section := range requiredSections {
		if !strings.Contains(contentStr, section) {
			return fmt.Errorf("changelog missing required section: %s", section)
		}
	}
	
	fmt.Println("Changelog format is valid")
	return nil
}

// Release prepares a release by updating version in changelog
func (c Changelog) Release() error {
	// Get version from command line args or environment
	version := os.Getenv("VERSION")
	if version == "" {
		return fmt.Errorf("VERSION environment variable is required")
	}
	
	fmt.Printf("Preparing release %s...\n", version)
	
	// Read existing changelog
	content, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
	}
	
	// Replace [Unreleased] with version and date
	date := time.Now().Format("2006-01-02")
	oldContent := string(content)
	newContent := strings.Replace(oldContent, "## [Unreleased]", fmt.Sprintf("## [%s] - %s", version, date), 1)
	
	// Add new [Unreleased] section at the top
	unreleasedSection := `
## [Unreleased]

### Added
- 

### Changed
- 

### Deprecated
- 

### Removed
- 

### Fixed
- 

### Security
- 

`
	
	// Insert new unreleased section after the first version
	lines := strings.Split(newContent, "\n")
	var finalLines []string
	inserted := false
	
	for _, line := range lines {
		finalLines = append(finalLines, line)
		
		// Insert after the newly created version header
		if !inserted && strings.Contains(line, fmt.Sprintf("## [%s] - %s", version, date)) {
			finalLines = append(finalLines, strings.Split(unreleasedSection, "\n")...)
			inserted = true
		}
	}
	
	// Write updated changelog
	finalContent := strings.Join(finalLines, "\n")
	if err := os.WriteFile("CHANGELOG.md", []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
	}
	
	fmt.Printf("Release %s prepared successfully\n", version)
	return nil
}