package interceptors

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"

	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc"
)

func APIKeyInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Internal, "Failed to extract metadata from context")
		}

		apiKeys := md.Get("Api-key")
		if len(apiKeys) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "API key is missing")
		}

		apiKey := apiKeys[0]
		if apiKey != "your-api-key" {
			return nil, status.Errorf(codes.PermissionDenied, "Invalid API key")
		}

		// Gọi handler tiếp theo
		return handler(ctx, req)
	}
}

func LoggingInterceptor(logger core.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Debugf("GRPC request: %s", info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Errorf("GRPC request failed: %v", err)
		} else {
			logger.Debugf("Completed GRPC request: %s", info.FullMethod)
		}
		return resp, err
	}
}

func RecoverInterceptor(logger core.Logger, isDebug bool, timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isDebug {
			logger.Debugf("GRPC request: %s", info.FullMethod)
		}
		timeoutCh := time.After(timeout)
		var respChan = make(chan interface{}, 1)
		var errChan = make(chan error, 1)
		defer func() {
			if r := recover(); r != nil {
				errChan <- status.Errorf(codes.Internal, "Internal server error")
			}
		}()

		go func() {
			resp, err := handler(ctx, req)
			if err != nil {
				errChan <- err
			} else {
				respChan <- resp
			}
		}()
		select {
		case resp := <-respChan:
			if isDebug {
				logger.Debugf("GRPC response: %s", info.FullMethod)
			}
			return resp, nil
		case err := <-errChan:
			if isDebug {
				logger.Errorf("GRPC request failed: %v", err)
			}
			return nil, err
		case <-timeoutCh:
			errTimeout := status.Errorf(codes.DeadlineExceeded, "GRPC request timed out")
			if isDebug {
				logger.Errorf("GRPC request time out: %v", errTimeout)
			}
			return nil, errTimeout
		}
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
}
