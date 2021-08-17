// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: telemetry.proto

package telemetry

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Telemetry service

type TelemetryService interface {
	GetTelemetryService(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricDefinitionCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricReportDefinitionCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricReportCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetTriggerCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricDefinition(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricReportDefinition(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetMetricReport(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	GetTrigger(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
	UpdateTrigger(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error)
}

type telemetryService struct {
	c    client.Client
	name string
}

func NewTelemetryService(name string, c client.Client) TelemetryService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "telemetry"
	}
	return &telemetryService{
		c:    c,
		name: name,
	}
}

func (c *telemetryService) GetTelemetryService(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetTelemetryService", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricDefinitionCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricDefinitionCollection", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricReportDefinitionCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricReportDefinitionCollection", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricReportCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricReportCollection", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetTriggerCollection(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetTriggerCollection", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricDefinition(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricDefinition", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricReportDefinition(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricReportDefinition", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetMetricReport(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetMetricReport", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) GetTrigger(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.GetTrigger", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryService) UpdateTrigger(ctx context.Context, in *TelemetryRequest, opts ...client.CallOption) (*TelemetryResponse, error) {
	req := c.c.NewRequest(c.name, "Telemetry.UpdateTrigger", in)
	out := new(TelemetryResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Telemetry service

type TelemetryHandler interface {
	GetTelemetryService(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricDefinitionCollection(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricReportDefinitionCollection(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricReportCollection(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetTriggerCollection(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricDefinition(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricReportDefinition(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetMetricReport(context.Context, *TelemetryRequest, *TelemetryResponse) error
	GetTrigger(context.Context, *TelemetryRequest, *TelemetryResponse) error
	UpdateTrigger(context.Context, *TelemetryRequest, *TelemetryResponse) error
}

func RegisterTelemetryHandler(s server.Server, hdlr TelemetryHandler, opts ...server.HandlerOption) error {
	type telemetry interface {
		GetTelemetryService(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricDefinitionCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricReportDefinitionCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricReportCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetTriggerCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricDefinition(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricReportDefinition(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetMetricReport(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		GetTrigger(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
		UpdateTrigger(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error
	}
	type Telemetry struct {
		telemetry
	}
	h := &telemetryHandler{hdlr}
	return s.Handle(s.NewHandler(&Telemetry{h}, opts...))
}

type telemetryHandler struct {
	TelemetryHandler
}

func (h *telemetryHandler) GetTelemetryService(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetTelemetryService(ctx, in, out)
}

func (h *telemetryHandler) GetMetricDefinitionCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricDefinitionCollection(ctx, in, out)
}

func (h *telemetryHandler) GetMetricReportDefinitionCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricReportDefinitionCollection(ctx, in, out)
}

func (h *telemetryHandler) GetMetricReportCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricReportCollection(ctx, in, out)
}

func (h *telemetryHandler) GetTriggerCollection(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetTriggerCollection(ctx, in, out)
}

func (h *telemetryHandler) GetMetricDefinition(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricDefinition(ctx, in, out)
}

func (h *telemetryHandler) GetMetricReportDefinition(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricReportDefinition(ctx, in, out)
}

func (h *telemetryHandler) GetMetricReport(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetMetricReport(ctx, in, out)
}

func (h *telemetryHandler) GetTrigger(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.GetTrigger(ctx, in, out)
}

func (h *telemetryHandler) UpdateTrigger(ctx context.Context, in *TelemetryRequest, out *TelemetryResponse) error {
	return h.TelemetryHandler.UpdateTrigger(ctx, in, out)
}
