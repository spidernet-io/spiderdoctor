// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// EchoRes echo request
//
// swagger:model EchoRes
type EchoRes struct {

	// client source ip
	ClientIP string `json:"clientIp,omitempty"`

	// other  information
	OtherDetail map[string]string `json:"otherDetail,omitempty"`

	// request header
	RequestHeader map[string]string `json:"requestHeader,omitempty"`

	// request url
	RequestURL string `json:"requestUrl,omitempty"`

	// server host name
	ServerName string `json:"serverName,omitempty"`
}

// Validate validates this echo res
func (m *EchoRes) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this echo res based on context it is used
func (m *EchoRes) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *EchoRes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EchoRes) UnmarshalBinary(b []byte) error {
	var res EchoRes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
