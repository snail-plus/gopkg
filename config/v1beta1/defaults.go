// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/onex.
//

package v1beta1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RecommendedDefaultMySQLConfiguration defaults a pointer to a
// MySQLConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultMySQLConfiguration(obj *MySQLConfiguration) {
	if obj.Host == "" {
		obj.Host = "127.0.0.1:3306"
	}

	if obj.MaxIdleConnections == 0 {
		obj.MaxIdleConnections = 100
	}

	if obj.MaxOpenConnections == 0 {
		obj.MaxOpenConnections = 100
	}

	zero := metav1.Duration{}
	if obj.MaxConnectionLifeTime == zero {
		obj.MaxConnectionLifeTime = metav1.Duration{Duration: 10 * time.Second}
	}
}

// RecommendedDefaultRedisConfiguration defaults a pointer to a
// RedisConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultRedisConfiguration(obj *RedisConfiguration) {
	if obj.Addr == "" {
		obj.Addr = "127.0.0.1:6379"
	}

	if obj.MaxRetries == 0 {
		obj.MaxRetries = 3
	}

	zero := metav1.Duration{}
	if obj.Timeout == zero {
		obj.Timeout = metav1.Duration{Duration: 10 * time.Second}
	}
}

/*
// RecommendedDefaultClientConnectionConfiguration defaults a pointer to a
// ClientConnectionConfiguration struct. This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDefaultClientConnectionConfiguration(obj *ClientConnectionConfiguration) {
	if len(obj.ContentType) == 0 {
		obj.ContentType = "application/vnd.kubernetes.protobuf"
	}
	if obj.QPS == 0.0 {
		obj.QPS = 50.0
	}
	if obj.Burst == 0 {
		obj.Burst = 100
	}
}

// RecommendedDebuggingConfiguration defaults profiling and debugging configuration.
// This will set the recommended default
// values, but they may be subject to change between API versions. This function
// is intentionally not registered in the scheme as a "normal" `SetDefaults_Foo`
// function to allow consumers of this type to set whatever defaults for their
// embedded configs. Forcing consumers to use these defaults would be problematic
// as defaulting in the scheme is done as part of the conversion, and there would
// be no easy way to opt-out. Instead, if you want to use this defaulting method
// run it in your wrapper struct of this type in its `SetDefaults_` method.
func RecommendedDebuggingConfiguration(obj *DebuggingConfiguration) {
	if obj.EnableProfiling == nil {
		obj.EnableProfiling = utilpointer.Bool(true) // profile debugging is cheap to have exposed and standard on kube binaries
	}
}

// NewRecommendedDebuggingConfiguration returns the current recommended DebuggingConfiguration.
// This may change between releases as recommendations shift.
func NewRecommendedDebuggingConfiguration() *DebuggingConfiguration {
	ret := &DebuggingConfiguration{}
	RecommendedDebuggingConfiguration(ret)
	return ret
}
*/
