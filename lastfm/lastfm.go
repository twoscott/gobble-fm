package lastfm

import (
	"encoding/xml"
	"maps"
	"regexp"
	"slices"
	"time"
)

const (
	BaseURL = "https://www.last.fm"
	APIURL  = BaseURL + "/api"
	AuthURL = APIURL + "/auth"
)

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

type Duration time.Duration

// Unwrap returns the duration as a time.Duration.
func (d Duration) Unwrap() time.Duration {
	return time.Duration(d)
}

// String returns the duration as a string.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// UnmarshalXML implements the xml.Unmarshaler interface for Duration.
func (d *Duration) UnmarshalXML(dc *xml.Decoder, start xml.StartElement) error {
	var sec int64
	if err := dc.DecodeElement(&sec, &start); err != nil {
		return err
	}

	*d = Duration(time.Duration(sec) * time.Second)
	return nil
}

type ImgSize string

const (
	// 34x34
	ImgSizeSmall ImgSize = "small"
	// 64x64
	ImgSizeMedium ImgSize = "medium"
	// 174x174
	ImgSizeLarge ImgSize = "large"
	// 300x300
	ImgSizeExtralarge ImgSize = "extralarge"
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
	case ImgSizeExtralarge:
		return "300x300/"
	case ImgSizeOriginal:
		return ""
	default:
		return ImgSizeExtralarge.PathSize()
	}
}

func (s ImgSize) intSize() int {
	switch s {
	case ImgSizeSmall:
		return 1
	case ImgSizeMedium:
		return 2
	case ImgSizeLarge:
		return 3
	case ImgSizeExtralarge:
		return 4
	default:
		return 0
	}
}

// Compare returns -1 if s < other, 0 if s == other, and 1 if s > other.
func (s ImgSize) Compare(other ImgSize) int {
	return s.intSize() - other.intSize()
}

type Image map[ImgSize]string

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

	var url string
	if err := dc.DecodeElement(&url, &start); err != nil {
		return err
	}

	(*i)[size] = url
	return nil
}

// String returns the string representation of the Image URL.
func (i Image) String() string {
	return i.URL()
}

// URL returns the URL of the image in its extra large size.
// This is the same as calling SizedURL(ImgSizeExtralarge).
func (i Image) URL() string {
	if len(i) == 0 {
		return ""
	}

	if url, ok := i[ImgSizeExtralarge]; ok {
		return url
	}

	sizes := slices.SortedFunc(maps.Keys(i), func(a, b ImgSize) int { return a.Compare(b) })
	return i[sizes[len(sizes)-1]]
}

// SizedURL returns the URL of the image with the specified size.
func (i Image) SizedURL(size ImgSize) string {
	if url, ok := i[size]; ok {
		return url
	}

	return i.resizeURL(size)
}

// OriginalURL returns the URL of the image in original size.
func (i Image) OriginalURL() string {
	return i.resizeURL(ImgSizeOriginal)
}

func (i Image) resizeURL(size ImgSize) string {
	sizeRegex, err := regexp.Compile(`i/u/\d+(s|x\d+)/`)
	if err != nil {
		return i.URL()
	}

	return sizeRegex.ReplaceAllString(i.URL(), "i/u/"+size.PathSize())
}
