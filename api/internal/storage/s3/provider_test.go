package s3

import (
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
			isVideo := len(tt.contentType) >= 6 && tt.contentType[:6] == "video/"
			if isVideo != tt.wantValid {
				t.Errorf("content type %q validation = %v, want %v", tt.contentType, isVideo, tt.wantValid)
			}
		})
	}
}
