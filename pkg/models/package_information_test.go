package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPackageInformation_MarshalUnmarshal(t *testing.T) {
	// Create a sample package information
	createdAt, _ := time.Parse(time.RFC3339, "2023-05-24T19:21:28.229Z")
	pkg := PackageInformation{
		Name:             "test-package",
		Downloads:        12345,
		Version:          "1.0.0",
		VersionCreatedAt: createdAt,
		VersionDownloads: 1000,
		Platform:         "ruby",
		Authors:          "Test Author",
		Info:             "Test package information",
		Licenses:         []string{"MIT"},
		Metadata: Metadata{
			DocumentationURI: "https://example.com/docs",
			BugTrackerURI:    "https://example.com/bugs",
		},
		Yanked:           false,
		Sha:              "abcdef1234567890",
		ProjectURI:       "https://rubygems.org/gems/test-package",
		GemURI:           "https://rubygems.org/gems/test-package-1.0.0.gem",
		HomepageURI:      "https://example.com",
		DocumentationURI: "https://example.com/docs",
		MailingListURI:   "https://example.com/mailing-list",
		SourceCodeURI:    "https://github.com/example/test-package",
		BugTrackerURI:    "https://github.com/example/test-package/issues",
		ChangelogURI:     "https://github.com/example/test-package/blob/master/CHANGELOG.md",
		Dependencies: Dependencies{
			Development: []*Dependency{
				{
					Name:         "test-dev-dependency",
					Requirements: ">= 0.1.0",
				},
			},
			Runtime: []*Dependency{
				{
					Name:         "test-runtime-dependency",
					Requirements: "= 2.0.0",
				},
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(pkg)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledPkg PackageInformation
	err = json.Unmarshal(jsonData, &unmarshaledPkg)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, pkg.Name, unmarshaledPkg.Name)
	assert.Equal(t, pkg.Downloads, unmarshaledPkg.Downloads)
	assert.Equal(t, pkg.Version, unmarshaledPkg.Version)
	assert.Equal(t, pkg.VersionCreatedAt.Format(time.RFC3339), unmarshaledPkg.VersionCreatedAt.Format(time.RFC3339))
	assert.Equal(t, pkg.Platform, unmarshaledPkg.Platform)
	assert.Equal(t, pkg.Authors, unmarshaledPkg.Authors)
	assert.Equal(t, pkg.Licenses, unmarshaledPkg.Licenses)
	assert.Equal(t, pkg.Metadata.DocumentationURI, unmarshaledPkg.Metadata.DocumentationURI)
	assert.Equal(t, pkg.Metadata.BugTrackerURI, unmarshaledPkg.Metadata.BugTrackerURI)
	assert.Equal(t, len(pkg.Dependencies.Development), len(unmarshaledPkg.Dependencies.Development))
	assert.Equal(t, len(pkg.Dependencies.Runtime), len(unmarshaledPkg.Dependencies.Runtime))
	assert.Equal(t, pkg.Dependencies.Development[0].Name, unmarshaledPkg.Dependencies.Development[0].Name)
	assert.Equal(t, pkg.Dependencies.Runtime[0].Requirements, unmarshaledPkg.Dependencies.Runtime[0].Requirements)
}

func TestPackageInformation_JsonUnmarshal(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"name": "rails",
		"downloads": 436090160,
		"version": "7.0.5",
		"version_created_at": "2023-05-24T19:21:28.229Z",
		"version_downloads": 54428,
		"platform": "ruby",
		"authors": "David Heinemeier Hansson",
		"info": "Ruby on Rails is a full-stack web framework",
		"licenses": ["MIT"],
		"metadata": {
			"documentation_uri": "https://api.rubyonrails.org/v7.0.5/",
			"bug_tracker_uri": "https://github.com/rails/rails/issues",
			"source_code_uri": "https://github.com/rails/rails/tree/v7.0.5"
		},
		"yanked": false,
		"sha": "57ef2baa4a1f5f954bc6e5a019b1fac8486ece36f79c1cf366e6de33210637fe",
		"project_uri": "https://rubygems.org/gems/rails",
		"gem_uri": "https://rubygems.org/gems/rails-7.0.5.gem",
		"homepage_uri": "https://rubyonrails.org",
		"dependencies": {
			"development": [],
			"runtime": [
				{
					"name": "actioncable",
					"requirements": "= 7.0.5"
				},
				{
					"name": "actionmailbox",
					"requirements": "= 7.0.5"
				}
			]
		}
	}`

	var pkg PackageInformation
	err := json.Unmarshal([]byte(jsonData), &pkg)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, "rails", pkg.Name)
	assert.Equal(t, 436090160, pkg.Downloads)
	assert.Equal(t, "7.0.5", pkg.Version)
	assert.Equal(t, 54428, pkg.VersionDownloads)
	assert.Equal(t, "ruby", pkg.Platform)
	assert.Equal(t, "David Heinemeier Hansson", pkg.Authors)
	assert.Equal(t, []string{"MIT"}, pkg.Licenses)
	assert.Equal(t, "https://api.rubyonrails.org/v7.0.5/", pkg.Metadata.DocumentationURI)
	assert.Len(t, pkg.Dependencies.Runtime, 2)
	assert.Equal(t, "actioncable", pkg.Dependencies.Runtime[0].Name)
	assert.Equal(t, "= 7.0.5", pkg.Dependencies.Runtime[0].Requirements)
}
