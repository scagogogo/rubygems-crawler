package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_GetPackage(t *testing.T) {
	packageInformation, err := NewRepository().GetPackage(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)
	assert.Equal(t, "rails", packageInformation.Name)
}

func TestRepository_Search(t *testing.T) {
	search, err := NewRepository().Search(context.Background(), "rails", 1)
	assert.Nil(t, err)
	assert.True(t, len(search) > 0)
}

func TestRepository_Downloads(t *testing.T) {
	downloadsCount, err := NewRepository().Downloads(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, downloadsCount)
	assert.True(t, downloadsCount.TotalDownloads > 0)
}

func TestRepository_GetGemLatestVersion(t *testing.T) {
	latestVersion, err := NewRepository().GetGemLatestVersion(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, latestVersion)
	assert.NotEmpty(t, latestVersion.Version)
}

func TestRepository_GetGemVersions(t *testing.T) {
	versions, err := NewRepository().GetGemVersions(context.Background(), "rails")
	assert.Nil(t, err)
	assert.True(t, len(versions) > 0)
}

func TestRepository_LatestGems(t *testing.T) {
	gems, err := NewRepository().LatestGems(context.Background())
	assert.Nil(t, err)
	assert.True(t, len(gems) > 0)
}

func TestRepository_VersionDownloads(t *testing.T) {
	downloads, err := NewRepository().VersionDownloads(context.Background(), "rails", "v0.0.1")
	assert.Nil(t, err)
	assert.NotNil(t, downloads)
	assert.True(t, downloads.TotalDownloads > 0)
}
