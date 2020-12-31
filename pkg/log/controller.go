package log

import (
	l "log"
)

type Controller struct {
	trace *Facility
	debug *Facility
	info  *Facility
	warn  *Facility
	error *Facility
	fatal *Facility
}

func NewController(verbosity int, logger *l.Logger) *Controller {
	switch verbosity {
	case 0:
		return NewOffController()
	case 1:
		return NewFatalController(logger)
	case 2:
		return NewErrorController(logger)
	case 3:
		return NewWarnController(logger)
	case 4:
		return NewInfoController(logger)
	case 5:
		return NewDebugController(logger)
	}

	return NewTraceController(logger)
}

func NewTraceController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: enabled,
		debug: enabled,
		info:  enabled,
		warn:  enabled,
		error: enabled,
		fatal: enabled,
	}

	return controller
}

func NewDebugController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: DisabledFacility,
		debug: enabled,
		info:  enabled,
		warn:  enabled,
		error: enabled,
		fatal: enabled,
	}

	return controller
}

func NewInfoController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  enabled,
		warn:  enabled,
		error: enabled,
		fatal: enabled,
	}

	return controller
}

func NewWarnController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  DisabledFacility,
		warn:  enabled,
		error: enabled,
		fatal: enabled,
	}

	return controller
}

func NewErrorController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  DisabledFacility,
		warn:  DisabledFacility,
		error: enabled,
		fatal: enabled,
	}

	return controller
}

func NewFatalController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  DisabledFacility,
		warn:  DisabledFacility,
		error: DisabledFacility,
		fatal: enabled,
	}

	return controller
}

func NewOffController() *Controller {
	controller := &Controller{
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  DisabledFacility,
		warn:  DisabledFacility,
		error: DisabledFacility,
		fatal: DisabledFacility,
	}

	return controller
}

func (c *Controller) Trace() *Facility {
	return c.trace
}

func (c *Controller) Debug() *Facility {
	return c.debug
}

func (c *Controller) Info() *Facility {
	return c.info
}

func (c *Controller) Warn() *Facility {
	return c.warn
}

func (c *Controller) Error() *Facility {
	return c.error
}

func (c *Controller) Fatal() *Facility {
	return c.fatal
}
