package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcUserAgentHeader        = "user-agent"
	grpcGatewayClientIPHeader  = "x-forwarded-host"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}
	if mtdt, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := mtdt.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}
		if userAgents := mtdt.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}
		if clientIPs := mtdt.Get(grpcGatewayClientIPHeader); len(clientIPs) > 0 {
			md.UserAgent = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		md.ClientIP = p.Addr.String()
	}
	return md
}
