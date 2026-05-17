package authz

import (
	"context"
	"online-learning-platform/internal/shared/grpcerrors"

	"google.golang.org/grpc/metadata"
)

const RoleMetadataKey = "x-user-role"

func RoleFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(RoleMetadataKey)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func RequireRole(ctx context.Context, allowed ...string) error {
	role := RoleFromContext(ctx)
	if role == "" {
		return grpcerrors.Unauthenticated("x-user-role metadata is required")
	}
	for _, candidate := range allowed {
		if role == candidate {
			return nil
		}
	}
	return grpcerrors.PermissionDenied("permission denied for role " + role)
}
