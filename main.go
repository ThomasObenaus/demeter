package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/alertmanager"
	"github.com/thomasobenaus/sokar/api"
	"github.com/thomasobenaus/sokar/capacityPlanner"
	"github.com/thomasobenaus/sokar/logging"
	"github.com/thomasobenaus/sokar/nomadConnector"
	"github.com/thomasobenaus/sokar/scaleEventAggregator"
	"github.com/thomasobenaus/sokar/scaler"
	"github.com/thomasobenaus/sokar/sokar"
)

func main() {

	// parse commandline args and consume environment variables
	parsedArgs := parseArgs()
	if !parsedArgs.validateArgs() {
		os.Exit(1)
	}

	jobname := parsedArgs.JobName
	jobMinCount := parsedArgs.JobMinCount
	jobMaxCount := parsedArgs.JobMaxCount
	scaleBy := parsedArgs.ScaleBy
	localPort := 11000

	// set up logging
	lCfg := logging.Config{
		UseStructuredLogging:       parsedArgs.StructuredLogging,
		UseUnixTimestampForLogging: parsedArgs.UseUnixTimestampForLogging,
	}
	loggingFactory := lCfg.New()
	logger := loggingFactory.NewNamedLogger("sokar")

	logger.Info().Msg("Set up the scaler ...")
	scaler, err := setupScaler(jobname, jobMinCount, jobMaxCount, parsedArgs.NomadServerAddr, loggingFactory)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed setting up the scaler")
	}
	logger.Info().Msg("Set up the scaler ... done")

	// OneShot mode
	if parsedArgs.OneShot {
		scaResult := scaler.ScaleBy(scaleBy)
		logger.Info().Msgf("Scale %s: %s", scaResult.State, scaResult.StateDescription)
		os.Exit(0)
	}

	logger.Info().Msg("Connecting components and setting up sokar ...")
	api := api.New(localPort, loggingFactory.NewNamedLogger("sokar.api"))

	var scaleAlertReceivers []scaleEventAggregator.ScaleAlertReceiver
	amCfg := alertmanager.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.amlertmanager"),
	}
	amConnector := amCfg.New()
	api.Router.POST("/alerts", amConnector.HandleScaleAlerts)
	scaleAlertReceivers = append(scaleAlertReceivers, amConnector)

	scaEvtAggCfg := scaleEventAggregator.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.scaEvtAggr"),
	}
	scaEvtAggr := scaEvtAggCfg.New(scaleAlertReceivers)
	api.Router.POST("/alert", scaEvtAggr.ScaleEvent)

	capaCfg := capacityPlanner.Config{
		Logger: loggingFactory.NewNamedLogger("sokar.capaPlanner"),
	}
	capaPlanner := capaCfg.New()

	sokarInst, err := setupSokar(scaEvtAggr, capaPlanner, scaler, api, logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed creating sokar.")
	}

	logger.Info().Msg("Connecting components and setting up sokar ... done")

	// Run all components
	sokarInst.Run()
	scaEvtAggr.Run()
	api.Run()

	// Install signal handler for shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signalChan
		logger.Info().Msgf("Received %v. Shutting down...", s)

		// Stop all components
		api.Stop()
		scaEvtAggr.Stop()
		sokarInst.Stop()
	}()

	// Wait till completion
	api.Join()
	scaEvtAggr.Join()
	sokarInst.Join()

	os.Exit(0)
}

func setupSokar(scaleEventAggregator sokar.ScaleEventAggregator, capacityPlanner sokar.CapacityPlanner, scaler sokar.Scaler, api api.API, logger zerolog.Logger) (*sokar.Sokar, error) {
	cfg := sokar.Config{
		Logger: logger,
	}
	sokarInst, err := cfg.New(scaleEventAggregator, capacityPlanner, scaler)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Registering http handlers ...")
	api.Router.POST(sokar.PathScaleBy, sokarInst.ScaleBy)
	logger.Info().Msg("Registering http handlers ... done")

	return sokarInst, err
}

func setupScaler(jobName string, min uint, max uint, nomadSrvAddr string, logF logging.LoggerFactory) (*scaler.Scaler, error) {

	// Set up the nomad connector
	nomadConnectorConfig := nomadConnector.NewDefaultConfig(nomadSrvAddr)
	nomadConnectorConfig.Logger = logF.NewNamedLogger("sokar.nomad")
	nomadConnector, err := nomadConnectorConfig.New()
	if err != nil {
		return nil, fmt.Errorf("Failed setting up nomad connector: %s.", err)
	}

	scaCfg := scaler.Config{
		JobName:  jobName,
		MinCount: min,
		MaxCount: max,
		Logger:   logF.NewNamedLogger("sokar.scaler"),
	}

	scaler, err := scaCfg.New(nomadConnector)
	if err != nil {
		return nil, fmt.Errorf("Failed setting up scaler: %s.", err)
	}

	return scaler, nil
}
