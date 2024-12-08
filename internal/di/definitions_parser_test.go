package di_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/strider2038/digen/internal/di"
)

var (
	//go:embed testdata/parsing/factories.txt
	testFactorySource string
)

func TestParseFactoryFromSource(t *testing.T) {
	factory, err := di.ParseFactoriesFromSource(testFactorySource)

	require.NoError(t, err)
	require.NotNil(t, factory)
	assert.NotNil(t, factory.Imports["usecase"])
	assert.NotNil(t, factory.Imports["domain"])
	assert.NotNil(t, factory.Imports["httpadapter"])
	assert.NotNil(t, factory.Imports["inmemory"])
	assert.Contains(t, factory.Factories, "EntityRepository")
	assert.Contains(t, factory.Factories, "UseCase")
	assert.Contains(t, factory.Factories, "Handler")
}
