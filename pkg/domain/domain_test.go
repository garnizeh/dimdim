package domain_test

import (
	"testing"

	"github.com/garnizeH/dimdim/pkg/domain"
)

func TestDomainIsDev(t *testing.T) {
	tests := []struct {
		name string
		d    domain.Domain
		want bool
	}{
		{
			name: "empty domain",
			d:    domain.Domain(""),
			want: false,
		},
		{
			name: "dev domain",
			d:    domain.Domain("localhost"),
			want: true,
		},
		{
			name: "prod domain",
			d:    domain.Domain("example.com"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsDev(); got != tt.want {
				t.Errorf("Domain.IsDev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainURL(t *testing.T) {
	tests := []struct {
		name   string
		domain domain.Domain
		path   []string
		want   string
	}{
		{
			name:   "dev domain - invalid path",
			domain: domain.Domain("localhost"),
			path:   nil,
			want:   "http://localhost/",
		},
		{
			name:   "dev domain - empty path",
			domain: domain.Domain("localhost"),
			path:   []string{},
			want:   "http://localhost/",
		},
		{
			name:   "dev domain",
			domain: domain.Domain("localhost"),
			path:   []string{"foo", "bar"},
			want:   "http://localhost/foo/bar",
		},
		{
			name:   "prod domain - invalid path",
			domain: domain.Domain("example.com"),
			path:   nil,
			want:   "https://example.com/",
		},
		{
			name:   "prod domain - empty path",
			domain: domain.Domain("example.com"),
			path:   []string{},
			want:   "https://example.com/",
		},
		{
			name:   "prod domain",
			domain: domain.Domain("example.com"),
			path:   []string{"foo", "bar"},
			want:   "https://example.com/foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.domain.URL(tt.path...); got != tt.want {
				t.Errorf("%q got %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
