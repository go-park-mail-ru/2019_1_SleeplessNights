package console

import "testing"

func TestMessage(t *testing.T) {
	Message("Hello world")
}

func TestSuccess(t *testing.T) {
	Success("Hello world")
}

func TestError(t *testing.T) {
	Error("Hello world")
}

func TestPredicate(t *testing.T) {
	Predicate(true, "Hello world")
}

func TestTitle(t *testing.T) {
	Title("Hello world")
}
