package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesConfiguredAndDefaultValues(t *testing.T) {
	t.Setenv(EnvAllowedRepositories, "github.com/CyberT33N/git-governance,github.example/acme/release")
	t.Setenv(EnvGitHubAppID, "42")
	t.Setenv(EnvGitHubInstallationID, "99")
	t.Setenv(EnvGitHubPrivateKeyPath, "/var/run/key.pem")
	t.Setenv(EnvPort, "")
	t.Setenv(EnvGitHubAPIBaseURL, "")
	t.Setenv(EnvRequestTimeout, "")
	t.Setenv(EnvMaxRequestBytes, "")
	t.Setenv(EnvMinimumTokenLifetime, "")

	configuration, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if configuration.Port != defaultPort {
		t.Fatalf("Port = %q, want %q", configuration.Port, defaultPort)
	}
	if !configuration.RepositoryAllowed(Repository{Host: "GITHUB.COM", Owner: "CyberT33N", Name: "git-governance"}) {
		t.Fatal("RepositoryAllowed() = false, want true")
	}
	if configuration.GitHubAPIBaseURL != defaultGitHubAPIBaseURL {
		t.Fatalf("GitHubAPIBaseURL = %q, want %q", configuration.GitHubAPIBaseURL, defaultGitHubAPIBaseURL)
	}
	if configuration.RequestTimeout != defaultRequestTimeout {
		t.Fatalf("RequestTimeout = %s, want %s", configuration.RequestTimeout, defaultRequestTimeout)
	}
	if configuration.MaxRequestBytes != defaultMaxRequestBytes {
		t.Fatalf("MaxRequestBytes = %d, want %d", configuration.MaxRequestBytes, defaultMaxRequestBytes)
	}
	if configuration.MinimumTokenLifetime != defaultMinimumTokenLifetime {
		t.Fatalf("MinimumTokenLifetime = %s, want %s", configuration.MinimumTokenLifetime, defaultMinimumTokenLifetime)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	valid := map[string]string{
		EnvAllowedRepositories:  "github.com/CyberT33N/git-governance",
		EnvGitHubAppID:          "1",
		EnvGitHubInstallationID: "2",
		EnvGitHubPrivateKeyPath: "/key.pem",
	}

	for name, value := range map[string]string{
		EnvPort:                 "70000",
		EnvAllowedRepositories:  "invalid",
		EnvGitHubAppID:          "0",
		EnvGitHubInstallationID: "-1",
		EnvGitHubPrivateKeyPath: " ",
		EnvGitHubAPIBaseURL:     "http://api.github.com",
		EnvRequestTimeout:       "-1s",
		EnvMaxRequestBytes:      "0",
		EnvMinimumTokenLifetime: "nonsense",
	} {
		t.Run(name, func(t *testing.T) {
			environment := cloneEnvironment(valid)
			environment[name] = value
			_, err := load(func(key string) string { return environment[key] })
			if err == nil {
				t.Fatalf("load() error = nil for %s=%q", name, value)
			}
		})
	}

	for _, name := range []string{
		EnvAllowedRepositories,
		EnvGitHubAppID,
		EnvGitHubInstallationID,
		EnvGitHubPrivateKeyPath,
	} {
		t.Run("missing "+name, func(t *testing.T) {
			environment := cloneEnvironment(valid)
			delete(environment, name)
			_, err := load(func(key string) string { return environment[key] })
			if err == nil || !strings.Contains(err.Error(), name) {
				t.Fatalf("load() error = %v, want error mentioning %s", err, name)
			}
		})
	}
}

func TestPortValue(t *testing.T) {
	for _, testCase := range []struct {
		name string
		raw  string
		want string
		ok   bool
	}{
		{name: "default", raw: "", want: defaultPort, ok: true},
		{name: "valid", raw: " 443 ", want: "443", ok: true},
		{name: "zero", raw: "0"},
		{name: "too high", raw: "65536"},
		{name: "non numeric", raw: "abc"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := portValue(testCase.raw)
			if (err == nil) != testCase.ok {
				t.Fatalf("portValue(%q) error = %v, want success %t", testCase.raw, err, testCase.ok)
			}
			if got != testCase.want {
				t.Fatalf("portValue(%q) = %q, want %q", testCase.raw, got, testCase.want)
			}
		})
	}
}

func TestRepositoryParsing(t *testing.T) {
	repository, err := parseRepository(" GitHub.COM /CyberT33N/git-governance ")
	if err != nil {
		t.Fatalf("parseRepository() error = %v", err)
	}
	if repository != (Repository{Host: "github.com", Owner: "CyberT33N", Name: "git-governance"}) {
		t.Fatalf("parseRepository() = %#v", repository)
	}

	for _, raw := range []string{
		"",
		"one/two",
		"one/two/three/four",
		"github.com//repo",
		"github.com/owner/",
		"github .com/owner/repo",
		"github.com/owner name/repo",
		"github.com/owner/repo name",
	} {
		if _, err := parseRepository(raw); err == nil {
			t.Errorf("parseRepository(%q) error = nil", raw)
		}
	}

	repositories, err := repositoriesValue("github.com/CyberT33N/git-governance,github.com/CyberT33N/git-governance")
	if err != nil {
		t.Fatalf("repositoriesValue() error = %v", err)
	}
	if len(repositories) != 1 {
		t.Fatalf("len(repositories) = %d, want 1", len(repositories))
	}
	if _, err := repositoriesValue(" "); err == nil {
		t.Fatal("repositoriesValue(empty) error = nil")
	}
}

func TestValueHelpers(t *testing.T) {
	if got, err := positiveIntegerValue("TEST", " 42 "); err != nil || got != "42" {
		t.Fatalf("positiveIntegerValue() = %q, %v", got, err)
	}
	for _, raw := range []string{"", "0", "-1", "x"} {
		if _, err := positiveIntegerValue("TEST", raw); err == nil {
			t.Errorf("positiveIntegerValue(%q) error = nil", raw)
		}
	}

	if got, err := httpsURLValue("URL", "https://api.example.test/base/", ""); err != nil || got != "https://api.example.test/base" {
		t.Fatalf("httpsURLValue() = %q, %v", got, err)
	}
	if got, err := httpsURLValue("URL", "", "https://fallback.example"); err != nil || got != "https://fallback.example" {
		t.Fatalf("httpsURLValue(fallback) = %q, %v", got, err)
	}
	for _, raw := range []string{
		"http://example.test",
		"https://user@example.test",
		"https://example.test?query=yes",
		"https://example.test#fragment",
		"::invalid::",
	} {
		if _, err := httpsURLValue("URL", raw, ""); err == nil {
			t.Errorf("httpsURLValue(%q) error = nil", raw)
		}
	}

	if got, err := durationValue("DURATION", "3s", time.Second); err != nil || got != 3*time.Second {
		t.Fatalf("durationValue() = %s, %v", got, err)
	}
	if got, err := durationValue("DURATION", "", time.Second); err != nil || got != time.Second {
		t.Fatalf("durationValue(fallback) = %s, %v", got, err)
	}
	for _, raw := range []string{"0s", "-1s", "bad"} {
		if _, err := durationValue("DURATION", raw, time.Second); err == nil {
			t.Errorf("durationValue(%q) error = nil", raw)
		}
	}

	if got, err := positiveInt64Value("SIZE", "2", 1); err != nil || got != 2 {
		t.Fatalf("positiveInt64Value() = %d, %v", got, err)
	}
	if got, err := positiveInt64Value("SIZE", "", 1); err != nil || got != 1 {
		t.Fatalf("positiveInt64Value(fallback) = %d, %v", got, err)
	}
	for _, raw := range []string{"0", "-1", "bad"} {
		if _, err := positiveInt64Value("SIZE", raw, 1); err == nil {
			t.Errorf("positiveInt64Value(%q) error = nil", raw)
		}
	}
}

func TestConfigRepositoryAllowedAndListenAddress(t *testing.T) {
	configuration := Config{
		Port: "8080",
		AllowedRepositories: map[Repository]struct{}{
			{Host: "github.com", Owner: "CyberT33N", Name: "git-governance"}: {},
		},
	}
	if !configuration.RepositoryAllowed(Repository{Host: "GITHUB.COM", Owner: "CyberT33N", Name: "git-governance"}) {
		t.Fatal("RepositoryAllowed() = false, want true")
	}
	if configuration.RepositoryAllowed(Repository{Host: "github.com", Owner: "CyberT33N", Name: "other"}) {
		t.Fatal("RepositoryAllowed() = true, want false")
	}
	if got := configuration.ListenAddress(); got != ":8080" {
		t.Fatalf("ListenAddress() = %q, want %q", got, ":8080")
	}
}

func cloneEnvironment(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}
