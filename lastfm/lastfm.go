package lastfm

import (
	"encoding/xml"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	BaseURL = "https://www.last.fm"
	APIURL  = BaseURL + "/api"
	AuthURL = APIURL + "/auth"
)

type Params struct {
	shouldSign bool
	values     url.Values
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

type Image string

// String returns the string representation of the Image URL.
func (i Image) String() string {
	return string(i)
}

// URL returns the URL of the image in its extra large size.
// This is the same as calling SizedURL(ImgSizeExtralarge).
func (i Image) URL() string {
	return i.SizedURL(ImgSizeExtralarge)
}

// SizedURL returns the URL of the image with the specified size.
func (i Image) SizedURL(size ImgSize) string {
	sizeRegex := regexp.MustCompile(`i/u/\d+(s|x\d+)/`)
	return sizeRegex.ReplaceAllString(i.String(), "i/u/"+size.PathSize())
}

type DateTime time.Time

// Time returns the time.Time representation of the DateTime.
func (dt DateTime) Time() time.Time {
	return time.Time(dt)
}

// Sring returns the string representation of the DateTime in DateTime format.
func (dt DateTime) String() string {
	return dt.Format(time.DateTime)
}

// Format returns the string representation of the DateTime in the
// specified format.
func (dt DateTime) Format(format string) string {
	return dt.Time().Format(format)
}

// MarshalJSON implements the json.Marshaler interface for DateTime.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	return time.Time(dt).MarshalJSON()
}

// UnmarshalXML implements the xml.Unmarshaler interface for DateTime.
func (dt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var uts string
	for _, attr := range start.Attr {
		if attr.Name.Local == "uts" || attr.Name.Local == "unixtime" {
			uts = attr.Value
			break
		}
	}

	if uts != "" {
		sec, err := strconv.ParseInt(uts, 10, 64)
		if err != nil {
			return err
		}
		*dt = DateTime(time.Unix(sec, 0))
	}

	var discard string
	if err := d.DecodeElement(&discard, &start); err != nil {
		return err
	}

	return nil
}
