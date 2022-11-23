// Code generated by go-swagger; DO NOT EDIT.

// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package healthy

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetHealthyReadinessHandlerFunc turns a function with the right signature into a get healthy readiness handler
type GetHealthyReadinessHandlerFunc func(GetHealthyReadinessParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetHealthyReadinessHandlerFunc) Handle(params GetHealthyReadinessParams) middleware.Responder {
	return fn(params)
}

// GetHealthyReadinessHandler interface for that can handle valid get healthy readiness params
type GetHealthyReadinessHandler interface {
	Handle(GetHealthyReadinessParams) middleware.Responder
}

// NewGetHealthyReadiness creates a new http.Handler for the get healthy readiness operation
func NewGetHealthyReadiness(ctx *middleware.Context, handler GetHealthyReadinessHandler) *GetHealthyReadiness {
	return &GetHealthyReadiness{Context: ctx, Handler: handler}
}

/*
	GetHealthyReadiness swagger:route GET /healthy/readiness healthy getHealthyReadiness

# Readiness probe

pod readiness probe for agent and controller pod
*/
type GetHealthyReadiness struct {
	Context *middleware.Context
	Handler GetHealthyReadinessHandler
}

func (o *GetHealthyReadiness) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetHealthyReadinessParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}