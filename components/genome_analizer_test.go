package components

import (
	"testing"

	"github.com/lube/mutantes/models"
	"github.com/stretchr/testify/assert"
)

func TestGenomeAnalizer(t *testing.T) {
	ga := NewGenomeAnalizer()

	// Secuencias horizontales y verticales
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

	// Sin secuencias validas
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

	// Secuencias oblicuas
	res = ga.IsMutant(&models.Genome{
		Bases: []string{
			"ATGCGG",
			"GATTGC",
			"AAGGGA",
			"AGGGAG",
			"CCGCTA",
			"TGACTG",
		},
	})

	// 3 x 3
	assert.Equal(t, res, true)
	res = ga.IsMutant(&models.Genome{
		Bases: []string{
			"CGG",
			"TGC",
			"GGA",
		},
	})

	assert.Equal(t, res, false)
}
