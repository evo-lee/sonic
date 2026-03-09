package main

import (
	"bytes"
	"context"
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

func TestFrameworkParity(t *testing.T) {
	if os.Getenv("SONIC_ENABLE_FRAMEWORK_PARITY") != "1" {
		t.Skip("set SONIC_ENABLE_FRAMEWORK_PARITY=1 to run framework parity integration test")
	}
	bin := buildTestBinary(t)
	baseline := runFrameworkSnapshot(t, bin, "gin")
	candidate := runFrameworkSnapshot(t, bin, "hertz")

	for path, want := range baseline {
		got := candidate[path]
		if got.status != want.status {
			t.Fatalf("status mismatch for %s: gin=%d hertz=%d", path, want.status, got.status)
		}
		if got.contentType != want.contentType {
			t.Fatalf("content-type mismatch for %s: gin=%q hertz=%q", path, want.contentType, got.contentType)
		}
		if got.body != want.body {
			t.Fatalf("body mismatch for %s", path)
		}
		if want.requestID != "" && got.requestID == "" {
			t.Fatalf("missing request id for %s in hertz response", path)
		}
		if got.locale != want.locale {
			t.Fatalf("locale mismatch for %s: gin=%q hertz=%q", path, want.locale, got.locale)
		}
		if got.allowOrigin != want.allowOrigin {
			t.Fatalf("cors origin mismatch for %s: gin=%q hertz=%q", path, want.allowOrigin, got.allowOrigin)
		}
	}
}

func buildTestBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "sonic-test-bin")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Env = append(os.Environ(),
		"GOTOOLCHAIN=local",
		"GOPROXY=https://proxy.golang.org,direct",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build binary: %v\n%s", err, out)
	}
	return bin
}

func runFrameworkSnapshot(t *testing.T, bin, framework string) map[string]httpSnapshot {
	t.Helper()
	port := freePort(t)
	cfg := writeTempConfig(t, framework, port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "-config", cfg)
	cmd.Dir = repoRoot(t)
	var log bytes.Buffer
	cmd.Stdout = &log
	cmd.Stderr = &log
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", framework, err)
	}
	defer func() {
		cancel()
		_ = cmd.Wait()
	}()
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
	out := make(map[string]httpSnapshot, len(paths))
	for name, req := range paths {
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("%s %s request failed: %v\n%s", framework, name, err, log.String())
		}
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		out[name] = httpSnapshot{
			status:      resp.StatusCode,
			body:        string(body),
			contentType: resp.Header.Get("Content-Type"),
			requestID:   resp.Header.Get("X-Request-ID"),
			locale:      resp.Header.Get("Content-Language"),
			allowOrigin: resp.Header.Get("Access-Control-Allow-Origin"),
		}
	}
	return out
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

func writeTempConfig(t *testing.T, framework string, port int) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repoRoot(t), "conf", "config.dev.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	content = strings.Replace(content, "framework: gin", "framework: "+framework, 1)
	content = strings.Replace(content, "port: 8080", fmt.Sprintf("port: %d", port), 1)
	cfg := filepath.Join(t.TempDir(), framework+".yaml")
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
