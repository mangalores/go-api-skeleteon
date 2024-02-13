package utils

import (
	log "github.com/sirupsen/logrus"
	"math"
)

type Delta struct {
	last    float64
	setLast bool
}

func NewDelta() Delta {
	return Delta{}
}

func (d *Delta) Calc(value float64) *float64 {
	if !d.setLast {
		d.last = value
		d.setLast = true

		return nil
	}

	delta := value - d.last
	d.last = value
	if math.IsNaN(delta) {
		log.Warn("not a number!")
	}
	return &delta
}

type DeltaP struct {
	Delta
}

func NewDeltaP() DeltaP {
	return DeltaP{}
}

func (d *DeltaP) Calc(value float64) *float64 {
	last := d.last
	delta := d.Delta.Calc(value)
	if delta == nil || last == 0 {
		return nil
	}

	deltaP := *delta / last
	if math.IsNaN(deltaP) {
		log.Warn("not a number!")
	}
	return &deltaP
}

type SMA struct {
	size   int
	values []float64
	sum    float64
	count  int
}

func NewSMA(size int) SMA {
	return SMA{
		size,
		make([]float64, 0),
		0,
		0,
	}
}

func (m *SMA) add(value float64) {
	if m.count >= m.size {
		v := m.values[0]
		m.values = m.values[1:]
		m.sum -= v
	}

	m.values = append(m.values, value)
	m.sum += value
	m.count = len(m.values)
}

func (m *SMA) calc() float64 {
	if m.count < m.size {
		return 0
	}

	return m.sum / float64(m.count)
}

func (m *SMA) Calc(value float64) *float64 {
	m.add(value)
	if m.count < m.size {
		return nil
	}

	v := m.calc()

	if math.IsNaN(v) {
		log.Warn("not a number!")
	}

	return &v
}

type EMA struct {
	SMA
	last       float64
	multiplier float64
}

func NewEMA(size int, smoothing int) EMA {
	return EMA{
		NewSMA(size),
		0,
		multiplier(smoothing, size),
	}
}

func (m *EMA) Calc(value float64) *float64 {
	m.SMA.add(value)

	var ema float64
	if m.count < m.size {
		return nil
	} else if m.count == m.size {
		ema = m.SMA.calc()
	} else {
		ema = m.multiplier*(value-m.last) + m.last
	}

	m.last = ema

	return &ema
}

func multiplier(smoothing int, size int) float64 {
	return float64(smoothing) / float64(size+1)
}
