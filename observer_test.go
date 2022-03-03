package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitObserver(t *testing.T) {
	assert.NotNil(t, observerProvider)

	pi, found := Instance().Find(TypeLogger)
	assert.True(t, found)
	assert.NotNil(t, pi)

	pi, found = Instance().Find(TypeErrorReporter)
	assert.True(t, found)
	assert.NotNil(t, pi)

	pi, found = Instance().Find(TypeMetrics)
	assert.True(t, found)
	assert.NotNil(t, pi)
}
