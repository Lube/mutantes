package models

import (
	stderrors "errors"
	"strings"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/lube/mutantes/errors"
)

// Genome represents a genome record.
type Genome struct {
	Bases []string `json:"dna" db:"dna"`
}

// GetKey for Redis insertion
func (g Genome) GetKey() string {
	return strings.Join(g.Bases, "")
}

// Validate validates the Genome fields.
func (g Genome) Validate() error {
	err := validation.Validate(g.Bases, validation.By(checkNxN))

	if err != nil {
		return errors.InvalidData(map[string]error{"dna": err})
	}

	err = validation.ValidateStruct(&g,
		validation.Field(&g.Bases, validation.Required.Error("is required")))

	if err != nil {
		return errors.InvalidData(map[string]error{"dna": err})
	}

	return nil
}

func checkNxN(value interface{}) error {
	s, _ := value.([]string)
	n := len(s)
	for _, base := range s {
		if len(base) != n {
			return stderrors.New("Genome matrix must be N x N")
		}
	}

	return nil
}

// Stats represents an stat record.
type Stats struct {
	CountMutantDNA int64   `json:"count_mutant_dna"`
	CountHumanDNA  int64   `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

// NewStats generate a Stats Model
func NewStats(countHumans int64, countMutants int64) *Stats {
	var ratio float64

	if countHumans+countMutants == 0 {
		ratio = 0
	} else {
		ratio = float64(countMutants) / (float64(countMutants) + float64(countHumans))
	}

	return &Stats{
		CountHumanDNA:  countHumans,
		CountMutantDNA: countMutants,
		Ratio:          float64(int(ratio*100)) / 100,
	}
}
