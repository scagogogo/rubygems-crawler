package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryDownloadCount_MarshalUnmarshal(t *testing.T) {
	// Create a sample RepositoryDownloadCount
	downloadCount := RepositoryDownloadCount{
		TotalDownloads: 1000000,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(downloadCount)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledDownloadCount RepositoryDownloadCount
	err = json.Unmarshal(jsonData, &unmarshaledDownloadCount)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, downloadCount.TotalDownloads, unmarshaledDownloadCount.TotalDownloads)
}

func TestRepositoryDownloadCount_JsonUnmarshal(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"total": 436090160
	}`

	var downloadCount RepositoryDownloadCount
	err := json.Unmarshal([]byte(jsonData), &downloadCount)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, 436090160, downloadCount.TotalDownloads)
}

func TestVersionDownloadCount_MarshalUnmarshal(t *testing.T) {
	// Create a sample VersionDownloadCount
	versionDownloadCount := VersionDownloadCount{
		VersionDownloads: 54428,
		TotalDownloads:   436090160,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(versionDownloadCount)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledVersionDownloadCount VersionDownloadCount
	err = json.Unmarshal(jsonData, &unmarshaledVersionDownloadCount)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, versionDownloadCount.VersionDownloads, unmarshaledVersionDownloadCount.VersionDownloads)
	assert.Equal(t, versionDownloadCount.TotalDownloads, unmarshaledVersionDownloadCount.TotalDownloads)
}

func TestVersionDownloadCount_JsonUnmarshal(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"version_downloads": 54428,
		"total_downloads": 436090160
	}`

	var versionDownloadCount VersionDownloadCount
	err := json.Unmarshal([]byte(jsonData), &versionDownloadCount)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, 54428, versionDownloadCount.VersionDownloads)
	assert.Equal(t, 436090160, versionDownloadCount.TotalDownloads)
}
