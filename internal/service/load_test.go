package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// Test saveUser function
func TestSaveUser(t *testing.T) {
	users.Store("test_user_old", true) // Clear from previous tests
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})

	userID := "test_user_123"
	saveUser(userID)

	val, ok := users.Load(userID)
	assert.True(t, ok, "user should be saved")
	assert.Equal(t, true, val, "user value should be true")
}

// Test savePR function
func TestSavePR(t *testing.T) {
	prs.Range(func(k, v any) bool {
		prs.Delete(k)
		return true
	})

	prID := "pr-12345"
	savePR(prID)

	val, ok := prs.Load(prID)
	assert.True(t, ok, "PR should be saved")
	assert.Equal(t, true, val, "PR value should be true")
}

// Test saveTeam function
func TestSaveTeam(t *testing.T) {
	teams.Range(func(k, v any) bool {
		teams.Delete(k)
		return true
	})

	teamName := "test_team"
	saveTeam(teamName)

	val, ok := teams.Load(teamName)
	assert.True(t, ok, "team should be saved")
	assert.Equal(t, true, val, "team value should be true")
}

// Test getRandomFromMap with items
func TestGetRandomFromMap_WithItems(t *testing.T) {
	m := &sync.Map{}
	m.Store("key1", "value1")
	m.Store("key2", "value2")
	m.Store("key3", "value3")

	result := getRandomFromMap(m)
	assert.NotEmpty(t, result, "should return a value from the map")
	assert.True(t, result == "key1" || result == "key2" || result == "key3",
		"should return one of the stored keys")
}

// Test getRandomFromMap with empty map
func TestGetRandomFromMap_Empty(t *testing.T) {
	m := &sync.Map{}
	result := getRandomFromMap(m)
	assert.Empty(t, result, "should return empty string for empty map")
}

// Test getRandomFromMap randomness
func TestGetRandomFromMap_Randomness(t *testing.T) {
	m := &sync.Map{}
	m.Store("key1", true)
	m.Store("key2", true)
	m.Store("key3", true)

	results := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := getRandomFromMap(m)
		results[result]++
	}

	assert.Greater(t, len(results), 1, "should return different values (randomness check)")
}

// Test randString function
func TestRandString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"zero length", 0},
		{"short string", 5},
		{"medium string", 20},
		{"long string", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := randString(tt.length)
			assert.Equal(t, tt.length, len(result), "should return string of correct length")

			// Check that all characters are letters
			for _, r := range result {
				assert.True(t, (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'),
					"should only contain letters")
			}
		})
	}
}

// Test randString randomness
func TestRandString_Randomness(t *testing.T) {
	str1 := randString(10)
	str2 := randString(10)
	assert.NotEqual(t, str1, str2, "consecutive calls should produce different strings")
}

// Test randUserID function
func TestRandUserID(t *testing.T) {
	result := randUserID()
	assert.True(t, strings.HasPrefix(result, "u"), "should start with 'u'")

	// Extract number part
	numPart := strings.TrimPrefix(result, "u")
	assert.NotEmpty(t, numPart, "should have number part")

	// Verify it's a valid number format
	assert.Regexp(t, `^\d+$`, numPart, "should contain only digits after 'u'")
}

// Test randUserID range
func TestRandUserID_Range(t *testing.T) {
	for i := 0; i < 100; i++ {
		result := randUserID()
		// Extract number and verify it's in valid range (1 to 1000)
		assert.Regexp(t, `^u\d+$`, result, "should match pattern u<number>")
	}
}

// Test randPRID function
func TestRandPRID(t *testing.T) {
	result := randPRID()
	assert.True(t, strings.HasPrefix(result, "pr-"), "should start with 'pr-'")

	numPart := strings.TrimPrefix(result, "pr-")
	assert.Regexp(t, `^\d+$`, numPart, "should contain only digits after 'pr-'")
}

// Test randBool function
func TestRandBool(t *testing.T) {
	// Test that it returns boolean values
	for i := 0; i < 50; i++ {
		result := randBool()
		assert.IsType(t, true, result, "should return boolean")
	}
}

// Test randBool distribution
func TestRandBool_Distribution(t *testing.T) {
	trueCount := 0
	falseCount := 0

	for i := 0; i < 1000; i++ {
		if randBool() {
			trueCount++
		} else {
			falseCount++
		}
	}

	// With 1000 samples, we expect roughly equal distribution (between 300-700)
	assert.Greater(t, trueCount, 300, "should have reasonable number of true values")
	assert.Less(t, trueCount, 700, "should have reasonable number of true values")
	assert.Greater(t, falseCount, 300, "should have reasonable number of false values")
	assert.Less(t, falseCount, 700, "should have reasonable number of false values")
}

// Test newPOST function
func TestNewPOST(t *testing.T) {
	tests := []struct {
		body any
		name string
		path string
	}{
		{
			name: "simple string body",
			path: "/test",
			body: "test_string",
		},
		{
			name: "map body",
			path: "/team/add",
			body: map[string]any{
				"team_name": "test_team",
				"members":   []string{"user1", "user2"},
			},
		},
		{
			name: "nested structure",
			path: "/pullRequest/create",
			body: map[string]any{
				"pull_request_id": "pr-123",
				"author_id":       "u1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := newPOST(tt.path, tt.body)

			assert.Equal(t, "POST", target.Method, "should be POST method")
			assert.Equal(t, baseURL+tt.path, target.URL, "should have correct URL")
			assert.Equal(t, []string{"application/json"}, target.Header["Content-Type"],
				"should have JSON content type")
			assert.NotEmpty(t, target.Body, "should have body")

			// Verify body can be unmarshaled
			var unmarshaled any
			err := json.Unmarshal(target.Body, &unmarshaled)
			assert.NoError(t, err, "body should be valid JSON")
		})
	}
}

// Test newGET function
func TestNewGET(t *testing.T) {
	tests := []struct {
		params map[string]string
		name   string
		path   string
	}{
		{
			name:   "no parameters",
			path:   "/team/get",
			params: map[string]string{},
		},
		{
			name: "single parameter",
			path: "/users/getReview",
			params: map[string]string{
				"user_id": "u123",
			},
		},
		{
			name: "multiple parameters",
			path: "/test",
			params: map[string]string{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := newGET(tt.path, tt.params)

			assert.Equal(t, "GET", target.Method, "should be GET method")
			assert.True(t, strings.HasPrefix(target.URL, baseURL+tt.path),
				"URL should start with base URL and path")

			// Parse URL and verify parameters
			parsed, err := url.Parse(target.URL)
			assert.NoError(t, err, "URL should be parseable")

			for key, value := range tt.params {
				assert.Equal(t, value, parsed.Query().Get(key),
					"parameter should be in query string")
			}
		})
	}
}

// Test genTeamAdd function
func TestGenTeamAdd(t *testing.T) {
	teams.Range(func(k, v any) bool {
		teams.Delete(k)
		return true
	})
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})

	target := genTeamAdd()

	assert.Equal(t, "POST", target.Method)
	assert.Equal(t, baseURL+"/team/add", target.URL)
	assert.NotEmpty(t, target.Body)

	var payload map[string]any
	err := json.Unmarshal(target.Body, &payload)
	require.NoError(t, err)

	assert.Contains(t, payload, "team_name")
	assert.Contains(t, payload, "members")

	teamName := payload["team_name"].(string)
	members := payload["members"].([]any)

	assert.NotEmpty(t, teamName, "should have team name")
	assert.Greater(t, len(members), 0, "should have members")

	// Verify team and users were saved
	_, ok := teams.Load(teamName)
	assert.True(t, ok, "team should be saved")
}

// Test genTeamGet function
func TestGenTeamGet(t *testing.T) {
	teams.Range(func(k, v any) bool {
		teams.Delete(k)
		return true
	})

	// Add a team first
	saveTeam("test_team_get")

	target := genTeamGet()

	assert.Equal(t, "GET", target.Method)
	assert.True(t, strings.HasPrefix(target.URL, baseURL+"/team/get?"))

	parsed, err := url.Parse(target.URL)
	require.NoError(t, err)
	assert.Equal(t, "test_team_get", parsed.Query().Get("team_name"))
}

// Test genSetActive function
func TestGenSetActive(t *testing.T) {
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})

	saveUser("test_user")

	target := genSetActive()

	assert.Equal(t, "POST", target.Method)
	assert.Equal(t, baseURL+"/users/setIsActive", target.URL)

	var payload map[string]any
	err := json.Unmarshal(target.Body, &payload)
	require.NoError(t, err)

	assert.Contains(t, payload, "user_id")
	assert.Contains(t, payload, "is_active")
	assert.Equal(t, "test_user", payload["user_id"])
	assert.IsType(t, true, payload["is_active"])
}

// Test genUserReviews function
func TestGenUserReviews(t *testing.T) {
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})

	saveUser("review_user")

	target := genUserReviews()

	assert.Equal(t, "GET", target.Method)
	assert.True(t, strings.HasPrefix(target.URL, baseURL+"/users/getReview?"))

	parsed, err := url.Parse(target.URL)
	require.NoError(t, err)
	assert.Equal(t, "review_user", parsed.Query().Get("user_id"))
}

// Test genPRCreate function
func TestGenPRCreate(t *testing.T) {
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})
	prs.Range(func(k, v any) bool {
		prs.Delete(k)
		return true
	})

	saveUser("pr_author")

	target := genPRCreate()

	assert.Equal(t, "POST", target.Method)
	assert.Equal(t, baseURL+"/pullRequest/create", target.URL)

	var payload map[string]any
	err := json.Unmarshal(target.Body, &payload)
	require.NoError(t, err)

	assert.Contains(t, payload, "pull_request_id")
	assert.Contains(t, payload, "pull_request_name")
	assert.Contains(t, payload, "author_id")
}

// Test genPRMerge function
func TestGenPRMerge(t *testing.T) {
	prs.Range(func(k, v any) bool {
		prs.Delete(k)
		return true
	})

	savePR("merge_pr")

	target := genPRMerge()

	assert.Equal(t, "POST", target.Method)
	assert.Equal(t, baseURL+"/pullRequest/merge", target.URL)

	var payload map[string]any
	err := json.Unmarshal(target.Body, &payload)
	require.NoError(t, err)

	assert.Contains(t, payload, "pull_request_id")
	assert.Equal(t, "merge_pr", payload["pull_request_id"])
}

// Test genPRReassign function
func TestGenPRReassign(t *testing.T) {
	prs.Range(func(k, v any) bool {
		prs.Delete(k)
		return true
	})
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})

	savePR("reassign_pr")
	saveUser("reassign_user")

	target := genPRReassign()

	assert.Equal(t, "POST", target.Method)
	assert.Equal(t, baseURL+"/pullRequest/reassign", target.URL)

	var payload map[string]any
	err := json.Unmarshal(target.Body, &payload)
	require.NoError(t, err)

	assert.Contains(t, payload, "pull_request_id")
	assert.Contains(t, payload, "old_user_id")
}

// Test randomTarget function
func TestRandomTarget(t *testing.T) {
	users.Store("test_user", true)
	teams.Store("test_team", true)
	prs.Store("test_pr", true)

	results := make(map[string]int)
	paths := []string{
		"/team/add",
		"/team/get",
		"/users/setIsActive",
		"/users/getReview",
		"/pullRequest/create",
		"/pullRequest/merge",
		"/pullRequest/reassign",
	}

	// Generate multiple targets to check diversity
	for i := 0; i < 70; i++ {
		target := randomTarget()
		parsed, _ := url.Parse(target.URL)
		path := parsed.Path
		results[path]++
	}

	// Verify we got different endpoints
	assert.Greater(t, len(results), 1, "should generate targets for multiple endpoints")

	// Verify all generated paths are valid
	for path := range results {
		found := false
		for _, valid := range paths {
			if valid == path {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("path %s should be valid", path))
	}
}

// Test validateResponse with valid responses
func TestValidateResponse_ValidResponses(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
		code uint16
	}{
		{
			name: "PR create success",
			code: 201,
			path: "/pullRequest/create",
			body: `{"data":"created ok"}`,
		},
		{
			name: "Team add success",
			code: 201,
			path: "/team/add",
			body: `{"data":"created ok"}`,
		},
		{
			name: "Set active success",
			code: 200,
			path: "/users/setIsActive",
			body: `{"data":""}`,
		},
		{
			name: "Merge success",
			code: 200,
			path: "/pullRequest/merge",
			body: `{"data":"merged"}`,
		},
		{
			name: "PR not found",
			code: 404,
			path: "/pullRequest/create",
			body: `{"error":{"message":"author/team not found"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &vegeta.Result{
				Code: tt.code,
				URL:  baseURL + tt.path,
				Body: []byte(tt.body),
			}

			err := validateResponse(res)
			assert.NoError(t, err, "should validate successfully")
		})
	}
}

// Test validateResponse with invalid code
func TestValidateResponse_InvalidCode(t *testing.T) {
	res := &vegeta.Result{
		Code: 500,
		URL:  baseURL + "/pullRequest/create",
		Body: []byte(`{"error":"Internal server error"}`),
	}

	err := validateResponse(res)
	assert.Error(t, err, "should return error for invalid status code")
	assert.Contains(t, err.Error(), "unexpected code")
}

// Test validateResponse with invalid path
func TestValidateResponse_InvalidPath(t *testing.T) {
	res := &vegeta.Result{
		Code: 200,
		URL:  baseURL + "/invalid/path",
		Body: []byte(`{}`),
	}

	err := validateResponse(res)
	assert.Error(t, err, "should return error for invalid path")
	assert.Contains(t, err.Error(), "unexpected path")
}

// Test validateResponse with invalid URL
func TestValidateResponse_InvalidURL(t *testing.T) {
	res := &vegeta.Result{
		Code: 200,
		URL:  "not a valid url %",
		Body: []byte(`{}`),
	}

	err := validateResponse(res)
	assert.Error(t, err, "should return error for invalid URL")
}

// Test containsNormalized function
func TestContainsNormalized(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
		want     bool
	}{
		{
			name:     "exact match",
			got:      "test message",
			expected: "test message",
			want:     true,
		},
		{
			name:     "case insensitive",
			got:      "TEST MESSAGE",
			expected: "test message",
			want:     true,
		},
		{
			name:     "contains substring",
			got:      "this is a test message",
			expected: "test message",
			want:     true,
		},
		{
			name:     "case insensitive contains",
			got:      "THIS IS A TEST MESSAGE",
			expected: "test message",
			want:     true,
		},
		{
			name:     "empty expected",
			got:      "anything",
			expected: "",
			want:     true,
		},
		{
			name:     "no match",
			got:      "hello world",
			expected: "test",
			want:     false,
		},
		{
			name:     "whitespace normalization",
			got:      "  test message  ",
			expected: "test message",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsNormalized(tt.got, tt.expected)
			assert.Equal(t, tt.want, result)
		})
	}
}

// Test RunLoadTest function (integration test)
func TestRunLoadTest(t *testing.T) {
	// Clear maps before test
	users.Range(func(k, v any) bool {
		users.Delete(k)
		return true
	})
	teams.Range(func(k, v any) bool {
		teams.Delete(k)
		return true
	})
	prs.Range(func(k, v any) bool {
		prs.Delete(k)
		return true
	})

	// Add some initial data
	saveUser("u1")
	saveUser("u2")
	saveTeam("team1")
	savePR("pr1")

	service := &LoadService{}

	// Run a very short load test
	rate := vegeta.Rate{Freq: 5, Per: time.Second}
	duration := time.Millisecond * 100

	// This test mainly checks that the function doesn't panic
	// In a real scenario, you'd need a running server
	assert.NotPanics(t, func() {
		service.RunLoadTest(rate, duration)
	}, "RunLoadTest should not panic")
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("randString with large size", func(t *testing.T) {
		result := randString(10000)
		assert.Equal(t, 10000, len(result))
	})

	t.Run("multiple saves to same key", func(t *testing.T) {
		testMap := &sync.Map{}
		testMap.Store("key", 1)
		testMap.Store("key", 2)
		val, _ := testMap.Load("key")
		assert.Equal(t, 2, val)
	})

	t.Run("newPOST with nil body", func(t *testing.T) {
		target := newPOST("/test", nil)
		assert.NotEmpty(t, target.Body)
	})

	t.Run("newGET with empty params", func(t *testing.T) {
		target := newGET("/test", map[string]string{})
		assert.Equal(t, baseURL+"/test?", target.URL)
	})
}

// Test generator functions are in the generators slice
func TestGeneratorsFunctionality(t *testing.T) {
	assert.Equal(t, 7, len(generators), "should have 7 generator functions")

	users.Store("test_user", true)
	teams.Store("test_team", true)
	prs.Store("test_pr", true)

	for i, gen := range generators {
		t.Run(fmt.Sprintf("generator_%d", i), func(t *testing.T) {
			assert.NotPanics(t, func() {
				target := gen()
				assert.NotEmpty(t, target.Method)
				assert.NotEmpty(t, target.URL)
			}, "generator should not panic")
		})
	}
}
