package ginmiddleware

import "testing"

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name        string
		requestPath string
		rulePath    string
		want        bool
	}{
		{
			name:        "exact match root",
			requestPath: "/",
			rulePath:    "/",
			want:        true,
		},
		{
			name:        "exact match simple path",
			requestPath: "/api/v1/users",
			rulePath:    "/api/v1/users",
			want:        true,
		},
		{
			name:        "no match different path",
			requestPath: "/api/v1/users",
			rulePath:    "/api/v1/accounts",
			want:        false,
		},
		{
			name:        "simple wildcard prefix match",
			requestPath: "/api/v1/users/123",
			rulePath:    "/api/v1/users/*",
			want:        true,
		},
		{
			name:        "simple wildcard prefix match without extra segments",
			requestPath: "/api/v1/users/",
			rulePath:    "/api/v1/users/*",
			want:        true,
		},
		{
			name:        "simple wildcard prefix non-match",
			requestPath: "/api/v1/admin/123",
			rulePath:    "/api/v1/users/*",
			want:        false,
		},
		{
			name:        "wildcard only star matches anything",
			requestPath: "/anything/here",
			rulePath:    "*",
			want:        true,
		},
		{
			name:        "empty rule only matches empty request",
			requestPath: "",
			rulePath:    "",
			want:        true,
		},
		{
			name:        "empty rule does not match non-empty request",
			requestPath: "/something",
			rulePath:    "",
			want:        false,
		},
		{
			name:        "rule shorter than request without wildcard",
			requestPath: "/api/v1/users/123",
			rulePath:    "/api/v1",
			want:        false,
		},
		{
			name:        "prefix equal but no wildcard",
			requestPath: "/api/v1/usersX",
			rulePath:    "/api/v1/users",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPath(tt.requestPath, tt.rulePath)
			if got != tt.want {
				t.Fatalf("matchPath(%q, %q) = %v, want %v",
					tt.requestPath, tt.rulePath, got, tt.want)
			}
		})
	}
}
