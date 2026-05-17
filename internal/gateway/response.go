package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"online-learning-platform/internal/shared/authz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func decode(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return false
	}
	return true
}

func requestContext(r *http.Request) context.Context {
	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	return ctx
}

func grpcContext(r *http.Request) context.Context {
	ctx := requestContext(r)
	role := r.Header.Get("X-User-Role")
	userID := r.Header.Get("X-User-Id")
	pairs := make([]string, 0, 4)
	if role != "" {
		pairs = append(pairs, authz.RoleMetadataKey, role)
	}
	if userID != "" {
		pairs = append(pairs, "x-user-id", userID)
	}
	if len(pairs) == 0 {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, pairs...)
}

func writeGRPC(w http.ResponseWriter, payload any, err error) {
	if err != nil {
		writeJSON(w, grpcHTTPStatus(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func grpcHTTPStatus(err error) int {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadGateway
	}
}

//func clientsAvailable(w http.ResponseWriter, clients *Clients) bool {
//	if clients == nil {
//		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "grpc clients unavailable"})
//		return false
//	}
//	return true
//}

func methodAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func pathSegments(path, prefix string) []string {
	trimmed := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}
