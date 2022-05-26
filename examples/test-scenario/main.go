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
		Add("oneWithChaos", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosExperiments))

	f1Scenarios.Execute()
}

func scenarioOne(t *testing.T) testing.RunFn {
	runFn := func(t *testing.T) {}
	return runFn
}

func scenarioOneChaosExperiments(b *chaosmesh.ChaosExperimentsBuilder) {
	b.
		WithExistingNetworkChaos("default", "scenario-two").
		WithNetworkChaos(&v1alpha1.NetworkChaos{
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
		}).
		WithExistingDNSChaos("chaos-testing", "dns-chaos-example")
}
