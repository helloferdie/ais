package handler

import (
	"ais/lib/libresponse"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockArticle = map[string]interface{}{
	"author": "O'Reilly Media, Inc",
	"title":  "Introducing Go",
	"body":   "Perfect for beginners familiar with programming basics, this hands-on guide provides an easy introduction to Go, the general-purpose programming language from Google. Author Caleb Doxsey covers the language’s core features with step-by-step instructions and exercises in each chapter to help you practice what you learn. By the time you finish this book, not only will you be able to write real Go programs, you'll be ready to tackle advanced techniques.",
}

var mockArticleJSON, _ = json.Marshal(mockArticle)

type mockArticleResponse struct {
	ID        int64  `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type mockListArticleResponse struct {
	Items      []mockArticleResponse `json:"items"`
	TotalItems int64                 `json:"total_items"`
	TotalPages int64                 `json:"total_pages"`
}

// Load environment variable for test
func init() {
	godotenv.Load("../.env.test")
}

// testArticleCreate - Reuse test to insert mock data to database
func testArticleCreate(t *testing.T) int64 {
	// Setup
	mockJSON, err := json.Marshal(mockArticle)
	if err != nil {
		t.Error("Cannot create json from mock data")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(mockJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	require.NoError(t, ArticlePost(c))
	return testVerifyArticleResponse(t, rec, mockArticle)
}

// testVerifyArticleResponse - Reuse test to verify response body with expected data
func testVerifyArticleResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedData map[string]interface{}) int64 {
	require.Equal(t, http.StatusOK, rec.Code)

	// Check valid response format
	resp := new(libresponse.Response)
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Error("Not valid response format")
		t.FailNow()
	}

	// Compare mock data with response data
	respData := new(mockArticleResponse)
	b, _ := json.Marshal(resp.Data)
	err = json.Unmarshal(b, &respData)
	if err != nil {
		t.Error("Not valid mock response format")
		t.FailNow()
	}

	respMapData := resp.Data.(map[string]interface{})
	for k, v := range expectedData {
		assert.EqualValues(t, v, respMapData[k], "Key :"+k)
	}
	if t.Failed() {
		t.Error("Mock data does not match with response data")
		t.FailNow()
	}

	return respData.ID
}

// testArticleList - Reuse test to list articles
func testArticleList(t *testing.T, params map[string]string) *mockListArticleResponse {
	// Setup
	q := make(url.Values)
	for k, v := range params {
		q.Set(k, v)
	}

	// Set default page if not exist
	_, ok := params["page"]
	if !ok {
		q.Set("page", "1")
	}

	// Set default items_per_page if not exist
	_, ok = params["items_per_page"]
	if !ok {
		q.Set("items_per_page", "100")
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	require.NoError(t, ArticleList(c))
	require.Equal(t, http.StatusOK, rec.Code)

	// Check valid response format
	resp := new(libresponse.Response)
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Error("Not valid response format")
		t.FailNow()
	}

	respData := new(mockListArticleResponse)
	b, _ := json.Marshal(resp.Data)
	err = json.Unmarshal(b, &respData)
	if err != nil {
		t.Error("Not valid mock response format")
		t.FailNow()
	}

	t.Logf("Test list params %v return %v items", params, respData.TotalItems)
	return respData
}

func TestArticlePost(t *testing.T) {
	testArticleCreate(t)
}

func TestArticleView(t *testing.T) {
	// Setup
	insertedID := testArticleCreate(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(insertedID, 10))

	// Test
	require.NoError(t, ArticleView(c))
	responseID := testVerifyArticleResponse(t, rec, mockArticle)

	// Make sure article ID is correct
	require.Equal(t, insertedID, responseID)
}

func TestArticleList(t *testing.T) {
	// Search articles from author "ImMockAuthor"
	params := map[string]string{
		"author": "ImMockAuthor",
	}
	resp := testArticleList(t, params)
	assert.Equal(t, int64(0), resp.TotalItems, params) // Test for no data

	// Search articles from author "Reilly"
	testArticleCreate(t)
	params = map[string]string{
		"author": "Reilly",
	}
	resp = testArticleList(t, params)
	assert.Greater(t, resp.TotalItems, int64(0), params) // Test for return more than 1 data
	if len(resp.Items) > 1 {
		// Test for data sort by recent created_at
		tsArticle0, err := time.Parse(time.RFC3339, resp.Items[0].CreatedAt)
		if err != nil {
			t.Errorf("Not valid timestamp format %s", resp.Items[0].CreatedAt)
			t.FailNow()
		}

		tsArticle1, err := time.Parse(time.RFC3339, resp.Items[1].CreatedAt)
		if err != nil {
			t.Errorf("Not valid timestamp format %s", resp.Items[1].CreatedAt)
			t.FailNow()
		}
		assert.GreaterOrEqual(t, tsArticle0, tsArticle1)
	}

	// Search articles from with keywords "Reilly"
	params = map[string]string{
		"query": "core practice programming",
	}
	resp = testArticleList(t, params)
	assert.Greater(t, resp.TotalItems, int64(0), params) // Test for return more than 1 data
	if len(resp.Items) > 1 {
		// Test for data sort by recent created_at
		tsArticle0, err := time.Parse(time.RFC3339, resp.Items[0].CreatedAt)
		if err != nil {
			t.Errorf("Not valid timestamp format %s", resp.Items[0].CreatedAt)
			t.FailNow()
		}

		tsArticle1, err := time.Parse(time.RFC3339, resp.Items[1].CreatedAt)
		if err != nil {
			t.Errorf("Not valid timestamp format %s", resp.Items[1].CreatedAt)
			t.FailNow()
		}
		assert.GreaterOrEqual(t, tsArticle0, tsArticle1)
	}
}

func TestArticleUpdate(t *testing.T) {
	// Setup
	insertedID := testArticleCreate(t)

	e := echo.New()
	mockUpdate := map[string]interface{}{
		"id":     insertedID,
		"title":  "New title",
		"author": "Gopher",
		"body":   "Hello world",
	}
	upJSON, _ := json.Marshal(mockUpdate)
	upReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(upJSON)))
	upReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	upRec := httptest.NewRecorder()
	upCtx := e.NewContext(upReq, upRec)

	// Test
	require.NoError(t, ArticleUpdate(upCtx))

	// Verify changes updated
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(insertedID, 10))
	require.NoError(t, ArticleView(c))
	testVerifyArticleResponse(t, rec, mockUpdate)
}

func TestArticleDelete(t *testing.T) {
	// Setup
	insertedID := testArticleCreate(t)

	// Instance for view
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(insertedID, 10))

	// Instance for delete
	delJSON, _ := json.Marshal(map[string]interface{}{
		"id": insertedID,
	})
	delReq := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(delJSON)))
	delReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	delRec := httptest.NewRecorder()
	delCtx := e.NewContext(delReq, delRec)

	// Test
	require.NoError(t, ArticleView(c))

	// Do delete
	require.NoError(t, ArticleDelete(delCtx))
	require.Equal(t, http.StatusOK, delRec.Code)

	// Test view again, expected 404
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(insertedID, 10))
	require.NoError(t, ArticleView(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
}
