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

package dcim

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewDcimPlatformsDeleteParams creates a new DcimPlatformsDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewDcimPlatformsDeleteParams() *DcimPlatformsDeleteParams {
	return &DcimPlatformsDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewDcimPlatformsDeleteParamsWithTimeout creates a new DcimPlatformsDeleteParams object
// with the ability to set a timeout on a request.
func NewDcimPlatformsDeleteParamsWithTimeout(timeout time.Duration) *DcimPlatformsDeleteParams {
	return &DcimPlatformsDeleteParams{
		timeout: timeout,
	}
}

// NewDcimPlatformsDeleteParamsWithContext creates a new DcimPlatformsDeleteParams object
// with the ability to set a context for a request.
func NewDcimPlatformsDeleteParamsWithContext(ctx context.Context) *DcimPlatformsDeleteParams {
	return &DcimPlatformsDeleteParams{
		Context: ctx,
	}
}

// NewDcimPlatformsDeleteParamsWithHTTPClient creates a new DcimPlatformsDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewDcimPlatformsDeleteParamsWithHTTPClient(client *http.Client) *DcimPlatformsDeleteParams {
	return &DcimPlatformsDeleteParams{
		HTTPClient: client,
	}
}

/*
DcimPlatformsDeleteParams contains all the parameters to send to the API endpoint

	for the dcim platforms delete operation.

	Typically these are written to a http.Request.
*/
type DcimPlatformsDeleteParams struct {

	/* ID.

	   A unique integer value identifying this platform.
	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the dcim platforms delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DcimPlatformsDeleteParams) WithDefaults() *DcimPlatformsDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the dcim platforms delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *DcimPlatformsDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) WithTimeout(timeout time.Duration) *DcimPlatformsDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) WithContext(ctx context.Context) *DcimPlatformsDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) WithHTTPClient(client *http.Client) *DcimPlatformsDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) WithID(id int64) *DcimPlatformsDeleteParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the dcim platforms delete params
func (o *DcimPlatformsDeleteParams) SetID(id int64) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *DcimPlatformsDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}