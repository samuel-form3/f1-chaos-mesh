package chaosmesh

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1/testing"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	yamlUtil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type experimentsConfigurator struct {
	experiments *chaosExperiments
	kubeCli     client.Client
	t           *testing.T
}

func newExperimentsConfigurator(t *testing.T, kubeCli client.Client, experiments *chaosExperiments) *experimentsConfigurator {
	return &experimentsConfigurator{
		experiments: experiments,
		kubeCli:     kubeCli,
		t:           t,
	}
}

func (c *experimentsConfigurator) ConfigureExperiments() error {
	c.t.Logger.Info("Setting up chaos experiments")

	for gvk, cc := range c.experiments.chaos {
		for _, ccc := range cc {
			err := c.createChaos(gvk, ccc)
			if err != nil {
				c.t.Error(err)
				return err
			}
		}
	}

	for gvk, exps := range c.experiments.chaosFromFiles {
		for _, filePath := range exps {
			err := c.createChaosFromFile(gvk, filePath)
			if err != nil {
				return err
			}
		}
	}

	for gvk, exps := range c.experiments.chaosFromYaml {
		for _, yaml := range exps {
			err := c.createChaosFromYaml(gvk, yaml)
			if err != nil {
				return err
			}
		}
	}

	for _, wf := range c.experiments.chaosWorkflows {
		err := c.createChaosWorkflow(wf)
		if err != nil {
			return err
		}
	}

	for _, filePath := range c.experiments.chaosWorkflowsFromFiles {
		err := c.createChaosWorkflowFromFile(filePath)
		if err != nil {
			return err
		}
	}

	for _, yaml := range c.experiments.chaosWorkflowsFromYaml {
		err := c.createChaosWorkflowFromYaml(yaml)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *experimentsConfigurator) CleanupExperiments() error {
	c.t.Logger.Info("Cleaning up chaos experiments")

	for gvk, cc := range c.experiments.chaos {
		for _, ccc := range cc {
			err := c.deleteChaos(gvk, ccc)
			if err != nil {
				c.t.Logger.Error(err)
			}
		}
	}

	for gvk, cc := range c.experiments.chaosFromFiles {
		for _, ccc := range cc {
			err := c.deleteChaosFromFile(gvk, ccc)
			if err != nil {
				c.t.Logger.Error(err)
			}
		}
	}

	for gvk, cc := range c.experiments.chaosFromYaml {
		for _, ccc := range cc {
			err := c.deleteChaosFromYaml(gvk, ccc)
			if err != nil {
				c.t.Logger.Error(err)
			}
		}
	}

	for _, wf := range c.experiments.chaosWorkflows {
		err := c.deleteChaosWorkflow(wf)
		if err != nil {
			c.t.Logger.Error(err)
		}
	}

	for _, wf := range c.experiments.chaosWorkflowsFromFiles {
		err := c.deleteChaosWorkflowFromFile(wf)
		if err != nil {
			c.t.Logger.Error(err)
		}
	}

	for _, wf := range c.experiments.chaosWorkflowsFromYaml {
		err := c.deleteChaosWorkflowFromYaml(wf)
		if err != nil {
			c.t.Logger.Error(err)
		}
	}

	return nil
}

func (c *experimentsConfigurator) createChaos(gvk schema.GroupVersionKind, obj *unstructured.Unstructured) error {
	expFriendlyName := generateExperimentFriendlyName(gvk.Kind, obj.GetNamespace(), obj.GetName())

	c.t.Logger.Infof("Setting up chaos experiment %s", expFriendlyName)
	err := c.kubeCli.Create(context.Background(), obj, &client.CreateOptions{})
	if err != nil {
		c.t.Logger.Errorf("Error setting up chaos experiment %s", expFriendlyName)
		return err
	}

	err = c.waitForExperimentToBeInjected(gvk, obj, expFriendlyName)
	if err != nil {
		c.t.Logger.Errorf("Chaos experiment %s was not injected, err: %s", expFriendlyName, err)
		return err
	}

	return nil
}

func (c *experimentsConfigurator) waitForExperimentToBeInjected(gvk schema.GroupVersionKind, obj *unstructured.Unstructured, expFriendlyName string) error {
	return wait.PollImmediate(2*time.Second, 1*time.Minute, func() (bool, error) {
		updObj := &unstructured.Unstructured{}
		updObj.SetGroupVersionKind(gvk)

		err := c.kubeCli.Get(
			context.Background(),
			types.NamespacedName{Namespace: obj.GetNamespace(), Name: obj.GetName()},
			updObj)

		if err != nil {
			c.t.Logger.Infof("Could not get chaos experiment %s, err: %s", expFriendlyName, err)
			return false, nil
		}

		var status chaosmeshv1alpha1.ChaosStatus
		unstructuredStatus := updObj.Object["status"].(map[string]interface{})
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredStatus, &status)
		if err != nil {
			c.t.Logger.Infof("Could not get chaos experiment status %s, err: %s", expFriendlyName, err)
			return false, nil
		}

		return experimentHasCondition(status, chaosmeshv1alpha1.ConditionAllInjected), nil
	})
}

func (c *experimentsConfigurator) createChaosFromFile(gvk schema.GroupVersionKind, filePath string) error {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	err := unmarshalFile(filePath, obj)
	if err != nil {
		return err
	}

	return c.createChaos(gvk, obj)
}

func (c *experimentsConfigurator) createChaosFromYaml(gvk schema.GroupVersionKind, yaml string) error {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	err := unmarshalYaml(yaml, obj)
	if err != nil {
		return err
	}

	return c.createChaos(gvk, obj)
}

func (c *experimentsConfigurator) deleteChaos(gvk schema.GroupVersionKind, obj *unstructured.Unstructured) error {
	expFriendlyName := generateExperimentFriendlyName(gvk.Kind, obj.GetNamespace(), obj.GetName())
	c.t.Logger.Infof("Cleaning up chaos experiment %s", expFriendlyName)
	err := c.kubeCli.Delete(context.Background(), obj, &client.DeleteOptions{})
	if err != nil {
		c.t.Logger.Errorf("Error cleaning up chaos experiment %s", expFriendlyName)
		return err
	}
	return nil
}

func (c *experimentsConfigurator) deleteChaosFromFile(gvk schema.GroupVersionKind, filePath string) error {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	err := unmarshalFile(filePath, obj)
	if err != nil {
		return err
	}

	return c.deleteChaos(gvk, obj)
}

func (c *experimentsConfigurator) deleteChaosFromYaml(gvk schema.GroupVersionKind, yaml string) error {
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)
	err := unmarshalYaml(yaml, obj)
	if err != nil {
		return err
	}

	return c.deleteChaos(gvk, obj)
}

func (c *experimentsConfigurator) createChaosWorkflow(wf *chaosmeshv1alpha1.Workflow) error {
	expFriendlyName := generateExperimentFriendlyName("Workflow", wf.GetNamespace(), wf.GetName())

	c.t.Logger.Infof("Setting up chaos workflow %s", expFriendlyName)
	err := c.kubeCli.Create(context.Background(), wf, &client.CreateOptions{})
	if err != nil {
		c.t.Logger.Errorf("Error setting up chaos workflow %s, err : %s", expFriendlyName, err)
		return err
	}

	err = wait.PollImmediate(2*time.Second, 1*time.Minute, func() (bool, error) {
		var updWf chaosmeshv1alpha1.Workflow
		err := c.kubeCli.Get(
			context.Background(),
			types.NamespacedName{Namespace: wf.GetNamespace(), Name: wf.GetName()},
			&updWf)

		if err != nil {
			c.t.Logger.Infof("Could not get chaos workflow %s, err: %s", expFriendlyName, err)
			return false, nil
		}

		for _, sc := range updWf.Status.Conditions {
			if sc.Type == chaosmeshv1alpha1.WorkflowConditionScheduled && sc.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})

	if err != nil {
		c.t.Logger.Errorf("Chaos experiment %s was not injected, err: %s", expFriendlyName, err)
		return err
	}

	return nil
}

func (c *experimentsConfigurator) createChaosWorkflowFromFile(filePath string) error {
	wf := &chaosmeshv1alpha1.Workflow{}
	err := unmarshalFile(filePath, wf)
	if err != nil {
		return err
	}

	return c.createChaosWorkflow(wf)
}

func (c *experimentsConfigurator) createChaosWorkflowFromYaml(yaml string) error {
	wf := &chaosmeshv1alpha1.Workflow{}
	err := unmarshalYaml(yaml, wf)
	if err != nil {
		return err
	}

	return c.createChaosWorkflow(wf)
}

func (c *experimentsConfigurator) deleteChaosWorkflow(wf *chaosmeshv1alpha1.Workflow) error {
	expFriendlyName := generateExperimentFriendlyName("Workflow", wf.GetNamespace(), wf.GetName())

	c.t.Logger.Infof("Deleting chaos workflow %s", expFriendlyName)
	err := c.kubeCli.Delete(context.Background(), wf, &client.DeleteOptions{})
	if err != nil {
		c.t.Logger.Errorf("Error deleting up chaos workflow %s", expFriendlyName)
		return err
	}

	return nil
}

func (c *experimentsConfigurator) deleteChaosWorkflowFromFile(filePath string) error {
	wf := &chaosmeshv1alpha1.Workflow{}
	err := unmarshalFile(filePath, wf)
	if err != nil {
		return err
	}

	return c.deleteChaosWorkflow(wf)
}

func (c *experimentsConfigurator) deleteChaosWorkflowFromYaml(yaml string) error {
	wf := &chaosmeshv1alpha1.Workflow{}
	err := unmarshalYaml(yaml, wf)
	if err != nil {
		return err
	}

	return c.deleteChaosWorkflow(wf)
}

func unmarshalFile(filePath string, obj interface{}) error {
	f, err := os.Open(filePath)
	if err != nil {
		return errors.Wrapf(err, "error opening file %s", filePath)
	}
	err = yamlUtil.NewYAMLOrJSONDecoder(f, 100).Decode(obj)
	if err != nil {
		return errors.Wrapf(err, "error decoding yaml from file %s", filePath)
	}

	return nil
}

func unmarshalYaml(yaml string, obj interface{}) error {
	yr := strings.NewReader(yaml)
	err := yamlUtil.NewYAMLOrJSONDecoder(yr, 100).Decode(obj)
	if err != nil {
		return err
	}

	return nil
}

func experimentHasCondition(status chaosmeshv1alpha1.ChaosStatus, condition chaosmeshv1alpha1.ChaosConditionType) bool {
	for _, sc := range status.Conditions {
		if sc.Type == condition && sc.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func generateExperimentFriendlyName(experimentType string, namespace string, name string) string {
	return fmt.Sprintf("[%s]::%s/%s", experimentType, namespace, name)
}
