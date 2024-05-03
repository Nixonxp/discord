// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UpdateUserRequest update user request
//
// swagger:model UpdateUserRequest
type UpdateUserRequest struct {

	// avatar photo url
	AvatarPhotoURL string `json:"avatar_photo_url,omitempty"`

	// email
	// Example: user123@mail.ru
	Email string `json:"email,omitempty"`

	// login
	// Example: user_login
	Login string `json:"login,omitempty"`

	// name
	// Example: user123
	Name string `json:"name,omitempty"`

	// password
	Password string `json:"password,omitempty"`
}

// Validate validates this update user request
func (m *UpdateUserRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this update user request based on context it is used
func (m *UpdateUserRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UpdateUserRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateUserRequest) UnmarshalBinary(b []byte) error {
	var res UpdateUserRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
