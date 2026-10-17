package regru

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// setup starts a test server; handlers are registered on mux with paths like "/domain/check".
func setup(t *testing.T) (*Client, *http.ServeMux) {
	t.Helper()
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	c := NewClient(nil).WithAuthPassword("test", "test")
	c.BaseURL, _ = url.Parse(srv.URL + "/")
	return c, mux
}

// handle registers path, asserts it is a POST and that input_data equals wantInput
// (compared as JSON; nil skips the check), and replies with answer wrapped in a success envelope.
func handle(t *testing.T, mux *http.ServeMux, path string, wantInput any, answer string) {
	t.Helper()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("%s: method = %s, want POST", path, r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if wantInput != nil {
			var got, want any
			if err := json.Unmarshal([]byte(r.PostForm.Get("input_data")), &got); err != nil {
				t.Errorf("%s: input_data %q: %v", path, r.PostForm.Get("input_data"), err)
			}
			b, _ := json.Marshal(wantInput)
			_ = json.Unmarshal(b, &want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: input_data = %s, want %s", path, r.PostForm.Get("input_data"), b)
			}
		}
		fmt.Fprintf(w, `{"result":"success","answer":%s}`, answer)
	})
}

func TestNop(t *testing.T) {
	c, mux := setup(t)
	mux.HandleFunc("/nop", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		for k, want := range map[string]string{"username": "test", "password": "test", "output_format": "json"} {
			if got := r.PostForm.Get(k); got != want {
				t.Errorf("%s = %q, want %q", k, got, want)
			}
		}
		if r.PostForm.Has("input_data") {
			t.Error("input_data sent for nil params")
		}
		fmt.Fprint(w, `{"answer":{"login":"test","user_id":"0"},"result":"success"}`)
	})
	got, _, err := c.Nop(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got.Login != "test" || got.UserID != 0 {
		t.Errorf("got %+v", got)
	}
}

func TestGetServiceID(t *testing.T) {
	c, mux := setup(t)
	handle(t, mux, "/get_service_id", map[string]any{"domain_name": "qqq.ru"}, `{"service_id":"123456"}`)
	got, _, err := c.GetServiceID(context.Background(), &ServiceIdent{DomainName: "qqq.ru"})
	if err != nil {
		t.Fatal(err)
	}
	if got.ServiceID != 123456 {
		t.Errorf("ServiceID = %d", got.ServiceID)
	}
}

func TestErrorResponse(t *testing.T) {
	c, mux := setup(t)
	mux.HandleFunc("/nop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"error_code":"PASSWORD_AUTH_FAILED","error_text":"Username/password incorrect","result":"error"}`)
	})
	_, _, err := c.Nop(context.Background())
	var e *ErrorResponse
	if !errors.As(err, &e) || e.Code != ErrCodePasswordAuth || ErrorCode(err) != ErrCodePasswordAuth {
		t.Fatalf("err = %v", err)
	}
}

func TestNonJSONResponse(t *testing.T) {
	c, mux := setup(t)
	mux.HandleFunc("/nop", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	})
	if _, _, err := c.Nop(context.Background()); err == nil || ErrorCode(err) != "" {
		t.Fatalf("err = %v", err)
	}
}

func TestItemResult(t *testing.T) {
	var items []ServiceItem
	err := json.Unmarshal([]byte(`[
		{"dname":"test.ru","result":"success","service_id":12345},
		{"error_code":"INVALID_SERVICE_ID","error_text":"service_id is invalid","result":"error","service_id":"22bug22"}
	]`), &items)
	if err != nil || len(items) != 2 || items[0].Err() != nil || items[0].ServiceID != 12345 || ErrorCode(items[1].Err()) != "INVALID_SERVICE_ID" {
		t.Fatalf("items = %+v", items)
	}
	var item ServiceItem
	_ = json.Unmarshal([]byte(`{"dname":"ya.ru","error_code":"DOMAIN_ALREADY_EXISTS","result":"Domain already exists"}`), &item)
	if ErrorCode(item.Err()) != "DOMAIN_ALREADY_EXISTS" || item.Err().Error() != "regru: DOMAIN_ALREADY_EXISTS: Domain already exists" {
		t.Errorf("Err = %v", item.Err())
	}
}

func TestFlexTypes(t *testing.T) {
	var v struct {
		A, B, C FlexInt
		D, E    FlexFloat
		F, G    FlexString
		H, I    Bool
	}
	err := json.Unmarshal([]byte(`{"A":12,"B":"34","C":"","D":"1.5","E":2,"F":7,"G":"x","H":"1","I":0}`), &v)
	if err != nil {
		t.Fatal(err)
	}
	if v.A != 12 || v.B != 34 || v.C != 0 || v.D != 1.5 || v.E != 2 || v.F != "7" || v.G != "x" || !bool(v.H) || bool(v.I) {
		t.Errorf("got %+v", v)
	}
	if b, _ := json.Marshal(struct{ X Bool }{true}); string(b) != `{"X":1}` {
		t.Errorf("Bool marshal = %s", b)
	}
}

// Example from section 2.2.4.1 of the docs.
func TestSigText(t *testing.T) {
	form := url.Values{"username": {"test"}, "sig": {"ignored"}, "input_format": {"0"}}
	input := []byte(`{"password":"test","output_content_type":"plain","show_input_params":1,"domain_name":"qqq.ru","leftdata":[1,2,3],"zero":0,"empty":"","off":false}`)
	got, err := sigText(form, input)
	if err != nil {
		t.Fatal(err)
	}
	if want := "1;1;2;3;plain;qqq.ru;test;test"; got != want {
		t.Errorf("sigText = %q, want %q", got, want)
	}
}

func TestSignatureRequest(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(nil).WithSignature("user", key)
	req, err := c.NewRequest(context.Background(), "domain/nop", map[string]string{"dname": "qqq.ru"})
	if err != nil {
		t.Fatal(err)
	}
	_ = req.ParseForm()
	if req.PostForm.Has("password") {
		t.Error("password sent with signature auth")
	}
	sig, _ := base64.StdEncoding.DecodeString(req.PostForm.Get("sig"))
	sum := sha512.Sum512([]byte("json;json;qqq.ru;user"))
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA512, sum[:], sig); err != nil {
		t.Errorf("signature does not verify: %v", err)
	}
}

func TestItemResultErrorWithoutCode(t *testing.T) {
	r := ItemResult{Result: "error", ErrorText: "DNSSEC is not available"}
	if err := r.Err(); err == nil || err.Error() != "regru: : DNSSEC is not available" {
		t.Errorf("Err = %v", err)
	}
}
