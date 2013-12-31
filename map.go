// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// map.go [created: Wed,  5 Jun 2013]

package enumstr

import (
	"fmt"
)

type Map map[interface{}]string

func (m Map) invert() (map[string]interface{}, error) {
	rev := make(map[string]interface{}, len(m))
	for v, str := range m {
		if _, ok := rev[str]; ok {
			return nil, fmt.Errorf("duplicate string %q", str)
		}
		rev[str] = v
	}
	return rev, nil
}

type enumMap struct {
	m   Map
	rev map[string]interface{}
}

func EnumMap(m Map) Enum {
	rev, err := m.invert()
	if err != nil {
		panic(err)
	}
	return &enumMap{m, rev}
}

func (enum *enumMap) Format(v interface{}) (string, error) {
	// TODO catch panic
	if str, ok := enum.m[v]; ok {
		return str, nil
	}
	return "", ErrUnknown
}

func (enum *enumMap) Parse(str string) (interface{}, error) {
	for v, vstr := range enum.m {
		if vstr == str {
			return v, nil
		}
	}
	return "", ErrUnknown
}

type enumMapDefault struct {
	def interface{}
	m   *enumMap
}

func EnumMapDefault(def interface{}, m Map) EnumDefault {
	if _, ok := m[def]; ok {
		return &enumMapDefault{def, EnumMap(m).(*enumMap)}
	}
	panic("unknown default")
}

func (enum *enumMapDefault) Format(v interface{}) string {
	str, err := enum.m.Format(v)
	if err == nil {
		return str
	}
	str, _ = enum.m.Format(enum.def) // membership guaranteed
	return str
}

func (enum *enumMapDefault) Parse(str string) interface{} {
	v, err := enum.m.Parse(str)
	if err == nil {
		return v
	}
	return enum.def
}
