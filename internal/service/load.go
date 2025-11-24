package service

import (
	"encoding/json"
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"
)

type LoadService struct{}

const baseURL = "http://localhost:8080"

var (
	users sync.Map // key: user_id string, value: bool
	teams sync.Map // key: team_name string, value: bool
	prs   sync.Map // key: pr_id string, value: bool
)

var expected = map[string]map[int][]string{
	"/pullRequest/create": {
		201: {"created ok"},
		404: {"author/team not found"},
	},
	"/team/add": {
		201: {"created ok"},
	},
	"/users/setIsActive": {
		200: {""},
	},
	"/pullRequest/merge": {
		200: {"merged"},
		404: {"pr not found"},
	},
	"/pullRequest/reassign": {
		200: {"reassigned"},
		404: {"pr/user not found"},
	},
	"/team/get": {
		200: {"ok"},
		404: {"team not found"},
	},
	"/users/getReview": {
		200: {"ok"},
		404: {"user not found"},
	},
}

func saveUser(id string)   { users.Store(id, true) }
func savePR(id string)     { prs.Store(id, true) }
func saveTeam(name string) { teams.Store(name, true) }

func getRandomFromMap(m *sync.Map) string {
	keys := make([]string, 0)
	m.Range(func(k, v any) bool {
		keys = append(keys, k.(string))
		return true
	})
	if len(keys) == 0 {
		return ""
	}
	return keys[rand.Intn(len(keys))]
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randUserID() string {
	return fmt.Sprintf("u%d", rand.Intn(1000)+1)
}

func randPRID() string {
	return fmt.Sprintf("pr-%d", rand.Intn(100000))
}

func randBool() bool {
	return rand.Intn(2) == 1
}

func newPOST(path string, body any) vegeta.Target {
	b, _ := json.Marshal(body)
	return vegeta.Target{
		Method: "POST",
		URL:    baseURL + path,
		Header: map[string][]string{"Content-Type": {"application/json"}},
		Body:   b,
	}
}

func newGET(path string, params map[string]string) vegeta.Target {
	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	return vegeta.Target{
		Method: "GET",
		URL:    baseURL + path + "?" + q.Encode(),
	}
}

func genTeamAdd() vegeta.Target {
	n := rand.Intn(4) + 2
	members := make([]map[string]any, n)

	teamName := "team_" + randString(5)
	saveTeam(teamName)

	for i := 0; i < n; i++ {
		uid := randUserID()
		saveUser(uid)

		members[i] = map[string]any{
			"user_id":   uid,
			"username":  randString(6),
			"is_active": randBool(),
		}
	}

	body := map[string]any{
		"team_name": teamName,
		"members":   members,
	}
	return newPOST("/team/add", body)
}

func genTeamGet() vegeta.Target {
	team := getRandomFromMap(&teams)

	return newGET("/team/get", map[string]string{
		"team_name": team,
	})
}

func genSetActive() vegeta.Target {
	user := getRandomFromMap(&users)

	return newPOST("/users/setIsActive", map[string]any{
		"user_id":   user,
		"is_active": randBool(),
	})
}

func genUserReviews() vegeta.Target {
	user := getRandomFromMap(&users)

	return newGET("/users/getReview", map[string]string{
		"user_id": user,
	})
}

func genPRCreate() vegeta.Target {
	user := getRandomFromMap(&users)
	prID := randPRID()
	savePR(prID)

	body := map[string]any{
		"pull_request_id":   prID,
		"pull_request_name": "pr_" + randString(7),
		"author_id":         user,
	}
	return newPOST("/pullRequest/create", body)
}

func genPRMerge() vegeta.Target {
	pr := getRandomFromMap(&prs)

	return newPOST("/pullRequest/merge", map[string]any{
		"pull_request_id": pr,
	})
}

func genPRReassign() vegeta.Target {
	pr := getRandomFromMap(&prs)
	user := getRandomFromMap(&users)

	return newPOST("/pullRequest/reassign", map[string]any{
		"pull_request_id": pr,
		"old_user_id":     user,
	})
}

var generators = []func() vegeta.Target{
	genTeamAdd,
	genTeamGet,
	genSetActive,
	genUserReviews,
	genPRCreate,
	genPRMerge,
	genPRReassign,
}

func randomTarget() vegeta.Target {
	return generators[rand.Intn(len(generators))]()
}

func validateResponse(res *vegeta.Result) error {
	parsed, err := url.Parse(res.URL)
	if err != nil {
		return fmt.Errorf("bad url: %v", err)
	}

	path := parsed.Path
	allows, ok := expected[path]
	if !ok {
		return fmt.Errorf("unexpected path %s", path)
	}

	code := int(res.Code)
	msgs, ok := allows[code]
	if !ok {
		return fmt.Errorf("[unexpected code] %d for %s; body=%s", code, path, string(res.Body))
	}

	var payload map[string]any
	var msg string

	if err := json.Unmarshal(res.Body, &payload); err == nil {
		if errObj, ok := payload["error"].(map[string]any); ok {
			if m, ok := errObj["message"].(string); ok {
				msg = m
			}
		}
	}

	if msg == "" {
		msg = string(res.Body)
	}

	for _, expectedMsg := range msgs {
		if expectedMsg == "" || containsNormalized(msg, expectedMsg) {
			return nil
		}
	}

	return fmt.Errorf("[unexpected msg] %d %s → '%s' (expected one of: %v)",
		code, path, msg, msgs)
}

func containsNormalized(got, expected string) bool {
	g := strings.ToLower(strings.TrimSpace(got))
	e := strings.ToLower(strings.TrimSpace(expected))
	return strings.Contains(g, e)
}

func (s *LoadService) RunLoadTest(rate vegeta.Rate, duration time.Duration) {
	targeter := vegeta.Targeter(func(t *vegeta.Target) error {
		*t = randomTarget()

		if len(t.Body) > 0 {
			fmt.Printf("REQUEST: %s %s\nBody: %s\n\n", t.Method, t.URL, string(t.Body))
		} else {
			fmt.Printf("REQUEST: %s %s\n\n", t.Method, t.URL)
		}

		return nil
	})

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, duration, "mixed-load") {
		fmt.Printf("RESPONSE [%d] %s\nBody: %s\n\n",
			res.Code, res.URL, string(res.Body))

		if err := validateResponse(res); err != nil {
			fmt.Printf("❌ Validation error: %v\n\n", err)
		}

		metrics.Add(res)
	}

	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latency P99: %s\n", metrics.Latencies.P99)
}
