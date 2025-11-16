package services

import (
	"math/rand"
	"time"
)

type Randomizer interface {
	Perm(n int) []int
	Intn(n int) int
}

type defaultRandomizer struct {
	r *rand.Rand
}

func NewDefaultRandomizer() Randomizer {
	return &defaultRandomizer{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (d *defaultRandomizer) Perm(n int) []int { return d.r.Perm(n) }
func (d *defaultRandomizer) Intn(n int) int   { return d.r.Intn(n) }
