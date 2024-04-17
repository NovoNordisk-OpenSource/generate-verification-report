package inputs

import (
	"errors"
	"strings"
)

type Parameters struct {
	ProjectRepositoryPath                                 string `json:"project_repository_path"`
	TestResultsPath                                       string `json:"test_results_path"`
	EnvironmentName                                       string `json:"environment_name"`
	ExtractRequirementsNameToIdMappingPyLocation          string `json:"extract_requirements_name_to_id_mapping_py_location"`
	FeatureFilesPath                                      string `json:"feature_files_path"`
	GetPullRequestIdPyLocation                            string `json:"get_pull_request_id_py_location"`
	RenderItChangeIssueLinkGithubPyLocation               string `json:"render_it_change_issue_link_github_py_location"`
	GetPullRequestDetailsForHashGithubShLocation          string `json:"get_pull_request_details_for_hash_github_sh_location"`
	GetPullRequestMergedTimestampGithubShLocation         string `json:"get_pull_request_merged_timestamp_github_sh_location"`
	GetPullRequestUrlGithubShLocation                     string `json:"get_pull_request_url_github_sh_location"`
	GetPullRequestItChangeIssueGithubShLocation           string `json:"get_pull_request_it_change_issue_github_sh_location"`
	GetVerificationReportArtifactNameForContextShLocation string `json:"get_verification_report_artifact_name_for_context_sh_location"`
	GetVerificationReportFilenameForContextShLocation     string `json:"get_verification_report_filename_for_context_sh_location"`
	ItSolutionName                                        string `json:"it_solution_name"`
	PipelineRunId                                         string `json:"pipeline_run_id"`
	ProjectName                                           string `json:"project_name"`
	ReadyFor                                              string `json:"ready_for"`
	RenderConfigurationSpecificationsPyLocation           string `json:"render_configuration_specifications_py_location"`
	RenderDesignSpecificationsPyLocation                  string `json:"render_design_specifications_py_location"`
	RenderJsonTestResultPyLocation                        string `json:"render_json_test_result_py_location"`
	RenderReplacePyLocation                               string `json:"render_replace_py_location"`
	RenderRequirementsPyLocation                          string `json:"render_requirements_py_location"`
	StageName                                             string `json:"stage_name"`
	SystemConfigurationSpecificationPath                  string `json:"system_configuration_specification_path"`
	SystemDesignSpecificationPath                         string `json:"system_design_specification_path"`
	TestResultsArtifactName                               string `json:"test_results_artifact_name"`
	TestResultsFormat                                     string `json:"test_results_format"`
	VerificationReportTemplateLocation                    string `json:"verification_report_template_location"`
	Github                                                Github `json:"github"`
	Azure                                                 Azure  `json:"azure"`
}

type Github struct {
}

type Azure struct {
}

func (p *Parameters) IsValid() (bool, error) {
	// Guard clauses
	if p.PipelineRunId == "" {
		return false, errors.New("pipeline_run_id is required")
	}
	if p.EnvironmentName == "" {
		return false, errors.New("environment_name is required")
	}
	if strings.Contains(p.EnvironmentName, " ") {
		return false, errors.New("environment_name cannot contain spaces")
	}

	return true, nil
}
