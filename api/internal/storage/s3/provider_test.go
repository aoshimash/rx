package s3

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestValidateObjectKey(t *testing.T) {
	provider := &Provider{}

	tests := []struct {
		name      string
		objectKey string
		want      bool
	}{
		{
			name:      "valid object key",
			objectKey: "videos/user123/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      true,
		},
		{
			name:      "valid object key with mov extension",
			objectKey: "videos/user-abc/550e8400-e29b-41d4-a716-446655440000.mov",
			want:      true,
		},
		{
			name:      "valid object key with underscore in user ID",
			objectKey: "videos/user_123/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      true,
		},
		{
			name:      "invalid - missing videos prefix",
			objectKey: "user123/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      false,
		},
		{
			name:      "invalid - wrong prefix",
			objectKey: "images/user123/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      false,
		},
		{
			name:      "invalid - missing extension",
			objectKey: "videos/user123/550e8400-e29b-41d4-a716-446655440000",
			want:      false,
		},
		{
			name:      "invalid - empty string",
			objectKey: "",
			want:      false,
		},
		{
			name:      "invalid - path traversal attempt",
			objectKey: "videos/../secrets/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      false,
		},
		{
			name:      "invalid - special characters in user ID",
			objectKey: "videos/user@123/550e8400-e29b-41d4-a716-446655440000.mp4",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := provider.ValidateObjectKey(tt.objectKey)
			if got != tt.want {
				t.Errorf("ValidateObjectKey(%q) = %v, want %v", tt.objectKey, got, tt.want)
			}
		})
	}
}

func TestNormalizeUserID_HashesSpecialCharacters(t *testing.T) {
	// User IDs with special characters must be SHA256-hashed for safe object key prefix
	userID := "user@example.com"
	got := normalizeUserID(userID)
	if got == "" {
		t.Fatal("normalizeUserID returned empty for user ID with special chars")
	}
	if got == userID {
		t.Error("expected hashed value for special chars, got same user ID")
	}
	// SHA256 hex is 64 chars
	if len(got) != 64 {
		t.Errorf("expected 64-char hex, got len %d", len(got))
	}
	h := sha256.Sum256([]byte(userID))
	want := hex.EncodeToString(h[:])
	if got != want {
		t.Errorf("normalizeUserID(%q) = %q, want %q (SHA256 hex)", userID, got, want)
	}
}

func TestUploadURLRequest_ContentTypeValidation(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		wantValid   bool
	}{
		{"video/mp4", "video/mp4", true},
		{"video/quicktime", "video/quicktime", true},
		{"video/webm", "video/webm", true},
		{"image/png", "image/png", false},
		{"application/json", "application/json", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isVideoContentType(tt.contentType)
			if got != tt.wantValid {
				t.Errorf("isVideoContentType(%q) = %v, want %v", tt.contentType, got, tt.wantValid)
			}
		})
	}
}
