package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListHostsEmpty(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/hosts", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var hosts []interface{}
	json.Unmarshal(w.Body.Bytes(), &hosts)

	if len(hosts) != 0 {
		t.Errorf("Expected empty array, got %d hosts", len(hosts))
	}
}

func TestCreateHost(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	body := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	req, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var host map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &host)

	if host["name"] != "test-server" {
		t.Errorf("Expected name 'test-server', got '%v'", host["name"])
	}
	if host["hostname"] != "192.168.1.100" {
		t.Errorf("Expected hostname '192.168.1.100', got '%v'", host["hostname"])
	}
}

func TestListHostsWithData(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	createReq, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	req, _ := http.NewRequest("GET", "/hosts", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var hosts []interface{}
	json.Unmarshal(w.Body.Bytes(), &hosts)

	if len(hosts) != 1 {
		t.Errorf("Expected 1 host, got %d hosts", len(hosts))
	}
}

func TestGetHostByID(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	createReq, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdHost map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdHost)
	hostID := int(createdHost["id"].(float64))

	req, _ := http.NewRequest("GET", "/hosts/1", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var host map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &host)

	if host["id"] != float64(hostID) {
		t.Errorf("Expected id %d, got %v", hostID, host["id"])
	}
}

func TestGetHostByIDNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/hosts/999", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateHost(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	createReq, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	updateBody := `{"name":"updated-server","hostname":"192.168.1.200","username":"admin","ssh_key_path":"/home/admin/.ssh/id_rsa"}`
	req, _ := http.NewRequest("PUT", "/hosts/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var host map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &host)

	if host["name"] != "updated-server" {
		t.Errorf("Expected name 'updated-server', got '%v'", host["name"])
	}
	if host["hostname"] != "192.168.1.200" {
		t.Errorf("Expected hostname '192.168.1.200', got '%v'", host["hostname"])
	}
}

func TestDeleteHost(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	createBody := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	createReq, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	req, _ := http.NewRequest("DELETE", "/hosts/1", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	getReq, _ := http.NewRequest("GET", "/hosts/1", nil)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", getW.Code)
	}
}

func TestDeleteHostNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("DELETE", "/hosts/999", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}
