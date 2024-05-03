// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Friend friend
//
// swagger:model Friend
type Friend struct {

	// email
	// Example: user123@mail.ru
	Email string `json:"email,omitempty"`

	// login
	// Example: user_login
	Login string `json:"login,omitempty"`

	// name
	// Example: user123
	Name string `json:"name,omitempty"`

	// user id
	// Example: 1
	UserID int64 `json:"user_id,omitempty"`
}

// Validate validates this friend
func (m *Friend) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this friend based on context it is used
func (m *Friend) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Friend) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Friend) UnmarshalBinary(b []byte) error {
	var res Friend
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}