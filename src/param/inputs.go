package param

type Parameters struct {
	PipelineId                                            string `json:"pipeline_id"`
	EnvironmentName                                       string `json:"environment_name"`
	ExtractRequirementsNameToIdMappingPyLocation          string `json:"extract_requirements_name_to_id_mapping_py_location"`
	FeatureFilesPath                                      string `json:"feature_files_path"`
	GetPullRequestIdPyLocation                            string `json:"get_pull_request_id_py_location"`
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
	SystemConfigurationsPath                              string `json:"system_configurations_path"`
	SystemDesignPath                                      string `json:"system_design_path"`
	TestResultsArtifactName                               string `json:"test_results_artifact_name"`
	TestResultsFormat                                     string `json:"test_results_format"`
	VerificationReportTemplateLocation                    string `json:"verification_report_template_location"`
	// TemplateRepo                                          string `json:"template_repo"` // no longer required in the Dagger version
}
