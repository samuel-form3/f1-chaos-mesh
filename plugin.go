package chaosmesh

import (
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1/scenarios"
	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type ChaosPlugin struct {
	kubeCli client.Client
	initErr error
}

func NewChaosPlugin() *ChaosPlugin {
	cp := &ChaosPlugin{}

	cliConfig, err := config.GetConfig()
	if err != nil {
		cp.initErr = err
		return cp
	}

	scheme := runtime.NewScheme()
	cl, err := client.New(cliConfig, client.Options{Scheme: scheme})
	if err != nil {
		cp.initErr = err
		return cp
	}

	err = chaosmeshv1alpha1.AddToScheme(scheme)
	if err != nil {
		cp.initErr = err
		return cp
	}

	cp.kubeCli = cl

	return cp
}

func (cp *ChaosPlugin) WithExperiments(cfn ChaosExperimentsConfigureFn) scenarios.ScenarioOption {
	return func(s *scenarios.Scenario) {

		s.ScenarioFn = cp.wrapScenarioWithExperiments(s.ScenarioFn, cfn)
	}
}

func (cp *ChaosPlugin) wrapScenarioWithExperiments(s testing.ScenarioFn, cfn ChaosExperimentsConfigureFn) testing.ScenarioFn {
	return func(t *testing.T) testing.RunFn {
		if cp.initErr != nil {
			t.Fatalf("Could not initialize chaos plugin correctly: %s", cp.initErr)
		}

		experimentsBuilder := newChaosExperimentsBuilder()
		cfn(experimentsBuilder)
		experiments := experimentsBuilder.build()

		ec := newExperimentsConfigurator(t, cp.kubeCli, experiments)

		t.Cleanup(func() {
			err := ec.CleanupExperiments()
			t.Require.NoError(err)
		})

		// Configure Experiments
		err := ec.ConfigureExperiments()
		if err != nil {
			ec.CleanupExperiments()
			t.FailNow()
		}

		return s(t)
	}
}
