package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestHost(t *testing.T, ts *testServer) string {
	body := `{"name":"test-server","hostname":"192.168.1.100","username":"root","ssh_key_path":"/home/user/.ssh/id_rsa"}`
	req, _ := http.NewRequest("POST", "/hosts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	var host map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &host)
	return host["id"].(string)
}

func createTestPackage(t *testing.T, ts *testServer) string {
	body := `{"name":"test-package","compose_content":"version: '3'\nservices:\n  web:\n    image: nginx"}`
	req, _ := http.NewRequest("POST", "/packages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	var pkg map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pkg)
	return pkg["id"].(string)
}

func TestListDeploymentsEmpty(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/deployments", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var deployments []interface{}
	json.Unmarshal(w.Body.Bytes(), &deployments)

	if len(deployments) != 0 {
		t.Errorf("Expected empty array, got %d deployments", len(deployments))
	}
}

func TestCreateDeployment(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	body := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	req, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var deployment map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deployment)

	if deployment["host_id"] != hostID {
		t.Errorf("Expected host_id %s, got %v", hostID, deployment["host_id"])
	}
	if deployment["package_id"] != packageID {
		t.Errorf("Expected package_id %s, got %v", packageID, deployment["package_id"])
	}
}

func TestCreateDeploymentInvalidHost(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	packageID := createTestPackage(t, ts)

	body := `{"host_id":"550e8400-e29b-41d4-a716-446655440000","package_id":"` + packageID + `"}`
	req, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateDeploymentInvalidPackage(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)

	body := `{"host_id":"` + hostID + `","package_id":"550e8400-e29b-41d4-a716-446655440000"}`
	req, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListDeploymentsWithData(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	createBody := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	createReq, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	req, _ := http.NewRequest("GET", "/deployments", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var deployments []interface{}
	json.Unmarshal(w.Body.Bytes(), &deployments)

	if len(deployments) != 1 {
		t.Errorf("Expected 1 deployment, got %d deployments", len(deployments))
	}
}

func TestGetDeploymentByID(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	createBody := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	createReq, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdDeployment map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdDeployment)
	deploymentID := createdDeployment["id"].(string)

	req, _ := http.NewRequest("GET", "/deployments/"+deploymentID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var deployment map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deployment)

	if deployment["id"] != deploymentID {
		t.Errorf("Expected id %s, got %v", deploymentID, deployment["id"])
	}
}

func TestGetDeploymentByIDNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("GET", "/deployments/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteDeployment(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	createBody := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	createReq, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdDeployment map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdDeployment)
	deploymentID := createdDeployment["id"].(string)

	req, _ := http.NewRequest("DELETE", "/deployments/"+deploymentID, nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	getReq, _ := http.NewRequest("GET", "/deployments/"+deploymentID, nil)
	getW := httptest.NewRecorder()
	ts.router.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Errorf("Expected 404 after delete, got %d", getW.Code)
	}
}

func TestDeleteDeploymentNotFound(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	req, _ := http.NewRequest("DELETE", "/deployments/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestRunDeployment(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	createBody := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	createReq, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdDeployment map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdDeployment)
	deploymentID := createdDeployment["id"].(string)

	req, _ := http.NewRequest("POST", "/deployments/"+deploymentID+"/run", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var deployment map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deployment)

	if deployment["status"] != "running" {
		t.Errorf("Expected status 'running', got '%v'", deployment["status"])
	}
}

func TestGetDeploymentLogs(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.close()

	hostID := createTestHost(t, ts)
	packageID := createTestPackage(t, ts)

	createBody := `{"host_id":"` + hostID + `","package_id":"` + packageID + `"}`
	createReq, _ := http.NewRequest("POST", "/deployments", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	ts.router.ServeHTTP(createW, createReq)

	var createdDeployment map[string]interface{}
	json.Unmarshal(createW.Body.Bytes(), &createdDeployment)
	deploymentID := createdDeployment["id"].(string)

	req, _ := http.NewRequest("GET", "/deployments/"+deploymentID+"/logs", nil)
	w := httptest.NewRecorder()
	ts.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var logs []interface{}
	json.Unmarshal(w.Body.Bytes(), &logs)

	if len(logs) != 0 {
		t.Errorf("Expected empty logs array initially, got %d logs", len(logs))
	}
}
