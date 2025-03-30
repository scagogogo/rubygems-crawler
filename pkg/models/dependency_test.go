package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDependencyInfo_MarshalUnmarshal(t *testing.T) {
	// Create a sample dependency info
	dep := DependencyInfo{
		Name:          "rails",
		DependentName: "activerecord",
		Requirements:  ">= 5.0.0",
		DependentType: "runtime",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(dep)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledDep DependencyInfo
	err = json.Unmarshal(jsonData, &unmarshaledDep)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, dep.Name, unmarshaledDep.Name)
	assert.Equal(t, dep.DependentName, unmarshaledDep.DependentName)
	assert.Equal(t, dep.Requirements, unmarshaledDep.Requirements)
	assert.Equal(t, dep.DependentType, unmarshaledDep.DependentType)
}

func TestDependencyInfo_JsonUnmarshal(t *testing.T) {
	// Sample JSON data
	jsonData := `{
		"name": "rails",
		"dependent_name": "activerecord",
		"requirements": ">= 5.0.0",
		"dependent_type": "runtime"
	}`

	var dep DependencyInfo
	err := json.Unmarshal([]byte(jsonData), &dep)
	assert.NoError(t, err)

	// Verify parsed data
	assert.Equal(t, "rails", dep.Name)
	assert.Equal(t, "activerecord", dep.DependentName)
	assert.Equal(t, ">= 5.0.0", dep.Requirements)
	assert.Equal(t, "runtime", dep.DependentType)
}

func TestDependency_MarshalUnmarshal(t *testing.T) {
	// Create a sample dependency
	dep := Dependency{
		Name:         "rails",
		Requirements: ">= 5.0.0",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(dep)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledDep Dependency
	err = json.Unmarshal(jsonData, &unmarshaledDep)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, dep.Name, unmarshaledDep.Name)
	assert.Equal(t, dep.Requirements, unmarshaledDep.Requirements)
}

func TestDependencies_MarshalUnmarshal(t *testing.T) {
	// Create a sample dependencies struct
	deps := Dependencies{
		Development: []*Dependency{
			{
				Name:         "rspec",
				Requirements: ">= 3.0.0",
			},
		},
		Runtime: []*Dependency{
			{
				Name:         "activesupport",
				Requirements: "= 6.0.0",
			},
			{
				Name:         "activerecord",
				Requirements: "= 6.0.0",
			},
		},
	}

	// Convert to JSON
	jsonData, err := json.Marshal(deps)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Convert back from JSON
	var unmarshaledDeps Dependencies
	err = json.Unmarshal(jsonData, &unmarshaledDeps)
	assert.NoError(t, err)

	// Check if fields match
	assert.Equal(t, len(deps.Development), len(unmarshaledDeps.Development))
	assert.Equal(t, len(deps.Runtime), len(unmarshaledDeps.Runtime))
	assert.Equal(t, deps.Development[0].Name, unmarshaledDeps.Development[0].Name)
	assert.Equal(t, deps.Development[0].Requirements, unmarshaledDeps.Development[0].Requirements)
	assert.Equal(t, deps.Runtime[0].Name, unmarshaledDeps.Runtime[0].Name)
	assert.Equal(t, deps.Runtime[1].Requirements, unmarshaledDeps.Runtime[1].Requirements)
}
