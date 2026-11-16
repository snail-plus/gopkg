// Copyright 2024 eve.  All rights reserved.

package app

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/snail-plus/gopkg/log"
	"github.com/snail-plus/gopkg/util/homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const configFlagName = "config"

// CliOptions abstracts configuration options for reading parameters from the
// command line.
//
// 为兼容 options 包中已有的 Validate() []error 实现（见 options.IOptions），
// 这里不强制 Validate 的具体签名，运行时支持以下组合：
//   - 可选 Complete() error，加 Validate() error
//   - Validate() []error（Complete 可选）
type CliOptions interface{}

func completeOptions(opts CliOptions) error {
	if c, ok := opts.(interface{ Complete() error }); ok {
		return c.Complete()
	}
	return nil
}

func validateOptions(opts CliOptions) error {
	switch v := opts.(type) {
	case interface{ Validate() error }:
		return v.Validate()
	case interface{ Validate() []error }:
		return errors.Join(v.Validate()...)
	default:
		return nil
	}
}

var cfgFile string

// AddConfigFlag adds flags for a specific server to the specified FlagSet object.
// It also sets a passed functions to read values from configuration file into viper
// when each cobra command's Execute method is called.
func AddConfigFlag(fs *pflag.FlagSet, name string) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	// Enable viper's automatic environment variable parsing. This means
	// that viper will automatically read values corresponding to viper
	// variables from environment variables.
	viper.AutomaticEnv()
	// Set the environment variable prefix. Use the strings.ReplaceAll function
	// to replace hyphens with underscores in the name, and use strings.ToUpper
	// to convert the name to uppercase, then set it as the prefix for environment variables.
	viper.SetEnvPrefix(strings.ReplaceAll(strings.ToUpper(name), "-", "_"))
	// Set the replacement rules for environment variable keys. Use the
	// strings.NewReplacer function to specify replacing periods and hyphens with underscores.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")

			if names := strings.Split(name, "-"); len(names) > 1 {
				viper.AddConfigPath(filepath.Join(homedir.HomeDir(), "."+names[0]))
				viper.AddConfigPath(filepath.Join("/etc", names[0]))
			}

			viper.SetConfigName(name)
		}

		if err := viper.ReadInConfig(); err != nil {
			log.Debugw("Failed to read configuration file", "file", cfgFile, "err", err)
		}

	})
}

// sensitiveConfigKey 判断配置项是否含敏感信息（口令、密钥、令牌等）。
func sensitiveConfigKey(key string) bool {
	lower := strings.ToLower(key)
	for _, part := range []string{"password", "passwd", "secret", "token", "credential", "private-key"} {
		if strings.Contains(lower, part) {
			return true
		}
	}
	for _, seg := range strings.FieldsFunc(lower, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	}) {
		if seg == "key" {
			return true
		}
	}
	return false
}

// PrintConfig 以 debug 级别打印配置，敏感项以 *** 代替，避免口令/密钥落盘。
func PrintConfig() {
	for _, key := range viper.AllKeys() {
		value := any("***")
		if !sensitiveConfigKey(key) {
			value = viper.Get(key)
		}
		log.Debugw(fmt.Sprintf("CFG: %s=%v", key, value))
	}
}

func init() {
	pflag.StringVarP(&cfgFile, configFlagName, "c", cfgFile, "Read configuration from specified `FILE`, "+
		"support JSON, TOML, YAML, HCL, or Java properties formats.")
}
