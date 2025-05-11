package lastfm

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// The string format Last.fm uses to represent dates and times.
const TimeFormat = "02 Jan 2006, 15:04"

// IntBool wraps a boolean and represents a Last.fm integer boolean.
type IntBool bool

func (b IntBool) Bool() bool {
	return bool(b)
}

// UnmarshalXML implements the xml.Unmarshaler interface for IntBool. Unmarshals
// an integer value into a boolean. 1 is true, 0 is false.
func (b *IntBool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var val int
	if err := d.DecodeElement(&val, &start); err != nil {
		return err
	}

	switch val {
	case 1:
		*b = true
	case 0:
		*b = false
	default:
		return fmt.Errorf("invalid IntBool value: %d", val)
	}

	return nil
}

// DateTime wraps time.Time and represents a Last.fm DateTime.
type DateTime time.Time

// Unix returns the Unix timestamp of the DateTime.
func (dt DateTime) Unix() int64 {
	return dt.Time().Unix()
}

// Time returns the time.Time representation of the DateTime.
func (dt DateTime) Time() time.Time {
	return time.Time(dt)
}

// String returns the string representation of the DateTime in DateTime format.
func (dt DateTime) String() string {
	return dt.Format(time.DateTime)
}

// Format returns the string representation of the DateTime in the
// specified format.
func (dt DateTime) Format(format string) string {
	return dt.Time().Format(format)
}

// UnmarshalXML implements the xml.Unmarshaler interface for DateTime.
func (dt *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content string
	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	var uts string
	for _, attr := range start.Attr {
		if attr.Name.Local == "uts" || attr.Name.Local == "unixtime" {
			uts = attr.Value
			break
		}
	}

	if uts != "" {
		sec, err := strconv.ParseInt(uts, 10, 64)
		if err == nil {
			*dt = DateTime(time.Unix(sec, 0))
			return nil
		}
	}

	if content != "" {
		sec, err := strconv.ParseInt(content, 10, 64)
		if err == nil {
			*dt = DateTime(time.Unix(sec, 0))
			return nil
		}

		t, err := time.ParseInLocation(TimeFormat, content, time.UTC)
		if err == nil {
			*dt = DateTime(t)
			return nil
		}
	}

	return nil
}

// UnmarshalXMLAttr implements the xml.UnmarshalerAttr interface for DateTime.
func (dt *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {
	sec, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return nil
	}
	*dt = DateTime(time.Unix(sec, 0))
	return nil
}

// Duration wraps a time.Duration in seconds.
type Duration time.Duration

const (
	DurationHour   = Duration(time.Hour)
	DurationMinute = Duration(time.Minute)
	DurationSecond = Duration(time.Second)
)

// DurationMinSec returns a Duration from minutes and seconds.
func DurationMinSec(minutes, sec int) Duration {
	return (Duration(minutes) * DurationMinute) + (Duration(sec) * DurationSecond)
}

// DurationSeconds returns a Duration from seconds.
func DurationSeconds(seconds int) Duration {
	return Duration(seconds) * DurationSecond
}

// EncodeValues implements the url.ValuesEncoder interface for Duration.
func (d Duration) EncodeValues(key string, v *url.Values) error {
	sec := strconv.FormatFloat(time.Duration(d).Seconds(), 'f', 0, 64)
	v.Set(key, sec)
	return nil
}

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
	var s string
	if err := dc.DecodeElement(&s, &start); err != nil {
		return err
	}

	sec, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// sometimes field isn't a number (e.g., "userdata: NULL")
		return nil
	}

	*d = Duration(time.Duration(sec) * time.Second)
	return nil
}

// DurationMilli wraps a time.Duration in milliseconds.
type DurationMilli time.Duration

// Unwrap returns the duration as a time.Duration.
func (d DurationMilli) Unwrap() time.Duration {
	return Duration(d).Unwrap()
}

// String returns the duration as a string.
func (d DurationMilli) String() string {
	return Duration(d).String()
}

// UnmarshalXML implements the xml.Unmarshaler interface for Duration.
func (d *DurationMilli) UnmarshalXML(dc *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := dc.DecodeElement(&s, &start); err != nil {
		return err
	}

	mil, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// sometimes field isn't a number (e.g., "userdata: NULL")
		return nil
	}

	*d = DurationMilli(time.Duration(mil) * time.Millisecond)
	return nil
}
