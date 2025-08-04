package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// GrpcLoggerAdapter grpc-ecosystemのloggingインターセプターを使用するためのアダプター
func GrpcLoggerAdapter() grpc.UnaryServerInterceptor {
	logger := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		// zerologのレベルに変換
		var logEvent *zerolog.Event
		switch lvl {
		case logging.LevelDebug:
			logEvent = log.Debug()
		case logging.LevelInfo:
			logEvent = log.Info()
		case logging.LevelWarn:
			logEvent = log.Warn()
		case logging.LevelError:
			logEvent = log.Error()
		default:
			logEvent = log.Info()
		}

		// フィールドを追加
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				logEvent = logEvent.Interface(fields[i].(string), fields[i+1])
			}
		}

		logEvent.Msg(msg)
	})

	return logging.UnaryServerInterceptor(logger)
}

// HttpLoggerMiddleware grpc-gatewayのruntime.WithMiddlewaresで使用するHTTPログミドルウェア
func HttpLoggerMiddleware() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		startTime := time.Now()

		// リクエスト処理後にログを出力するため、コンテキストに開始時間を保存
		ctx = context.WithValue(ctx, "request_start_time", startTime)
		ctx = context.WithValue(ctx, "request_method", req.Method)
		ctx = context.WithValue(ctx, "request_path", req.RequestURI)

		return metadata.New(nil)
	})
}

// HttpLoggerResponseInterceptor HTTPレスポンスのログ出力を行うインターセプター
func HttpLoggerResponseInterceptor() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
		startTime, ok := ctx.Value("request_start_time").(time.Time)
		if !ok {
			return nil
		}

		method, _ := ctx.Value("request_method").(string)
		path, _ := ctx.Value("request_path").(string)
		duration := time.Since(startTime)

		// レスポンスのステータスコードを取得
		statusCode := http.StatusOK
		if ww, ok := w.(interface{ Status() int }); ok {
			statusCode = ww.Status()
		}

		logger := log.Info()
		if statusCode >= 400 {
			logger = log.Error()
		}

		logger.Str("protocol", "http").
			Str("method", method).
			Str("path", path).
			Int("status_code", statusCode).
			Str("status_text", http.StatusText(statusCode)).
			Dur("duration", duration).
			Msg("HTTP request processed (grpc-gateway)")

		return nil
	})
}

// HttpLogger grpc-gatewayのruntime.WithMiddlewaresで使用するHTTPログミドルウェアの組み合わせ
func HttpLogger() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		HttpLoggerMiddleware(),
		HttpLoggerResponseInterceptor(),
	}
}

// 従来のHTTPログミドルウェア（後方互換性のため残す）
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func LegacyHttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode >= 400 {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Str("remote_addr", req.RemoteAddr).
			Str("user_agent", req.UserAgent()).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("HTTP request processed")
	})
}
