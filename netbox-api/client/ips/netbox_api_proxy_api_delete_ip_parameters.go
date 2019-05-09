// Code generated by go-swagger; DO NOT EDIT.

package ips

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "git.f-i-ts.de/cloud-native/metal/metal-api/netbox-api/models"
)

// NewNetboxAPIProxyAPIDeleteIPParams creates a new NetboxAPIProxyAPIDeleteIPParams object
// with the default values initialized.
func NewNetboxAPIProxyAPIDeleteIPParams() *NetboxAPIProxyAPIDeleteIPParams {
	var ()
	return &NetboxAPIProxyAPIDeleteIPParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewNetboxAPIProxyAPIDeleteIPParamsWithTimeout creates a new NetboxAPIProxyAPIDeleteIPParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewNetboxAPIProxyAPIDeleteIPParamsWithTimeout(timeout time.Duration) *NetboxAPIProxyAPIDeleteIPParams {
	var ()
	return &NetboxAPIProxyAPIDeleteIPParams{

		timeout: timeout,
	}
}

// NewNetboxAPIProxyAPIDeleteIPParamsWithContext creates a new NetboxAPIProxyAPIDeleteIPParams object
// with the default values initialized, and the ability to set a context for a request
func NewNetboxAPIProxyAPIDeleteIPParamsWithContext(ctx context.Context) *NetboxAPIProxyAPIDeleteIPParams {
	var ()
	return &NetboxAPIProxyAPIDeleteIPParams{

		Context: ctx,
	}
}

// NewNetboxAPIProxyAPIDeleteIPParamsWithHTTPClient creates a new NetboxAPIProxyAPIDeleteIPParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewNetboxAPIProxyAPIDeleteIPParamsWithHTTPClient(client *http.Client) *NetboxAPIProxyAPIDeleteIPParams {
	var ()
	return &NetboxAPIProxyAPIDeleteIPParams{
		HTTPClient: client,
	}
}

/*NetboxAPIProxyAPIDeleteIPParams contains all the parameters to send to the API endpoint
for the netbox api proxy api delete ip operation typically these are written to a http.Request
*/
type NetboxAPIProxyAPIDeleteIPParams struct {

	/*Request
	  The ip deletion body

	*/
	Request *models.CIDR

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) WithTimeout(timeout time.Duration) *NetboxAPIProxyAPIDeleteIPParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) WithContext(ctx context.Context) *NetboxAPIProxyAPIDeleteIPParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) WithHTTPClient(client *http.Client) *NetboxAPIProxyAPIDeleteIPParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithRequest adds the request to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) WithRequest(request *models.CIDR) *NetboxAPIProxyAPIDeleteIPParams {
	o.SetRequest(request)
	return o
}

// SetRequest adds the request to the netbox api proxy api delete ip params
func (o *NetboxAPIProxyAPIDeleteIPParams) SetRequest(request *models.CIDR) {
	o.Request = request
}

// WriteToRequest writes these params to a swagger request
func (o *NetboxAPIProxyAPIDeleteIPParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Request != nil {
		if err := r.SetBodyParam(o.Request); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}