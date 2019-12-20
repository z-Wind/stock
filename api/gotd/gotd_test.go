package gotd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/z-Wind/stock/api/gotd/api"
	"github.com/z-Wind/stock/instance"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var (
	apiKey      = instance.TdAPIKey
	redirectURL = instance.TdURL
)

func init() {
	path, err := filepath.Abs(filepath.Join("../../instance/", "secrets.properties"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Secrets Path:%s\n", path)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var td struct {
		AccessToken  string `json:"AccessToken"`
		RefreshToken string `json:"RefreshToken"`
	}

	json.Unmarshal(b, &td)
	if apiKey == "" {
		panic("td is null")
	}
}

func TestNewTD(t *testing.T) {
	path, err := filepath.Abs("../../instance/")
	if err != nil {
		panic(err)
	}
	td, err := NewTD(apiKey, redirectURL, path)
	if err != nil {
		panic(err)
	}
	t.Logf("%+v", td)
}

func TestAuth(t *testing.T) {
	t.Run("OpenAuthServer", testOpenAuthServer)
	//t.Run("TDAmeritrade_Authentication", OpenAuthServer)
}

func testOpenAuthServer(t *testing.T) {
	td := TDAmeritrade{
		apiKey:   apiKey,
		certPath: filepath.Join("../../instance/", "cert.pem"),
		keyPath:  filepath.Join("../../instance/", "key.pem"),
	}
	code := "123"
	go func() {
		//time.Sleep(time.Second * 1)
		// disable https security verify
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		time.Sleep(time.Second * 3)
		res, err := http.Get("https://localhost:8080/?code=" + code)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		t.Logf("%s", content)
	}()

	err := td.openAuthServer()
	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if td.code != code {
		t.Errorf("Get: %s, Want: %s ", td.code, code)
	}
}

func testTDAmeritrade_Authentication(t *testing.T) {
	td := TDAmeritrade{
		apiKey:   apiKey,
		certPath: filepath.Join("../../instance/", "cert.pem"),
		keyPath:  filepath.Join("../../instance/", "key.pem"),
	}

	err := td.Authentication(redirectURL)
	if err != nil {
		panic(err)
	}

	t.Logf("%+v", td)
	got, err := td.GetAccounts([]string{api.AccountFieldPositions, api.AccountFieldOrders})
	if err != nil {
		t.Errorf("%s", err)
		panic(err)
	}
	content, _ := json.MarshalIndent(got, "", "    ")
	t.Log(string(content))

	err = td.RefreshAccessToken()
	if err != nil {
		t.Errorf("%s", err)
		panic(err)
	}
	ws, err := td.GetWatchlistsforMultipleAccounts()
	if err != nil {
		t.Errorf("%s", err)
		panic(err)
	}
	content, _ = json.MarshalIndent(ws, "", "    ")
	t.Log(string(content))

	if td.IsAccessTokenFail() {
		t.Error("IsAccessTokenExpired should be false")
	}

	if td.IsRefreshTokenExpired() {
		t.Error("IsRefreshTokenExpired should be false")
	}
}

func TestTDAmeritrade_SaveLoad(t *testing.T) {
	td := TDAmeritrade{
		apiKey:      apiKey,
		certPath:    filepath.Join("../../instance/", "cert.pem"),
		keyPath:     filepath.Join("../../instance/", "key.pem"),
		code:        "code",
		redirectURI: "redirectURI",
		Token: Token{
			AccessToken:       "AccessToken",
			AccessTokenStart:  time.Now(),
			RefreshToken:      "RefreshToken",
			RefreshTokenStart: time.Now().AddDate(0, 1, 0),
		},
	}
	tdCopy := td

	os.Mkdir("test", os.ModePerm)

	err := td.SaveAuth("test")
	if err != nil {
		panic(err)
	}
	err = td.LoadAuth("test")
	if err != nil {
		panic(err)
	}
	if !cmp.Equal(td, tdCopy, cmp.AllowUnexported(td)) {
		t.Errorf("Get %+v want %+v", td, tdCopy)
	}

	os.RemoveAll("test")
}

func TestTDAmeritrade_GetAccountIDs(t *testing.T) {
	tests := []struct {
		name    string
		notWant []int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{"Test", []int64{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td, err := NewTD(apiKey, redirectURL, "../../instance")
			if err != nil {
				panic(err)
			}
			got, err := td.GetAccountIDs()
			if (err != nil) != tt.wantErr {
				t.Errorf("TDAmeritrade.GetAccountIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(got, tt.notWant) {
				t.Errorf("TDAmeritrade.GetAccountIDs() = %v, notWant %v", got, tt.notWant)
			}
		})
	}
}
