package services

import (
	"github.com/lube/mutantes/app"
	"github.com/lube/mutantes/errors"
	"github.com/lube/mutantes/models"
)

type genomeDAO interface {
	Insert(rs app.RequestScope, genome *models.Genome, isMutant bool) error
	Stats(rs app.RequestScope) (int64, int64, error)
}

type genomeAnalizer interface {
	IsMutant(genome *models.Genome) bool
}

// GenomeService provides services related with genomes.
type GenomeService struct {
	dao      genomeDAO
	analizer genomeAnalizer
}

// NewGenomeService creates a new GenomeService.
func NewGenomeService(dao genomeDAO, analizer genomeAnalizer) *GenomeService {
	return &GenomeService{dao, analizer}
}

// Analize creates a new genome.
func (s *GenomeService) Analize(rs app.RequestScope, genome *models.Genome) (bool, error) {
	if err := genome.Validate(); err != nil {
		return false, err
	}

	if s.analizer.IsMutant(genome) {
		s.dao.Insert(rs, genome, true)
		return true, nil
	}

	s.dao.Insert(rs, genome, false)
	return false, errors.NotAMutant()
}

// Stats returns the number of genomes.
func (s *GenomeService) Stats(rs app.RequestScope) (*models.Stats, error) {
	countHumans, countMutants, err := s.dao.Stats(rs)

	if err != nil {
		return &models.Stats{}, err
	}

	return models.NewStats(countHumans, countMutants), nil
}
