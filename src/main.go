package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"dagger.io/dagger"
	"github.com/rs/zerolog"
)

//var Logger zerolog.Logger

func main() {
	//log := initLogger()
	//ctx := context.WithValue(context.Background(), "log", log)
	ctx := context.Background()

	if err := CollectParameters(ctx, log); err != nil {
		panic(err)
	}
	if err := VerifyParameters(ctx, log); err != nil {
		panic(err)
	}
	// TODO: Port of stages
	if err := GenerateVerificationReport(ctx, log); err != nil {
		panic(err)
	}
}

func VerifyParameters(ctx context.Context, log zerolog.Logger) error {
	log.Info().Msg("Verifying parameters")

	return nil
}

func GenerateVerificationReport(ctx context.Context, log zerolog.Logger) error {
	log.Info().Msg("Generating verification report")

	return nil
}

func CollectParameters(ctx context.Context, log zerolog.Logger) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log))
	if err != nil {
		return err
	}
	defer client.Close()

	// // Run commands in container environment
	// out, err := client.
	// 	Container().
	// 	From("alpine:3.17").
	// 	WithExec([]string{"apk", "add", "curl"}).
	// 	WithExec([]string{"ls", "-la"}).
	// 	Stdout(ctx)
	// if err != nil {
	// 	return err
	// }
	// // Print result
	// fmt.Println(out)

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

	log := zerolog.New(multiWriter).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("Logger initialized")
	return log
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
