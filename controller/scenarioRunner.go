package scencontroller

import (
	fr "github.com/Dharitri-org/drtg-scenario/fileresolver"
	mjparse "github.com/Dharitri-org/drtg-scenario/json/parse"
	mj "github.com/Dharitri-org/drtg-scenario/model"
)

// ScenarioRunner describes a component that can run a VM scenario.
type ScenarioRunner interface {
	// Reset clears state/world.
	Reset()

	// RunScenario executes the scenario and checks if it passed. Failure is signaled by returning an error.
	// The FileResolver helps with resolving external steps.
	// TODO: group into a "execution context" param.
	RunScenario(*mj.Scenario, fr.FileResolver) error
}

// ScenarioController is a component that can run json scenarios, using a provided executor.
type ScenarioController struct {
	Executor    ScenarioRunner
	RunsNewTest bool
	Parser      mjparse.Parser
}

// NewScenarioController creates new ScenarioController instance.
func NewScenarioController(executor ScenarioRunner, fileResolver fr.FileResolver) *ScenarioController {
	return &ScenarioController{
		Executor: executor,
		Parser:   mjparse.NewParser(fileResolver),
	}
}
