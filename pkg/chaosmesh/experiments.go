package chaosmesh

import (
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
)

type ChaosExperimentsBuilder struct {
	experiments *ChaosExperiments
}

func newChaosExperimentsBuilder() *ChaosExperimentsBuilder {
	return &ChaosExperimentsBuilder{
		experiments: &ChaosExperiments{
			networkChaos: []*chaosmeshv1alpha1.NetworkChaos{},
		},
	}
}

type ChaosExperimentsConfigureFn func(b *ChaosExperimentsBuilder)

func (b *ChaosExperimentsBuilder) WithNetworkChaos(nc *chaosmeshv1alpha1.NetworkChaos) *ChaosExperimentsBuilder {
	b.experiments.networkChaos = append(b.experiments.networkChaos, nc)
	return b
}

func (b *ChaosExperimentsBuilder) build() *ChaosExperiments {
	return b.experiments
}

type ChaosExperiments struct {
	networkChaos []*chaosmeshv1alpha1.NetworkChaos
}
