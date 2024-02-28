package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/nn-dma/generate-verification-report/inputs"

	"dagger.io/dagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InputDir  = "input"
	OutputDir = "output"
)

var (
	parameters inputs.Parameters
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
	if err := CreateVerificationReportFilename(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}
	if err := CreateVerificationReportArtifactName(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}
	// TODO: Port of stages
	if err := GenerateVerificationReport(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln("Error:", err))
		panic(err)
	}

	// Collect and verify parameters

	// Collect and verify test results

	// Checkout the repository (or provide a path for it? locally)

	// Preprocess
	// - Run scripts that collect GitHub/ADO information via API
	// - Run scripts that render/generate HTML
	// - Run scripts that generate report filename and artifact name

	// Generate verification report
}

func VerifyParameters(ctx context.Context) error {
	log.Info().Msg("Verifying parameters")

	return nil
}

func CreateVerificationReportFilename(ctx context.Context) error {
	log.Info().Msg("Creating verification report filename")

	// - bash: echo "##vso[task.setvariable variable=verification_report_file]$(${{ parameters.get_verification_report_filename_for_context_sh_location }} "${{ parameters.environment_name }}" "$(Build.BuildId)" "${{ parameters.ready_for }}").html"
	// displayName: Generate verification report filename

	return nil
}

func CreateVerificationReportArtifactName(ctx context.Context) error {
	// - bash: echo "##vso[task.setvariable variable=verification_report_artifact]$(${{ parameters.get_verification_report_artifact_name_for_context_sh_location }} "${{ parameters.ready_for }}")"
	// displayName: Generate verification report artifact name

	// log.Info().Msg("Creating verification report artifact name")

	// // Initialize Dagger client
	// client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Logger))
	// if err != nil {
	// 	return err
	// }
	// defer client.Close()

	// // Execute the bash script and set the output as a variable value
	// output, err := client.Container().From("alpine:latest").
	// 	WithExec([]string{"sh", "-c", fmt.Sprintf("get_verification_report_artifact_name_for_context.sh %s", parameters.ready_for)}).
	// 	Export(ctx, "")
	// if err != nil {
	// 	return err
	// }

	// // Set the output as a variable value on the Dagger context
	// ctx = context.WithValue(ctx, "verification_report_artifact", string(output))

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

	hostdir := OutputDir
	//reportTemplateFile := client.Host().File("template/VerificationReportTemplate.html")

	_, err = client.Container().From("alpine:latest").
		WithDirectory(OutputDir, client.Directory().WithFile("report.html", client.Host().File("template/VerificationReportTemplate.html"))).
		WithWorkdir(".").
		WithExec([]string{"ls", "-la", OutputDir}).
		Directory(OutputDir).
		Export(ctx, hostdir)
	if err != nil {
		return err
	}

	// NOTE: Logging file size is for debugging purposes for now——may be removed in the future unless having it in the logs is useful
	reportTemplateFile := path.Join(OutputDir, "report.html")
	generatedReportFile := client.Host().File(reportTemplateFile)
	size, err := generatedReportFile.Size(ctx)
	if err != nil {
		return err
	}
	log.Info().Msgf("Verification report generated: %s/report.html is %d bytes", OutputDir, size)

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
	entries, err := client.Host().Directory(InputDir).Entries(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	for _, entry := range entries {
		if entry == "parameters.json" {
			entryPath := path.Join(InputDir, entry)
			log.Info().Msg(fmt.Sprintf("Found parameters file: '%s'", entryPath))
			log.Info().Msg(fmt.Sprintf("Reading '%s'", entryPath))
			parameters, err = readParameters(entryPath)
			if err != nil {
				log.Error().Msg(fmt.Sprintln("Error:", err))
			} else {
				log.Info().Msg(fmt.Sprintf("Parsed parameters: %+v", parameters))
			}
		}
	}

	// Check if parameters are valid
	if valid, err := parameters.IsValid(); !valid {
		return err
	}
	log.Info().Msg("Parameters are valid")

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
func readParameters(fileName string) (inputs.Parameters, error) {
	// Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		return inputs.Parameters{}, err
	}
	defer file.Close()

	// Read the content of the file
	data, err := io.ReadAll(file)
	if err != nil {
		return inputs.Parameters{}, err
	}
	log.Info().Msg(fmt.Sprintf("Raw parameters: \n%s", data))

	// Create a map to unmarshal JSON data
	var parameters inputs.Parameters
	err = json.Unmarshal([]byte(data), &parameters)
	if err != nil {
		return inputs.Parameters{}, err
	}

	return parameters, nil
}
