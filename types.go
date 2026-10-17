package regru

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

// FlexInt is an integer the API may encode either as a JSON number or a string.
// A non-numeric string decodes to 0: error items in list answers echo invalid IDs
// back (e.g. "service_id":"22bug22") and must not fail the whole answer.
type FlexInt int64

func (n *FlexInt) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	if b[0] == '{' || b[0] == '[' {
		return &json.UnmarshalTypeError{Value: string(b), Type: reflect.TypeOf(*n)}
	}
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		f, _ := strconv.ParseFloat(string(b), 64)
		v = int64(f)
	}
	*n = FlexInt(v)
	return nil
}

// FlexFloat is a number the API may encode either as a JSON number or a string (prices, balances).
type FlexFloat float64

func (n *FlexFloat) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	v, err := strconv.ParseFloat(string(b), 64)
	*n = FlexFloat(v)
	return err
}

// FlexString is a string the API may encode as a JSON number.
type FlexString string

func (s *FlexString) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var v string
		err := json.Unmarshal(b, &v)
		*s = FlexString(v)
		return err
	}
	if string(b) == "null" {
		*s = ""
		return nil
	}
	*s = FlexString(b)
	return nil
}

// Bool is a flag sent to the API as 0/1.
type Bool bool

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	s := string(bytes.Trim(data, `"`))
	*b = Bool(s == "1" || s == "true")
	return nil
}

// ServiceIdent identifies a domain or service (section 2.2.4.3).
// Set ServiceID (recommended), or DName/DomainName, optionally with ServType/SubType,
// or UplinkServiceID with ServType/SubType.
type ServiceIdent struct {
	ServiceID       int64  `json:"service_id,omitempty"`
	UserServID      string `json:"user_servid,omitempty"`
	DName           string `json:"dname,omitempty"`
	DomainName      string `json:"domain_name,omitempty"`
	ServType        string `json:"servtype,omitempty"`
	SubType         string `json:"subtype,omitempty"`
	UplinkServiceID int64  `json:"uplink_service_id,omitempty"`
}

// FolderIdent identifies a folder (section 2.2.4.5). FolderID is recommended.
type FolderIdent struct {
	FolderID   int64  `json:"folder_id,omitempty"`
	FolderName string `json:"folder_name,omitempty"`
}

// PaymentOptions are common to functions that order or renew services (section 2.2.4.6).
type PaymentOptions struct {
	PayType     string `json:"pay_type,omitempty"` // bank, pbank, prepay (default), yacard
	OkIfNoMoney Bool   `json:"ok_if_no_money,omitempty"`
}

// ServiceItem is a per-service entry in list answers.
type ServiceItem struct {
	ItemResult
	DName      string     `json:"dname,omitempty"`
	ServiceID  FlexInt    `json:"service_id,omitempty"`
	ServType   string     `json:"servtype,omitempty"`
	SubType    string     `json:"subtype,omitempty"`
	State      string     `json:"state,omitempty"`
	UserServID FlexString `json:"user_servid,omitempty"`
}
