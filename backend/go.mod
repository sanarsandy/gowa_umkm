module gowa-backend

go 1.24.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.0

	// Gemini AI SDK
	github.com/google/generative-ai-go v0.20.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/jmoiron/sqlx v1.4.0
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/labstack/echo/v4 v4.11.4
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/redis/go-redis/v9 v9.3.0

	// WhatsApp dependencies - version will be resolved by go mod tidy
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	golang.org/x/crypto v0.44.0
	golang.org/x/oauth2 v0.21.0
	google.golang.org/api v0.186.0
)

require (
	go.mau.fi/whatsmeow v0.0.0-20251127132918-b9ac3d51d746
	google.golang.org/protobuf v1.36.10
)

require (
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/ai v0.8.0 // indirect
	cloud.google.com/go/auth v0.6.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/longrunning v0.5.7 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beeper/argo-go v1.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elliotchance/orderedmap/v3 v3.1.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/petermattis/goid v0.0.0-20250904145737-900bdf8bb490 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vektah/gqlparser/v2 v2.5.27 // indirect
	go.mau.fi/libsignal v0.2.1 // indirect
	go.mau.fi/util v0.9.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.51.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.51.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20251113190631-e25ba8c21ef6 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/grpc v1.64.1 // indirect
)
