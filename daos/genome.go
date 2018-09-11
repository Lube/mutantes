package daos

import (
	"github.com/lube/mutantes/app"
	"github.com/lube/mutantes/models"
)

// GenomeDAO persists genome data in database
type GenomeDAO struct {
}

// NewGenomeDAO creates a new GenomeDAO
func NewGenomeDAO() *GenomeDAO {
	return &GenomeDAO{}
}

// Insert saves a new genome record in the database.
func (dao *GenomeDAO) Insert(rs app.RequestScope, genome *models.Genome, isMutant bool) error {

	var key string
	if isMutant {
		key = "mutantes"
	} else {
		key = "humanos"
	}
	return rs.DB().SAdd(key, genome.GetKey()).Err()
}

// Stats retrieves the genome records and generates a new Stats model.
func (dao *GenomeDAO) Stats(rs app.RequestScope) (int64, int64, error) {
	countMutants, err := rs.DB().SCard("mutantes").Result()
	if err != nil {
		return 0, 0, err
	}
	countHumans, err := rs.DB().SCard("humanos").Result()
	if err != nil {
		return 0, 0, err
	}
	return countHumans, countMutants, nil
}
