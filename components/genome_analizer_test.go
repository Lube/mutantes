package components

import (
	"testing"

	"github.com/lube/mutantes/models"
	"github.com/stretchr/testify/assert"
)

func TestGenomeAnalizer(t *testing.T) {
	ga := NewGenomeAnalizer()

	res := ga.IsMutant(&models.Genome{
		Bases: []string{
			"ATGCGA",
			"AAGTGC",
			"AAAAAA",
			"AGAAGG",
			"CCTCTA",
			"TCACTG",
		},
	})

	assert.Equal(t, res, true)

	res = ga.IsMutant(&models.Genome{
		Bases: []string{
			"ATGCGA",
			"GAGTGC",
			"AAGGAA",
			"AGAAGG",
			"CCTCTA",
			"TCACTG",
		},
	})

	assert.Equal(t, res, false)
}
