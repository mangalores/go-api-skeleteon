package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"testing"
)

func TestDelta(t *testing.T) {
	delta := Delta{}

	var actual *float64

	assert.Nil(t, delta.Calc(1))

	actual = delta.Calc(1.5)
	require.NotNil(t, actual)
	assert.Equal(t, float64(0.5), *actual)

	actual = delta.Calc(1.5)
	require.NotNil(t, actual)
	assert.Equal(t, float64(0), *actual)

	actual = delta.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(-0.5), *actual)
}

func TestDeltaP(t *testing.T) {
	delta := DeltaP{}

	var actual *float64

	assert.Nil(t, delta.Calc(2))

	actual = delta.Calc(2.5)
	require.NotNil(t, actual)
	assert.Equal(t, float64(0.25), *actual)

	actual = delta.Calc(1.25)
	require.NotNil(t, actual)
	assert.Equal(t, float64(-0.5), *actual)

	actual = delta.Calc(1.25)
	require.NotNil(t, actual)
	assert.Equal(t, float64(0), *actual)
}

func TestCalcSMA(t *testing.T) {
	sma := NewSMA(3)

	var actual *float64

	assert.Nil(t, sma.Calc(1))

	assert.Nil(t, sma.Calc(1))

	actual = sma.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(1), *actual)
}

func TestCalcSMASum(t *testing.T) {
	sma := NewSMA(3)

	var actual *float64

	assert.Nil(t, sma.Calc(1))
	assert.Nil(t, sma.Calc(2))

	actual = sma.Calc(3)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)

	actual = sma.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)

	actual = sma.Calc(2)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)

	actual = sma.Calc(3)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)
}

func TestCalcEMA(t *testing.T) {
	ema := NewEMA(3, 2)

	var actual *float64

	assert.Nil(t, ema.Calc(1))
	assert.Nil(t, ema.Calc(1))

	actual = ema.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(1), *actual)

	actual = ema.Calc(4)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)

	actual = ema.Calc(4)
	require.NotNil(t, actual)
	assert.Equal(t, float64(3), *actual)

	actual = ema.Calc(4)
	require.NotNil(t, actual)
	assert.Equal(t, float64(4), *actual)

	actual = ema.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(3), *actual)

	actual = ema.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(2), *actual)

	actual = ema.Calc(1)
	require.NotNil(t, actual)
	assert.Equal(t, float64(1), *actual)
}
