package codegen

import (
	"fmt"
	"log/slog"
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

func Logger(logger *slog.Logger) PipelineOption {
	return func(pipeline *Pipeline) {
		pipeline.logger = logger
	}
}

func Debug(enabled bool) PipelineOption {
	return func(pipeline *Pipeline) {
		pipeline.Debug = enabled
	}
}
