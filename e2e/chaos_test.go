package e2e

import "testing"

func TestStructExperiment(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_struct_chaos_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}

func TestYamlExperiment(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_yaml_chaos_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}

func TestFileExperiment(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_file_chaos_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}

func TestStructWorkflow(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_struct_chaos_workflow_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}

func TestFileWorkflow(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_file_chaos_workflow_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}

func TestYamlWorkflow(t *testing.T) {
	given, when, then := newF1ScenarioStage(t)

	given.
		f1_is_configured_to_run_a_scenario_with_a_yaml_chaos_workflow_experiment()

	when.
		the_f1_scenario_is_executed()

	then.
		the_chaos_experiments_are_created().
		and().
		the_f1_scenario_succeeds().
		and().
		the_chaos_experiments_are_cleaned_up()
}
