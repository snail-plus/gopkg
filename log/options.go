// Copyright 2024 eve.  All rights reserved.

package log

import (
	"github.com/spf13/pflag"
	"go.uber.org/zap/zapcore"
)

// Options contains configuration options for logging.
type Options struct {
	// DisableCaller specifies whether to include caller information in the log.
	DisableCaller bool `json:"disable-caller,omitempty" yaml:"disable-caller"`
	// DisableStacktrace specifies whether to record a stack trace for all messages at or above panic level.
	DisableStacktrace bool `json:"disable-stacktrace,omitempty" yaml:"disable-stacktrace"`
	// EnableColor specifies whether to output colored logs.
	EnableColor bool `json:"enable-color"       yaml:"enable-color"`
	// Level specifies the minimum log level. Valid values are: debug, info, warn, error, dpanic, panic, and fatal.
	Level string `json:"level,omitempty" yaml:"level"`
	// Format specifies the log output format. Valid values are: console and json.
	Format string `json:"format,omitempty" yaml:"format"`

	// OutputPaths specifies the output paths for the logs.
	OutputPaths []string `yaml:"output-paths"`
	MaxSize     int      `yaml:"max-size"`    // 单个日志文件最大大小（单位：MB）
	MaxBackups  int      `yaml:"max-backups"` // 保留的旧日志文件个数
	MaxAge      int      `yaml:"max-age"`     // 保留的旧日志文件最大天数（单位：天）
}

// NewOptions creates a new Options object with default values.
func NewOptions() *Options {
	return &Options{
		Level:       zapcore.InfoLevel.String(),
		Format:      "console",
		OutputPaths: []string{"stdout"},
	}
}

// Validate verifies flags passed to LogsOptions.
func (o *Options) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds command line flags for the configuration.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Level, "log.level", o.Level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace, ""+
		"Disable the log to record a stack trace for all messages at or above panic level.")
	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringVar(&o.Format, "log.format", o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.StringSliceVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "Output paths of log.")
}
