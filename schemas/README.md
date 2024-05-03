## Notes

The `expr.json` file is copied from https://github.com/grafana/grafana/blob/main/pkg/expr/query.request.schema.json

Modifications:

* occurrences of `https://json-schema.org/draft-04/schema` were replaced by `http://json-schema.org/draft-04/schema`
  * see: https://github.com/santhosh-tekuri/jsonschema/blob/534765aa6277d702a0aaf99c986492dda697d48f/compiler.go#L128-L138
  * could be fixed by upgrading santhosh-tekuri/jsonschema (we're using a really old version of the lib :/)
