package entry

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
)

// Coder is a interface for encoding and decoding data.
// It should be implemented by the specific game mode player, group, room, team, etc.
//
// Note: Do not implement this interface for the common entries or strategy entries.
type Coder interface {
	// Encode encodes current object to a byte array
	Encode() ([]byte, error)
	// Decode decodes a byte array to current object
	Decode(data []byte) error
}

// Encode encodes current object to a byte array with gob
func Encode(v any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode decodes a byte array to current object
func Decode(data []byte, v any) error {
	// v must be a pointer
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("v must be a pointer")
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}
