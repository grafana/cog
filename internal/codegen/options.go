package codegen

import (
	"fmt"
	"maps"
)

func StdoutReporter(msg string) {
	fmt.Println(msg)
}

func Parameters(extraParameters map[string]string) PipelineOption {
	return func(pipeline *Pipeline) {
		maps.Copy(pipeline.Parameters, extraParameters)
		pipeline.interpolateParameters()
	}
}

func Reporter(reporter ProgressReporter) PipelineOption {
	return func(pipeline *Pipeline) {
		pipeline.reporter = reporter
	}
}

func Debug(enabled bool) PipelineOption {
	return func(pipeline *Pipeline) {
		pipeline.Debug = enabled
	}
}
