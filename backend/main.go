package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	commonConfig "github.com/AntonioBR9998/go-common/config"
	"github.com/AntonioBR9998/gymapp/api"
	"github.com/AntonioBR9998/gymapp/internal/config"
	"github.com/AntonioBR9998/gymapp/internal/core"
	"github.com/AntonioBR9998/gymapp/internal/repository"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "0.0.0"
	Commit    = "I'm live!"
	BuildDate = "I don't remember exactly"
)

func main() {
	app := cli.App{
		Name:        "gymapp-backend",
		Usage:       "GymApp backend",
		Description: "GymApp backend service",
		Action:      startGymAppService,
		Version:     Version,
		Before:      BeforeFunc,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Required:  true,
				TakesFile: true,
				Aliases:   []string{"c"},
				Usage:     "load configuration from `FILE`",
				EnvVars:   []string{"GYMAPP_CONFIG_FILE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n%s", err)
	}
}

func startGymAppService(ctx *cli.Context) error {
	log.Infoln("starting GymApp")

	configFilePath := ctx.String("config")
	log.Debugf("the configuration file path is: %s", configFilePath)
	log.Infof("loading configuration from file '%s'", configFilePath)
	cfg := commonConfig.New(
		&config.Config{},
		func(cfg *config.Config) {
			cfg.ConfigPath = configFilePath
		},
	)

	log.Traceln("creating repository layer")
	repository := repository.NewRepository(*cfg)

	log.Traceln("creating service layer")
	service := core.NewCore(repository, *cfg)

	log.Traceln("creating REST API layer")
	s := api.NewAPI(*cfg, service)

	log.Infoln("the user server is on tap now: ", cfg.API.GetURL())
	if cfg.API.TLSEnabled {
		return http.ListenAndServeTLS(cfg.API.GetRelativeURL(), cfg.API.TLSCert, cfg.API.TLSKey, s.Router())
	} else {
		return http.ListenAndServe(cfg.API.GetRelativeURL(), s.Router())
	}
}

func BeforeFunc(ctx *cli.Context) error {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("shutting down GymApp backend service!")
		os.Exit(0)
	}()

	return nil
}
