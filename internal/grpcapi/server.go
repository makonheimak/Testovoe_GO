package grpcapi

import (
	"context"
	"errors"
	"fmt"
	"net"

	auditv1 "config-audit/api/gen/audit/v1"
	"config-audit/internal/app"
	"config-audit/internal/finding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type auditServer struct {
	auditv1.UnimplementedAuditServiceServer
	analyzer app.Analyzer
}

func NewServer(analyzer app.Analyzer) *grpc.Server {
	server := grpc.NewServer()
	auditv1.RegisterAuditServiceServer(server, auditServer{analyzer: analyzer})
	return server
}

func Run(ctx context.Context, addr string, analyzer app.Analyzer) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", addr, err)
	}

	server := NewServer(analyzer)
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func (server auditServer) Analyze(ctx context.Context, request *auditv1.AnalyzeRequest) (*auditv1.AnalyzeResponse, error) {
	source := request.GetFilename()
	if source == "" {
		source = "grpc-request"
	}

	findings, err := server.analyzer.AnalyzeBytes(ctx, request.GetConfig(), source)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &auditv1.AnalyzeResponse{Findings: toProtoFindings(findings)}, nil
}

func toProtoFindings(items []finding.Finding) []*auditv1.Finding {
	out := make([]*auditv1.Finding, 0, len(items))
	for _, item := range items {
		out = append(out, &auditv1.Finding{
			RuleId:         item.RuleID,
			Severity:       string(item.Severity),
			Message:        item.Message,
			Recommendation: item.Recommendation,
			Path:           item.Path,
			Source:         item.Source,
		})
	}
	return out
}
