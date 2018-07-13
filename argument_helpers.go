/*
 *   This file is part of GridWorker.
 *
 *   Copyright (c) 2018 Mocha Industries, LLC.
 *   All rights reserved.
 *
 *   GridWorker is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   GridWorker is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with GridWorker.  If not, see <https://www.gnu.org/licenses/>.
 */

package gridworker

import (
	"github.com/golang/protobuf/proto"
)

//
// Message Getters
//

// getInterface gets an object from the argument map
func (m *Message) getInterface(key string) *DynamicType {
	return m.arguments[key]
}

// GetString returns a string for a specific key in the argument map
func (m *Message) GetString(key string) string {
	return m.getInterface(key).GetStr()
}

// GetInt64 returns a int64 for a specific key in the argument map
func (m *Message) GetInt64(key string) int64 {
	return m.getInterface(key).GetInt()
}

// GetFloat64 returns a float64 for a specific key in the argument map
func (m *Message) GetFloat64(key string) float64 {
	return m.getInterface(key).GetFloat()
}

// GetBool returns a bool for a specific key in the argument map
func (m *Message) GetBool(key string) bool {
	return m.getInterface(key).GetBool()
}

// GetBytes returns a byte for a specific key in the argument map
func (m *Message) GetBytes(key string) []byte {
	return m.getInterface(key).GetBytes()
}

// GetMap returns a map for a specific key in the argument map
func (m *Message) GetMap(key string) map[string]interface{} {
	return m.getInterface(key).GetMap()
}

// GetSlice returns a slice for a specific key in the argument map
func (m *Message) GetSlice(key string) []interface{} {
	return m.getInterface(key).GetSlice()
}

//
// Setters
//

// setInterface sets an object in the argument map
func (m *Message) setInterface(key string, value *DynamicType) {
	m.arguments[key] = value
}

// SetString sets a string in the argument map
func (m *Message) SetString(key string, value string) {
	v := DynamicType{}
	v.Str = proto.String(value)
	v.SetType(DynamicTypeType_string)
	m.setInterface(key, &v)
}

// SetInt64 sets a int64 in the argument map
func (m *Message) SetInt64(key string, value int64) {
	v := DynamicType{}
	v.Int = proto.Int64(value)
	v.SetType(DynamicTypeType_int)
	m.setInterface(key, &v)
}

// SetFloat64 sets a float64 in the argument map
func (m *Message) SetFloat64(key string, value float64) {
	v := DynamicType{}
	v.Float = proto.Float64(value)
	v.SetType(DynamicTypeType_float)
	m.setInterface(key, &v)
}

// SetBool sets a bool in the argument map
func (m *Message) SetBool(key string, value bool) {
	v := DynamicType{}
	v.Bool = proto.Bool(value)
	v.SetType(DynamicTypeType_bool)
	m.setInterface(key, &v)
}

// SetBytes sets a bytes slice in the argument map
func (m *Message) SetBytes(key string, value []byte) {
	v := DynamicType{}
	v.Bytes = value
	v.SetType(DynamicTypeType_bytes)
	m.setInterface(key, &v)
}

// SetMap sets a map in the argument map
func (m *Message) SetMap(key string, value map[string]interface{}) {
	v := DynamicType{}
	v.SetMap(value)
	m.setInterface(key, &v)
}

// SetSlice sets a slice in the argument map
func (m *Message) SetSlice(key string, value []interface{}) {
	v := DynamicType{}
	v.SetSlice(value)
	m.setInterface(key, &v)
}

//
// Context Helpers
//

// GetString returns a string for a specific key in the context input map
func (c *Context) GetString(key string) string {
	return c.input.GetString(key)
}

// GetInt64 returns a int64 for a specific key in the context input map
func (c *Context) GetInt64(key string) int64 {
	return c.input.GetInt64(key)
}

// GetFloat64 returns a float64 for a specific key in the context input map
func (c *Context) GetFloat64(key string) float64 {
	return c.input.GetFloat64(key)
}

// GetBool returns a bool for a specific key in the context input map
func (c *Context) GetBool(key string) bool {
	return c.input.GetBool(key)
}

// GetBytes returns a byte slice for a specific key in the context input map
func (c *Context) GetBytes(key string) []byte {
	return c.input.GetBytes(key)
}

// GetMap returns a map for a specific key in the context input map
func (c *Context) GetMap(key string) map[string]interface{} {
	return c.input.GetMap(key)
}

// GetSlice returns a slice for a specific key in the context input map
func (c *Context) GetSlice(key string) []interface{} {
	return c.input.GetSlice(key)
}
