// Package lastfm provides a set of types and constants for working with the
// Last.fm API.
package lastfm

import (
	"encoding/xml"
	"maps"
	"regexp"
	"slices"
)

const (
	BaseURL = "https://www.last.fm"
	APIURL  = BaseURL + "/api"
	AuthURL = APIURL + "/auth"

	ImageHost    = "https://lastfm.freetls.fastly.net"
	BaseImageURL = ImageHost + "/i/u/"

	// NoArtistHash is the image hash for an artist with no image.
	NoArtistHash = "2a96cbd8b46e442fc41c2b86b821562f"
	// NoAlbumHash is the image hash for an album with no image.
	NoAlbumHash = "c6f59c1e5e7240a4c0d427abd71f3dbb"
	// NoTrackHash is the image hash for a track with no image.
	NoTrackHash = "4128a6eb29f94943c9d206c08e625904"
	// NoAvatarHash is the image hash for a user with no avatar.
	NoAvatarHash = "818148bf682d429dc215c1705eb27b98"

	// NoArtistImageURL is the image URL for an artist with no image.
	NoArtistImageURL ImageURL = BaseImageURL + NoArtistHash + ".png"
	// NoAlbumImageURL is the image URL for an album with no image.
	NoAlbumImageURL ImageURL = BaseImageURL + NoAlbumHash + ".png"
	// NoTrackImageURL is the image URL for a track with no image.
	NoTrackImageURL ImageURL = BaseImageURL + NoTrackHash + ".png"
	// NoAvatarImageURL is the image URL for a user with no avatar.
	NoAvatarImageURL ImageURL = BaseImageURL + NoAvatarHash + ".png"
)

// ImageURLSizeRegex is a regex to match the image URL size.
// It matches the following patterns:
// - i/u/34s/
// - i/u/64s/
// - i/u/174s/
// - i/u/300x300/
// - i/u/
var ImageURLSizeRegex = regexp.MustCompile(`i/u/(34s/|64s/|174s/|300x300/)?`)

// BuildImageURL builds the image URL for the given size and hash.
func BuildImageURL(size ImgSize, hash string) ImageURL {
	return ImageURL(BaseImageURL + size.PathSize() + hash + ".png")
}

type ImgSize string

const (
	// Used when an API response returns an image without a size attribute.
	ImgSizeUndefined ImgSize = "undefined"
	// 34x34
	ImgSizeSmall ImgSize = "small"
	// 64x64
	ImgSizeMedium ImgSize = "medium"
	// 174x174
	ImgSizeLarge ImgSize = "large"
	// 300x300
	ImgSizeExtraLarge ImgSize = "extralarge"
	// 300x300?
	ImgSizeMega ImgSize = "mega"
	// Original upload size
	ImgSizeOriginal ImgSize = "original"
)

// PathSize returns the path size string for the given ImgSize.
func (s ImgSize) PathSize() string {
	switch s {
	case ImgSizeSmall:
		return "34s/"
	case ImgSizeMedium:
		return "64s/"
	case ImgSizeLarge:
		return "174s/"
	case ImgSizeExtraLarge, ImgSizeMega:
		return "300x300/"
	case ImgSizeOriginal:
		return ""
	default:
		return ImgSizeExtraLarge.PathSize()
	}
}

// Compare returns -1 if s < other, 0 if s == other, and 1 if s > other.
func (s ImgSize) Compare(other ImgSize) int {
	return s.intSize() - other.intSize()
}

func (s ImgSize) intSize() int {
	switch s {
	case ImgSizeSmall:
		return 1
	case ImgSizeMedium:
		return 2
	case ImgSizeLarge:
		return 3
	case ImgSizeExtraLarge:
		return 4
	case ImgSizeMega:
		return 5
	default:
		return 0
	}
}

type ImageURL string

// String returns the string representation of the ImageURL.
func (i ImageURL) String() string {
	return string(i)
}

// Resize returns the resized image URL with the specified size.
func (i ImageURL) Resize(size ImgSize) string {
	if i == "" {
		return ""
	}

	return ImageURLSizeRegex.ReplaceAllString(i.String(), "i/u/"+size.PathSize())
}

type Image map[ImgSize]ImageURL

// UnmarshalXML implements the xml.Unmarshaler interface for Image.
func (i *Image) UnmarshalXML(dc *xml.Decoder, start xml.StartElement) error {
	if *i == nil {
		*i = make(Image)
	}

	var size ImgSize
	for _, attr := range start.Attr {
		if attr.Name.Local == "size" {
			size = ImgSize(attr.Value)
			break
		}
	}

	var url ImageURL
	if err := dc.DecodeElement(&url, &start); err != nil {
		return err
	}

	if url != "" {
		if size == "" {
			size = ImgSizeUndefined
		}
		(*i)[size] = url
	}

	return nil
}

// String returns the string representation of the Image URL.
func (i Image) String() string {
	return i.URL()
}

// URL returns the URL of the image in its extra large size.
// This is the same as calling SizedURL(ImgSizeExtralarge).
func (i Image) URL() string {
	return i.SizedURL(ImgSizeExtraLarge)
}

// SizedURL returns the URL of the image with the specified size.
func (i Image) SizedURL(size ImgSize) string {
	if url, ok := i[size]; ok {
		return url.String()
	}

	return i.url().Resize(size)
}

// OriginalURL returns the URL of the image in original size.
func (i Image) OriginalURL() string {
	return i.url().Resize(ImgSizeOriginal)
}

func (i Image) url() ImageURL {
	if len(i) == 0 {
		return ""
	}

	if url, ok := i[ImgSizeExtraLarge]; ok {
		return url
	}

	sizes := slices.SortedFunc(maps.Keys(i), func(a, b ImgSize) int {
		return a.Compare(b)
	})

	largest := sizes[len(sizes)-1]
	return i[largest]
}

type Period string

const (
	PeriodOverall Period = "overall"
	PeriodWeek    Period = "7day"
	PeriodMonth   Period = "1month"
	Period3Months Period = "3month"
	Period6Months Period = "6month"
	PeriodYear    Period = "12month"
)

type TagType string

const (
	TagTypeArtist TagType = "artist"
	TagTypeAlbum  TagType = "album"
	TagTypeTrack  TagType = "track"
)
