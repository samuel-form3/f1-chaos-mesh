#!/usr/bin/env bash

# create cluster
kind create cluster

# setup helm repositories
helm repo add chaos-mesh https://charts.chaos-mesh.org

# install chaos-mesh
kubectl create ns chaos-testing
helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-testing --set chaosDaemon.runtime=containerd --set chaosDaemon.socketPath=/run/containerd/containerd.sock