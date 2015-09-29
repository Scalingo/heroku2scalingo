// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package envconfig

import (
	"os"
	"testing"
)

type Specification struct {
	Debug bool
	Port  int
	Rate  float32
	User  string
}

type Specification2 struct {
	DebugLevel string
	ServerPort int
}

func TestProcess(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	err := Process("env_config", &s)
	if err != nil {
		t.Error(err.Error())
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.User != "Kelsey" {
		t.Errorf("expected %s, got %s", "Kelsey", s.User)
	}

	var s2 Specification2
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG_LEVEL", "info")
	os.Setenv("ENV_CONFIG_SERVER_PORT", "3000")
	err = Process("env_config", &s2)
	if err != nil {
		t.Error(err.Error())
	}
	if s2.DebugLevel != "info" {
		t.Errorf("expected %s, got %s", "info", s2.DebugLevel)
	}
	if s2.ServerPort != 3000 {
		t.Errorf("expected %d, got %d", 3000, s2.ServerPort)
	}
}

func testProcessWithoutPrefix(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("DEBUG", "true")
	os.Setenv("PORT", "8080")
	os.Setenv("RATE", "0.5")
	os.Setenv("USER", "Kelsey")
	err := Process("", &s)
	if err != nil {
		t.Error(err.Error())
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.User != "Kelsey" {
		t.Errorf("expected %s, got %s", "Kelsey", s.User)
	}
}

func TestParseErrorBool(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Debug" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

func TestParseErrorFloat32(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_RATE", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Rate" {
		t.Errorf("expected %s, got %v", "Rate", v.FieldName)
	}
	if s.Rate != 0 {
		t.Errorf("expected %v, got %v", 0, s.Rate)
	}
}

func TestParseErrorInt(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_PORT", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Port" {
		t.Errorf("expected %s, got %v", "Port", v.FieldName)
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}
}

func TestErrInvalidSpecification(t *testing.T) {
	m := make(map[string]string)
	err := Process("env_config", &m)
	if err != ErrInvalidSpecification {
		t.Errorf("expected %v, got %v", ErrInvalidSpecification, err)
	}
}

func TestToUnderscoreCase(t *testing.T) {
	// In golang an exported field name is capitalized camelcase
	input := "SampleFieldName"
	output := toUnderscoreCase(input)
	expectedOutput := "Sample_Field_Name"
	if output != expectedOutput {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}
