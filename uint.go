// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// uint.go [created: Wed,  5 Jun 2013]

package enumstr

import (
	"fmt"
)

func dupCheck(strs []string) {
	set := make(map[string]int)
	for _, s := range strs {
		set[s] += 1
	}
	dups := make([]string, 0, len(set))
	for k, count := range set {
		if count > 1 {
			dups = append(dups, k)
		}
	}
	if len(dups) > 0 {
		panic(fmt.Errorf("duplicate strings: %v", dups))
	}
}

type uintEnum []string

func (enum uintEnum) Format(v interface{}) (string, error) {
	if i, ok := v.(uint); ok {
		if int(i) < len(enum) {
			return enum[i], nil
		}
		return "", ErrUnknown
	}
	return "", fmt.Errorf("expected uint but received %T", v)
}

func (enum uintEnum) Parse(str string) (interface{}, error) {
	for i := range enum {
		if enum[i] == str {
			return uint(i), nil
		}
	}
	return 0, ErrUnknown
}

type uintEnumDefault struct {
	invalid uint
	uintEnum
}

// like UintEnum but implements EnumDefault.
// does not handle non-consective sets well.
// does not handle non-empty sets excluding zero well.
func UintEnumDefault(invalid uint, m []string) EnumDefault {
	dupCheck(m)
	if int(invalid) < len(m) {
		return &uintEnumDefault{invalid, uintEnum(m)}
	}
	panic("unknown default")
}

func (enum *uintEnumDefault) Format(v interface{}) string {
	str, err := enum.uintEnum.Format(v)
	if err != nil {
		return enum.uintEnum[enum.invalid]
	}
	return str
}

func (enum *uintEnumDefault) Parse(str string) interface{} {
	v, err := enum.uintEnum.Parse(str)
	if err != nil {
		return enum.invalid
	}
	return v
}

// does not handle non-consective sets well.
// does not handle non-empty sets excluding zero well.
func UintEnum(m []string) Enum {
	dupCheck(m)
	enum := make(uintEnum, len(m))
	copy(enum, m)
	return enum
}
