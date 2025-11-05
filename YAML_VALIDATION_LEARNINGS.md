# YAML Validation - Key Learnings & Implementation Notes

## Critical Learning: YAML Validation Fail-Fast Implementation

### **Current Problem Analysis**

**Root Cause**: YAML validation errors are silently ignored at startup, leading to production issues
**Impact**: Invalid configurations with unknown root keys (like `component:`) are accepted, causing runtime failures

### **Key Technical Details**

#### 1. **Current Validation Logic (PROBLEMATIC)**
**File**: `pkg/shared/yaml_parser.go:57` (ParseYAMLMetadata caller)
```go
// CURRENT PROBLEMATIC PATTERN
func ParseYAMLMetadata(filePath string) (bool, *ConfigFile, error) {
    // ... validation logic ...
    err := validateRootKeys(rawConfig)
    if err != nil {
        return true, nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
    }
    // Continues execution despite validation errors!
}
```

#### 2. **Validation Function** (CORRECT)
**File**: `pkg/shared/yaml_parser.go:85` (validateRootKeys)
```go
func validateRootKeys(rawConfig map[string]interface{}) error {
    allowedKeys := map[string]bool{
        "i18n":     true,
        "auth":     true,
        "metadata": true,
        "error":    true,
        "dynamic":  true,
    }
    // ... validation logic - this works correctly
}
```

#### 3. **Runtime Loading (CORRECT)**
**File**: `pkg/router/middleware/template_middleware.go` (template_middleware.go)
```go
// Runtime loading SHOULD gracefully fallback (this is correct)
func (tm *templateMiddleware) loadComponentMetadata(templateFile string) *shared.ConfigFile {
    configFileFound, config, err := shared.ParseYAMLMetadata(yamlPath)
    if err != nil {
        // This is correct behavior for runtime
        return nil // No metadata available
    }
    return config
}
```

### **Implementation Strategy**

#### **Phase 1: Identify Startup vs Runtime Scenarios**

**Startup Scenarios (should be fatal):**
- Application bootstrap
- Configuration loading
- Service initialization

**Runtime Scenarios (should be graceful):**
- Template middleware loading
- Component metadata loading
- Dynamic template discovery

#### **Phase 2: Enhanced Error Handling**

**For Startup Fatal Errors:**
```go
// In startup code
if configFileFound {
    // Log detailed error information
    logger.Error("FATAL: YAML validation failed during startup",
        zap.String("file", filePath),
        zap.Error("validation_error", err))

    // Exit application immediately
    os.Exit(1)
}
```

#### **Phase 3: Improve Error Messages**

**Current**: `"unknown root key 'component' - allowed keys are: i18n, auth, metadata, layout, error, dynamic"`

**Enhanced**:
```
ERROR: Invalid YAML structure in file 'demo/app/components/footer.templ.yaml'
Unknown root key 'component' found at line 45, column 3

Allowed root keys are:
- i18n: Internationalization translations (en, de, fr, etc.)
- auth: Authentication and authorization settings
- metadata: Template metadata and configuration
- error: Error handling and display configuration
- dynamic: Dynamic parameter validation rules

Example of correct structure:
i18n:
  en:
    welcome: "Welcome"
  de:
    welcome: "Willkommen"

metadata:
  title: "My App"
  version: "1.0.0"
```

### **Implementation Files & Changes**

#### **Primary Changes Needed:**

1. **pkg/shared/yaml_parser.go**
   - Add `FatalParseYAMLMetadata()` function for startup scenarios
   - Keep existing `ParseYAMLMetadata()` for runtime scenarios
   - Enhance `validateRootKeys()` error messages

2. **Application Bootstrap**
   - Update startup code to use `FatalParseYAMLMetadata()`
   - Add clear error logging and application exit

3. **Template Middleware**
   - Keep existing graceful fallback behavior
   - This is already correct - no changes needed

### **Testing Strategy**

#### **Test Cases for Fatal Failures:**
```go
// Should cause application exit
func TestInvalidYAMLStartup() {
    // Create YAML with unknown root key
    // Attempt to start application
    // Verify application exits with code 1
}
```

#### **Test Cases for Runtime Graceful Fallback:**
```go
// Should not cause application exit
func TestInvalidYAMLRuntime() {
    // Create YAML with unknown root key
    // Load via template middleware
    // Verify nil return and graceful handling
}
```

### **Validation Checklist**

- [ ] Startup validation errors cause application exit
- [ ] Runtime validation errors are handled gracefully
- [ ] Error messages include file paths and line numbers
- [ ] Error messages provide examples of valid structure
- [ ] No breaking changes to existing runtime behavior
- [ ] All existing tests still pass
- [ ] New tests cover both startup and runtime scenarios

### **Production Benefits**

1. **Fast Fail Principle**: Errors caught early, not in production
2. **Clear Error Messages**: Developers can quickly fix configuration issues
3. **Maintains Graceful Runtime**: Runtime loading still works as expected
4. **Backward Compatibility**: No breaking changes to existing functionality

### **Risk Mitigation**

1. **Gradual Rollout**: Implement in stages with proper testing
2. **Backward Compatibility**: Ensure existing runtime behavior unchanged
3. **Extensive Testing**: Cover all YAML loading scenarios
4. **Documentation**: Provide migration guide for affected code