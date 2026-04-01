package jsonschema

import (
	"context"

	"github.com/grafana/cog/internal/httputil"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type HTTPURLLoader struct{}

func (loader *HTTPURLLoader) Load(url string) (any, error) {
	body, err := httputil.LoadURL(context.Background(), url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return jsonschema.UnmarshalJSON(body)
}
