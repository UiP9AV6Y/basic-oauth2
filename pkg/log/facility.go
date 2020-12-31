package log

import (
	"io/ioutil"
	l "log"
)

type Facility struct {
	enabled bool
	logger  *l.Logger
}

var DisabledFacility *Facility = disabledFacility()

func NewFacility(logger *l.Logger) *Facility {
	facility := &Facility{
		enabled: true,
		logger:  logger,
	}

	return facility
}

func disabledFacility() *Facility {
	logger := l.New(ioutil.Discard, "", l.LstdFlags)
	facility := &Facility{
		enabled: false,
		logger:  logger,
	}

	return facility
}

func (f *Facility) Enabled() bool {
	return f.enabled
}

func (f *Facility) Logger() *l.Logger {
	return f.logger
}

func (f *Facility) Print(v ...interface{}) {
	f.logger.Print(v...)
}

func (f *Facility) Println(v ...interface{}) {
	f.logger.Println(v...)
}

func (f *Facility) Printf(format string, v ...interface{}) {
	f.logger.Printf(format, v...)
}
