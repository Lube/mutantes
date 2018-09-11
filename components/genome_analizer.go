package components

import (
	"regexp"
	"strings"

	"github.com/lube/mutantes/models"
)

// MinimumSequenceLength of mutant dna
const MinimumSequenceLength = 4

// GenomeAnalizer provides component related to logic analyzing genomes.
type GenomeAnalizer struct {
	mutantChecker *regexp.Regexp
}

// NewGenomeAnalizer creates a new GenomeAnalizer.
func NewGenomeAnalizer() *GenomeAnalizer {
	var mutantChecker = regexp.MustCompile(`[A]{4,}|[C]{4,}|[T]{4,}|[G]{4,}`)
	return &GenomeAnalizer{mutantChecker}
}

// IsMutant Revisar documentación complementaria ./IsMutant.md
func (component *GenomeAnalizer) IsMutant(m *models.Genome) bool {
	sequences := 0
	shutdown := make(chan struct{})
	matches := make(chan bool)

	if len(m.Bases) < 4 {
		return false
	}

	// rows
	go component.checkBases(m, getRows, len(m.Bases), matches, shutdown)

	// columns
	go component.checkBases(m, getColumns, len(m.Bases), matches, shutdown)

	// Para toda matriz de NxN existen al menos
	// N - (longitud minima de diagonal) * 2 diagonales + 1
	qtyDiagonals := (len(m.Bases)-MinimumSequenceLength)*2 + 1
	go component.checkBases(m, getDiagonalsNESO, qtyDiagonals, matches, shutdown)
	go component.checkBases(m, getDiagonalsNOSE, qtyDiagonals, matches, shutdown)

	// Procesamos secuencias de diagonales, filas y columnas
	for i := 0; i < qtyDiagonals*2+len(m.Bases)*2; i++ {
		if <-matches {
			sequences = sequences + 1
		}

		if sequences == 2 {
			close(shutdown)
			return true
		}
	}

	return false
}

func (component *GenomeAnalizer) checkBases(m *models.Genome, getter func(m *models.Genome, i int) string, qtyBases int, matches chan<- bool, shutdown <-chan struct{}) {
	for i := 0; i < qtyBases; i++ {
		matches <- component.mutantChecker.MatchString(getter(m, i))
		select {
		case <-shutdown:
			return
		default:
		}
	}
}

func getDiagonalsNESO(g *models.Genome, i int) string {
	N := len(g.Bases)
	p := MinimumSequenceLength + i - 1

	var diagonal string
	for q := max(0, p-N+1); q <= min(p, N-1); q++ {
		diagonal = diagonal + string(g.Bases[p-q][q])
	}

	return diagonal
}

func getDiagonalsNOSE(g *models.Genome, i int) string {
	N := len(g.Bases)
	p := MinimumSequenceLength + i - 1

	var diagonal string
	for q := max(0, p-N+1); q <= min(p, N-1); q++ {
		diagonal = diagonal + string(g.Bases[p-q][N-1-q])
	}

	return diagonal
}

func getRows(g *models.Genome, i int) string {
	return g.Bases[i]
}

func getColumns(g *models.Genome, i int) string {
	return strings.Join(fmap(g.Bases, func(b string) string {
		return b[i : i+1]
	}), "")
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
