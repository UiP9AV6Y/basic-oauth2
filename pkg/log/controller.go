package log

import (
	l "log"
)

type Controller struct {
	level Level
	trace *Facility
	debug *Facility
	info  *Facility
	warn  *Facility
	error *Facility
	fatal *Facility
}

func NewController(level Level, logger *l.Logger) *Controller {
	switch level {
	case LevelOff:
		return NewOffController()
	case LevelFatal:
		return NewFatalController(logger)
	case LevelError:
		return NewErrorController(logger)
	case LevelWarn:
		return NewWarnController(logger)
	case LevelInfo:
		return NewInfoController(logger)
	case LevelDebug:
		return NewDebugController(logger)
	}

	return NewTraceController(logger)
}

func NewTraceController(logger *l.Logger) *Controller {
	enabled := NewFacility(logger)
	controller := &Controller{
		level: LevelTrace,
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
		level: LevelDebug,
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
		level: LevelInfo,
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
		level: LevelWarn,
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
		level: LevelError,
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
		level: LevelFatal,
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
		level: LevelOff,
		trace: DisabledFacility,
		debug: DisabledFacility,
		info:  DisabledFacility,
		warn:  DisabledFacility,
		error: DisabledFacility,
		fatal: DisabledFacility,
	}

	return controller
}

func (c *Controller) Level() Level {
	return c.level
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
