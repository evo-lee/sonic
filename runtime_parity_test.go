package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type httpSnapshot struct {
	status      int
	body        string
	contentType string
	requestID   string
	locale      string
	allowOrigin string
}

func TestHertzRuntimeSmoke(t *testing.T) {
	if os.Getenv("SONIC_ENABLE_RUNTIME_SMOKE") != "1" {
		t.Skip("set SONIC_ENABLE_RUNTIME_SMOKE=1 to run the hertz runtime smoke test")
	}
	bin := buildTestBinary(t)
	snapshot := runHertzSmoke(t, bin)

	assertStatus(t, "GET /ping", snapshot["GET /ping"], http.StatusOK)
	assertStatus(t, "GET /", snapshot["GET /"], http.StatusOK)
	assertStatus(t, "GET /api/admin/is_installed", snapshot["GET /api/admin/is_installed"], http.StatusOK)
	assertStatus(t, "GET /api/content/options/comment", snapshot["GET /api/content/options/comment"], http.StatusOK)
	assertStatus(t, "GET /admin_random/", snapshot["GET /admin_random/"], http.StatusOK)
	assertStatus(t, "GET /css/app.a231e5ba.css", snapshot["GET /css/app.a231e5ba.css"], http.StatusOK)
	assertStatus(t, "OPTIONS /api/admin/login", snapshot["OPTIONS /api/admin/login"], http.StatusNoContent)
	assertStatus(t, "POST /api/admin/login", snapshot["POST /api/admin/login"], http.StatusOK)
	assertStatus(t, "GET /api/admin/users/profiles", snapshot["GET /api/admin/users/profiles"], http.StatusOK)
}

func assertStatus(t *testing.T, name string, snap httpSnapshot, status int) {
	t.Helper()
	if snap.status != status {
		t.Fatalf("status mismatch for %s: got=%d want=%d body=%s", name, snap.status, status, snap.body)
	}
	if snap.requestID == "" && strings.HasPrefix(name, "GET /api/admin/") {
		t.Fatalf("missing request id for %s", name)
	}
}

func buildTestBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "sonic-test-bin")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local", "GOPROXY=https://proxy.golang.org,direct")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build binary: %v\n%s", err, out)
	}
	return bin
}

func runHertzSmoke(t *testing.T, bin string) map[string]httpSnapshot {
	t.Helper()
	port := freePort(t)
	cfg := writeTempConfig(t, port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "-config", cfg)
	cmd.Dir = repoRoot(t)
	var log bytes.Buffer
	cmd.Stdout = &log
	cmd.Stderr = &log
	if err := cmd.Start(); err != nil {
		t.Fatalf("start hertz: %v", err)
	}
	defer func() { cancel(); _ = cmd.Wait() }()
	waitForReady(t, port, &log)

	client := &http.Client{Timeout: 5 * time.Second}
	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	paths := map[string]*http.Request{
		"GET /ping":                        mustRequest(t, http.MethodGet, base+"/ping", nil),
		"GET /":                            mustRequest(t, http.MethodGet, base+"/", nil),
		"GET /api/admin/is_installed":      mustRequest(t, http.MethodGet, base+"/api/admin/is_installed", nil),
		"GET /api/content/options/comment": mustRequest(t, http.MethodGet, base+"/api/content/options/comment", nil),
		"GET /admin_random/":               mustRequest(t, http.MethodGet, base+"/admin_random/", nil),
		"GET /css/app.a231e5ba.css":        mustRequest(t, http.MethodGet, base+"/css/app.a231e5ba.css", nil),
		"OPTIONS /api/admin/login":         mustCORSRequest(t, base+"/api/admin/login"),
	}
	out := make(map[string]httpSnapshot, len(paths)+2)
	for name, req := range paths {
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("%s request failed: %v\n%s", name, err, log.String())
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		out[name] = snapshotFromResponse(resp, body)
	}

	username := envOrDefault("SONIC_RUNTIME_SMOKE_USERNAME", "litang")
	password := envOrDefault("SONIC_RUNTIME_SMOKE_PASSWORD", "Ll3313222")
	loginBody := fmt.Sprintf(`{"username":%q,"password":%q}`, username, password)
	loginReq := mustRequest(t, http.MethodPost, base+"/api/admin/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := client.Do(loginReq)
	if err != nil {
		t.Fatalf("POST /api/admin/login request failed: %v\n%s", err, log.String())
	}
	loginRespBody, _ := io.ReadAll(loginResp.Body)
	_ = loginResp.Body.Close()
	out["POST /api/admin/login"] = snapshotFromResponse(loginResp, loginRespBody)
	accessToken := extractAccessToken(t, loginRespBody)

	profileReq := mustRequest(t, http.MethodGet, base+"/api/admin/users/profiles", nil)
	profileReq.Header.Set("Admin-Authorization", accessToken)
	profileResp, err := client.Do(profileReq)
	if err != nil {
		t.Fatalf("GET /api/admin/users/profiles request failed: %v\n%s", err, log.String())
	}
	profileRespBody, _ := io.ReadAll(profileResp.Body)
	_ = profileResp.Body.Close()
	out["GET /api/admin/users/profiles"] = snapshotFromResponse(profileResp, profileRespBody)
	return out
}

func snapshotFromResponse(resp *http.Response, body []byte) httpSnapshot {
	return httpSnapshot{
		status:      resp.StatusCode,
		body:        string(body),
		contentType: resp.Header.Get("Content-Type"),
		requestID:   resp.Header.Get("X-Request-ID"),
		locale:      resp.Header.Get("Content-Language"),
		allowOrigin: resp.Header.Get("Access-Control-Allow-Origin"),
	}
}

func extractAccessToken(t *testing.T, body []byte) string {
	t.Helper()
	var payload struct {
		Status int `json:"status"`
		Data   struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("login response is not valid JSON: %v body=%s", err, body)
	}
	if payload.Status != http.StatusOK || payload.Data.AccessToken == "" {
		t.Fatalf("login did not return an access token body=%s", body)
	}
	return payload.Data.AccessToken
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func mustRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func mustCORSRequest(t *testing.T, url string) *http.Request {
	t.Helper()
	req := mustRequest(t, http.MethodOptions, url, nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	return req
}

func waitForReady(t *testing.T, port int, log *bytes.Buffer) {
	t.Helper()
	client := &http.Client{Timeout: 1 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/ping", port)
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("server did not become ready on port %d\n%s", port, log.String())
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox does not allow binding a local port for integration tests")
		}
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func writeTempConfig(t *testing.T, port int) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "conf", "config.dev.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	content := strings.Replace(string(data), "port: 8080", fmt.Sprintf("port: %d", port), 1)
	cfg := filepath.Join(t.TempDir(), "hertz.yaml")
	if err := os.WriteFile(cfg, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return wd
}
