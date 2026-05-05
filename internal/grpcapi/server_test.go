package grpcapi

import (
	"context"
	"net"
	"strings"
	"testing"

	auditv1 "config-audit/api/gen/audit/v1"
	"config-audit/internal/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestAnalyze(t *testing.T) {
	client, cleanup := newTestClient(t)
	defer cleanup()

	response, err := client.Analyze(context.Background(), &auditv1.AnalyzeRequest{
		Config:   []byte(`{"storage":{"digest-algorithm":"MD5"}}`),
		Filename: "config.json",
	})
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if len(response.GetFindings()) != 1 {
		t.Fatalf("expected one finding, got %#v", response.GetFindings())
	}

	if !strings.Contains(response.GetFindings()[0].GetRuleId(), "weak-algorithm") {
		t.Fatalf("expected weak-algorithm finding, got %#v", response.GetFindings()[0])
	}
}

func TestAnalyzeRejectsInvalidConfig(t *testing.T) {
	client, cleanup := newTestClient(t)
	defer cleanup()

	_, err := client.Analyze(context.Background(), &auditv1.AnalyzeRequest{
		Config:   []byte(`{"storage":`),
		Filename: "config.json",
	})
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func newTestClient(t *testing.T) (auditv1.AuditServiceClient, func()) {
	t.Helper()

	listener := bufconn.Listen(1024 * 1024)
	server := NewServer(app.NewDefaultAnalyzer())

	go func() {
		_ = server.Serve(listener)
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("DialContext returned error: %v", err)
	}

	cleanup := func() {
		_ = conn.Close()
		server.Stop()
	}

	return auditv1.NewAuditServiceClient(conn), cleanup
}
