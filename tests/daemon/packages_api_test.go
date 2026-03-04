package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPackagesEmpty(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/packages", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var packages []interface{}
	json.Unmarshal(w.Body.Bytes(), &packages)

	if len(packages) != 0 {
		t.Errorf("Expected empty array, got %d packages", len(packages))
	}
}

func TestCreatePackage(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	body := `{"name":"test-package","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	req, _ := http.NewRequest("POST", "/packages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var pkg map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pkg)

	if pkg["name"] != "test-package" {
		t.Errorf("Expected name 'test-package', got '%v'", pkg["name"])
	}
}

func TestListPackagesWithData(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-package","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/packages", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	req, _ := http.NewRequest("GET", "/packages", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var packages []interface{}
	json.Unmarshal(w.Body.Bytes(), &packages)

	if len(packages) != 1 {
		t.Errorf("Expected 1 package, got %d packages", len(packages))
	}
}

func TestGetPackageByID(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-package","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/packages", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdPkg map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdPkg)
	packageID := createdPkg["id"].(string)

	req, _ := http.NewRequest("GET", "/packages/"+packageID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var pkg map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pkg)

	if pkg["id"] != packageID {
		t.Errorf("Expected id %s, got %v", packageID, pkg["id"])
	}
	if pkg["compose_content"] == nil {
		t.Error("Expected compose_content to be present")
	}
}

func TestGetPackageByIDNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/packages/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeletePackage(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-package","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/packages", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdPkg map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdPkg)
	packageID := createdPkg["id"].(string)

	req, _ := http.NewRequest("DELETE", "/packages/"+packageID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	getReq, _ := http.NewRequest("GET", "/packages/"+packageID, nil)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", getW.Code)
	}
}

func TestDeletePackageNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("DELETE", "/packages/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}
