package regru

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

// sign implements section 2.2.4.1 "Аутентификация по сигнатуре": all non-empty,
// non-zero scalar values of the request (form fields and input_data contents,
// except sig) are sorted, joined with ";" and signed with RSA-SHA512.
func sign(key *rsa.PrivateKey, form url.Values, input []byte) (string, error) {
	text, err := sigText(form, input)
	if err != nil {
		return "", err
	}
	sum := sha512.Sum512([]byte(text))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA512, sum[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func sigText(form url.Values, input []byte) (string, error) {
	var vals []string
	for k, vs := range form {
		if k == "sig" || k == "input_data" {
			continue
		}
		for _, v := range vs {
			vals = appendSigValue(vals, v)
		}
	}
	if len(input) > 0 {
		dec := json.NewDecoder(bytes.NewReader(input))
		dec.UseNumber()
		var data any
		if err := dec.Decode(&data); err != nil {
			return "", err
		}
		vals = collectSigValues(vals, data)
	}
	sort.Strings(vals)
	return strings.Join(vals, ";"), nil
}

func collectSigValues(vals []string, d any) []string {
	switch d := d.(type) {
	case map[string]any:
		for k, v := range d {
			if k != "sig" {
				vals = collectSigValues(vals, v)
			}
		}
	case []any:
		for _, v := range d {
			vals = collectSigValues(vals, v)
		}
	case string:
		vals = appendSigValue(vals, d)
	case json.Number:
		vals = appendSigValue(vals, d.String())
	case bool:
		if d {
			vals = append(vals, "1")
		}
	}
	return vals
}

// appendSigValue skips values that are false in Perl: "" and "0".
func appendSigValue(vals []string, v string) []string {
	if v == "" || v == "0" {
		return vals
	}
	return append(vals, v)
}
