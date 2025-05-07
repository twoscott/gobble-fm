package lastfm

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

const TimeFormat = "02 Jan 2006, 15:04"

type IntBool bool

func (b IntBool) Bool() bool {
	return bool(b)
}

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

type DateTime time.Time

// MarshalJSON implements the json.Marshaler interface for DateTime.
func (dt DateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dt))
}

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
	var dateString string
	if err := d.DecodeElement(&dateString, &start); err != nil {
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
		if err != nil {
			return err
		}
		*dt = DateTime(time.Unix(sec, 0))
		return nil
	}

	if dateString != "" {
		t, err := time.ParseInLocation(TimeFormat, dateString, time.UTC)
		if err != nil {
			return err
		}

		*dt = DateTime(t)
	}

	return nil
}

// UnmarshalXMLAttr implements the xml.UnmarshalerAttr interface for DateTime.
func (dt *DateTime) UnmarshalXMLAttr(attr xml.Attr) error {
	sec, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return err
	}
	*dt = DateTime(time.Unix(sec, 0))
	return nil
}

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
