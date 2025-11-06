# Production Deployment

**Comprehensive guide to deploying Templ Router applications in production environments.**

## Overview

This guide covers production deployment strategies for Templ Router applications, including Docker containerization, Kubernetes orchestration, monitoring, security best practices, and operational considerations.

**Key Topics:**
- Docker deployment strategies
- Kubernetes configuration
- Security and authentication
- Performance optimization
- Monitoring and logging
- CI/CD pipelines
- Environment management
- Backup and recovery

## Configuration Prefix Notice

**Important:** All environment variables in this documentation use the default prefix `TR_`. This prefix is **configurable** when you set up your application:

```go
// Default configuration (uses TR_ prefix)
container.RegisterRouterServices("TR")

// Custom prefix configuration
container.RegisterRouterServices("MYAPP")  // Environment variables will use MYAPP_ prefix
```

**Examples:**
- Default: `TR_SERVER_HOST=localhost`
- Custom: `MYAPP_SERVER_HOST=localhost`
- Multiple apps: `APP1_SERVER_HOST=localhost` and `APP2_SERVER_HOST=localhost`

All examples below use the default `TR_` prefix, but you can replace `TR` with your custom prefix in all environment variable names.

## Docker Deployment

### Multi-stage Dockerfile

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install trgen tool
RUN go install github.com/denkhaus/templ-router/cmd/trgen@latest

# Generate templates
RUN templ generate
RUN trgen --scan-path=app --module-name=github.com/youruser/yourproject

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy templates and static files
COPY --from=builder /app/app ./app
COPY --from=builder /app/generated ./generated

# Set permissions
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Start application
CMD ["./main"]
```

### Docker Compose for Production

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      # Server Configuration
      - TR_SERVER_HOST=0.0.0.0
      - TR_SERVER_PORT=8080
      - TR_SERVER_BASE_URL=https://yourapp.com

      # Database Configuration
      - TR_DATABASE_HOST=postgres
      - TR_DATABASE_PORT=5432
      - TR_DATABASE_USER=${POSTGRES_USER}
      - TR_DATABASE_PASSWORD=${POSTGRES_PASSWORD}
      - TR_DATABASE_NAME=${POSTGRES_DB}
      - TR_DATABASE_SSL_MODE=require

      # Authentication
      - TR_AUTH_SESSION_EXPIRY=1h
      - TR_AUTH_SESSION_COOKIE_SECURE=true
      - TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true
      - TR_AUTH_SESSION_COOKIE_SAME_SITE=Strict

      # Internationalization
      - TR_I18N_SUPPORTED_LOCALES=en-US,de-DE,fr-FR
      - TR_I18N_DEFAULT_LOCALE=en-US
      - TR_I18N_URL_PREFIX=true

      # Security
      - TR_SECURITY_ENABLE_MIDDLEWARE=true
      - TR_SECURITY_ENABLE_RATE_LIMIT=true
      - TR_SECURITY_RATE_LIMIT_REQUESTS=100
      - TR_SECURITY_ENABLE_CSRF=true

      # Logging
      - TR_LOGGING_LEVEL=info
      - TR_LOGGING_FORMAT=json
      - TR_LOGGING_ENABLE_FILE=true
      - TR_LOGGING_FILE_PATH=/var/log/app/app.log

      # Environment
      - TR_ENVIRONMENT_KIND=production

      # Monitoring
      - TR_ENVIRONMENT_METRICS_ENABLED=true
      - TR_ENVIRONMENT_HEALTH_CHECK_ENABLED=true

    volumes:
      - ./logs:/var/log/app

    depends_on:
      - postgres
      - redis

    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    restart: unless-stopped
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-db:/docker-entrypoint-initdb.d
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:
  logs:

networks:
  app-network:
    driver: bridge
```

### Environment Configuration

```bash
# .env.production
POSTGRES_DB=yourapp_prod
POSTGRES_USER=yourapp_user
POSTGRES_PASSWORD=your-secure-password

# Security (generate strong passwords)
TR_SECURITY_CSRF_SECRET=$(openssl rand -base64 32)
TR_DATABASE_PASSWORD=$(openssl rand -base64 32)

# SSL/TLS Configuration
TR_SERVER_READ_TIMEOUT=15s
TR_SERVER_WRITE_TIMEOUT=15s
TR_SERVER_IDLE_TIMEOUT=60s
TR_SERVER_SHUTDOWN_TIMEOUT=10s

# Session Security
TR_AUTH_SESSION_EXPIRY=1h
TR_AUTH_SESSION_COOKIE_NAME=session_id
TR_AUTH_SESSION_COOKIE_SECURE=true
TR_AUTH_SESSION_COOKIE_HTTP_ONLY=true
TR_AUTH_SESSION_COOKIE_SAME_SITE=Strict
TR_AUTH_SESSION_COOKIE_PATH=/

# Rate Limiting
TR_SECURITY_RATE_LIMIT_REQUESTS=60
TR_SECURITY_RATE_LIMIT_WINDOW=1m

# Logging
TR_LOGGING_LEVEL=warn
TR_LOGGING_FORMAT=json
TR_LOGGING_FILE_MAX_SIZE=100MB
TR_LOGGING_FILE_MAX_BACKUPS=5
```

## Kubernetes Deployment

### Kubernetes Manifests

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: yourapp
  labels:
    name: yourapp
    environment: production

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: yourapp-config
  namespace: yourapp
data:
  server.conf: |
    TR_SERVER_HOST=0.0.0.0
    TR_SERVER_PORT=8080
    TR_SERVER_READ_TIMEOUT=15s
    TR_SERVER_WRITE_TIMEOUT=15s
  logging.conf: |
    TR_LOGGING_LEVEL=info
    TR_LOGGING_FORMAT=json
    TR_LOGGING_ENABLE_FILE=true
  i18n.conf: |
    TR_I18N_SUPPORTED_LOCALES=en-US,de-DE,fr-FR
    TR_I18N_DEFAULT_LOCALE=en-US
    TR_I18N_URL_PREFIX=true

---
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: yourapp-secrets
  namespace: yourapp
type: Opaque
data:
  database-password: <base64-encoded-password>
  csrf-secret: <base64-encoded-csrf-secret>
  jwt-secret: <base64-encoded-jwt-secret>

---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: yourapp
  namespace: yourapp
  labels:
    app: yourapp
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: yourapp
      version: v1
  template:
    metadata:
      labels:
        app: yourapp
        version: v1
    spec:
      containers:
      - name: yourapp
        image: yourapp/yourapp:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: TR_DATABASE_HOST
          value: "postgres-service"
        - name: TR_DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: yourapp-secrets
              key: database-password
        - name: TR_SECURITY_CSRF_SECRET
          valueFrom:
            secretKeyRef:
              name: yourapp-secrets
              key: csrf-secret
        - name: TR_ENVIRONMENT_KIND
          value: "production"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

---
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: yourapp-service
  namespace: yourapp
spec:
  selector:
    app: yourapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP

---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: yourapp-ingress
  namespace: yourapp
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - yourapp.com
    secretName: yourapp-tls
  rules:
  - host: yourapp.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: yourapp-service
            port:
              number: 80

---
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: yourapp-hpa
  namespace: yourapp
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: yourapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Database StatefulSet

```yaml
# k8s/postgres.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: yourapp
spec:
  serviceName: postgres
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          value: "yourapp_prod"
        - name: POSTGRES_USER
          value: "yourapp_user"
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: yourapp-secrets
              key: database-password
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
  volumeClaimTemplates:
  - metadata:
      name: postgres-storage
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 10Gi
      storageClassName: fast-ssd
```

## CI/CD Pipeline

### GitHub Actions Workflow

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    branches: [main]
  tags:
    - 'v*'

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'

    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-

    - name: Install dependencies
      run: go mod download

    - name: Install tools
      run: |
        go install github.com/a-h/templ/cmd/templ@latest
        go install github.com/denkhaus/templ-router/cmd/trgen@latest

    - name: Generate templates
      run: templ generate

    - name: Generate template registry
      run: trgen --scan-path=app --module-name=github.com/youruser/yourproject

    - name: Run tests
      run: go test ./...

    - name: Run E2E tests
      run: |
        # Start application in background
        ./main &
        APP_PID=$!

        # Wait for application to start
        sleep 10

        # Run E2E tests
        go test ./tests/e2e/...

        # Cleanup
        kill $APP_PID

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'

    steps:
    - uses: actions/checkout@v3

    - name: Log in to GitHub Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: |
          ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
          ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ meta.version }}
        labels: ${{ steps.meta.outputs.labels }}

    - name: Deploy to Kubernetes
      run: |
        # Set up kubectl
        curl -LO "https://dl.k8s.io/release/v1.29.0/bin/linux/amd64/kubectl"
        chmod +x kubectl
        sudo mv kubectl /usr/local/bin/

        # Configure kubectl
        echo "${{ secrets.KUBE_CONFIG }}" | base64 -d > kubeconfig
        export KUBECONFIG=kubeconfig

        # Apply Kubernetes manifests
        kubectl apply -f k8s/

        # Wait for rollout
        kubectl rollout status deployment/yourapp -n yourapp
        kubectl rollout status deployment/postgres -n yourapp
```

### ArgoCD GitOps

```yaml
# argocd/application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: yourapp
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/youruser/yourapp.git
    targetRevision: HEAD
    path: k8s
  destination:
      server: https://kubernetes.default.svc
      namespace: yourapp
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
  ignoreDifferences:
  - group: apps
    kind: Deployment
    jsonPointers:
    - path: /spec/template/spec/containers/0/image
  - group: apps
    kind: StatefulSet
    jsonPointers:
    - path: /spec/template/spec/containers/0/image
```

## Security Configuration

### SSL/TLS Configuration

```yaml
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    # Redirect HTTP to HTTPS
    server {
        listen 80;
        server_name yourapp.com;
        return 301 https://$server_name$request_uri;
    }

    # HTTPS configuration
    server {
        listen 443 ssl http2;
        server_name yourapp.com;

        ssl_certificate /etc/nginx/ssl/cert.pem;
        ssl_certificate_key /etc/ssl/key.pem;
        ssl_session_timeout 1d;
        ssl_session_cache shared:SSL:50m;
        ssl_session_tickets off;

        # Modern SSL configuration
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-SHA384;
        ssl_prefer_server_ciphers off;

        # Security headers
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options DENY;
        add_header X-Content-Type-Options nosniff;
        add_header X-XSS-Protection "1; mode=block";
        add_header Referrer-Policy "strict-origin-when-cross-origin";

        location / {
            proxy_pass http://yourapp-service:80;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # API endpoints
        location /api/ {
            proxy_pass http://yourapp-service:80;
            proxy_set_header Host $host;

            # CORS headers
            add_header Access-Control-Allow-Origin "$http_origin" always;
            add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
            add_header Access-Control-Allow-Headers "Authorization, Content-Type" always;
            add_header Access-Control-Allow-Credentials "true" always;
        }
    }
}
```

### Security Headers Middleware

```go
// pkg/middleware/security_middleware.go
package middleware

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

type SecurityMiddleware struct {
    config *SecurityConfig
}

func NewSecurityMiddleware(config *SecurityConfig) *SecurityMiddleware {
    return &SecurityMiddleware{config: config}
}

func (m *SecurityMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        m.setSecurityHeaders(w)

        // CSP header
        if m.config.CSPEnabled {
            m.setCSPHeader(w)
        }

        next.ServeHTTP(w, r)
    })
}

func (m *SecurityMiddleware) setSecurityHeaders(w http.ResponseWriter) {
    headers := map[string]string{
        "X-Frame-Options":           "DENY",
        "X-Content-Type-Options":   "nosniff",
        "X-XSS-Protection":        "1; mode=block",
        "Referrer-Policy":          "strict-origin-when-cross-origin",
        "Permissions-Policy":       "interest-cohort=()",
    }

    for key, value := range headers {
        w.Header().Set(key, value)
    }
}

func (m *SecurityMiddleware) setCSPHeader(w http.ResponseWriter) {
    csp := "default-src 'self'; " +
          "script-src 'self' 'unsafe-inline'; " +
          "style-src 'self' 'unsafe-inline'; " +
          "img-src 'self' data: https:; " +
          "font-src 'self'; " +
          "connect-src 'self'; " +
          "frame-ancestors 'self'; " +
          "form-action 'self';"

    w.Header().Set("Content-Security-Policy", csp)
}
```

## Monitoring and Observability

### Prometheus Metrics

```go
// pkg/metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/samber/do/v2"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "status", "path"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)

type MetricsMiddleware struct {
    injector *do.Injector
}

func NewMetricsMiddleware(injector *do.Injector) *MetricsMiddleware {
    return &MetricsMiddleware{injector: injector}
}

func (m *MetricsMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Increment active connections
        activeConnections.Inc()
        defer activeConnections.Dec()

        next.ServeHTTP(w, r)

        // Record metrics
        duration := time.Since(start)
        httpRequestsTotal.WithLabelValues(
            r.Method,
            "200", // Should be captured from response
            r.URL.Path,
        ).Inc()

        httpRequestDuration.WithLabelValues(
            r.Method,
            r.URL.Path,
        ).Observe(duration.Seconds())
    })
}

func InitMetrics() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(activeConnections)

    // Expose metrics endpoint
    http.Handle("/metrics", promhttp.Handler())
}
```

### Health Checks

```go
// pkg/health/health.go
package health

import (
    "encoding/json"
    "net/http"
    "time"
)

type HealthStatus struct {
    Status    string                 `json:"status"`
    Timestamp time.Time              `json:"timestamp"`
    Version   string                 `json:"version"`
    Checks    map[string]HealthCheck   `json:"checks"`
}

type HealthCheck struct {
    Status  string `json:"status"`
    Message string `json:"message,omitempty"`
}

type HealthChecker struct {
    checks map[string]func() HealthCheck
    config *HealthConfig
}

func NewHealthChecker(config *HealthConfig) *HealthChecker {
    return &HealthChecker{
        checks: make(map[string]func() HealthCheck),
        config: config,
    }
}

func (h *HealthChecker) AddCheck(name string, check func() HealthCheck) {
    h.checks[name] = check
}

func (h *HealthChecker) CheckHealth() HealthStatus {
    status := HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   h.config.Version,
        Checks:    make(map[string]HealthCheck),
    }

    allHealthy := true
    for name, check := range h.checks {
        result := check()
        status.Checks[name] = result

        if result.Status != "healthy" {
            allHealthy = false
        }
    }

    if !allHealthy {
        status.Status = "unhealthy"
    }

    return status
}

func (h *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
    status := h.CheckHealth()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(getStatusCode(status.Status))
    json.NewEncoder(w).Encode(status)
}

func getStatusCode(status string) int {
    switch status {
    case "healthy":
        return 200
    case "degraded":
        return 200
    case "unhealthy":
        return 503
    default:
        return 500
    }
}
```

### Structured Logging

```go
// pkg/logging/logger.go
package logging

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Logger struct {
    *zap.Logger
}

func NewLogger(config *LogConfig) (*Logger, error) {
    var zapConfig zap.Config

    switch config.Level {
    case "debug":
        zapConfig = zap.NewDevelopmentConfig()
    case "info":
        zapConfig = zap.NewProductionConfig()
    case "warn":
        zapConfig = zap.NewProductionConfig()
        zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
    case "error":
        zapConfig = zap.NewProductionConfig()
        zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
    default:
        zapConfig = zap.NewProductionConfig()
    }

    // Add file output if enabled
    if config.EnableFile && config.FilePath != "" {
        zapConfig.OutputPaths = []string{config.FilePath}
        zapConfig.EncoderConfig = zapcore.NewJSONEncoder(
            zapcore.EncoderConfig{
                TimeKey:        "timestamp",
                LevelKey:       "level",
                NameKey:        "logger",
                CallerKey:      "caller",
                MessageKey:     "msg",
                StacktraceKey:    "stacktrace",
                LineEnding:     zapcore.DefaultLineEnding,
                EncodeLevel:    zapcore.LowercaseLevelEncoder,
                EncodeTime:     zapcore.ISO8601TimeEncoder,
                EncodeDuration: zapcore.SecondsDurationEncoder,
                EncodeCaller:    zapcore.ShortCallerEncoder,
            },
        )
    }

    logger, err := zapConfig.Build()
    if err != nil {
        return nil, err
    }

    return &Logger{Logger: logger}, nil
}

func (l *Logger) Request(method, path, status int, duration time.Duration, userID string) {
    l.Logger.Info("HTTP Request",
        zap.String("method", method),
        zap.String("path", path),
        zap.Int("status", status),
        zap.Duration("duration", duration),
        zap.String("user_id", userID),
        zap.String("remote_addr", "client_ip"),
    )
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
    l.Logger.Error(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.Logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
    l.Logger.Debug(msg, fields...)
}
```

## Performance Optimization

### Database Connection Pooling

```go
// pkg/database/connection_pool.go
package database

import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/lib/pq"
)

type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime  time.Duration
}

func NewConnectionPool(config *DatabaseConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host,
        config.Port,
        config.User,
        config.Password,
        config.Database,
        config.SSLMode,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)
    db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}
```

### Template Caching

```go
// pkg/cache/template_cache.go
package cache

import (
    "crypto/sha256"
    "fmt"
    "time"
    "github.com/denkhaus/templ-router/pkg/interfaces"
)

type TemplateCache struct {
    cache   map[string]*CacheEntry
    ttl     time.Duration
    maxSize int
}

type CacheEntry struct {
    Content   interface{}
    CreatedAt time.Time
    ExpiresAt time.Time
}

func NewTemplateCache(ttl time.Duration, maxSize int) *TemplateCache {
    return &TemplateCache{
        cache:   make(map[string]*CacheEntry),
        ttl:     ttl,
        maxSize: maxSize,
    }
}

func (c *TemplateCache) Get(key string) (interface{}, bool) {
    if entry, exists := c.cache[key]; exists {
        if time.Now().Before(entry.ExpiresAt) {
            delete(c.cache, key)
            return nil, false
        }
        return entry.Content, true
    }
    return nil, false
}

func (c *TemplateCache) Set(key string, content interface{}) {
    // Implement LRU eviction if cache is full
    if len(c.cache) >= c.maxSize {
        c.evictOldest()
    }

    c.cache[key] = &CacheEntry{
        Content:   content,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(c.ttl),
    }
}

func (c *TemplateCache) evictOldest() {
    var oldestKey string
    var oldestTime time.Time

    for key, entry := range c.cache {
        if oldestTime.IsZero() || entry.CreatedAt.Before(oldestTime) {
            oldestTime = entry.CreatedAt
            oldestKey = key
        }
    }

    if oldestKey != "" {
        delete(c.cache, oldestKey)
    }
}
```

## Backup and Recovery

### Database Backup Strategy

```bash
#!/bin/bash
# backup.sh - Database backup script

set -e

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="yourapp_user"
DB_NAME="yourapp_prod"
BACKUP_DIR="/backups/database"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Create database backup
echo "Creating database backup: $BACKUP_FILE"
PGPASSWORD="$DB_PASSWORD" pg_dump \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    --no-owner \
    --no-privileges \
    --verbose \
    --file="$BACKUP_FILE"

# Compress backup
echo "Compressing backup: $BACKUP_FILE.gz"
gzip "$BACKUP_FILE"

# Clean up old backups (keep last 7 days)
echo "Cleaning up old backups..."
find "$BACKUP_DIR" -name "*.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

### Application Configuration Backup

```bash
#!/bin/bash
# backup-config.sh - Configuration backup script

set -e

CONFIG_DIR="/etc/yourapp"
BACKUP_DIR="/backups/config"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup environment configuration
echo "Backing up environment variables..."
env | grep "^TR_" > "$BACKUP_DIR/env_$DATE.env"

# Backup SSL certificates
if [ -d "/etc/ssl/certs" ]; then
    echo "Backing up SSL certificates..."
    tar -czf "$BACKUP_DIR/ssl_certs_$DATE.tar.gz" -C /etc/ssl/certs .
fi

# Backup application configuration
if [ -d "$CONFIG_DIR" ]; then
    echo "Backing up application configuration..."
    tar -czf "$BACKUP_DIR/config_$DATE.tar.gz" -C "$CONFIG_DIR" .
fi

echo "Configuration backup completed"
```

### Disaster Recovery Plan

```markdown
# DISASTER RECOVERY PLAN

## 1. Immediate Response (0-1 hour)
- Alert stakeholders
- Assess damage
- Initiate incident response
- Preserve evidence

## 2. Service Restoration (1-4 hours)
- Restore from most recent backup
- Verify data integrity
- Restore application services
- Test critical functionality

## 3. Full Recovery (4-24 hours)
- Restore all data from backups
- Rebuild infrastructure as needed
- Verify all systems are functional
- Conduct post-incident review

## 4. Post-Recovery (24+ hours)
- Monitor system stability
- Implement improvements
- Update recovery procedures
- Document lessons learned

## Recovery Procedures

### Database Recovery
1. Stop application services
2. Restore database from backup
3. Verify data integrity
4. Restart services
5. Test critical functionality

### Application Recovery
1. Restore application files
2. Rebuild Docker images
3. Update Kubernetes manifests
4. Redeploy services
5. Test all functionality

### Infrastructure Recovery
1. Recreate cloud resources
2. Restore configurations
3. Update DNS records
4. Reconfigure monitoring
5. Test connectivity
```

## Next Steps

- **[Getting Started](GETTING-STARTED.md)** - Basic deployment setup
- **[Configuration](CONFIGURATION.md)** - Production configuration
- **[Security](AUTHENTICATION.md)** - Security best practices
- **[Monitoring](../CONTRIBUTING.md#monitoring)** - Set up monitoring
- **[Architecture](ARCHITECTURE.md)** - System architecture

## Need Help?

- **[Discussions](https://github.com/denkhaus/templ-router/discussions)** - Community discussions
- **[Issues](https://github.com/denkhaus/templ-router/issues)** - Report bugs or request features
- **[Kubernetes Documentation](https://kubernetes.io/docs/)** - Kubernetes deployment
- **[Docker Documentation](https://docs.docker.com/)** - Docker deployment