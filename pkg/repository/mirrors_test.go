package repository

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRepositoryMirrorURLs tests that the mirror repository factory functions
// correctly set up repositories with the appropriate base URLs
func TestRepositoryMirrorURLs(t *testing.T) {
	// Test RubyChina repository URL
	t.Run("RubyChina URL", func(t *testing.T) {
		repo := NewRubyChinaRepository()
		assert.NotNil(t, repo)

		// Type assertion to access the internal options
		repoImpl, ok := repo.(*RepositoryImpl)
		if !ok {
			t.Fatalf("Expected *RepositoryImpl, got %T", repo)
		}
		assert.Equal(t, ServerURLRubyChina, repoImpl.options.ServerURL)
	})

	// Test TSingHua repository URL
	t.Run("TSingHua URL", func(t *testing.T) {
		repo := NewTSingHuaRepository()
		assert.NotNil(t, repo)

		repoImpl, ok := repo.(*RepositoryImpl)
		if !ok {
			t.Fatalf("Expected *RepositoryImpl, got %T", repo)
		}
		assert.Equal(t, ServerURLTSingHua, repoImpl.options.ServerURL)
	})

	// Test AliYun repository URL
	t.Run("AliYun URL", func(t *testing.T) {
		repo := NewAliYunRepository()
		assert.NotNil(t, repo)

		repoImpl, ok := repo.(*RepositoryImpl)
		if !ok {
			t.Fatalf("Expected *RepositoryImpl, got %T", repo)
		}
		assert.Equal(t, ServerURLAliYun, repoImpl.options.ServerURL)
	})
}

// Live API tests are skipped by default as they require network access
// and may fail due to rate limiting, authentication requirements, etc.
// To run these tests, use: go test -v -run TestLiveAPI ./pkg/repository/...

func TestLiveAPI_MirrorFunctionality(t *testing.T) {
	// Skip these tests by default
	t.Skip("Skipping live API tests as they require network access")

	// Test RubyChina mirror
	t.Run("RubyChina GetPackage", func(t *testing.T) {
		repo := NewRubyChinaRepository()
		packageInfo, err := repo.GetPackage(context.Background(), "rails")

		// If the error is about authentication or access, just note it but don't fail the test
		if err != nil && (strings.Contains(err.Error(), "403") ||
			strings.Contains(err.Error(), "401") ||
			strings.Contains(err.Error(), "authentication")) {
			t.Logf("API access issue (likely restricted access): %v", err)
			return
		}

		assert.NoError(t, err)
		if err == nil {
			assert.NotNil(t, packageInfo)
			assert.Equal(t, "rails", packageInfo.Name)
		}
	})

	// Test TSingHua mirror
	t.Run("TSingHua GetPackage", func(t *testing.T) {
		repo := NewTSingHuaRepository()
		packageInfo, err := repo.GetPackage(context.Background(), "rails")

		// If the error is about authentication or access, just note it but don't fail the test
		if err != nil && (strings.Contains(err.Error(), "403") ||
			strings.Contains(err.Error(), "401") ||
			strings.Contains(err.Error(), "authentication")) {
			t.Logf("API access issue (likely restricted access): %v", err)
			return
		}

		assert.NoError(t, err)
		if err == nil {
			assert.NotNil(t, packageInfo)
			assert.Equal(t, "rails", packageInfo.Name)
		}
	})

	// Test AliYun mirror
	t.Run("AliYun GetPackage", func(t *testing.T) {
		repo := NewAliYunRepository()
		packageInfo, err := repo.GetPackage(context.Background(), "rails")

		// If the error is about authentication or access, just note it but don't fail the test
		if err != nil && (strings.Contains(err.Error(), "403") ||
			strings.Contains(err.Error(), "401") ||
			strings.Contains(err.Error(), "authentication")) {
			t.Logf("API access issue (likely restricted access): %v", err)
			return
		}

		assert.NoError(t, err)
		if err == nil {
			assert.NotNil(t, packageInfo)
			assert.Equal(t, "rails", packageInfo.Name)
		}
	})
}
