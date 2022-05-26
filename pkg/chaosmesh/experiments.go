package chaosmesh

import (
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type chaosExperiments struct {
	chaos         map[schema.GroupVersionKind][]*unstructured.Unstructured
	existingChaos map[schema.GroupVersionKind][]types.NamespacedName
}

type ChaosExperimentsBuilder struct {
	experiments *chaosExperiments
}

func newChaosExperimentsBuilder() *ChaosExperimentsBuilder {
	return &ChaosExperimentsBuilder{
		experiments: &chaosExperiments{
			chaos:         map[schema.GroupVersionKind][]*unstructured.Unstructured{},
			existingChaos: map[schema.GroupVersionKind][]types.NamespacedName{},
		},
	}
}

type ChaosExperimentsConfigureFn func(b *ChaosExperimentsBuilder)

func (b *ChaosExperimentsBuilder) WithAWSChaos(c *chaosmeshv1alpha1.AWSChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("AWSChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingAWSChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("AWSChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithDNSChaos(c *chaosmeshv1alpha1.DNSChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("DNSChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingDNSChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("DNSChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithGCPChaos(c *chaosmeshv1alpha1.GCPChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("GCPChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingGCPChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("GCPChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithHTTPChaos(c *chaosmeshv1alpha1.HTTPChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("HTTPChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingHTTPChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("HTTPChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithIOChaos(c *chaosmeshv1alpha1.IOChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("IOChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingIOChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("IOChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithJVMChaos(c *chaosmeshv1alpha1.JVMChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("JVMChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingJVMChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("JVMChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithKernelChaos(c *chaosmeshv1alpha1.KernelChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("KernelChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingKernelChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("KernelChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithNetworkChaos(c *chaosmeshv1alpha1.NetworkChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingNetworkChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithPhysicalMachineChaos(c *chaosmeshv1alpha1.PhysicalMachine) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PhysicalMachine"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingPhysicalMachine(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PhysicalMachine"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithPodChaos(c *chaosmeshv1alpha1.PodChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingPodChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithPodHTTPChaos(c *chaosmeshv1alpha1.PodHttpChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodHttpChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingPodHttpChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodHttpChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithPodIOChaos(c *chaosmeshv1alpha1.PodIOChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodIOChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingPodIOChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodIOChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithPodNetworkChaos(c *chaosmeshv1alpha1.PodNetworkChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodNetworkChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingPodNetworkChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodNetworkChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithStressChaos(c *chaosmeshv1alpha1.StressChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("StressChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingStressChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("StressChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) WithTimeChaos(c *chaosmeshv1alpha1.TimeChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("TimeChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithExistingTimeChaos(namespace string, name string) *ChaosExperimentsBuilder {
	return b.withExistingChaos(chaosmeshv1alpha1.GroupVersion.WithKind("TimeChaos"), namespace, name)
}

func (b *ChaosExperimentsBuilder) withChaos(gvk schema.GroupVersionKind, c interface{}) *ChaosExperimentsBuilder {
	_, ok := b.experiments.chaos[gvk]
	if !ok {
		b.experiments.chaos[gvk] = []*unstructured.Unstructured{}
	}

	obj, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(c)
	uc := &unstructured.Unstructured{Object: obj}
	uc.SetGroupVersionKind(gvk)

	b.experiments.chaos[gvk] = append(b.experiments.chaos[gvk], uc)

	return b
}

func (b *ChaosExperimentsBuilder) withExistingChaos(gvk schema.GroupVersionKind, namespace string, experimentName string) *ChaosExperimentsBuilder {
	_, ok := b.experiments.existingChaos[gvk]
	if !ok {
		b.experiments.existingChaos[gvk] = []types.NamespacedName{}
	}

	b.experiments.existingChaos[gvk] = append(b.experiments.existingChaos[gvk], types.NamespacedName{
		Name:      experimentName,
		Namespace: namespace,
	})

	return b
}

func (b *ChaosExperimentsBuilder) build() *chaosExperiments {
	return b.experiments
}
