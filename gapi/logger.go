package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}
	logger := log.Info()

	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protcol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_string", statusCode.String()).
		Str("request duration", duration.String()).
		Msg("Received gRPC request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.Body = b
	return r.ResponseWriter.Write(b)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		duration := time.Since(startTime)

		logger := log.Info()

		rec := &ResponseRecorder{
			ResponseWriter: res,
			statusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)

		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", req.Proto).
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.statusCode).
			Str("status_string", http.StatusText(rec.statusCode)).
			Str("request duration", duration.String()).
			Msg("Received HTTP request")

	})
}
