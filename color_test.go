package pixel

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorConversions(t *testing.T) {
	c0 := RGBA{64, 127, 128, 255}
	c1 := FloatsToColor(ColorToFloats[float32](c0))
	require.Equal(t, c0, c1)
}

func TestToRGBA(t *testing.T) {
	type args struct {
		c color.Color
	}
	tests := []struct {
		name string
		args args
		want RGBA
	}{
		{
			name: "pixel.rgba",
			args: args{c: RGBA{64, 127, 128, 255}},
			want: RGBA{64, 127, 128, 255},
		},
		{
			name: "color.rgba",
			args: args{c: color.RGBA{64, 127, 128, 255}},
			want: RGBA{64, 127, 128, 255},
		},
		{
			name: "color.nrgba",
			args: args{c: color.NRGBA{64, 127, 128, 255}},
			want: RGBA{64, 127, 128, 255},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToRGBA(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRGBA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGB(t *testing.T) {
	rgba := RGB(0.25, 0.5, 0.75)
	require.Equal(t, RGBA{63, 127, 191, 255}, rgba)
}
