package e2e

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	chaosmeshv1alpha1 "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	k8sClient client.Client
)

func TestMain(m *testing.M) {
	err := StartEnv()
	if err != nil {
		StopEnv()
		panic(err)
	}
	code := m.Run()

	StopEnv()
	os.Exit(code)
}

func StartEnv() error {
	err := createKindCluster()
	if err != nil {
		return err
	}

	err = installChaosMesh()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	err = chaosmeshv1alpha1.AddToScheme(scheme)
	if err != nil {
		return err
	}

	k8sClient, err = client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return err
	}

	return nil
}

func StopEnv() error {

	return deleteKindCluster()
}

func createKindCluster() error {
	clusterList, err := exec.Command("kind", "get", "clusters").Output()
	if err != nil {
		return err
	}
	if strings.Contains(string(clusterList), "f1-chaos-mesh-tests") {
		return nil
	}

	cmd := exec.Command("kind", "create", "cluster", "--name", "f1-chaos-mesh-tests")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func installChaosMesh() error {
	cmd := exec.Command("helm", "repo", "add", "chaos-mesh", "https://charts.chaos-mesh.org")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command(
		"helm", "upgrade", "-i",
		"chaos-mesh", "chaos-mesh/chaos-mesh",
		"--set", "chaosDaemon.runtime=containerd",
		"--set", "chaosDaemon.socketPath=/run/containerd/containerd.sock",
		"--set", "controllerManager.replicaCount=1",
		"--wait")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

func deleteKindCluster() error {
	reuseCluster := os.Getenv("REUSE_CLUSTER")
	if reuseCluster == "true" {
		return nil
	}

	cmd := exec.Command("kind", "delete", "cluster", "--name", "f1-chaos-mesh-tests")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
