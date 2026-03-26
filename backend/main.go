package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	commonConfig "github.com/AntonioBR9998/go-common/config"
	"github.com/AntonioBR9998/gymapp/internal/config"
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

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{ServerName: cfg.ServerName, InsecureSkipVerify: true}

	log.Traceln("creating repository layer")
	// TODO

	log.Traceln("creating service layer")
	// TODO

	log.Traceln("creating REST API layer")
	// TODO
	// s := server.NewAPI(*cfg, service)

	// log.Infoln("the user server is on tap now: ", cfg.API.GetURL())
	// return http.ListenAndServe(cfg.API.GetRelativeURL(), s.Router())

	return nil // Delete this line when the function is implemented
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
