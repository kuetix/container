package test

import (
	"testing"

	"github.com/kuetix/container"
)

// helper to reset all global containers between tests
func resetAll() {
	container.Reset()
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	deferred := false
	func() {
		deferred = true
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic, got none")
			}
		}()
		fn()
	}()
	if !deferred {
		t.Fatalf("panic checker did not run")
	}
}

func TestToFetchAndFetch(t *testing.T) {
	resetAll()
	container.ToFetch("answer", 42)
	if got := container.Fetch("answer"); got != 42 {
		t.Fatalf("Fetch() = %v, want 42", got)
	}

	// Fetch missing should panic
	mustPanic(t, func() { _ = container.Fetch("missing") })
}

func TestToFetchFunc(t *testing.T) {
	resetAll()
	container.ToFetchFunc("now", func() interface{} { return "ok" })
	if got := container.Fetch("now"); got != "ok" {
		t.Fatalf("Fetch(ToFetchFunc) = %v, want ok", got)
	}
}

func TestToResolveAndResolve(t *testing.T) {
	resetAll()
	count := 0
	container.ToResolve("gen", func() interface{} {
		count++
		return count
	})

	// Resolve should call factory each time
	if got := container.Resolve("gen"); got != 1 {
		t.Fatalf("Resolve() first = %v, want 1", got)
	}
	if got := container.Resolve("gen"); got != 2 {
		t.Fatalf("Resolve() second = %v, want 2", got)
	}

	// Resolve missing should panic
	mustPanic(t, func() { _ = container.Resolve("missing") })
}

func TestGetPrefersSingletonThenFactory(t *testing.T) {
	resetAll()
	container.ToFetch("name", "singleton")
	container.ToResolve("name", func() interface{} { return "factory" })

	if got := container.Get("name"); got != "singleton" {
		t.Fatalf("Get() with both present = %v, want singleton", got)
	}
}

func TestGetFromFactoryAndMissing(t *testing.T) {
	resetAll()
	container.ToResolve("x", func() interface{} { return 7 })
	if got := container.Get("x"); got != 7 {
		t.Fatalf("Get() factory = %v, want 7", got)
	}

	mustPanic(t, func() { _ = container.Get("missing") })
}

func TestParameters(t *testing.T) {
	resetAll()
	container.ToParameter("env", "prod")
	if got := container.Parameter("env"); got != "prod" {
		t.Fatalf("Parameter() = %v, want prod", got)
	}
	mustPanic(t, func() { _ = container.Parameter("missing") })
}

func TestHasAndCanChecks(t *testing.T) {
	resetAll()
	// none present
	if container.Has("a") {
		t.Fatalf("Has() = true, want false")
	}
	if container.CanFetch("a") {
		t.Fatalf("CanFetch() = true, want false")
	}
	if container.CanResolve("a") {
		t.Fatalf("CanResolve() = true, want false")
	}
	if container.HasParameter("a") {
		t.Fatalf("HasParameter() = true, want false")
	}

	container.ToFetch("s", 1)
	container.ToResolve("f", func() interface{} { return 2 })
	container.ToParameter("p", 3)

	if !container.Has("s") || !container.Has("f") || !container.Has("p") {
		t.Fatalf("Has() did not detect present keys")
	}
	if !container.CanFetch("s") {
		t.Fatalf("CanFetch() = false, want true")
	}
	if !container.CanResolve("f") {
		t.Fatalf("CanResolve() = false, want true")
	}
	if !container.HasParameter("p") {
		t.Fatalf("HasParameter() = false, want true")
	}
}
