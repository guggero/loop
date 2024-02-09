// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package swapserverrpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// InstantSwapServerClient is the client API for InstantSwapServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InstantSwapServerClient interface {
	// RequestInstantLoopOut initiates an instant loop out swap.
	RequestInstantLoopOut(ctx context.Context, in *InstantLoopOutRequest, opts ...grpc.CallOption) (*InstantLoopOutResponse, error)
	// PollPaymentAccepted polls the server to see if the payment has been
	// accepted.
	PollPaymentAccepted(ctx context.Context, in *PollPaymentAcceptedRequest, opts ...grpc.CallOption) (*PollPaymentAcceptedResponse, error)
	// InitHtlcSig is called by the client to initiate the htlc sig exchange.
	InitHtlcSig(ctx context.Context, in *InitHtlcSigRequest, opts ...grpc.CallOption) (*InitHtlcSigResponse, error)
	// PushHtlcSig is called by the client to push the htlc sigs to the server.
	PushHtlcSig(ctx context.Context, in *PushHtlcSigRequest, opts ...grpc.CallOption) (*PushHtlcSigResponse, error)
	// PushPreimage is called by the client to push the preimage to the server.
	// This returns the musig2 signatures that the client needs to sweep the
	// htlc.
	PushPreimage(ctx context.Context, in *PushPreimageRequest, opts ...grpc.CallOption) (*PushPreimageResponse, error)
	// CancelInstantSwap tries to cancel the instant swap. This can only be
	// called if the swap has not been accepted yet.
	CancelInstantSwap(ctx context.Context, in *CancelInstantSwapRequest, opts ...grpc.CallOption) (*CancelInstantSwapResponse, error)
	// GetInstantOutQuote returns the absolute fee in satoshis for the swap and
	// the pubkey to query the route to estimate offchain payment fees.
	GetInstantOutQuote(ctx context.Context, in *GetInstantOutQuoteRequest, opts ...grpc.CallOption) (*GetInstantOutQuoteResponse, error)
}

type instantSwapServerClient struct {
	cc grpc.ClientConnInterface
}

func NewInstantSwapServerClient(cc grpc.ClientConnInterface) InstantSwapServerClient {
	return &instantSwapServerClient{cc}
}

func (c *instantSwapServerClient) RequestInstantLoopOut(ctx context.Context, in *InstantLoopOutRequest, opts ...grpc.CallOption) (*InstantLoopOutResponse, error) {
	out := new(InstantLoopOutResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/RequestInstantLoopOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) PollPaymentAccepted(ctx context.Context, in *PollPaymentAcceptedRequest, opts ...grpc.CallOption) (*PollPaymentAcceptedResponse, error) {
	out := new(PollPaymentAcceptedResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/PollPaymentAccepted", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) InitHtlcSig(ctx context.Context, in *InitHtlcSigRequest, opts ...grpc.CallOption) (*InitHtlcSigResponse, error) {
	out := new(InitHtlcSigResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/InitHtlcSig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) PushHtlcSig(ctx context.Context, in *PushHtlcSigRequest, opts ...grpc.CallOption) (*PushHtlcSigResponse, error) {
	out := new(PushHtlcSigResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/PushHtlcSig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) PushPreimage(ctx context.Context, in *PushPreimageRequest, opts ...grpc.CallOption) (*PushPreimageResponse, error) {
	out := new(PushPreimageResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/PushPreimage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) CancelInstantSwap(ctx context.Context, in *CancelInstantSwapRequest, opts ...grpc.CallOption) (*CancelInstantSwapResponse, error) {
	out := new(CancelInstantSwapResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/CancelInstantSwap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *instantSwapServerClient) GetInstantOutQuote(ctx context.Context, in *GetInstantOutQuoteRequest, opts ...grpc.CallOption) (*GetInstantOutQuoteResponse, error) {
	out := new(GetInstantOutQuoteResponse)
	err := c.cc.Invoke(ctx, "/looprpc.InstantSwapServer/GetInstantOutQuote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InstantSwapServerServer is the server API for InstantSwapServer service.
// All implementations must embed UnimplementedInstantSwapServerServer
// for forward compatibility
type InstantSwapServerServer interface {
	// RequestInstantLoopOut initiates an instant loop out swap.
	RequestInstantLoopOut(context.Context, *InstantLoopOutRequest) (*InstantLoopOutResponse, error)
	// PollPaymentAccepted polls the server to see if the payment has been
	// accepted.
	PollPaymentAccepted(context.Context, *PollPaymentAcceptedRequest) (*PollPaymentAcceptedResponse, error)
	// InitHtlcSig is called by the client to initiate the htlc sig exchange.
	InitHtlcSig(context.Context, *InitHtlcSigRequest) (*InitHtlcSigResponse, error)
	// PushHtlcSig is called by the client to push the htlc sigs to the server.
	PushHtlcSig(context.Context, *PushHtlcSigRequest) (*PushHtlcSigResponse, error)
	// PushPreimage is called by the client to push the preimage to the server.
	// This returns the musig2 signatures that the client needs to sweep the
	// htlc.
	PushPreimage(context.Context, *PushPreimageRequest) (*PushPreimageResponse, error)
	// CancelInstantSwap tries to cancel the instant swap. This can only be
	// called if the swap has not been accepted yet.
	CancelInstantSwap(context.Context, *CancelInstantSwapRequest) (*CancelInstantSwapResponse, error)
	// GetInstantOutQuote returns the absolute fee in satoshis for the swap and
	// the pubkey to query the route to estimate offchain payment fees.
	GetInstantOutQuote(context.Context, *GetInstantOutQuoteRequest) (*GetInstantOutQuoteResponse, error)
	mustEmbedUnimplementedInstantSwapServerServer()
}

// UnimplementedInstantSwapServerServer must be embedded to have forward compatible implementations.
type UnimplementedInstantSwapServerServer struct {
}

func (UnimplementedInstantSwapServerServer) RequestInstantLoopOut(context.Context, *InstantLoopOutRequest) (*InstantLoopOutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestInstantLoopOut not implemented")
}
func (UnimplementedInstantSwapServerServer) PollPaymentAccepted(context.Context, *PollPaymentAcceptedRequest) (*PollPaymentAcceptedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PollPaymentAccepted not implemented")
}
func (UnimplementedInstantSwapServerServer) InitHtlcSig(context.Context, *InitHtlcSigRequest) (*InitHtlcSigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitHtlcSig not implemented")
}
func (UnimplementedInstantSwapServerServer) PushHtlcSig(context.Context, *PushHtlcSigRequest) (*PushHtlcSigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushHtlcSig not implemented")
}
func (UnimplementedInstantSwapServerServer) PushPreimage(context.Context, *PushPreimageRequest) (*PushPreimageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushPreimage not implemented")
}
func (UnimplementedInstantSwapServerServer) CancelInstantSwap(context.Context, *CancelInstantSwapRequest) (*CancelInstantSwapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelInstantSwap not implemented")
}
func (UnimplementedInstantSwapServerServer) GetInstantOutQuote(context.Context, *GetInstantOutQuoteRequest) (*GetInstantOutQuoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInstantOutQuote not implemented")
}
func (UnimplementedInstantSwapServerServer) mustEmbedUnimplementedInstantSwapServerServer() {}

// UnsafeInstantSwapServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InstantSwapServerServer will
// result in compilation errors.
type UnsafeInstantSwapServerServer interface {
	mustEmbedUnimplementedInstantSwapServerServer()
}

func RegisterInstantSwapServerServer(s grpc.ServiceRegistrar, srv InstantSwapServerServer) {
	s.RegisterService(&InstantSwapServer_ServiceDesc, srv)
}

func _InstantSwapServer_RequestInstantLoopOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InstantLoopOutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).RequestInstantLoopOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/RequestInstantLoopOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).RequestInstantLoopOut(ctx, req.(*InstantLoopOutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_PollPaymentAccepted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PollPaymentAcceptedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).PollPaymentAccepted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/PollPaymentAccepted",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).PollPaymentAccepted(ctx, req.(*PollPaymentAcceptedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_InitHtlcSig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitHtlcSigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).InitHtlcSig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/InitHtlcSig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).InitHtlcSig(ctx, req.(*InitHtlcSigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_PushHtlcSig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushHtlcSigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).PushHtlcSig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/PushHtlcSig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).PushHtlcSig(ctx, req.(*PushHtlcSigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_PushPreimage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushPreimageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).PushPreimage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/PushPreimage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).PushPreimage(ctx, req.(*PushPreimageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_CancelInstantSwap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelInstantSwapRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).CancelInstantSwap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/CancelInstantSwap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).CancelInstantSwap(ctx, req.(*CancelInstantSwapRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InstantSwapServer_GetInstantOutQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInstantOutQuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InstantSwapServerServer).GetInstantOutQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/looprpc.InstantSwapServer/GetInstantOutQuote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InstantSwapServerServer).GetInstantOutQuote(ctx, req.(*GetInstantOutQuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InstantSwapServer_ServiceDesc is the grpc.ServiceDesc for InstantSwapServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InstantSwapServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "looprpc.InstantSwapServer",
	HandlerType: (*InstantSwapServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestInstantLoopOut",
			Handler:    _InstantSwapServer_RequestInstantLoopOut_Handler,
		},
		{
			MethodName: "PollPaymentAccepted",
			Handler:    _InstantSwapServer_PollPaymentAccepted_Handler,
		},
		{
			MethodName: "InitHtlcSig",
			Handler:    _InstantSwapServer_InitHtlcSig_Handler,
		},
		{
			MethodName: "PushHtlcSig",
			Handler:    _InstantSwapServer_PushHtlcSig_Handler,
		},
		{
			MethodName: "PushPreimage",
			Handler:    _InstantSwapServer_PushPreimage_Handler,
		},
		{
			MethodName: "CancelInstantSwap",
			Handler:    _InstantSwapServer_CancelInstantSwap_Handler,
		},
		{
			MethodName: "GetInstantOutQuote",
			Handler:    _InstantSwapServer_GetInstantOutQuote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "instantout.proto",
}
