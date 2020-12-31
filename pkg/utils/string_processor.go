package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type StringProcessor interface {
	Process(string) (string, error)
}

type FunctorProcessor func(string) string

func (p FunctorProcessor) Process(s string) (string, error) {
	return p(s), nil
}

type ExpandingProcessor struct {
	callback  func(string) string
	lastError error
}

func NewExpandingProcessor(callback StringProcessor) *ExpandingProcessor {
	processor := &ExpandingProcessor{}
	processor.callback = func(s string) string {
		v, err := callback.Process(s)
		if err != nil {
			processor.lastError = err
		}
		return v
	}

	return processor
}

func (p *ExpandingProcessor) Process(s string) (string, error) {
	p.lastError = nil
	v := os.Expand(s, p.callback)
	err := p.lastError
	p.lastError = nil
	return v, err
}

type FileReferenceProcessor struct {
	Indicator byte
	Directory string
}

func NewRelativeFileReferenceProcessor(indicator byte, file string) *FileReferenceProcessor {
	directory := filepath.Dir(file)
	processor := &FileReferenceProcessor{
		Indicator: indicator,
		Directory: directory,
	}

	return processor
}

func (p *FileReferenceProcessor) Process(s string) (string, error) {
	if len(s) <= 1 || s[0] != p.Indicator {
		return s, nil
	}

	file := ResolvePath(p.Directory, s[1:])
	b, err := ioutil.ReadFile(file)
	if b != nil {
		return string(b), err
	}

	return "", err
}

type ChainProcessor []StringProcessor

func NewChainProcessor(p ...StringProcessor) ChainProcessor {
	return ChainProcessor(p)
}

func (p ChainProcessor) Process(s string) (v string, err error) {
	v = s

	for _, c := range p {
		v, err = c.Process(v)
		if err != nil {
			return
		}
	}

	return
}

var (
	EnvProcessor StringProcessor = FunctorProcessor(os.Getenv)
)
