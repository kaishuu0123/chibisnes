package chibisnes

type Controller struct {
	console        *Console
	controllerType byte
	// latchline
	latchLine bool
	// for controller
	currentState uint16 // actual state
	latchedState uint16
}

const (
	ButtonB = iota
	ButtonY
	ButtonSelect
	ButtonStart
	ButtonUp
	ButtonDown
	ButtonLeft
	ButtonRight
	ButtonA
	ButtonX
	ButtonL
	ButtonR
)

func NewController(console *Console) *Controller {
	return &Controller{
		console: console,
	}
}

func (controller *Controller) Reset() {
	controller.latchLine = false
	controller.latchedState = 0
}

func (controller *Controller) Cycle() {
	if controller.latchLine {
		controller.latchedState = controller.currentState
	}
}

func (controller *Controller) Read() byte {
	var ret byte = byte(controller.latchedState) & 1
	controller.latchedState >>= 1
	controller.latchedState |= 0x8000
	return ret
}
