package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_GetPackage(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	packageInformation, err := NewRepository().GetPackage(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)
	if packageInformation != nil {
		assert.Equal(t, "rails", packageInformation.Name)
	}
}

func TestRepository_Search(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	search, err := NewRepository().Search(context.Background(), "rails", 1)
	assert.Nil(t, err)
	assert.True(t, len(search) > 0)
}

func TestRepository_Downloads(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	downloadsCount, err := NewRepository().Downloads(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, downloadsCount)
	if downloadsCount != nil {
		assert.True(t, downloadsCount.TotalDownloads > 0)
	}
}

func TestRepository_GetGemLatestVersion(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	latestVersion, err := NewRepository().GetGemLatestVersion(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, latestVersion)
	if latestVersion != nil {
		assert.NotEmpty(t, latestVersion.Version)
	}
}

func TestRepository_GetGemVersions(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	versions, err := NewRepository().GetGemVersions(context.Background(), "rails")
	assert.Nil(t, err)
	assert.True(t, len(versions) > 0)
}

func TestRepository_LatestGems(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	gems, err := NewRepository().LatestGems(context.Background())
	assert.Nil(t, err)
	assert.True(t, len(gems) > 0)
}

func TestRepository_VersionDownloads(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	downloads, err := NewRepository().VersionDownloads(context.Background(), "rails", "v0.0.1")
	assert.Nil(t, err)
	assert.NotNil(t, downloads)
	if downloads != nil {
		assert.True(t, downloads.TotalDownloads > 0)
	}
}
