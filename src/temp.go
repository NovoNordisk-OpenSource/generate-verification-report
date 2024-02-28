package main

import (
	"context"
	"path"

	"dagger.io/dagger"
	"github.com/rs/zerolog/log"
)

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

func PrintCollectedTestResults(ctx context.Context) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(log.Logger))
	if err != nil {
		return err
	}
	defer client.Close()

	log.Info().Msg("Listing collected test results")
	_, err = client.Container().From("alpine:latest").
		WithWorkdir(".").
		WithExec([]string{"ls", "-la", path.Join("input", "testresults")}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}
