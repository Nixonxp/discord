// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CreateServerResponse create server response
//
// swagger:model CreateServerResponse
type CreateServerResponse struct {

	// id
	// Example: 1
	ID int64 `json:"id,omitempty"`

	// name
	// Example: server name
	Name string `json:"name,omitempty"`
}

// Validate validates this create server response
func (m *CreateServerResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create server response based on context it is used
func (m *CreateServerResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreateServerResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateServerResponse) UnmarshalBinary(b []byte) error {
	var res CreateServerResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
