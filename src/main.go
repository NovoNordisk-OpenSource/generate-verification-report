package main

import (
	"context"
	"dagger.io/dagger"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func main() {
	log.Logger = initLogger()
	log.Info().Msg("Logger initialized")

	ctx := context.Background()

	if err := CollectParameters(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}
	if err := VerifyParameters(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}
	// TODO: Port of stages
	if err := GenerateVerificationReport(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}
}

func VerifyParameters(ctx context.Context) error {
	log.Info().Msg("Verifying parameters")

	return nil
}

func GenerateVerificationReport(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Logger))
	if err != nil {
		return err
	}
	defer client.Close()

	log.Info().Msg("Generating verification report")

	hostdir := "output"

	_, err = client.Container().From("alpine:latest").
		WithDirectory("output", client.Directory().WithNewFile("report.html", "This is a test verification report generated from a Dagger workflow")).
		WithWorkdir(".").
		WithExec([]string{"ls", "-la", "output"}).
		WithExec([]string{"cat", "output/report.html"}).
		Directory("output").
		Export(ctx, hostdir)
	if err != nil {
		return err
	}

	return nil
}

func CollectParameters(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Logger))
	if err != nil {
		return err
	}
	defer client.Close()

	log.Info().Msg("Collecting parameters")
	entries, err := client.Host().Directory(".").Entries(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	for _, entry := range entries {
		if entry == "parameters.json" {
			json, err := readJSON(entry)
			if err != nil {
				log.Error().Msg(fmt.Sprintln("Error:", err))
			} else {
				log.Info().Msg(json)
			}
		}
	}

	return nil
}

func initLogger() zerolog.Logger {
	logFile, _ := os.OpenFile(
		"run.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)

	multi := zerolog.New(multiWriter).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return multi
}

// NOTE: For debugging purposes for now
func readJSON(fileName string) (string, error) {
	// Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the content of the file
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Create a map to unmarshal JSON data
	var jsonData map[string]interface{}

	// Unmarshal JSON data into the map
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", err
	}

	// Print the content
	json := fmt.Sprintf("JSON Data:\n%+v\n", jsonData)

	return json, nil
}

type Parameters struct {
	PipelineId  string `json:"pipeline_id"`
	ProjectName string `json:"project_name"`
}
