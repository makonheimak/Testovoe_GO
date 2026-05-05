package auditv1

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const AuditServiceName = "audit.v1.AuditService"

type AuditServiceClient interface {
	Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error)
}

type auditServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuditServiceClient(cc grpc.ClientConnInterface) AuditServiceClient {
	return &auditServiceClient{cc: cc}
}

func (client *auditServiceClient) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := client.cc.Invoke(ctx, "/audit.v1.AuditService/Analyze", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type AuditServiceServer interface {
	Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error)
}

type UnimplementedAuditServiceServer struct{}

func (UnimplementedAuditServiceServer) Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Analyze not implemented")
}

func RegisterAuditServiceServer(server grpc.ServiceRegistrar, service AuditServiceServer) {
	server.RegisterService(&AuditService_ServiceDesc, service)
}

func _AuditService_Analyze_Handler(service any, ctx context.Context, decoder func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(AnalyzeRequest)
	if err := decoder(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return service.(AuditServiceServer).Analyze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     service,
		FullMethod: "/audit.v1.AuditService/Analyze",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return service.(AuditServiceServer).Analyze(ctx, req.(*AnalyzeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var AuditService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: AuditServiceName,
	HandlerType: (*AuditServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Analyze",
			Handler:    _AuditService_Analyze_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/audit/v1/audit.proto",
}
