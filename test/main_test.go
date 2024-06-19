package test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}
	enableTest := os.Getenv("ENABLE_TEST")
	if enableTest == "true" {
		t.Run("TestAdd", testAdd)
		t.Run("TestSubtract", testSubtract)
	} else {
		t.Log("Test is disabled")
	}
}
