// Package lastfm provides a set of types and constants for working with the
// Last.fm API.
package lastfm

import "testing"

func TestImageURL_Resize(t *testing.T) {
	cases := []struct {
		name string
		i    ImageURL
		size ImgSize
		want string
	}{
		{
			name: "Small resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeSmall,
			want: "https://lastfm.freetls.fastly.net/i/u/34s/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Medium resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeMedium,
			want: "https://lastfm.freetls.fastly.net/i/u/64s/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Large resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeLarge,
			want: "https://lastfm.freetls.fastly.net/i/u/174s/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "ExtraLarge resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeExtraLarge,
			want: "https://lastfm.freetls.fastly.net/i/u/300x300/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Mega resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeMega,
			want: "https://lastfm.freetls.fastly.net/i/u/300x300/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Original resize from original",
			i:    "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeOriginal,
			want: "https://lastfm.freetls.fastly.net/i/u/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Resize from existing size",
			i:    "https://lastfm.freetls.fastly.net/i/u/ar0/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeMedium,
			want: "https://lastfm.freetls.fastly.net/i/u/64s/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Resize from invalid size",
			i:    "https://lastfm.freetls.fastly.net/i/u/invalid_size/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeMedium,
			want: "https://lastfm.freetls.fastly.net/i/u/64s/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
		{
			name: "Resize from invalid path",
			i:    "https://lastfm.freetls.fastly.net/58ac84b9c6b2978cad0a0099a38bc530.jpg",
			size: ImgSizeLarge,
			want: "https://lastfm.freetls.fastly.net/58ac84b9c6b2978cad0a0099a38bc530.jpg",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.i.Resize(c.size); got != c.want {
				t.Errorf("ImageURL.Resize() = %v, want %v", got, c.want)
			}
		})
	}
}
