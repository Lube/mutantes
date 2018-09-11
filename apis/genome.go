package apis

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/lube/mutantes/app"
	"github.com/lube/mutantes/models"
)

type (
	genomeService interface {
		Analize(rs app.RequestScope, model *models.Genome) (bool, error)
		Stats(rs app.RequestScope) (*models.Stats, error)
	}

	genomeResource struct {
		service genomeService
	}
)

// ServeGenomeResource sets up the routing of endpoints and the corresponding handlers.
func ServeGenomeResource(rg *routing.Router, service genomeService) {
	r := &genomeResource{service}
	rg.Get("/stats", r.stats)
	rg.Post("/mutant", r.analize)
}

func (r *genomeResource) analize(c *routing.Context) error {
	var model models.Genome
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.Analize(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *genomeResource) stats(c *routing.Context) error {
	stats, err := r.service.Stats(app.GetRequestScope(c))
	if err != nil {
		return err
	}

	return c.Write(stats)
}
