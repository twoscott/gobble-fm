package api

import (
	"encoding/xml"
	"errors"
)

type LFMWrapper struct {
	XMLName  xml.Name `xml:"lfm"`
	Status   string   `xml:"status,attr"`
	InnerXML []byte   `xml:",innerxml"`
}

// Empty checks if the InnerXML is empty.
func (lf *LFMWrapper) Empty() bool {
	return len(lf.InnerXML) == 0
}

// StatusOK checks if the status is "ok".
func (lf *LFMWrapper) StatusOK() bool {
	return lf.Status == "ok"
}

// StatusFailed checks if the status is "failed".
func (lf *LFMWrapper) StatusFailed() bool {
	return lf.Status == "failed"
}

// UnwrapError attempts to extract and return a LastFMError from the LFMWrapper
// instance. If the status of the LFMWrapper indicates success, it returns nil
// for both the error and LastFMError. Otherwise, it unmarshals the InnerXML
// field into a LastFMError structure. If the unmarshaling succeeds and the
// LastFMError contains an error code, it returns the LastFMError. If no error
// code is present or unmarshaling fails, it returns an appropriate error.
func (lf *LFMWrapper) UnwrapError() (*LastFMError, error) {
	if lf.StatusOK() {
		return nil, nil
	}

	var lferr LastFMError
	if err := lf.UnmarshalInnerXML(&lferr); err != nil {
		return nil, err
	}
	if lferr.HasErrorCode() {
		return &lferr, nil
	}

	return nil, errors.New("no error code in response")
}

// UnmarshalInnerXML unmarshals the InnerXML field into the provided destination
// variable.
func (lf *LFMWrapper) UnmarshalInnerXML(dest any) error {
	return xml.Unmarshal(lf.InnerXML, dest)
}
