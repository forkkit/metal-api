// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NetworkAllocateRequest network allocate request
// swagger:model NetworkAllocateRequest
type NetworkAllocateRequest struct {

	// The description of this prefix
	Description string `json:"description,omitempty"`

	// The id of the partition
	Partition string `json:"partition,omitempty"`

	// The prefix to split out a child from
	// Required: true
	// Min Length: 1
	Prefix *string `json:"prefix"`

	// The length of the child prefix to split out
	// Maximum: 32
	// Minimum: 0
	PrefixLength *float64 `json:"prefix_length,omitempty"`

	// The name of the project to assign this prefix to
	Project string `json:"project,omitempty"`

	// The name of the tenant to assign this prefix to
	Tenant string `json:"tenant,omitempty"`
}

// Validate validates this network allocate request
func (m *NetworkAllocateRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePrefix(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePrefixLength(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkAllocateRequest) validatePrefix(formats strfmt.Registry) error {

	if err := validate.Required("prefix", "body", m.Prefix); err != nil {
		return err
	}

	if err := validate.MinLength("prefix", "body", string(*m.Prefix), 1); err != nil {
		return err
	}

	return nil
}

func (m *NetworkAllocateRequest) validatePrefixLength(formats strfmt.Registry) error {

	if swag.IsZero(m.PrefixLength) { // not required
		return nil
	}

	if err := validate.Minimum("prefix_length", "body", float64(*m.PrefixLength), 0, false); err != nil {
		return err
	}

	if err := validate.Maximum("prefix_length", "body", float64(*m.PrefixLength), 32, false); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *NetworkAllocateRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetworkAllocateRequest) UnmarshalBinary(b []byte) error {
	var res NetworkAllocateRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}