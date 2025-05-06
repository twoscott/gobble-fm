package lastfm

import (
	"encoding/xml"
	"fmt"
)

type IntBool bool

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
