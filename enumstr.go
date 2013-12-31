// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// enumstr.go [created: Tue,  4 Jun 2013]

/*
Package enumstr creates string mappings for serializing constant numeric values.

Basic usage

The most common way to use enumstr is as way to help generate String() methods.

	type Fruit uint
	const (
		Banana Fruit = iota
		Apple
		Potato
	)
	var fruitEnum = enumstr.UintEnum([]string{
		Banana: "naners",
		Apple:  "apps",
		Potato: "tots",
	})
	func (f Fruit) String() string {
		return enumstr.String(fruitEnum, "wut?")
	}

Lookup from string

Values can also be looked up by their string representation.

	func FruitString(str string) (Fruit, error) {
		enum, err := fruitEnum.Parse(str)
		if err == nil {
			return enum.(Fruit), nil
		}
		return 0, fmt.Errorf("invalid fruit: %#v", str)
	}

JSON serialization

Values can easily be serialized as JSON.

	func (f Fruit) MarshalJSON() ([]byte, error) {
		return MarshalJSON(fruitEnum, f)
	}
	func (f *Fruit) UnmarshalJSON(p []byte) error {
		return UnmarshalJSON(fruitEnum, p, f)
	}

Default values

If a value represents an invalid value, it makes writing methods even easier

	type Stooge uint
	const (
		Invalid Stooge = iota
		Larry
		Moe
		Curly
		Shemp
	)
	var stoogeEnum = enumstr.UintEnumDefault(Default, []string{
		Invalid: "INVALID_STOOGE",
		Larry:   "Larry",
		Moe:     "Moe",
		Curly:   "Curly",
		Shemp:   "Shemp",
	})
	func StoogeString(str string) (Stooge, error) {
		return stoogeEnum.Parse(str).(Stooge), nil
	}
	func (f Stooge) String() string {
		return stoogeEnum.Format(f)
	}
	func (f Stooge) MarshalJSON() ([]byte, error) {
		return MarshalJSON(stoogeEnum, f)
	}
	func (f *Stooge) UnmarshalJSON(p []byte) error {
		return UnmarshalJSON(stoogeEnum, p, f)
	}
*/
package enumstr

import (
	"encoding/json"
	"fmt"
	"reflect"
)

var ErrUnknown = fmt.Errorf("unknown")

type Enum interface {
	Format(interface{}) (string, error)
	Parse(string) (interface{}, error)
}

// an enumeration with an element that represents an invalid value.
// on invalid input maps should return values the invalid element.
type EnumDefault interface {
	Format(interface{}) string
	Parse(string) interface{}
}

func String(enum Enum, val interface{}, invalid string) string {
	str, err := enum.Format(val)
	if err != nil {
		return invalid
	}
	return str
}

func MarshalJSON(m Enum, val interface{}) ([]byte, error) {
	str, err := m.Format(val)
	if err != nil {
		return nil, err
	}
	return json.Marshal(str)
}

func UnmarshalJSON(m Enum, p []byte, val interface{}) (err error) {
	var str string
	err = json.Unmarshal(p, &str)
	if err != nil {
		return err
	}
	enum, err := m.Parse(str)
	if err != nil {
		return err
	}
	// assign to val...
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("can't assign %T to %T", enum, val)
		}
	}()
	reflect.Indirect(reflect.ValueOf(val)).Set(reflect.ValueOf(enum))
	return nil
}

func StringDefault(enum EnumDefault, val interface{}) string {
	return enum.Format(val)
}

func MarshalJSONDefault(m EnumDefault, val interface{}) ([]byte, error) {
	return json.Marshal(m.Format(val))
}

func UnmarshalJSONDefault(m EnumDefault, p []byte, val interface{}) (err error) {
	var str string
	err = json.Unmarshal(p, &str)
	if err != nil {
		return err
	}
	enum := m.Parse(str)
	// assign to val...
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("can't assign %T to %T", enum, val)
		}
	}()
	reflect.Indirect(reflect.ValueOf(val)).Set(reflect.ValueOf(enum))
	return nil
}
