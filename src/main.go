package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/nn-dma/generate-verification-report/color"
	"github.com/nn-dma/generate-verification-report/inputs"

	"dagger.io/dagger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	InputDir        = "input"
	OutputDir       = "output"
	ScriptDir       = "script"
	ArtifactDir     = "artifact"
	RequirementsDir = "requirements"
	RepositoryDir   = "repository"
)

var (
	parameters inputs.Parameters
	workDir    string
)

func init() {
	log.Logger = initLogger()
	log.Info().Msg("---- NEW RUN, Logger initialized ----")

	currentDir, err := os.Getwd()
	if err != nil {
		log.Error().Msg(fmt.Sprintln(err))
		panic(err)
	}
	workDir = strings.Split(currentDir, "/src")[0]
	log.Info().Msg(fmt.Sprintf("Working directory: '%s'", workDir))
}

func main() {
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
	// TODO: Simplify by moving this to the python container
	log.Info().Msg("Collecting test results")
	collector := client.Container().From("alpine:latest").
		WithDirectory("input/testresults", client.Host().Directory(path.Join(InputDir, "testresults"))).
		WithExec([]string{"sh", "-c", "echo 'number of test results (.json files):' $(ls -1 input/testresults | grep .json | wc -l)"})
	log.Info().Msg("Test results collected")

	// 2. Generate verification report
	log.Info().Msg("Generating verification report")
	// Define local variables and secrets required for the verification report generation
	GITHUB_TOKEN := client.SetSecret("GITHUB_TOKEN", os.Getenv("GITHUB_TOKEN"))

	log.Info().Msg("Preparing state with parameters and test results and outputting debug information")
	generator := client.Container().From("python:3.12.2-bookworm").
		WithEnvVariable("GITHUB_SHA", os.Getenv("GITHUB_SHA")).
		WithEnvVariable("GITHUB_REF_NAME", os.Getenv("GITHUB_REF_NAME")).
		WithEnvVariable("GITHUB_REPOSITORY", os.Getenv("GITHUB_REPOSITORY")).
		WithEnvVariable("GITHUB_RUN_ID", os.Getenv("GITHUB_RUN_ID")).
		WithSecretVariable("GITHUB_TOKEN", GITHUB_TOKEN).
		WithDirectory(ScriptDir, client.Host().Directory(path.Join("src", ScriptDir))).
		WithDirectory(RequirementsDir, client.Host().Directory(parameters.FeatureFilesPath)).
		WithDirectory(RepositoryDir, client.Host().Directory(parameters.ProjectRepositoryPath)).
		WithDirectory("input/testresults", collector.Directory("input/testresults")).
		WithDirectory(OutputDir, client.Directory().WithFile("report.html", client.Host().File("src/template/VerificationReportTemplate.html"))).
		WithExec([]string{"mkdir", ArtifactDir}).
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
	// TODO: Port to GitHub API + write tests
	/*
		echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_id"
		prId=$(python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_id)
		echo $prId
		sed -i "s|<var>PULL_REQUEST_LINK</var>|$(System.CollectionUri)$(System.TeamProject)/_git/$(Build.Repository.Name)/pullrequest/$prId|g" ${{ parameters.verification_report_template_location }}
	*/

	log.Info().Msg("Extracting and rendering pull request closed timestamp")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering pull request closed timestamp'")})
	// TODO: Port to GitHub API + write tests
	/*
		echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_closed_timestamp"
		prClosedTimestamp=$(python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result pull_request_closed_timestamp)
		echo $prClosedTimestamp
		sed -i "s|<var>TIMESTAMP_PIPELINE_START</var>|$prClosedTimestamp|g" ${{ parameters.verification_report_template_location }}
	*/

	log.Info().Msg("Extracting and rendering related work items")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering related work items'")})
	// TODO: Port to GitHub API + write tests
	/*
		echo "python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result work_items > workItemsHtml.html"
		python3 ${{ parameters.get_pull_request_id_py_location }} -commit $COMMIT_HASH -accesstoken USE_ENV_VARIABLE -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) -result work_items > workItemsHtml.html
		cat workItemsHtml.html
		python3 ${{ parameters.render_replace_py_location }} -render ./workItemsHtml.html -template ${{ parameters.verification_report_template_location }} -placeholder "<var>WORK_ITEM_LINKS</var>"
		python3 ${{ parameters.render_replace_py_location }} -render ./workItemsHtml.html -template ${{ parameters.verification_report_template_location }} -placeholder "<kbd><var>CHANGE_ITEM</var></kbd>"
	*/

	// NOTE: This is probably not needed when using Dagger and is a remnant from the more sequential nature of the bash script.
	// log.Info().Msgf("Entering folder '%s' for correct script execution context", RepositoryDir)
	// generator = generator.
	// 	WithWorkdir(".").
	// 	WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Entering folder '$(Build.Repository.Name)' for correct script execution context'")})
	// // TODO: Port
	// /*
	// 	cd $(Build.Repository.Name)
	// */

	log.Info().Msg("Extracting and mapping feature names with unique tags")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and mapping feature names with unique tags'")}).
		WithExec([]string{"ls", "-la"}).
		WithExec([]string{"ls", "-la", ScriptDir}).
		WithExec([]string{"ls", "-la", RequirementsDir}).
		WithExec([]string{"python", parameters.ExtractRequirementsNameToIdMappingPyLocation, "-folder", RequirementsDir}, dagger.ContainerWithExecOpts{RedirectStdout: path.Join(ArtifactDir, "requirementsNameToIdMapping.dict")}).
		WithExec([]string{"ls", "-la", ArtifactDir}).
		WithExec([]string{"cat", path.Join(ArtifactDir, "requirementsNameToIdMapping.dict")})
	/*
		python3 ../${{ parameters.extract_requirements_name_to_id_mapping_py_location }} -folder ${{ parameters.feature_files_path }} > ../requirementsNameToIdMapping.dict
	*/

	log.Info().Msg("Extracting and rendering requirements")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering requirements'")})
	// TODO: Port to GitHub format + write tests
	/*
		python3 ../${{ parameters.render_requirements_py_location }} -folder ${{ parameters.feature_files_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfRequirementsHtml.html
		python3 ../${{ parameters.render_replace_py_location }} -render ./listOfRequirementsHtml.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_REQUIREMENTS</var>"
	*/

	log.Info().Msg("Extracting and rendering design specifications")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering design specifications'")})
	// TODO: Port to GitHub format + write tests
	/*
		python3 ../${{ parameters.render_design_specifications_py_location }} -folder ${{ parameters.system_design_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfDesignSpecifications.html
		python3 ../${{ parameters.render_replace_py_location }} -render ./listOfDesignSpecifications.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_DESIGN_SPECIFICATIONS</var>"
	*/

	log.Info().Msg("Extracting and rendering configuration specifications")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering configuration specifications'")})
	// TODO: Port to GitHub format + write tests
	/*
		python3 ../${{ parameters.render_configuration_specifications_py_location }} -folder ${{ parameters.system_configuration_path }} -branch origin/release/$(Build.SourceBranchName) -organization novonordiskit -project '$(System.TeamProject)' -repository $(Build.Repository.Name) > listOfConfigurationSpecifications.html
		python3 ../${{ parameters.render_replace_py_location }} -render ./listOfConfigurationSpecifications.html -template ../${{ parameters.verification_report_template_location }} -placeholder "<var>LIST_OF_CONFIGURATION_SPECIFICATIONS</var>"
	*/

	// NOTE: This is probably not needed when using Dagger and is a remnant from the more sequential nature of the bash script.
	// log.Info().Msg("Exiting folder '$(Build.Repository.Name)' for correct script execution context")
	// generator = generator.
	// 	WithWorkdir(".").
	// 	WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Exiting folder '$(Build.Repository.Name)' for correct script execution context'")})
	// // TODO: Port (not sure if this is needed)
	// /*
	// 	cd ..
	// */

	log.Info().Msg("Extracting and rendering test results")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Extracting and rendering test results'")}).
		WithExec([]string{"python", parameters.RenderJsonTestResultPyLocation, "-folder", path.Join(InputDir, "testresults"), "-mapping", path.Join(ArtifactDir, "requirementsNameToIdMapping.dict")}, dagger.ContainerWithExecOpts{RedirectStdout: path.Join(ArtifactDir, "renderJsonTestResults.html")}).
		WithExec([]string{"python", parameters.RenderReplacePyLocation, "-render", path.Join(ArtifactDir, "renderJsonTestResults.html"), "-template", "output/report.html", "-placeholder", "<var>TESTCASE_RESULTS</var>"})
	/*
		python3 ${{ parameters.render_json_test_result_py_location }} -folder $(Pipeline.Workspace)/${{ parameters.test_results_artifact_name }} -mapping ./requirementsNameToIdMapping.dict > testResultsHtml.html
		python3 ${{ parameters.render_replace_py_location }} -render ./testResultsHtml.html -template ${{ parameters.verification_report_template_location }} -placeholder "<var>TESTCASE_RESULTS</var>"
	*/

	log.Info().Msg("Rendering IT solution name")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering IT solution name'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>IT_SOLUTION_NAME</var>|" + parameters.ItSolutionName + "|g' output/report.html"})
	/*
		sed -i 's|<var>IT_SOLUTION_NAME</var>|${{ parameters.it_solution_name }}|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Make sure the parameter is set to either ADO or GitHub pipeline/workflow run ID
	log.Info().Msg("Rendering pipeline run ID")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering pipeline run ID'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>PIPELINE_RUN_ID</var>|" + parameters.PipelineRunId + "|g' output/report.html"})
	/*
		sed -i 's|<var>PIPELINE_RUN_ID</var>|$(Build.BuildId)|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Make sure the parameter is set
	log.Info().Msg("Rendering target environment name")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering target environment name'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>ENVIRONMENT</var>|" + parameters.EnvironmentName + "|g' output/report.html"})
	/*
		sed -i 's|<var>ENVIRONMENT</var>|${{ parameters.environment_name }}|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Make sure the parameter is set
	// TODO: Update the placeholder name to be generic (not ADO or GitHub specific)
	log.Info().Msg("Rendering GitHub project name")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering GitHub project name'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>ADO_PROJECT_NAME</var>|" + parameters.ProjectName + "|g' output/report.html"})
	/*
		sed -i 's|<var>ADO_PROJECT_NAME</var>|$(System.TeamProject)|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Write tests
	// TODO: Make sure the parameter is set
	log.Info().Msg("Rendering 'ready for' (production/use) value")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering 'ready for' (production/use) value'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>IS_READY_FOR</var>|" + parameters.ReadyFor + "|g' output/report.html"})
	/*
		sed -i 's|<var>IS_READY_FOR</var>|${{ parameters.ready_for }}|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Write tests
	// TODO: Update the placeholder name to be generic (not ADO or GitHub specific)
	log.Info().Msg("Rendering pipeline run link")
	pipelineRunLink, err := generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering pipeline run link'")}).
		WithExec([]string{"sh", "-c", "echo https://github.com/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID}"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	pipelineRunLink = strings.TrimSpace(pipelineRunLink)
	log.Info().Msgf("Pipeline run link: %s", pipelineRunLink)

	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering pipeline run link'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>ADO_PIPELINE_RUN_LINK</var>|" + pipelineRunLink + "|g' output/report.html"})
	/*
		sed -i 's|<var>ADO_PIPELINE_RUN_LINK</var>|$(System.CollectionUri)$(System.TeamProject)/_build/results?buildId=$(Build.BuildId)\&view=results|g' ${{ parameters.verification_report_template_location }}
	*/

	// NOTE: For GitHub, the pipeline run link is the same as the link to artifacts as these are not different pages (unlike with ADO).
	log.Info().Msg("Rendering pipeline run artifacts link")
	generator = generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Rendering pipeline run artifacts link'")}).
		WithExec([]string{"sh", "-c", "sed -i 's|<var>ARTIFACTS_ADO_PIPELINE_LINK</var>|" + pipelineRunLink + "|g' output/report.html"})
	// TODO: Write tests
	// TODO: Update the placeholder name to be generic (not ADO or GitHub specific)
	/*
		sed -i 's|<var>ARTIFACTS_ADO_PIPELINE_LINK</var>|$(System.CollectionUri)$(System.TeamProject)/_build/results?buildId=$(Build.BuildId)\&view=artifacts\&pathAsName=false\&type=publishedArtifacts|g' ${{ parameters.verification_report_template_location }}
	*/

	// TODO: Write tests
	log.Info().Msg("Generate verification report filename")
	verificationReportFilename, err := generator.
		WithExec([]string{"sh", "-c", "echo '================> " + color.Purple("Generate verification report filename'")}).
		WithExec([]string{parameters.GetVerificationReportFilenameForContextShLocation, parameters.EnvironmentName, parameters.PipelineRunId, parameters.ReadyFor}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	verificationReportFilename = fmt.Sprintf("%s.html", strings.TrimSpace(verificationReportFilename))
	log.Info().Msgf("Verification report filename: %s", verificationReportFilename)

	// 3. Export the verification report to host 'output' directory
	// TODO: Simplify by moving this to the python container
	_, err = client.Container().From("alpine:latest").
		WithFile(fmt.Sprintf("output/%s", verificationReportFilename), generator.File("output/report.html")).
		Directory(OutputDir).
		Export(ctx, OutputDir)
	if err != nil {
		return err
	}

	// NOTE: Logging file size is for debugging purposes for now——may be removed in the future unless having it in the logs is useful
	reportTemplateFile := path.Join(OutputDir, verificationReportFilename)
	generatedReportFile := client.Host().File(reportTemplateFile)
	size, err := generatedReportFile.Size(ctx)
	if err != nil {
		return err
	}
	log.Info().Msgf("Verification report generated: %s/%s is %d bytes", OutputDir, verificationReportFilename, size)

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
			entryPath := path.Join(workDir, InputDir, entry)
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
