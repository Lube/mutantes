package apis

import (
	"net/http"
	"testing"

	"github.com/lube/mutantes/components"
	"github.com/lube/mutantes/daos"
	"github.com/lube/mutantes/services"
	"github.com/lube/mutantes/testUtils"
)

func TestGenome(t *testing.T) {
	testUtils.ResetDB()
	router := newRouter()
	ServeGenomeResource(router, services.NewGenomeService(daos.NewGenomeDAO(), components.NewGenomeAnalizer()))

	notFoundError := `{"error_code":"NOT_FOUND", "message":"NOT_FOUND"}`
	dnaInvalidError := `{"error_code":"INVALID_DATA", "message":"INVALID_DATA", "details" : [{"field":"dna", "error":"Genome matrix must be N x N"}]}`
	dnaRequiredError := `{ "error_code":"INVALID_DATA", "message":"INVALID_DATA", "details" : [{"error":"dna: is required.", "field":"dna"}]}`

	runAPITests(t, router, []apiTestCase{
		{"t1 - try to get a non declared route", "GET", "/genomes/2", "", http.StatusNotFound, notFoundError},
		{"t2 - analyze a human genome", "POST", "/mutant", `{"dna":[
			"ATGCGA",
			"AAGTGC",
			"GGAAAG",
			"AGAAGG",
			"CCTCTA",
			"CCTCTA"
		]}`, http.StatusOK, `false`},
		{"t3 - analyze a mutant genome", "POST", "/mutant", `{"dna":[
			"ATGCGA",
			"AAGTGC",
			"AAAAAG",
			"AGAAGG",
			"CCTCTA",
			"CCTCTA"
		]}`, http.StatusOK, `true`},
		{"t4 - analyze a mutant genome", "POST", "/mutant", `{"dna":[
			"ATGAGA",
			"TAATGC",
			"AAAAAG",
			"AGAAGG",
			"CCTCTA",
			"CCTCTA"
		]}`, http.StatusOK, `true`},
		{"t5 - analyze an genome with validation error", "POST", "/mutant", `{"dna":[
			"ATGCGA"
		]}`, http.StatusBadRequest, dnaInvalidError},
		{"t6 - analyze an genome with validation error", "POST", "/mutant", `{}`, http.StatusBadRequest, dnaRequiredError},
		{"t7 - get stats", "GET", "/stats", ``, http.StatusOK, `{
			"count_mutant_dna": 2,
			"count_human_dna": 1,
			"ratio": 0.66
		}`},
	})
}
