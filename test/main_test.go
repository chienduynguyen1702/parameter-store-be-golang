package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		panic("Error loading .env.test file for test")
	}
	enableTest := os.Getenv("ENABLE_TEST")
	if enableTest != "true" {
		t.Log("Test is disabled")

	} else { // Test is enabled
		t.Run("TestAdd", testAdd)
		t.Run("TestSubtract", testSubtract)

	}
}
