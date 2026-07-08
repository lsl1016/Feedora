package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/feedora/backend/internal/dto"
	"github.com/gin-gonic/gin"
)

type testQueryRequest struct {
	dto.PageRequest
}

type testURIRequest struct {
	PostID int64 `uri:"postId" binding:"required"`
}

type testJSONRequest struct {
	Name string `json:"name" binding:"required"`
}

func TestNormalizePageRequestAppliesDefaultsAndLimit(t *testing.T) {
	req := dto.PageRequest{}

	normalizePageRequest(&req)

	if req.Page != 1 {
		t.Fatalf("expected default page to be 1, got %d", req.Page)
	}
	if req.PageSize != 10 {
		t.Fatalf("expected default pageSize to be 10, got %d", req.PageSize)
	}

	req = dto.PageRequest{Page: 3, PageSize: 999}
	normalizePageRequest(&req)
	if req.Page != 3 {
		t.Fatalf("expected page to remain 3, got %d", req.Page)
	}
	if req.PageSize != 100 {
		t.Fatalf("expected pageSize to be capped at 100, got %d", req.PageSize)
	}
}

func TestBindQueryReturnsParamErrorOnInvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test?page=abc", nil)

	var req testQueryRequest
	if bindQuery(c, &req) {
		t.Fatalf("expected bindQuery to fail for invalid query")
	}

	assertErrorParamInvalid(t, w, http.StatusBadRequest)
}

func TestBindURIReturnsParamErrorOnMissingPathValue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/posts", nil)

	var req testURIRequest
	if bindURI(c, &req) {
		t.Fatalf("expected bindURI to fail when required uri param is missing")
	}

	assertErrorParamInvalid(t, w, http.StatusBadRequest)
}

func TestBindJSONReturnsParamErrorOnInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", http.NoBody)
	c.Request.Header.Set("Content-Type", "application/json")

	var req testJSONRequest
	if bindJSON(c, &req) {
		t.Fatalf("expected bindJSON to fail for missing required field")
	}

	assertErrorParamInvalid(t, w, http.StatusBadRequest)
}

func assertErrorParamInvalid(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, w.Code)
	}

	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body.Code != 400 {
		t.Fatalf("expected error code 400, got %d", body.Code)
	}
}
