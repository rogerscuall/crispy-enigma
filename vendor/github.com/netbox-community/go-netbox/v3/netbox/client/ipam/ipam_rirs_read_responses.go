// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2020 The go-netbox Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ipam

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netbox-community/go-netbox/v3/netbox/models"
)

// IpamRirsReadReader is a Reader for the IpamRirsRead structure.
type IpamRirsReadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *IpamRirsReadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewIpamRirsReadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewIpamRirsReadDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewIpamRirsReadOK creates a IpamRirsReadOK with default headers values
func NewIpamRirsReadOK() *IpamRirsReadOK {
	return &IpamRirsReadOK{}
}

/*
IpamRirsReadOK describes a response with status code 200, with default header values.

IpamRirsReadOK ipam rirs read o k
*/
type IpamRirsReadOK struct {
	Payload *models.RIR
}

// IsSuccess returns true when this ipam rirs read o k response has a 2xx status code
func (o *IpamRirsReadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ipam rirs read o k response has a 3xx status code
func (o *IpamRirsReadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ipam rirs read o k response has a 4xx status code
func (o *IpamRirsReadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this ipam rirs read o k response has a 5xx status code
func (o *IpamRirsReadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this ipam rirs read o k response a status code equal to that given
func (o *IpamRirsReadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the ipam rirs read o k response
func (o *IpamRirsReadOK) Code() int {
	return 200
}

func (o *IpamRirsReadOK) Error() string {
	return fmt.Sprintf("[GET /ipam/rirs/{id}/][%d] ipamRirsReadOK  %+v", 200, o.Payload)
}

func (o *IpamRirsReadOK) String() string {
	return fmt.Sprintf("[GET /ipam/rirs/{id}/][%d] ipamRirsReadOK  %+v", 200, o.Payload)
}

func (o *IpamRirsReadOK) GetPayload() *models.RIR {
	return o.Payload
}

func (o *IpamRirsReadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RIR)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewIpamRirsReadDefault creates a IpamRirsReadDefault with default headers values
func NewIpamRirsReadDefault(code int) *IpamRirsReadDefault {
	return &IpamRirsReadDefault{
		_statusCode: code,
	}
}

/*
IpamRirsReadDefault describes a response with status code -1, with default header values.

IpamRirsReadDefault ipam rirs read default
*/
type IpamRirsReadDefault struct {
	_statusCode int

	Payload interface{}
}

// IsSuccess returns true when this ipam rirs read default response has a 2xx status code
func (o *IpamRirsReadDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ipam rirs read default response has a 3xx status code
func (o *IpamRirsReadDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ipam rirs read default response has a 4xx status code
func (o *IpamRirsReadDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ipam rirs read default response has a 5xx status code
func (o *IpamRirsReadDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ipam rirs read default response a status code equal to that given
func (o *IpamRirsReadDefault) IsCode(code int) bool {
	return o._statusCode == code
}

// Code gets the status code for the ipam rirs read default response
func (o *IpamRirsReadDefault) Code() int {
	return o._statusCode
}

func (o *IpamRirsReadDefault) Error() string {
	return fmt.Sprintf("[GET /ipam/rirs/{id}/][%d] ipam_rirs_read default  %+v", o._statusCode, o.Payload)
}

func (o *IpamRirsReadDefault) String() string {
	return fmt.Sprintf("[GET /ipam/rirs/{id}/][%d] ipam_rirs_read default  %+v", o._statusCode, o.Payload)
}

func (o *IpamRirsReadDefault) GetPayload() interface{} {
	return o.Payload
}

func (o *IpamRirsReadDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}