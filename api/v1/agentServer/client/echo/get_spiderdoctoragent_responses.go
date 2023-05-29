// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package echo

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/spidernet-io/spiderdoctor/api/v1/agentServer/models"
)

// GetSpiderdoctoragentReader is a Reader for the GetSpiderdoctoragent structure.
type GetSpiderdoctoragentReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSpiderdoctoragentReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetSpiderdoctoragentOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetSpiderdoctoragentOK creates a GetSpiderdoctoragentOK with default headers values
func NewGetSpiderdoctoragentOK() *GetSpiderdoctoragentOK {
	return &GetSpiderdoctoragentOK{}
}

/*
GetSpiderdoctoragentOK describes a response with status code 200, with default header values.

Success
*/
type GetSpiderdoctoragentOK struct {
	Payload *models.EchoRes
}

// IsSuccess returns true when this get spiderdoctoragent o k response has a 2xx status code
func (o *GetSpiderdoctoragentOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get spiderdoctoragent o k response has a 3xx status code
func (o *GetSpiderdoctoragentOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get spiderdoctoragent o k response has a 4xx status code
func (o *GetSpiderdoctoragentOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get spiderdoctoragent o k response has a 5xx status code
func (o *GetSpiderdoctoragentOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get spiderdoctoragent o k response a status code equal to that given
func (o *GetSpiderdoctoragentOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get spiderdoctoragent o k response
func (o *GetSpiderdoctoragentOK) Code() int {
	return 200
}

func (o *GetSpiderdoctoragentOK) Error() string {
	return fmt.Sprintf("[GET /spiderdoctoragent][%d] getSpiderdoctoragentOK  %+v", 200, o.Payload)
}

func (o *GetSpiderdoctoragentOK) String() string {
	return fmt.Sprintf("[GET /spiderdoctoragent][%d] getSpiderdoctoragentOK  %+v", 200, o.Payload)
}

func (o *GetSpiderdoctoragentOK) GetPayload() *models.EchoRes {
	return o.Payload
}

func (o *GetSpiderdoctoragentOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.EchoRes)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
