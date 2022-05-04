package chaosmesh

import (
	"context"
	"fmt"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1/scenarios"
	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"github.com/samuel-form3/f1-chaos-mesh/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type ChaosPlugin struct {
	kubeCli client.Client
	initErr error
}

func NewChaosPlugin(opts ...kubernetes.Option) *ChaosPlugin {
	cp := &ChaosPlugin{}

	kubeConfig := &kubernetes.Config{}
	for _, opt := range opts {
		opt(kubeConfig)
	}

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
			fmt.Println("wtf")
			t.Fatalf("Could not initialize chaos plugin correctly: %s", cp.initErr)
		}

		experimentsBuilder := newChaosExperimentsBuilder()
		cfn(experimentsBuilder)
		experiments := experimentsBuilder.build()

		// Configure Experiments
		err := cp.configureExperiments(t, experiments)
		t.Require.NoError(err)

		// Setup Cleanup Experiments
		t.Cleanup(func() {
			err := cp.cleanupExperiments(t, experiments)
			t.Require.NoError(err)
		})

		return s(t)
	}
}

func (cp *ChaosPlugin) configureExperiments(t *testing.T, exp *ChaosExperiments) error {
	t.Logger.Info("Setting up chaos experiments")
	for _, nc := range exp.networkChaos {
		t.Logger.Infof("Setting up chaos experiment [NetworkChaos]::%s", nc.Name)
		err := cp.kubeCli.Create(context.Background(), nc)
		if err != nil {
			t.Logger.Errorf("Error setting up chaos experiment [NetworkChaos]::%s : %s", nc.Name, err)
			return err
		}
	}

	return nil
}

func (cp *ChaosPlugin) cleanupExperiments(t *testing.T, exp *ChaosExperiments) error {
	t.Logger.Info("Cleaning up chaos experiments")
	for _, nc := range exp.networkChaos {
		t.Logger.Infof("Cleaning up chaos experiment [NetworkChaos]::%s", nc.Name)
		err := cp.kubeCli.Delete(context.Background(), nc)
		if err != nil {
			t.Logger.Errorf("Error setting up chaos experiment [NetworkChaos]::%s : %s", nc.Name, err)
		}
	}

	return nil
}
