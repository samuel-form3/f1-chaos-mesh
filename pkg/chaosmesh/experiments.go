package chaosmesh

import (
	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type chaosExperiments struct {
	chaos                   map[schema.GroupVersionKind][]*unstructured.Unstructured
	chaosFromFiles          map[schema.GroupVersionKind][]string
	chaosFromYaml           map[schema.GroupVersionKind][]string
	chaosWorkflows          []*chaosmeshv1alpha1.Workflow
	chaosWorkflowsFromFiles []string
	chaosWorkflowsFromYaml  []string
}

type ChaosExperimentsBuilder struct {
	experiments *chaosExperiments
}

func newChaosExperimentsBuilder() *ChaosExperimentsBuilder {
	return &ChaosExperimentsBuilder{
		experiments: &chaosExperiments{
			chaos:                   map[schema.GroupVersionKind][]*unstructured.Unstructured{},
			chaosFromFiles:          map[schema.GroupVersionKind][]string{},
			chaosFromYaml:           map[schema.GroupVersionKind][]string{},
			chaosWorkflows:          []*chaosmeshv1alpha1.Workflow{},
			chaosWorkflowsFromFiles: []string{},
			chaosWorkflowsFromYaml:  []string{},
		},
	}
}

type ChaosExperimentsConfigureFn func(b *ChaosExperimentsBuilder)

// AWS CHAOS

func (b *ChaosExperimentsBuilder) WithAWSChaos(c *chaosmeshv1alpha1.AWSChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("AWSChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithAWSChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("AWSChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithAWSChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("AWSChaos"), yaml)
}

// DNS Chaos

func (b *ChaosExperimentsBuilder) WithDNSChaos(c *chaosmeshv1alpha1.DNSChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("DNSChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithDNSChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("DNSChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithDNSChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("DNSChaos"), yaml)
}

// GCP Chaos

func (b *ChaosExperimentsBuilder) WithGCPChaos(c *chaosmeshv1alpha1.GCPChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("GCPChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithGCPChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("GCPChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithGCPChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("GCPChaos"), yaml)
}

// HTTP Chaos

func (b *ChaosExperimentsBuilder) WithHTTPChaos(c *chaosmeshv1alpha1.HTTPChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("HTTPChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithHTTPChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("HTTPChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithHTTPChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("HTTPChaos"), yaml)
}

// IO Chaos

func (b *ChaosExperimentsBuilder) WithIOChaos(c *chaosmeshv1alpha1.IOChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("IOChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithIOChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("IOChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithIOChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("IOChaos"), yaml)
}

// JVM Chaos

func (b *ChaosExperimentsBuilder) WithJVMChaos(c *chaosmeshv1alpha1.JVMChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("JVMChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithJVMChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("JVMChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithJVMChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("JVMChaos"), yaml)
}

// Kernel Chaos

func (b *ChaosExperimentsBuilder) WithKernelChaos(c *chaosmeshv1alpha1.KernelChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("KernelChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithKernelChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("KernelChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithKernelChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("KernelChaos"), yaml)
}

// Network Chaos

func (b *ChaosExperimentsBuilder) WithNetworkChaos(c *chaosmeshv1alpha1.NetworkChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithNetworkChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithNetworkChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("NetworkChaos"), yaml)
}

// Physical Machine Chaos

func (b *ChaosExperimentsBuilder) WithPhysicalMachineChaos(c *chaosmeshv1alpha1.PhysicalMachine) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PhysicalMachine"), c)
}

func (b *ChaosExperimentsBuilder) WithPhysicalMachineFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("PhysicalMachine"), filePath)
}

func (b *ChaosExperimentsBuilder) WithPhysicalMachineFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("PhysicalMachine"), yaml)
}

// Pod Chaos

func (b *ChaosExperimentsBuilder) WithPodChaos(c *chaosmeshv1alpha1.PodChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithPodChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("PodChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithPodChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("PodChaos"), yaml)
}

// Pod HTTP Chaos

func (b *ChaosExperimentsBuilder) WithPodHTTPChaos(c *chaosmeshv1alpha1.PodHttpChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodHttpChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithPodHttpChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("PodHttpChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithPodHttpChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("PodHttpChaos"), yaml)
}

// Pod IO Chaos

func (b *ChaosExperimentsBuilder) WithPodIOChaos(c *chaosmeshv1alpha1.PodIOChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodIOChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithPodIOChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("PodIOChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithPodIOChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("PodIOChaos"), yaml)
}

// Pod Network Chaos

func (b *ChaosExperimentsBuilder) WithPodNetworkChaos(c *chaosmeshv1alpha1.PodNetworkChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("PodNetworkChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithPodNetworkChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("PodNetworkChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithPodNetworkChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("PodNetworkChaos"), yaml)
}

// Stress Chaos

func (b *ChaosExperimentsBuilder) WithStressChaos(c *chaosmeshv1alpha1.StressChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("StressChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithStressChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("StressChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithStressChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("StressChaos"), yaml)
}

// Time Chaos

func (b *ChaosExperimentsBuilder) WithTimeChaos(c *chaosmeshv1alpha1.TimeChaos) *ChaosExperimentsBuilder {
	return b.withChaos(chaosmeshv1alpha1.GroupVersion.WithKind("TimeChaos"), c)
}

func (b *ChaosExperimentsBuilder) WithTimeChaosFromFile(filePath string) *ChaosExperimentsBuilder {
	return b.withChaosFromFile(chaosmeshv1alpha1.GroupVersion.WithKind("TimeChaos"), filePath)
}

func (b *ChaosExperimentsBuilder) WithTimeChaosFromYaml(yaml string) *ChaosExperimentsBuilder {
	return b.withChaosFromYaml(chaosmeshv1alpha1.GroupVersion.WithKind("TimeChaos"), yaml)
}

// Chaos Workflow

func (b *ChaosExperimentsBuilder) WithChaosWorkflow(c *chaosmeshv1alpha1.Workflow) *ChaosExperimentsBuilder {
	b.experiments.chaosWorkflows = append(b.experiments.chaosWorkflows, c)
	return b
}

func (b *ChaosExperimentsBuilder) WithChaosWorkflowFromFile(filePath string) *ChaosExperimentsBuilder {
	b.experiments.chaosWorkflowsFromFiles = append(b.experiments.chaosWorkflowsFromFiles, filePath)
	return b
}

func (b *ChaosExperimentsBuilder) WithChaosWorkflowFromYaml(yaml string) *ChaosExperimentsBuilder {
	b.experiments.chaosWorkflowsFromYaml = append(b.experiments.chaosWorkflowsFromYaml, yaml)
	return b
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

func (b *ChaosExperimentsBuilder) withChaosFromFile(gvk schema.GroupVersionKind, filePath string) *ChaosExperimentsBuilder {
	_, ok := b.experiments.chaosFromFiles[gvk]
	if !ok {
		b.experiments.chaosFromFiles[gvk] = []string{filePath}
	} else {
		b.experiments.chaosFromFiles[gvk] = append(b.experiments.chaosFromFiles[gvk], filePath)
	}
	return b
}

func (b *ChaosExperimentsBuilder) withChaosFromYaml(gvk schema.GroupVersionKind, yaml string) *ChaosExperimentsBuilder {
	_, ok := b.experiments.chaosFromYaml[gvk]
	if !ok {
		b.experiments.chaosFromYaml[gvk] = []string{yaml}
	} else {
		b.experiments.chaosFromYaml[gvk] = append(b.experiments.chaosFromYaml[gvk], yaml)
	}
	return b
}

func (b *ChaosExperimentsBuilder) build() *chaosExperiments {
	return b.experiments
}
