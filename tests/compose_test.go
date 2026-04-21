package pixel_test

import (
	"math"
	"testing"

	"github.com/gopxl/pixel/v2"
)

func approxEq(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func rgbaApproxEq(a, b pixel.RGBA, eps float64) bool {
	return approxEq(a.R, b.R, eps) &&
		approxEq(a.G, b.G, eps) &&
		approxEq(a.B, b.B, eps) &&
		approxEq(a.A, b.A, eps)
}

func TestComposeMultiply(t *testing.T) {
	const eps = 1e-9
	tests := []struct {
		name string
		a, b pixel.RGBA
		want pixel.RGBA
	}{
		{
			name: "both opaque - multiply channels",
			a:    pixel.RGBA{R: 0.5, G: 0.5, B: 0.5, A: 1},
			b:    pixel.RGBA{R: 0.8, G: 0.4, B: 1.0, A: 1},
			want: pixel.RGBA{R: 0.4, G: 0.2, B: 0.5, A: 1},
		},
		{
			name: "white source over opaque bg - unchanged",
			a:    pixel.RGBA{R: 1, G: 1, B: 1, A: 1},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
		},
		{
			name: "black source - all black",
			a:    pixel.RGBA{R: 0, G: 0, B: 0, A: 1},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 0, G: 0, B: 0, A: 1},
		},
		{
			name: "transparent source - backdrop unchanged",
			a:    pixel.RGBA{R: 0, G: 0, B: 0, A: 0},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pixel.ComposeMultiply.Compose(tt.a, tt.b)
			if !rgbaApproxEq(got, tt.want, eps) {
				t.Errorf("ComposeMultiply.Compose(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestComposeScreen(t *testing.T) {
	const eps = 1e-9
	tests := []struct {
		name string
		a, b pixel.RGBA
		want pixel.RGBA
	}{
		{
			name: "both opaque - screen formula: 1-(1-a)*(1-b)",
			a:    pixel.RGBA{R: 0.5, G: 0.5, B: 0.5, A: 1},
			b:    pixel.RGBA{R: 0.5, G: 0.5, B: 0.5, A: 1},
			// 1 - 0.5*0.5 = 0.75
			want: pixel.RGBA{R: 0.75, G: 0.75, B: 0.75, A: 1},
		},
		{
			name: "black source - backdrop unchanged",
			a:    pixel.RGBA{R: 0, G: 0, B: 0, A: 1},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
		},
		{
			name: "white source - output is white",
			a:    pixel.RGBA{R: 1, G: 1, B: 1, A: 1},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 1, G: 1, B: 1, A: 1},
		},
		{
			name: "transparent source - backdrop unchanged",
			a:    pixel.RGBA{R: 0, G: 0, B: 0, A: 0},
			b:    pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
			want: pixel.RGBA{R: 0.6, G: 0.3, B: 0.9, A: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pixel.ComposeScreen.Compose(tt.a, tt.b)
			if !rgbaApproxEq(got, tt.want, eps) {
				t.Errorf("ComposeScreen.Compose(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
