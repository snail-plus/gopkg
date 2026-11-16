// Copyright 2024 eve.  All rights reserved.

// Package verflag defines utility functions to handle command line flags
// related to version of Kubernetes.
package version

import (
	"fmt"
	"strconv"

	flag "github.com/spf13/pflag"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const strRawVersion string = "raw"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() any {
	return *v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// The type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// PrintVersionIfRequested 检查是否传入 --version，若传入则打印版本并返回 true。
// 不自行调用 os.Exit，由调用方决定进程退出行为。
func PrintVersionIfRequested(appName string) bool {
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		return true
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s %s\n", appName, Get().GitVersion)
		return true
	}
	return false
}
