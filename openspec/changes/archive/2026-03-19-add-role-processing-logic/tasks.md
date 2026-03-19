## 1. Configuration

- [x] 1.1 Add `Docker` section to `Config` struct in `pkg/config/config.go` with `Registry` and `ImageLoadTimeout` fields
- [x] 1.2 Update `Default()` function to provide sensible defaults (e.g., empty registry, 300s timeout)
- [x] 1.3 Update `config.example.yaml` with Docker configuration example

## 2. Image Processing Module

- [x] 2.1 Create `pkg/image/` package for Docker image processing logic
- [x] 2.2 Implement `ImageFile` function to validate path exists and has `.tar` extension
- [x] 2.3 Implement `LoadImage()` function using `docker load -i` with timeout support
- [x] 2.4 Implement `PushImage()` function to tag and push to registry
- [x] 2.5 Add error types for image processing failures

## 3. Image Validation in Role Handler

- [x] 3.1 Update role creation handler to validate `ImageFile` path when specified
- [x] 3.2 Add `ImageFile` field to role data model if not already present
- [x] 3.3 Integrate image processing into role creation flow after validation passes

## 4. Image Processing Integration

- [x] 4.1 Call image processing after role files are read and before database persistence
- [x] 4.2 Handle missing registry gracefully (skip image processing)
- [x] 4.3 Ensure role creation fails if image processing fails
- [x] 4.4 Log image processing progress and errors appropriately

## 5. Testing

- [x] 5.1 Add unit tests for image path validation
- [x] 5.2 Add unit tests for image load with timeout
- [x] 5.3 Add integration tests for full role creation with image processing
- [x] 5.4 Test error handling for missing files, invalid archives, registry failures
