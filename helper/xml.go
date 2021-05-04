package helper

import (
	"encoding/xml"
	"time"
)

// CustomTime is type data to help xml date time field to convert to time
type CustomTime struct {
	time.Time
}

// UnmarshalXML like override
func (c *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "2006-01-02 15:04:05" // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = CustomTime{parse}
	return nil
}
