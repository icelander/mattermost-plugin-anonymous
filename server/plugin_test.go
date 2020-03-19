package main

import (
	"bytes"
	"encoding/json"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PUB_KEY = []byte{1, 1, 2, 3}

func getFunctionalPlugin() Plugin {
	plugin := Plugin{}
	mockapi := &plugintest.API{}
	mockapi.On("KVSet", "user_1", PUB_KEY).Return(nil)
	mockapi.On("KVGet", "user_1").Return(PUB_KEY, nil)
	plugin.SetAPI(mockapi)
	return plugin
}

//test functionality of store and get key requests
func TestHandleGetPublicKey(t *testing.T) {
	tassert := assert.New(t)
	plugin := getFunctionalPlugin()
	w := httptest.NewRecorder()

	var data bytes.Buffer
	_ = json.NewEncoder(&data).Encode(SetPublicKeyRequest{PublicKey: PUB_KEY})

	r := httptest.NewRequest(http.MethodPost, "/api/pub_key/set", &data)

	r.Header.Add(USER_ID_HEADER_KEY, "user_1")
	plugin.ServeHTTP(nil, w, r)
	wr := w.Result()
	tassert.Equal(wr.StatusCode, 200)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/api/pub_key/get", nil)
	r.Header.Add(USER_ID_HEADER_KEY, "user_1")
	plugin.ServeHTTP(nil, w, r)
	tassert.NotNil(w.Result())
	tassert.Equal(w.Result().StatusCode, 200)
	bodyBytes, err := ioutil.ReadAll(w.Result().Body)
	tassert.Nil(err)

	type body struct {
		PublicKey []byte `json:"public_key"`
	}
	var b body
	_ = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&b)
	tassert.Equal(b.PublicKey, PUB_KEY)
}
