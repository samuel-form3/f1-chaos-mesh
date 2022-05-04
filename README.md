# f1-chaos-mesh
f1-chaos-mesh

# Installation

```bash
go get github.com/samuel-form3/f1-chaos-mesh
```

# How to use

See in the examples folder


```go
func main() {
	f1Chaos := chaosmesh.NewChaosPlugin()
	f1Scenarios := f1.Scenarios().
		Add("one", scenarioOne).
		Add("oneWithChaos", scenarioOne, f1Chaos.WithExperiments(scenarioOneChaosExperiments))

	f1Scenarios.Execute()
}
```