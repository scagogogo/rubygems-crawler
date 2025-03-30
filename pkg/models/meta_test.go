package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_MarshalUnmarshal(t *testing.T) {
	// Create a sample Metadata
	metadata := Metadata{
		DocumentationURI:    "https://api.rubyonrails.org/v7.0.5/",
		BugTrackerURI:       "https://github.com/rails/rails/issues",
		MailingListURI:      "https://discuss.rubyonrails.org/c/rubyonrails-talk",
		ChangelogURI:        "https://github.com/rails/rails/releases/tag/v7.0.5",
		SourceCodeURI:       "https://github.com/rails/rails/tree/v7.0.5",
		RubygemsMfaRequired: "true",
		WikiURI:             "https://github.com/rails/rails/wiki",
		HomepageURI:         "https://rubyonrails.org",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(metadata)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledMetadata Metadata
	err = json.Unmarshal(jsonData, &unmarshaledMetadata)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, metadata.DocumentationURI, unmarshaledMetadata.DocumentationURI)
	assert.Equal(t, metadata.BugTrackerURI, unmarshaledMetadata.BugTrackerURI)
	assert.Equal(t, metadata.MailingListURI, unmarshaledMetadata.MailingListURI)
	assert.Equal(t, metadata.ChangelogURI, unmarshaledMetadata.ChangelogURI)
	assert.Equal(t, metadata.SourceCodeURI, unmarshaledMetadata.SourceCodeURI)
	assert.Equal(t, metadata.RubygemsMfaRequired, unmarshaledMetadata.RubygemsMfaRequired)
	assert.Equal(t, metadata.WikiURI, unmarshaledMetadata.WikiURI)
	assert.Equal(t, metadata.HomepageURI, unmarshaledMetadata.HomepageURI)
}

func TestMetadata_JsonUnmarshal(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"documentation_uri": "https://api.rubyonrails.org/v7.0.5/",
		"bug_tracker_uri": "https://github.com/rails/rails/issues",
		"mailing_list_uri": "https://discuss.rubyonrails.org/c/rubyonrails-talk",
		"changelog_uri": "https://github.com/rails/rails/releases/tag/v7.0.5",
		"source_code_uri": "https://github.com/rails/rails/tree/v7.0.5",
		"rubygems_mfa_required": "true",
		"wiki_uri": "https://github.com/rails/rails/wiki",
		"homepage_uri": "https://rubyonrails.org"
	}`

	var metadata Metadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, "https://api.rubyonrails.org/v7.0.5/", metadata.DocumentationURI)
	assert.Equal(t, "https://github.com/rails/rails/issues", metadata.BugTrackerURI)
	assert.Equal(t, "https://discuss.rubyonrails.org/c/rubyonrails-talk", metadata.MailingListURI)
	assert.Equal(t, "https://github.com/rails/rails/releases/tag/v7.0.5", metadata.ChangelogURI)
	assert.Equal(t, "https://github.com/rails/rails/tree/v7.0.5", metadata.SourceCodeURI)
	assert.Equal(t, "true", metadata.RubygemsMfaRequired)
	assert.Equal(t, "https://github.com/rails/rails/wiki", metadata.WikiURI)
	assert.Equal(t, "https://rubyonrails.org", metadata.HomepageURI)
}

func TestMetadata_EmptyFields(t *testing.T) {
	// Sample JSON data with empty fields
	jsonData := `{
		"documentation_uri": "https://api.rubyonrails.org/v7.0.5/",
		"bug_tracker_uri": "",
		"mailing_list_uri": null,
		"changelog_uri": "https://github.com/rails/rails/releases/tag/v7.0.5",
		"source_code_uri": "https://github.com/rails/rails/tree/v7.0.5",
		"rubygems_mfa_required": "true"
	}`

	var metadata Metadata
	err := json.Unmarshal([]byte(jsonData), &metadata)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, "https://api.rubyonrails.org/v7.0.5/", metadata.DocumentationURI)
	assert.Equal(t, "", metadata.BugTrackerURI)
	assert.Equal(t, "", metadata.MailingListURI)
	assert.Equal(t, "https://github.com/rails/rails/releases/tag/v7.0.5", metadata.ChangelogURI)
	assert.Equal(t, "true", metadata.RubygemsMfaRequired)
	assert.Equal(t, "", metadata.WikiURI)
	assert.Equal(t, "", metadata.HomepageURI)
}
