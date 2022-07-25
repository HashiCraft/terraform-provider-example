package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupGetTestServer(t *testing.T) string {
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			e := json.NewEncoder(rw)

			br := &BlockRequest{
				X:        1,
				Y:        2,
				Z:        3,
				Material: "test",
				Facing:   "test",
				Half:     "test",
			}

			err := e.Encode(br)
			if err != nil {
				t.Fatal(err)
			}
		}),
	)

	return s.URL
}

func TestGETGetsBlock(t *testing.T) {
	url := setupGetTestServer(t)
	c := NewClient(url)

	b, err := c.GetBlock(1, 2, 3)
	require.NoError(t, err)

	require.Equal(t, 1, b.X)
	require.Equal(t, 2, b.Y)
	require.Equal(t, 3, b.Z)

}
