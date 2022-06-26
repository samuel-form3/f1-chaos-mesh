package main

import (
	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1"
	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"github.com/samuel-form3/f1-chaos-mesh/pkg/chaosmesh"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	f1Chaos := chaosmesh.NewChaosPlugin()
	f1Scenarios := f1.Scenarios().
		Add("one", scenarioOne).
		Add("oneWithChaos", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosExperiments)).
		Add("oneWithChaosFile", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosFromFile)).
		Add("oneWithChaosYaml", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosFromYaml)).
		Add("oneWithChaosWorkflow", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosWorkflow)).
		Add("oneWithChaosWorkflowFile", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosWorkflowFile)).
		Add("oneWithChaosWorkflowYaml", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosWorkflowYaml))

	f1Scenarios.Execute()
}

func scenarioOne(t *testing.T) testing.RunFn {
	runFn := func(t *testing.T) {}
	return runFn
}

func scenarioOneChaosExperiments(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithNetworkChaos(&v1alpha1.NetworkChaos{
		ObjectMeta: v1.ObjectMeta{
			Name:      "scenario-one",
			Namespace: "default",
		},
		Spec: v1alpha1.NetworkChaosSpec{
			Action: v1alpha1.DelayAction,
			PodSelector: v1alpha1.PodSelector{
				Mode: v1alpha1.AllMode,
				Selector: v1alpha1.PodSelectorSpec{
					GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
						Namespaces: []string{"default"},
					},
				},
			},
			TcParameter: v1alpha1.TcParameter{
				Delay: &v1alpha1.DelaySpec{
					Latency: "100ms",
				},
			},
		},
	})
}

func scenarioOneChaosFromFile(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithNetworkChaosFromFile("./env/existing-networkchaos.yaml")
}

func scenarioOneChaosFromYaml(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithNetworkChaosFromYaml(`
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: scenario-two
  namespace: default
spec:
  action: delay
  mode: one
  selector:
    namespaces:
      - default
    labelSelectors:
      'app': 'web-show'
  delay:
    latency: '10ms'
    correlation: '100'
    jitter: '0ms'
`)
}

func scenarioOneChaosWorkflow(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithChaosWorkflow(&v1alpha1.Workflow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "workflow-one",
			Namespace: "default",
		},
		Spec: v1alpha1.WorkflowSpec{
			Entry: "entry",
			Templates: []v1alpha1.Template{
				{
					Name:     "entry",
					Type:     v1alpha1.TypeSerial,
					Deadline: strPtr("5m"),
					Children: []string{
						"kill-coredns",
					},
				},
				{
					Name: "kill-coredns",
					Type: v1alpha1.TypePodChaos,
					Schedule: &v1alpha1.ChaosOnlyScheduleSpec{
						Schedule:          "@every 15s",
						ConcurrencyPolicy: v1alpha1.ForbidConcurrent,
						EmbedChaos: v1alpha1.EmbedChaos{
							PodChaos: &v1alpha1.PodChaosSpec{
								Action: v1alpha1.PodKillAction,
								ContainerSelector: v1alpha1.ContainerSelector{
									PodSelector: v1alpha1.PodSelector{
										Mode: v1alpha1.AllMode,
										Selector: v1alpha1.PodSelectorSpec{
											GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
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
}

func scenarioOneChaosWorkflowFile(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithChaosWorkflowFromFile("./env/existing-podchaosworkflow.yaml")
}

func scenarioOneChaosWorkflowYaml(b *chaosmesh.ChaosExperimentsBuilder) {
	b.WithChaosWorkflowFromYaml(`
apiVersion: chaos-mesh.org/v1alpha1
kind: Workflow
metadata:
  name: coredns-kill-workflow
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
      schedule: '@every 10s'
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
}

func strPtr(s string) *string { return &s }
