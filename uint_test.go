// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// uint.go [created: Wed,  5 Jun 2013]

package enumstr

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestStupidUnmarshalJSON(t *testing.T) {
	enum := UintEnum([]string{
		"abc",
		"def",
	})
	err := UnmarshalJSON(enum, []byte(`"abc"`), new(string))
	assert.NotEqual(t, nil, err)
}

func TestStupidUnmarshalJSONDefault(t *testing.T) {
	enum := UintEnumDefault(0, []string{
		"abc",
		"def",
	})
	err := UnmarshalJSONDefault(enum, []byte(`"abc"`), new(string))
	assert.NotEqual(t, nil, err)
}

func TestUintEnum(t *testing.T) {
	var recv interface{}
	defer func() { assert.Equal(t, nil, recv) }()
	defer func() { recv = recover() }()
	UintEnum([]string{
		"abc",
		"def",
	})
}

func TestUintEnumBadType(t *testing.T) {
	_, err := UintEnum([]string{
		"abc",
		"def",
	}).Format("not an int")
	assert.NotEqual(t, nil, err)
}

func TestUintEnumDuplicate(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	UintEnum([]string{
		"abc",
		"abc",
	})
}

func TestUintEnumDefault(t *testing.T) {
	var recv interface{}
	defer func() { assert.Equal(t, nil, recv) }()
	defer func() { recv = recover() }()
	UintEnumDefault(0, []string{
		"abc",
		"def",
	})
}

func TestUintEnumDefaultDuplicate(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	UintEnumDefault(0, []string{
		"abc",
		"abc",
	})
}

func TestUintEnumDefaultNotFound(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	UintEnumDefault(2, []string{
		"abc",
		"def",
	})
}

type Fruit uint

const (
	Banana Fruit = iota
	Apple
	Potato
)

var fruitEnum = UintEnum([]string{
	Banana: "naners",
	Apple:  "apps",
	Potato: "tots",
})

func ParseFruit(str string) (Fruit, error) {
	enum, err := fruitEnum.Parse(str)
	if err == nil {
		return Fruit(enum.(uint)), nil
	}
	return 0, fmt.Errorf("invalid fruit: %#v", str)
}

func (f Fruit) String() string {
	return String(fruitEnum, uint(f), nil, "wut?")
}

func (f Fruit) MarshalJSON() ([]byte, error) {
	return MarshalJSON(fruitEnum, uint(f))
}

func (f *Fruit) UnmarshalJSON(p []byte) error {
	return UnmarshalJSON(fruitEnum, p, (*uint)(f))
}

type Mode uint

const (
	Minvalid Mode = iota
	Mgo
	Mstop
	Mpotato
)

var modeEnumDefault = UintEnumDefault(uint(Minvalid), []string{
	Minvalid: "invalid",
	Mgo:      "go",
	Mstop:    "stop",
	Mpotato:  "potato",
})

func ParseMode(str string) Mode {
	return Mode(modeEnumDefault.Parse(str).(uint))
}

func (m Mode) String() string {
	return StringDefault(modeEnumDefault, uint(m))
}

func (m Mode) MarshalJSON() ([]byte, error) {
	return MarshalJSONDefault(modeEnumDefault, uint(m))
}

func (m *Mode) UnmarshalJSON(p []byte) error {
	return UnmarshalJSONDefault(modeEnumDefault, p, (*uint)(m))
}

func TestUintEnumString(t *testing.T) {
	assert.Equal(t, "naners", Banana.String())
	assert.Equal(t, "wut?", Fruit(5).String())

	assert.Equal(t, "go", Mgo.String())
	assert.Equal(t, "potato", Mpotato.String())
	assert.Equal(t, "invalid", Minvalid.String())
	assert.Equal(t, "invalid", Mode(1000).String())
}

func TestUintEnumParse(t *testing.T) {
	f, err := ParseFruit("tots")
	assert.Equal(t, nil, err)
	assert.Equal(t, Potato, f)

	_, err = ParseFruit("pants")
	assert.NotEqual(t, nil, err)

	assert.Equal(t, Mstop, ParseMode("stop"))
	assert.Equal(t, Minvalid, ParseMode("poppycock"))
}

func TestUintEnumMarshalJSON(t *testing.T) {
	p, err := Potato.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"tots"`, string(p))

	p, err = Fruit(100).MarshalJSON()
	assert.NotEqual(t, nil, err)

	p, err = Mpotato.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"potato"`, string(p))

	p, err = Minvalid.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"invalid"`, string(p))
}

func TestUintEnumUnmarshalJSON(t *testing.T) {
	var f Fruit

	err := f.UnmarshalJSON([]byte(`"apps"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, Apple, f)

	err = f.UnmarshalJSON([]byte(`"blah"`))
	assert.NotEqual(t, nil, err)

	err = f.UnmarshalJSON([]byte(`123`))
	assert.NotEqual(t, nil, err)

	var m Mode
	err = m.UnmarshalJSON([]byte(`"potato"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, Mpotato, m)

	err = m.UnmarshalJSON([]byte(`"blast off"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, Minvalid, m)

	// XXX is this the right behavior
	err = m.UnmarshalJSON([]byte(`123`))
	assert.NotEqual(t, nil, err)
}
