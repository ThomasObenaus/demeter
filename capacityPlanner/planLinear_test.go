package capacityPlanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PlanLinear(t *testing.T) {

	cfg := NewDefaultConfig()
	cap := cfg.New()

	assert.Equal(t, uint(10), cap.planLinear(0, 10))

	assert.Equal(t, uint(1), cap.planLinear(1, 0))
	assert.Equal(t, uint(20), cap.planLinear(1, 10))
	assert.Equal(t, uint(2), cap.planLinear(0.5, 1))
	assert.Equal(t, uint(15), cap.planLinear(0.5, 10))

	assert.Equal(t, uint(0), cap.planLinear(-1, 0))
	assert.Equal(t, uint(0), cap.planLinear(-1, 10))
	assert.Equal(t, uint(0), cap.planLinear(-0.5, 1))
	assert.Equal(t, uint(5), cap.planLinear(-0.5, 10))
}