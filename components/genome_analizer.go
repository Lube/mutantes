package components

import (
	"regexp"
	"strings"

	"github.com/lube/mutantes/models"
)

// GenomeAnalizer provides component related to logic analyzing genomes.
type GenomeAnalizer struct {
	mutantChecker *regexp.Regexp
}

// NewGenomeAnalizer creates a new GenomeAnalizer.
func NewGenomeAnalizer() *GenomeAnalizer {
	var mutantChecker = regexp.MustCompile(`[A]{4,}|[C]{4,}|[T]{4,}|[G]{4,}`)
	return &GenomeAnalizer{mutantChecker}
}

// IsMutant Stats returns the number of genomes.
func (component *GenomeAnalizer) IsMutant(m *models.Genome) bool {
	var columnBases []string
	for i := 0; i < len(m.Bases); i++ {
		columnBases = append(columnBases, strings.Join(fmap(m.Bases, func(b string) string {
			return b[i : i+1]
		}), ""))
	}

	var diagonalBases []string
	if len(m.Bases) >= 4 {
		ultimoIndice := len(m.Bases) - 1
		var diagonalPrincipal string
		var diagonalInvertida string
		// X   0   0   0   Y
		// 0   X   0   Y   0
		// 0   0  X/Y  0   0
		// 0   Y   0   X   0
		// Y   0   0   0   X
		for i := 0; i < len(m.Bases); i++ {
			diagonalPrincipal = diagonalPrincipal + string(m.Bases[i][i])
			diagonalInvertida = diagonalInvertida + string(m.Bases[i][ultimoIndice-i])
		}
		diagonalBases = append(diagonalBases, diagonalPrincipal, diagonalInvertida)

		// Para toda matriz de NxN existen al menos
		// N - (longitud minima de diagonal) * 2 diagonales
		// mas las diagonales principales
		for i := 0; i < (len(m.Bases) - 4); i++ {
			x := 1 + i
			y := 0

			var diagonalNOSETop string
			var diagonalNESOTop string
			var diagonalNOSEBottom string
			var diagonalNESOBottom string
			// j == Cantidad de elementos en una diagonal (N - (i+1) )
			// donde i es la enesima diagonal, o la coordenada x + 1 del
			// primer elemento de la diagonal
			for j := 0; j < len(m.Bases)-(i+1); j++ {
				diagonalNOSETop = diagonalNOSETop + string(m.Bases[y+j][x+j])
				diagonalNESOTop = diagonalNESOTop + string(m.Bases[y+j][ultimoIndice-(x+j)])
				diagonalNOSEBottom = diagonalNOSEBottom + string(m.Bases[x+j][y+j])
				diagonalNESOBottom = diagonalNESOBottom + string(m.Bases[x+j][ultimoIndice-(y+j)])
			}

			diagonalBases = append(diagonalBases, diagonalNOSETop, diagonalNOSEBottom, diagonalNESOTop, diagonalNESOBottom)
		}
	}

	sequences := []string{}
	sequences = append(sequences, m.Bases...)
	sequences = append(sequences, columnBases...)
	sequences = append(sequences, diagonalBases...)

	sequences = filter(sequences, func(b string) bool {
		return component.mutantChecker.MatchString(b)
	})

	return len(sequences) > 1
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func fmap(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
