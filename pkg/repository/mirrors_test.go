package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRubyChinaRepository(t *testing.T) {
	packageInformation, err := NewRubyChinaRepository().GetPackage(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)

}

func TestNewTSingHuaRepository(t *testing.T) {
	packageInformation, err := NewTSingHuaRepository().GetPackage(context.Background(), "rails")
	assert.Nil(t, err)
	assert.NotNil(t, packageInformation)
}