package gridworker

import (
	"encoding/json"
)

// SetMap sets a map for a DynamicType object
func (d *DynamicType) SetMap(m map[string]interface{}) {
	b, _ := json.Marshal(m)
	d.Bytes = b
	d.SetType(DynamicTypeType_json_map)
}

// GetMap gets a map from a DynamicType object
func (d *DynamicType) GetMap() map[string]interface{} {
	b := d.GetBytes()
	m := map[string]interface{}{}
	json.Unmarshal(b, m)

	return m
}

// SetSlice sets a slice for a DynamicType object
func (d *DynamicType) SetSlice(m []interface{}) {
	b, _ := json.Marshal(m)
	d.Bytes = b
	d.SetType(DynamicTypeType_json_slice)
}

// GetSlice extracts a slice out of a DynamicType object
func (d *DynamicType) GetSlice() []interface{} {
	b := d.GetBytes()
	m := []interface{}{}
	json.Unmarshal(b, m)

	return m
}

// SetType sets a DynamicTypeObject's type
func (d *DynamicType) SetType(t DynamicTypeType) {
	d.Type = &t
}
