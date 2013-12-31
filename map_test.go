// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// map_test.go [created: Mon, 30 Dec 2013]

package enumstr

import (
	"fmt"
	"github.com/bmizerany/assert"
	"testing"
)

func TestEnumMap(t *testing.T) {
	var recv interface{}
	defer func() { assert.Equal(t, nil, recv) }()
	defer func() { recv = recover() }()
	EnumMap(Map{
		0: "abc",
		1: "def",
	})
}

func TestEnumMapDuplicate(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	EnumMap(Map{
		0: "abc",
		1: "abc",
	})
}

func TestEnumMapDefault(t *testing.T) {
	var recv interface{}
	defer func() { assert.Equal(t, nil, recv) }()
	defer func() { recv = recover() }()
	EnumMapDefault(0, Map{
		0: "abc",
		1: "def",
	})
}

func TestEnumMapDefaultDuplicate(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	EnumMapDefault(0, Map{
		0: "abc",
		1: "abc",
	})
}

func TestEnumMapDefaultNotFound(t *testing.T) {
	var recv interface{}
	defer func() { assert.NotEqual(t, nil, recv) }()
	defer func() { recv = recover() }()
	EnumMapDefault(2, Map{
		0: "abc",
		1: "def",
	})
}

type REST struct {
	host     string
	method   string
	resource string
}

var puppyCreate = REST{"http://example.com", "POST", "/puppy"}
var puppyDelete = REST{"http://example.com", "DELETE", "/puppy/:id"}
var puppyShow = REST{"http://example.com", "GET", "/puppy/:id"}
var puppyUpdate = REST{"http://example.com", "PUT", "/puppy/:id"}

var puppyOpEnum = EnumMap(Map{
	puppyCreate: "Puppy CREATE",
	puppyDelete: "Puppy DELETE",
	puppyShow:   "Puppy SHOW",
	puppyUpdate: "Puppy UPDATE",
})

func ParseREST(str string) (REST, error) {
	enum, err := puppyOpEnum.Parse(str)
	if err == nil {
		return enum.(REST), nil
	}
	return REST{}, fmt.Errorf("invalid operation: %#v", str)
}

func (f REST) String() string {
	return String(puppyOpEnum, f, "REST(???)")
}

func (f REST) MarshalJSON() ([]byte, error) {
	return MarshalJSON(puppyOpEnum, f)
}

func (f *REST) UnmarshalJSON(p []byte) error {
	return UnmarshalJSON(puppyOpEnum, p, f)
}

func TestEnumMapString(t *testing.T) {
	assert.Equal(t, "Puppy CREATE", puppyCreate.String())
	assert.Equal(t, "REST(???)", REST{}.String())
}

func TestEnumMapParse(t *testing.T) {
	op, err := ParseREST("Puppy DELETE")
	assert.Equal(t, nil, err)
	assert.Equal(t, puppyDelete, op)

	_, err = ParseREST("hong kong FOOEY")
	assert.NotEqual(t, nil, err)
}

func TestEnumMapMarshalJSON(t *testing.T) {
	p, err := puppyUpdate.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"Puppy UPDATE"`, string(p))

	p, err = REST{}.MarshalJSON()
	assert.NotEqual(t, nil, err)
}

func TestEnumMapUnmarshalJSON(t *testing.T) {
	var op REST

	err := op.UnmarshalJSON([]byte(`"Puppy SHOW"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, puppyShow, op)

	err = op.UnmarshalJSON([]byte(`"blah"`))
	assert.NotEqual(t, nil, err)
}

type MIME struct {
	Name string
	Ext  string
}

var mimeInvalid = MIME{"invalid/mime", ""}
var mimeJSON = MIME{"application/json", ".json"}
var mimeCSV = MIME{"text/csv", ".csv"}
var mimeXML = MIME{"text/xml", ".xml"}

var mimeEnum = EnumMapDefault(mimeInvalid, Map{
	mimeInvalid: "INVALID",
	mimeJSON:    "JSON",
	mimeCSV:     "CSV",
	mimeXML:     "XML",
})

func ParseMIME(str string) MIME {
	return mimeEnum.Parse(str).(MIME)
}

func (f MIME) String() string {
	return StringDefault(mimeEnum, f)
}

func (f MIME) MarshalJSON() ([]byte, error) {
	return MarshalJSONDefault(mimeEnum, f)
}

func (f *MIME) UnmarshalJSON(p []byte) error {
	return UnmarshalJSONDefault(mimeEnum, p, f)
}

func TestEnumMapDefaultString(t *testing.T) {
	assert.Equal(t, "JSON", mimeJSON.String())
	assert.Equal(t, "INVALID", new(MIME).String())
}

func TestEnumMapDefaultParse(t *testing.T) {
	assert.Equal(t, mimeJSON, ParseMIME("JSON"))

	assert.Equal(t, mimeInvalid, ParseMIME("BLAH"))
}

func TestEnumMapDefaultMarshalJSON(t *testing.T) {
	p, err := mimeXML.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"XML"`, string(p))

	p, err = MIME{}.MarshalJSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, `"INVALID"`, string(p))
}

func TestEnumMapDefaultUnmarshalJSON(t *testing.T) {
	var m MIME

	err := m.UnmarshalJSON([]byte(`"CSV"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, mimeCSV, m)

	err = m.UnmarshalJSON([]byte(`"blah"`))
	assert.Equal(t, nil, err)
	assert.Equal(t, mimeInvalid, m)
}
