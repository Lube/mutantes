package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-redis/redis"
	"github.com/lube/mutantes/apis"
	"github.com/lube/mutantes/app"
	"github.com/lube/mutantes/components"
	"github.com/lube/mutantes/daos"
	"github.com/lube/mutantes/errors"
	"github.com/lube/mutantes/services"

	"github.com/go-ozzo/ozzo-routing/file"
)

func main() {
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	logger := logrus.New()

	db := redis.NewClient(&redis.Options{
		Network:  app.Config.Network,
		Addr:     app.Config.DSN,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := db.Ping().Result()
	if err != nil {
		panic(err)
	}

	http.Handle("/", buildRouter(logger, db))

	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func buildRouter(logger *logrus.Logger, db *redis.Client) *routing.Router {
	router := routing.New()

	router.Use(
		app.Init(logger, db),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
	)

	genomeDAO := daos.NewGenomeDAO()
	genomeAnalizer := components.NewGenomeAnalizer()
	apis.ServeGenomeResource(router, services.NewGenomeService(genomeDAO, genomeAnalizer))
	router.Get("/loaderio-6dc123f15883e72366e470f7e85b268e.html", file.Content("loaderio.html"))
	router.Options("/",
		func(c *routing.Context) error {
			return c.Write("OK")
		})

	return router
}
