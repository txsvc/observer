package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testType  ProviderType = 1
	otherType ProviderType = 2
)

type (
	TestProviderImpl struct {
	}
)

var (
	_ GenericProvider = (*TestProviderImpl)(nil)
)

func (tp *TestProviderImpl) Close() error {
	return nil
}

func newTestProvider() interface{} {
	return &TestProviderImpl{}
}

func TestWithProvider(t *testing.T) {
	opt := WithProvider("test", testType, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, testType, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestInitProvider(t *testing.T) {
	p, err := New()

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 0, len(p.providers))
}

func TestNewProvider(t *testing.T) {
	pc1 := WithProvider("test", testType, newTestProvider)
	p, err := New(pc1)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 1, len(p.providers))

	pi2, found := p.Find(testType)
	assert.True(t, found)
	assert.NotNil(t, pi2)

	pi3, found := p.Find(otherType)
	assert.False(t, found)
	assert.Nil(t, pi3)
}

func TestCloseProviders(t *testing.T) {
	pc1 := WithProvider("test", testType, newTestProvider)
	p, err := New(pc1)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 1, len(p.providers))

	assert.False(t, p.Close())
}

func TestDuplicateProvider(t *testing.T) {
	pc1 := WithProvider("test", testType, newTestProvider)
	pc2 := WithProvider("test", testType, newTestProvider)
	p, err := New(pc1, pc2)

	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestIgnoreDuplicateProvider(t *testing.T) {
	pc1 := WithProvider("test", testType, newTestProvider)
	pc2 := WithProvider("test", testType, newTestProvider)
	p, err := New(pc1)

	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, 1, len(p.providers))

	err = p.RegisterProviders(false, pc2)
	assert.Error(t, err)

	err = p.RegisterProviders(true, pc2)
	assert.NoError(t, err)
}
