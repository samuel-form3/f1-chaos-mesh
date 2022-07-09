package e2e

import (
	"context"
	"sync"
	"testing"
	"time"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1"
	f1Testing "github.com/form3tech-oss/f1/pkg/f1/testing"
	chaosmesh "github.com/samuel-form3/f1-chaos-mesh"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type f1ScenariosStage struct {
	t                   *testing.T
	runner              *f1.F1
	chaosPlugin         *chaosmesh.ChaosPlugin
	runErr              error
	expectedExperiments map[schema.GroupVersionKind][]types.NamespacedName
	watchExperimentsWg  *sync.WaitGroup
	watchExperimentsErr error
	k8sClient           client.Client
}

func newF1ScenarioStage(t *testing.T) (given, when, then *f1ScenariosStage) {
	s := &f1ScenariosStage{
		t:                   t,
		runner:              f1.Scenarios(),
		chaosPlugin:         chaosmesh.NewChaosPlugin(),
		expectedExperiments: map[schema.GroupVersionKind][]types.NamespacedName{},
		watchExperimentsWg:  &sync.WaitGroup{},
		k8sClient:           k8sClient,
	}

	s.watchExperimentsWg.Add(1)
	go s.watchExperiments()

	return s, s, s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_struct_chaos_experiment() *f1ScenariosStage {
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithNetworkChaos(&chaosmeshv1alpha1.NetworkChaos{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "scenario-struct",
					Namespace: "kube-system",
				},
				Spec: chaosmeshv1alpha1.NetworkChaosSpec{
					Action: chaosmeshv1alpha1.DelayAction,
					PodSelector: chaosmeshv1alpha1.PodSelector{
						Mode: chaosmeshv1alpha1.AllMode,
						Selector: chaosmeshv1alpha1.PodSelectorSpec{
							GenericSelectorSpec: chaosmeshv1alpha1.GenericSelectorSpec{
								Namespaces: []string{"kube-system"},
								LabelSelectors: map[string]string{
									"k8s-app": "kube-dns",
								},
							},
						},
					},
					TcParameter: chaosmeshv1alpha1.TcParameter{
						Delay: &chaosmeshv1alpha1.DelaySpec{
							Latency: "10ms",
						},
					},
				},
			})
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos")] = []types.NamespacedName{
		{Name: "scenario-struct", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_yaml_chaos_experiment() *f1ScenariosStage {
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithNetworkChaosFromYaml(`
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: scenario-yaml
  namespace: kube-system
spec:
  action: delay
  mode: one
  selector:
    namespaces:
      - kube-system
    labelSelectors:
      "k8s-app": "kube-dns"
  delay:
    latency: '10ms'
    correlation: '100'
    jitter: '0ms'
`)
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos")] = []types.NamespacedName{
		{Name: "scenario-yaml", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_file_chaos_experiment() *f1ScenariosStage {
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithNetworkChaosFromFile("./manifests/scenario-file.yaml")
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos")] = []types.NamespacedName{
		{Name: "scenario-file", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_struct_chaos_workflow_experiment() *f1ScenariosStage {
	deadline := "40s"
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithChaosWorkflow(&chaosmeshv1alpha1.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "workflow-struct",
					Namespace: "kube-system",
				},
				Spec: chaosmeshv1alpha1.WorkflowSpec{
					Entry: "entry",
					Templates: []chaosmeshv1alpha1.Template{
						{
							Name:     "entry",
							Type:     chaosmeshv1alpha1.TypeSerial,
							Deadline: strPtr("5m"),
							Children: []string{
								"kill-coredns",
							},
						},
						{
							Name:     "kill-coredns",
							Type:     chaosmeshv1alpha1.TemplateType(chaosmeshv1alpha1.TypeSchedule),
							Deadline: &deadline,
							Schedule: &chaosmeshv1alpha1.ChaosOnlyScheduleSpec{
								Schedule:          "@every 20s",
								Type:              chaosmeshv1alpha1.ScheduleTypePodChaos,
								ConcurrencyPolicy: chaosmeshv1alpha1.ForbidConcurrent,
								EmbedChaos: chaosmeshv1alpha1.EmbedChaos{
									PodChaos: &chaosmeshv1alpha1.PodChaosSpec{
										Action: chaosmeshv1alpha1.PodKillAction,
										ContainerSelector: chaosmeshv1alpha1.ContainerSelector{
											PodSelector: chaosmeshv1alpha1.PodSelector{
												Mode: chaosmeshv1alpha1.OneMode,
												Selector: chaosmeshv1alpha1.PodSelectorSpec{
													GenericSelectorSpec: chaosmeshv1alpha1.GenericSelectorSpec{
														Namespaces: []string{"kube-system"},
														LabelSelectors: map[string]string{
															"k8s-app": "kube-dns",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			})
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("Workflow")] = []types.NamespacedName{
		{Name: "workflow-struct", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_yaml_chaos_workflow_experiment() *f1ScenariosStage {
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithChaosWorkflowFromYaml(`
apiVersion: chaos-mesh.org/v1alpha1
kind: Workflow
metadata:
  name: workflow-yaml
  namespace: kube-system
spec:
  entry: entry
  templates:
  - name: entry
    templateType: Serial
    deadline: 240s
    children:
      - workflow-pod-chaos-schedule
  - name: workflow-pod-chaos-schedule
    templateType: Schedule
    deadline: 2m
    schedule:
      schedule: '@every 40s'
      concurrencyPolicy: Allow
      type: 'PodChaos'
      podChaos:
        action: pod-kill
        mode: one
        selector:
          namespaces:
            - kube-system
          labelSelectors:
            "k8s-app": "kube-dns"	
`)
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("Workflow")] = []types.NamespacedName{
		{Name: "workflow-yaml", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) f1_is_configured_to_run_a_scenario_with_a_file_chaos_workflow_experiment() *f1ScenariosStage {
	s.runner.Add(
		"exampleWithChaos",
		noopScenario,
		s.chaosPlugin.WithExperiments(func(b *chaosmesh.ChaosExperimentsBuilder) {
			b.WithChaosWorkflowFromFile("./manifests/workflow-file.yaml")
		}))

	s.expectedExperiments[chaosmeshv1alpha1.GroupVersion.WithKind("Workflow")] = []types.NamespacedName{
		{Name: "workflow-file", Namespace: "kube-system"},
	}

	return s
}

func (s *f1ScenariosStage) the_f1_scenario_is_executed() *f1ScenariosStage {
	s.runErr = s.runner.ExecuteWithArgs([]string{
		"run", "constant",
		"--rate", "1/s",
		"--max-duration", "5s",
		"--verbose",
		"exampleWithChaos",
	})

	return s
}

func (s *f1ScenariosStage) the_f1_scenario_succeeds() *f1ScenariosStage {
	require.NoError(s.t, s.runErr, "error executing scenarios")
	return s
}

func (s *f1ScenariosStage) the_chaos_experiments_are_created() *f1ScenariosStage {
	s.watchExperimentsWg.Wait()
	require.NoError(s.t, s.watchExperimentsErr, "error ensuring experiments were created")
	return s
}

func (s *f1ScenariosStage) the_chaos_experiments_are_cleaned_up() *f1ScenariosStage {
	err := wait.PollImmediate(1*time.Second, 20*time.Second, func() (bool, error) {
		for gvk, n := range s.expectedExperiments {
			for _, nn := range n {
				obj := &unstructured.Unstructured{}
				obj.SetGroupVersionKind(gvk)
				err := s.k8sClient.Get(context.Background(), nn, obj)
				if err == nil || !apierrors.IsNotFound(err) {
					return false, nil
				}
			}
		}
		return true, nil
	})

	require.NoError(s.t, err, "error ensuring experiments were cleaned up")
	return s
}

func (s *f1ScenariosStage) and() *f1ScenariosStage {
	return s
}

func (s *f1ScenariosStage) watchExperiments() {
	defer s.watchExperimentsWg.Done()

	s.watchExperimentsErr = wait.PollImmediate(1*time.Second, 30*time.Second, func() (bool, error) {
		for gvk, n := range s.expectedExperiments {
			for _, nn := range n {
				obj := &unstructured.Unstructured{}
				obj.SetGroupVersionKind(gvk)
				err := s.k8sClient.Get(context.Background(), nn, obj)
				if err != nil {
					return false, nil
				}
			}
		}
		return true, nil
	})
}

func noopScenario(t *f1Testing.T) f1Testing.RunFn {
	runFn := func(t *f1Testing.T) {}
	return runFn
}

func strPtr(s string) *string { return &s }
