package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/nn-dma/generate-verification-report/color"
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
	log.Info().Msg("---- NEW RUN ----")
	log.Info().Msg("Logger initialized")

	ctx := context.Background()

	// Collect and verify parameters
	if err := CollectParameters(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln(err))
		panic(err)
	}
	if err := VerifyParameters(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln(err))
		panic(err)
	}
	// Generate verification report
	// -Collect test results
	// -TODO: Verify test results (out of scope for now)
	// -Checkout the repository (or provide a path for it? locally)
	// - Preprocess
	// -- Run scripts that collect GitHub/ADO information via API
	// -- Run scripts that render/generate HTML
	// -- Run scripts that generate report filename and artifact name
	if err := GenerateVerificationReport(ctx); err != nil {
		log.Error().Msg(fmt.Sprintln(err))
		panic(err)
	}
}

func VerifyParameters(ctx context.Context) error {
	log.Info().Msg("Verifying parameters")
	if valid, err := parameters.IsValid(); !valid {
		return err
	}
	log.Info().Msg("Parameters are valid")

	return nil
}

func GenerateVerificationReport(ctx context.Context) error {
	// TODO: Port of stages

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Logger))
	if err != nil {
		return err
	}
	defer client.Close()

	// 1. Collect test results
	log.Info().Msg("Collecting test results")
	collector := client.Container().From("alpine:latest").
		WithWorkdir(".").
		WithDirectory("input/testresults", client.Host().Directory(path.Join(InputDir, "testresults"))).
		WithExec([]string{"sh", "-c", "echo 'number of test results (.json files):' $(ls -1 input/testresults | grep .json | wc -l)"})
	if err != nil {
		return err
	}
	log.Info().Msg("Test results collected")

	// 2. Generate verification report
	log.Info().Msg("Generating verification report")
	hostOutputDir := OutputDir

	log.Info().Msg("Preparing state with parameters and test results and outputting debug information")
	generator := client.Container().From("python:3.12.2-bookworm").
		WithEnvVariable("GITHUB_SHA", os.Getenv("GITHUB_SHA")).
		WithEnvVariable("GITHUB_REF_NAME", os.Getenv("GITHUB_REF_NAME")).
		WithDirectory("input/testresults", collector.Directory("input/testresults")).
		WithDirectory(OutputDir, client.Directory().WithFile("report.html", client.Host().File("template/VerificationReportTemplate.html"))).
		WithWorkdir(".").
		WithExec([]string{"ls", "-la", OutputDir}).
		WithExec([]string{"python", "--version"}).
		WithExec([]string{"sh", "-c", "git version"}).
		WithExec([]string{"sh", "-c", "echo current directory: $(pwd)"}).
		WithExec([]string{"sh", "-c", "echo branch: $(git branch --show-current)"}).
		WithExec([]string{"sh", "-c", "echo triggering commit hash: ${GITHUB_SHA}"}).
		WithExec([]string{"sh", "-c", "echo triggering branch: ${GITHUB_REF_NAME}"})

	log.Info().Msg("Extracting and rendering pull request links")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering pull request links'")})
		// TODO: Port to GitHub + write tests
		/*
			echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_id"
			prId=$(python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_id)
			echo $prId
			sed -i "s|<var>PULL_REQUEST_LINK</var>|$(System.CollectionUri)$(System.TeamProject)/_git/$(Build.Repository.Name)/pullrequest/$prId|g" ${{ parameters.verification_report_template_location }}
		*/

	log.Info().Msg("Extracting and rendering pull request closed timestamp")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering pull request closed timestamp'")})
		// TODO: Port to GitHub + write tests
		/*
			echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_closed_timestamp"
			prClosedTimestamp=$(python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_closed_timestamp)
			echo $prClosedTimestamp
			sed -i "s|<var>TIMESTAMP_PIPELINE_START</var>|$prClosedTimestamp|g" ${{ parameters.verification_report_template_location }}
		*/

	log.Info().Msg("Extracting and rendering related work items")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering related work items'")})
		// TODO: Port + write tests
		/*
			echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result work_items > workItemsHtml.html"
			python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result work_items > workItemsHtml.html
			cat workItemsHtml.html
			python3 ${{ parameters.render_replace_py_location }} -render ./workItemsHtml.html -template ${{ parameters.verification_report_template_location }} -placeholder "<var>WORK_ITEM_LINKS</var>"
			python3 ${{ parameters.render_replace_py_location }} -render ./workItemsHtml.html -template ${{ parameters.verification_report_template_location }} -placeholder "<kbd><var>CHANGE_ITEM</var></kbd>"
		*/

	log.Info().Msg("Entering folder '$(Build.Repository.Name)' for correct script execution context")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Entering folder '$(Build.Repository.Name)' for correct script execution context'")})
		// TODO: Port
		/*
			cd $(Build.Repository.Name)
		*/

	log.Info().Msg("Extracting and mapping feature names with unique tags")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and mapping feature names with unique tags'")})
		// TODO: Port + write tests
		/*
			python3 ../${{ parameters.extract_requirements_name_to_id_mapping_py_location }} -folder ${{ parameters.feature_files_path }} > ../requirementsNameToIdMapping.dict
		*/

	log.Info().Msg("Extracting and rendering requirements")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering requirements'")})
		// TODO: Port + write tests
		/*
			python3 ../${{ parameters.render_requirements_py_location }} -folder ${{ parameters.feature_files_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfRequirementsHtml.html
			python3 ../${{ parameters.render_replace_py_location }} -render ./listOfRequirementsHtml.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_REQUIREMENTS</var>"
		*/

	log.Info().Msg("Extracting and rendering design specifications")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering design specifications'")})
		// TODO: Port + write tests
		/*
			python3 ../${{ parameters.render_design_specifications_py_location }} -folder ${{ parameters.system_design_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfDesignSpecifications.html
			python3 ../${{ parameters.render_replace_py_location }} -render ./listOfDesignSpecifications.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_DESIGN_SPECIFICATIONS</var>"
		*/

	log.Info().Msg("Extracting and rendering configuration specifications")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering configuration specifications'")})
		// TODO: Port + write tests
		/*
			python3 ../${{ parameters.render_configuration_specifications_py_location }} -folder ${{ parameters.system_configuration_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfConfigurationSpecifications.html
			python3 ../${{ parameters.render_replace_py_location }} -render ./listOfConfigurationSpecifications.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_CONFIGURATION_SPECIFICATIONS</var>"
		*/

	log.Info().Msg("Exiting folder '$(Build.Repository.Name)' for correct script execution context")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Exiting folder '$(Build.Repository.Name)' for correct script execution context'")})
		// TODO: Port
		/*
			cd ..
		*/

	// 3. Export the verification report to host 'output' directory
	_, err = client.Container().From("alpine:latest").
		WithFile("output/report.html", generator.File("output/report.html")).
		Directory(OutputDir).
		Export(ctx, hostOutputDir)
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
	found := false
	for _, entry := range entries {
		if entry == "parameters.json" {
			found = true
			entryPath := path.Join(InputDir, entry)
			log.Info().Msg(fmt.Sprintf("Found parameters file: '%s'", entryPath))
			log.Info().Msg(fmt.Sprintf("Reading '%s'", entryPath))
			parameters, err = readParameters(entryPath) // Set the global parameters variable
			if err != nil {
				log.Error().Msg(fmt.Sprintln(err))
			} else {
				log.Info().Msg(fmt.Sprintf("Parsed parameters: %#v", parameters))
			}
		}
	}
	if !found {
		return fmt.Errorf("expected file 'parameters.json' not found in directory '%s'", InputDir)
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
	// NOTE: For debugging purposes for now
	log.Info().Msg(fmt.Sprintf("Raw parameters: \n%s", data))

	// Create a map to unmarshal JSON data
	var parameters inputs.Parameters
	err = json.Unmarshal([]byte(data), &parameters)
	if err != nil {
		return inputs.Parameters{}, err
	}

	return parameters, nil
}
