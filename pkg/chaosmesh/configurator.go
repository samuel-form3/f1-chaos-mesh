package chaosmesh

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"time"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/form3tech-oss/f1/pkg/f1/testing"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
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

	for gvk, exps := range c.experiments.existingChaos {
		for _, ce := range exps {
			err := c.enableExistingChaos(gvk, ce)
			if err != nil {
				return err
			}
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

	for gvk, exps := range c.experiments.existingChaos {
		for _, ce := range exps {
			err := c.pauseExistingChaos(gvk, ce)
			if err != nil {
				c.t.Logger.Error(err)
			}
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

	err = wait.PollImmediate(2*time.Second, 1*time.Minute, func() (bool, error) {
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

		return isExperimentAllInjected(status), nil
	})

	if err != nil {
		c.t.Logger.Errorf("Chaos experiment %s was not injected, err: %s", expFriendlyName, err)
		return err
	}

	return nil
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

func (c *experimentsConfigurator) enableExistingChaos(gvk schema.GroupVersionKind, nn types.NamespacedName) error {
	ce := &unstructured.Unstructured{}
	ce.SetGroupVersionKind(gvk)

	expFriendlyName := generateExperimentFriendlyName(gvk.Kind, nn.Namespace, nn.Name)
	c.t.Logger.Infof("Enabling chaos experiment %s", expFriendlyName)

	err := c.kubeCli.Get(context.Background(), nn, ce)
	if err != nil {
		c.t.Logger.Errorf("Error setting up chaos experiment %s, err: %s", expFriendlyName, err)
		return err
	}

	annotations := ce.GetAnnotations()
	delete(annotations, "experiment.chaos-mesh.org/pause")
	ce.SetAnnotations(annotations)
	err = c.kubeCli.Update(context.Background(), ce)
	if err != nil {
		c.t.Logger.Errorf("Could not update existing network chaos %s, err: %s", expFriendlyName, err)
		return err
	}

	return nil
}

func (c *experimentsConfigurator) pauseExistingChaos(gvk schema.GroupVersionKind, nn types.NamespacedName) error {
	ce := &unstructured.Unstructured{}
	ce.SetGroupVersionKind(gvk)

	expFriendlyName := generateExperimentFriendlyName(gvk.Kind, nn.Namespace, nn.Name)
	c.t.Logger.Infof("Pausing chaos experiment %s", expFriendlyName)

	err := c.kubeCli.Get(context.Background(), nn, ce)
	if err != nil {
		c.t.Logger.Errorf("Error getting chaos experiment %s, err: %s", expFriendlyName, err)
		return err
	}

	annotations := ce.GetAnnotations()
	annotations["experiment.chaos-mesh.org/pause"] = "true"
	ce.SetAnnotations(annotations)
	err = c.kubeCli.Update(context.Background(), ce)
	if err != nil {
		c.t.Logger.Errorf("Error updating chaos experiment %s", expFriendlyName)
		return err
	}

	return nil
}

func isExperimentAllInjected(status chaosmeshv1alpha1.ChaosStatus) bool {
	for _, sc := range status.Conditions {
		if sc.Type == chaosmeshv1alpha1.ConditionAllInjected && sc.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func generateExperimentFriendlyName(experimentType string, namespace string, name string) string {
	return fmt.Sprintf("[%s]::%s/%s", experimentType, namespace, name)
}
