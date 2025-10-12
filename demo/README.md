# Demo Application - Router Library Usage

## 🎯 Purpose
This demo shows how to use the feedback-router as a **standalone library** with **external template registry**.

## 🏗️ Architecture
```
Demo App
├── app/                    # Demo templates
├── generated/templates/    # Generated template registry
├── main.go                # Demo application entry point
└── cmd/generate-templates/ # Demo-specific template generator
```

## 🚀 Usage

### 1. Generate Templates
```bash
cd demo
go run cmd/generate-templates/main.go
```

### 2. Run Demo
```bash
go run main.go
```

### 3. Test Routes
```bash
curl http://localhost:8084/de/dashboard
curl http://localhost:8084/en/user/123
curl http://localhost:8084/de/product/laptop
```

## 📋 Key Features Demonstrated

### ✅ **Library Usage Pattern**
- Router imported as external library
- Template registry generated in application
- Dependency injection connects both

### ✅ **Clean Architecture**
- No hardcoded dependencies
- Interface-based template registry
- Pluggable store implementations

### ✅ **Production-Ready**
- 15 routes automatically discovered
- Multi-language support (DE/EN)
- Dynamic parameter handling
- Error handling with fallbacks

## 🔧 Configuration

The demo uses environment variables for configuration:
- `TEMPLATE_SCAN_PATH`: Path to scan for templates
- `TEMPLATE_OUTPUT_DIR`: Where to generate registry
- `TEMPLATE_TARGET_PACKAGE`: Specific package to generate for

## 📊 Results

When running, you should see:
```
🎯 Discovered 15 routes from template registry
🌐 Multi-language support: DE, EN
🔄 Dynamic routes: /locale_/user/id_, /locale_/product/id_
✅ Clean Architecture Demo Server on :8084
```