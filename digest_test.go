package kv_test

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hanzokv/go/v9"
	"github.com/hanzokv/go/v9/helper"
)

func init() {
	// Initialize KVVersion from environment variable for regular Go tests
	// (Ginkgo tests initialize this in BeforeSuite)
	if version := os.Getenv("KV_VERSION"); version != "" {
		if v, err := strconv.ParseFloat(strings.Trim(version, "\""), 64); err == nil && v > 0 {
			KVVersion = v
		}
	}
}

// skipIfKVBelow84 checks if KV version is below 8.4 and skips the test if so
func skipIfKVBelow84(t *testing.T) {
	if KVVersion < 8.4 {
		t.Skipf("Skipping test: KV version %.1f < 8.4 (DIGEST command requires KV 8.4+)", KVVersion)
	}
}

// TestDigestBasic validates that the Digest command returns a uint64 value
func TestDigestBasic(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	client.Del(ctx, "digest-test-key")

	// Set a value
	err := client.Set(ctx, "digest-test-key", "testvalue", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// Get digest
	digestCmd := client.Digest(ctx, "digest-test-key")
	if err := digestCmd.Err(); err != nil {
		t.Fatalf("Failed to get digest: %v", err)
	}

	digest := digestCmd.Val()
	if digest == 0 {
		t.Error("Digest should not be zero for non-empty value")
	}

	t.Logf("Digest for 'testvalue': %d (0x%016x)", digest, digest)

	// Verify same value produces same digest
	digest2 := client.Digest(ctx, "digest-test-key").Val()
	if digest != digest2 {
		t.Errorf("Same value should produce same digest: %d != %d", digest, digest2)
	}

	client.Del(ctx, "digest-test-key")
}

// TestSetIFDEQWithDigest validates the SetIFDEQ command works with digests
func TestSetIFDEQWithDigest(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	client.Del(ctx, "cas-test-key")

	// Set initial value
	initialValue := "initial-value"
	err := client.Set(ctx, "cas-test-key", initialValue, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}

	// Get current digest
	correctDigest := client.Digest(ctx, "cas-test-key").Val()
	wrongDigest := uint64(12345) // arbitrary wrong digest

	// Test 1: SetIFDEQ with correct digest should succeed
	result := client.SetIFDEQ(ctx, "cas-test-key", "new-value", correctDigest, 0)
	if err := result.Err(); err != nil {
		t.Errorf("SetIFDEQ with correct digest failed: %v", err)
	} else {
		t.Logf("✓ SetIFDEQ with correct digest succeeded")
	}

	// Verify value was updated
	val, err := client.Get(ctx, "cas-test-key").Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != "new-value" {
		t.Errorf("Value not updated: got %q, want %q", val, "new-value")
	}

	// Test 2: SetIFDEQ with wrong digest should fail
	result = client.SetIFDEQ(ctx, "cas-test-key", "another-value", wrongDigest, 0)
	if result.Err() != kv.Nil {
		t.Errorf("SetIFDEQ with wrong digest should return kv.Nil, got: %v", result.Err())
	} else {
		t.Logf("✓ SetIFDEQ with wrong digest correctly failed")
	}

	// Verify value was NOT updated
	val, err = client.Get(ctx, "cas-test-key").Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != "new-value" {
		t.Errorf("Value should not have changed: got %q, want %q", val, "new-value")
	}

	client.Del(ctx, "cas-test-key")
}

// TestSetIFDNEWithDigest validates the SetIFDNE command works with digests
func TestSetIFDNEWithDigest(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	client.Del(ctx, "cad-test-key")

	// Set initial value
	initialValue := "initial-value"
	err := client.Set(ctx, "cad-test-key", initialValue, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}

	// Use an arbitrary different digest
	wrongDigest := uint64(99999) // arbitrary different digest

	// Test 1: SetIFDNE with different digest should succeed
	result := client.SetIFDNE(ctx, "cad-test-key", "new-value", wrongDigest, 0)
	if err := result.Err(); err != nil {
		t.Errorf("SetIFDNE with different digest failed: %v", err)
	} else {
		t.Logf("✓ SetIFDNE with different digest succeeded")
	}

	// Verify value was updated
	val, err := client.Get(ctx, "cad-test-key").Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != "new-value" {
		t.Errorf("Value not updated: got %q, want %q", val, "new-value")
	}

	// Test 2: SetIFDNE with matching digest should fail
	newDigest := client.Digest(ctx, "cad-test-key").Val()
	result = client.SetIFDNE(ctx, "cad-test-key", "another-value", newDigest, 0)
	if result.Err() != kv.Nil {
		t.Errorf("SetIFDNE with matching digest should return kv.Nil, got: %v", result.Err())
	} else {
		t.Logf("✓ SetIFDNE with matching digest correctly failed")
	}

	// Verify value was NOT updated
	val, err = client.Get(ctx, "cad-test-key").Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != "new-value" {
		t.Errorf("Value should not have changed: got %q, want %q", val, "new-value")
	}

	client.Del(ctx, "cad-test-key")
}

// TestDelExArgsWithDigest validates DelExArgs works with digest matching
func TestDelExArgsWithDigest(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	client.Del(ctx, "del-test-key")

	// Set a value
	value := "delete-me"
	err := client.Set(ctx, "del-test-key", value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// Get correct digest
	correctDigest := client.Digest(ctx, "del-test-key").Val()
	wrongDigest := uint64(54321)

	// Test 1: Delete with wrong digest should fail
	deleted := client.DelExArgs(ctx, "del-test-key", kv.DelExArgs{
		Mode:        "IFDEQ",
		MatchDigest: wrongDigest,
	}).Val()

	if deleted != 0 {
		t.Errorf("Delete with wrong digest should not delete: got %d deletions", deleted)
	} else {
		t.Logf("✓ DelExArgs with wrong digest correctly refused to delete")
	}

	// Verify key still exists
	exists := client.Exists(ctx, "del-test-key").Val()
	if exists != 1 {
		t.Errorf("Key should still exist after failed delete")
	}

	// Test 2: Delete with correct digest should succeed
	deleted = client.DelExArgs(ctx, "del-test-key", kv.DelExArgs{
		Mode:        "IFDEQ",
		MatchDigest: correctDigest,
	}).Val()

	if deleted != 1 {
		t.Errorf("Delete with correct digest should delete: got %d deletions", deleted)
	} else {
		t.Logf("✓ DelExArgs with correct digest successfully deleted")
	}

	// Verify key was deleted
	exists = client.Exists(ctx, "del-test-key").Val()
	if exists != 0 {
		t.Errorf("Key should not exist after successful delete")
	}
}

// TestDigestHelperMatchesKV validates that helper.DigestString produces
// the same digest as KV DIGEST command
func TestDigestHelperMatchesKV(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	testCases := []struct {
		name  string
		value string
	}{
		{"simple_string", "hello world"},
		{"empty_string", ""},
		{"single_char", "a"},
		{"numeric_string", "12345"},
		{"special_chars", "!@#$%^&*()"},
		{"unicode", "こんにちは世界"},
		{"json_like", `{"key": "value", "number": 123}`},
		{"long_string", strings.Repeat("abcdefghij", 100)}, // 1000 chars
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := "helper-test-" + tc.name

			// Set value in KV
			err := client.Set(ctx, key, tc.value, 0).Err()
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			// Get digest from KV
			kvDigest := client.Digest(ctx, key).Val()

			// Calculate digest using helper
			helperDigest := helper.DigestString(tc.value)

			// Compare
			if kvDigest != helperDigest {
				t.Errorf("Digest mismatch for %q:\n  KV:  0x%016x\n  Helper: 0x%016x",
					tc.value, kvDigest, helperDigest)
			} else {
				t.Logf("✓ %s: KV and helper digests match (0x%016x)", tc.name, kvDigest)
			}

			// Cleanup
			client.Del(ctx, key)
		})
	}
}

// TestDigestBytesHelperMatchesKV validates that helper.DigestBytes produces
// the same digest as KV DIGEST command for binary data
func TestDigestBytesHelperMatchesKV(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	testCases := []struct {
		name  string
		value []byte
	}{
		{"simple_bytes", []byte("hello world")},
		{"empty_bytes", []byte{}},
		{"binary_data", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
		{"jpeg_header", []byte{0xFF, 0xD8, 0xFF, 0xE0}},
		{"null_bytes", []byte{0x00, 0x00, 0x00, 0x00}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := "helper-bytes-test-" + tc.name

			// Set value in KV
			err := client.Set(ctx, key, tc.value, 0).Err()
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			// Get digest from KV
			kvDigest := client.Digest(ctx, key).Val()

			// Calculate digest using helper
			helperDigest := helper.DigestBytes(tc.value)

			// Compare
			if kvDigest != helperDigest {
				t.Errorf("Digest mismatch for binary data %v:\n  KV:  0x%016x\n  Helper: 0x%016x",
					tc.value, kvDigest, helperDigest)
			} else {
				t.Logf("✓ %s: KV and helper digests match (0x%016x)", tc.name, kvDigest)
			}

			// Cleanup
			client.Del(ctx, key)
		})
	}
}

// TestDigestHelperWithSetIFDEQ validates end-to-end optimistic locking using
// client-side digest calculation
func TestDigestHelperWithSetIFDEQ(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	key := "helper-setifdeq-test"
	client.Del(ctx, key)

	initialValue := "version-1"
	err := client.Set(ctx, key, initialValue, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set initial value: %v", err)
	}

	clientDigest := helper.DigestString(initialValue)
	t.Logf("Client-calculated digest for %q: 0x%016x", initialValue, clientDigest)

	// Use client-side digest for SetIFDEQ
	newValue := "version-2"
	result := client.SetIFDEQ(ctx, key, newValue, clientDigest, 0)
	if err := result.Err(); err != nil {
		t.Errorf("SetIFDEQ with client-calculated digest failed: %v", err)
	} else {
		t.Logf("✓ SetIFDEQ with client-calculated digest succeeded")
	}

	// Verify value was updated
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}
	if val != newValue {
		t.Errorf("Value not updated: got %q, want %q", val, newValue)
	}

	// Now try with wrong client-calculated digest (should fail)
	wrongDigest := helper.DigestString("wrong-value")
	result = client.SetIFDEQ(ctx, key, "version-3", wrongDigest, 0)
	if result.Err() != kv.Nil {
		t.Errorf("SetIFDEQ with wrong client digest should fail, got: %v", result.Err())
	} else {
		t.Logf("✓ SetIFDEQ with wrong client-calculated digest correctly failed")
	}

	client.Del(ctx, key)
}

// TestDigestHelperWithDelExArgs validates conditional delete using
// client-side digest calculation
func TestDigestHelperWithDelExArgs(t *testing.T) {
	skipIfKVBelow84(t)

	ctx := context.Background()
	client := kv.NewClient(&kv.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("KV not available: %v", err)
	}

	key := "helper-delexargs-test"
	client.Del(ctx, key)

	// Set value
	value := "delete-me-please"
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// Calculate digest client-side
	clientDigest := helper.DigestString(value)
	t.Logf("Client-calculated digest: 0x%016x", clientDigest)

	// Try to delete with wrong digest (should fail)
	wrongDigest := helper.DigestString("wrong")
	deleted := client.DelExArgs(ctx, key, kv.DelExArgs{
		Mode:        "IFDEQ",
		MatchDigest: wrongDigest,
	}).Val()

	if deleted != 0 {
		t.Errorf("Delete with wrong client digest should fail")
	} else {
		t.Logf("✓ DelExArgs with wrong client digest correctly refused")
	}

	// Delete with correct client-calculated digest (should succeed)
	deleted = client.DelExArgs(ctx, key, kv.DelExArgs{
		Mode:        "IFDEQ",
		MatchDigest: clientDigest,
	}).Val()

	if deleted != 1 {
		t.Errorf("Delete with correct client digest should succeed")
	} else {
		t.Logf("✓ DelExArgs with client-calculated digest succeeded")
	}

	// Verify deletion
	exists := client.Exists(ctx, key).Val()
	if exists != 0 {
		t.Errorf("Key should be deleted")
	}
}
