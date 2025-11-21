package test

import (
	"testing"

	"github.com/kuetix/container"
)

func resetBootstrap() {
	container.ResetBootstrap()
}

func TestBootInitializesOnce(t *testing.T) {
	resetBootstrap()
	if container.Booted() {
		t.Fatalf("initOnce must start false")
	}
	container.Boot()
	if !container.Booted() {
		t.Fatalf("Boot() did not set initOnce true")
	}
	if container.DependencyInjection == nil {
		t.Fatalf("Boot() did not initialize DependencyInjection map")
	}

	// Add a sentinel, call Boot again; map should remain intact (idempotent)
	container.DependencyInjection["sentinel"] = func(string) {}
	container.Boot()
	if _, ok := container.DependencyInjection["sentinel"]; !ok {
		t.Fatalf("Boot() was not idempotent; DependencyInjection contents changed")
	}
}

func TestDependencyInjectionBootInvokesRegistered(t *testing.T) {
	resetBootstrap()
	container.Boot()
	called := make(map[string]int)

	container.DependencyInjection["one"] = func(name string) { called[name]++ }
	container.DependencyInjection["two"] = func(name string) { called[name]++ }

	container.DependencyInjectionBoot()

	if called["one"] != 1 || called["two"] != 1 {
		t.Fatalf("DependencyInjectionBoot() did not invoke all entries: got %#v", called)
	}
}
