package scenexec

import (
	"github.com/Dharitri-org/drtg-core/core/check"
	logger "github.com/Dharitri-org/drtg-logger"
	fr "github.com/Dharitri-org/drtg-scenario/scenario/expression/fileresolver"
	er "github.com/Dharitri-org/drtg-scenario/scenario/expression/reconstructor"
	scenio "github.com/Dharitri-org/drtg-scenario/scenario/io"
	scenmodel "github.com/Dharitri-org/drtg-scenario/scenario/model"
	worldmock "github.com/Dharitri-org/drtg-scenario/worldmock"
	vmcommon "github.com/Dharitri-org/drtg-vm-common"
)

var log = logger.GetOrCreate("scenario-exec")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ScenarioExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with VM.
type ScenarioExecutor struct {
	World             *worldmock.MockWorld
	vmBuilder         VMBuilder
	vm                VMInterface
	checkGas          bool
	scenarioTraceGas  []bool
	fileResolver      fr.FileResolver
	exprReconstructor er.ExprReconstructor
}

var _ scenio.ScenarioRunner = (*ScenarioExecutor)(nil)

// NewScenarioExecutor prepares a new VMTestExecutor instance.
func NewScenarioExecutor(vmBuilder VMBuilder) *ScenarioExecutor {
	world := vmBuilder.NewMockWorld()

	return &ScenarioExecutor{
		World:             world,
		vm:                nil,
		vmBuilder:         vmBuilder,
		checkGas:          true,
		scenarioTraceGas:  make([]bool, 0),
		fileResolver:      nil,
		exprReconstructor: er.ExprReconstructor{},
	}
}

// InitVM will initialize the VM and the builtin function container.
// Does nothing if the VM is already initialized.
func (ae *ScenarioExecutor) InitVM(scenGasSchedule scenmodel.GasSchedule) error {
	if ae.vm != nil {
		return nil
	}

	gasSchedule, err := ae.vmBuilder.GasScheduleMapFromScenarios(scenGasSchedule)
	if err != nil {
		return err
	}

	err = ae.World.InitBuiltinFunctions(gasSchedule)
	if err != nil {
		return err
	}

	ae.vm, err = ae.vmBuilder.NewVM(ae.World, gasSchedule)

	return err
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ScenarioExecutor) GetVM() vmcommon.VMExecutionHandler {
	return ae.vm
}

// Reset clears state/world.
// Is called in RunAllJSONScenariosInDirectory, but not in RunSingleJSONScenario.
func (ae *ScenarioExecutor) Reset() {
	if !check.IfNil(ae.vm) {
		ae.vm.Reset()
	}
	ae.World.Clear()
}

// Close will simply close the VM
func (ae *ScenarioExecutor) Close() {
	if !check.IfNil(ae.vm) {
		ae.vm.Reset()
	}
}

// PeekTraceGas returns the last position from the scenarioTraceGas, if existing
func (ae *ScenarioExecutor) PeekTraceGas() bool {
	length := len(ae.scenarioTraceGas)
	if length != 0 {
		return ae.scenarioTraceGas[length-1]
	}
	return false
}