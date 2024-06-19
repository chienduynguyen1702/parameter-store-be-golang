package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {
	err := godotenv.Load("../param.test")
	if err != nil {
		panic("Error loading .env.test file for test")
	}
	enableTest := os.Getenv("ENABLE_TEST")
	if enableTest != "true" {
		t.Log("Test is disabled")

	} else { // Test is enabled
		t.Run("TestAdd", testAdd)
		t.Run("TestSubtract", testSubtract)
		t.Run("TestMultiple", testMultiple)
		t.Run("TestCreate", testMultiple1)
		t.Run("TestUpdate", testMultiple2)
		t.Run("TestConnecttions", testMultiple3)
		t.Run("TestFunc", testMultiple4)
		t.Run("TestHappyCase", testMultiple5)
	}
}

func multiple(a, b int) int {
	return a * b
}
func testMultiple(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}

func testMultiple1(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}
func testMultiple2(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}
func testMultiple3(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}
func testMultiple4(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}
func testMultiple5(t *testing.T) {
	if multiple(3, 2) != 6 {
		t.Error("Expected 6")
	}
}
