package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRolesEmpty(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/roles", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var roles []interface{}
	json.Unmarshal(w.Body.Bytes(), &roles)

	if len(roles) != 0 {
		t.Errorf("Expected empty array, got %d roles", len(roles))
	}
}

func TestCreateRole(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	body := `{"name":"test-role","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	req, _ := http.NewRequest("POST", "/roles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var role map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &role)

	if role["name"] != "test-role" {
		t.Errorf("Expected name 'test-role', got '%v'", role["name"])
	}
}

func TestListRolesWithData(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-role","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/roles", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	req, _ := http.NewRequest("GET", "/roles", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var roles []interface{}
	json.Unmarshal(w.Body.Bytes(), &roles)

	if len(roles) != 1 {
		t.Errorf("Expected 1 role, got %d roles", len(roles))
	}
}

func TestGetRoleByID(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-role","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/roles", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdRole map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdRole)
	roleID := createdRole["id"].(string)

	req, _ := http.NewRequest("GET", "/roles/"+roleID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var role map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &role)

	if role["id"] != roleID {
		t.Errorf("Expected id %s, got %v", roleID, role["id"])
	}
	if role["compose_content"] == nil {
		t.Error("Expected compose_content to be present")
	}
}

func TestGetRoleByIDNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/roles/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteRole(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-role","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	createReq, _ := http.NewRequest("POST", "/roles", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdRole map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdRole)
	roleID := createdRole["id"].(string)

	req, _ := http.NewRequest("DELETE", "/roles/"+roleID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	getReq, _ := http.NewRequest("GET", "/roles/"+roleID, nil)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", getW.Code)
	}
}

func TestDeleteRoleNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("DELETE", "/roles/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}
