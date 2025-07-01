package codegen

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "syscall"

    "github.com/getkin/kin-openapi/openapi3"

    "github.com/grafana/cog/internal/ast"
    "github.com/grafana/cog/internal/openapi"
)

type OpenAPISources struct {
    InputBase `yaml:",inline"`

    URL          string   `yaml:"url"`
    Path         string   `yaml:"path"`
    PackageRegex string   `yaml:"package_regex,omitempty"`
    AuthBearer   string   `yaml:"auth_bearer,omitempty"` // To access to remote server
    Private      string   `yaml:"private,omitempty"`     // Generates public or private api
    SkipURLs     []string `yaml:"skip_urls,omitempty"`

    Validate bool `yaml:"validate,omitempty"`
}

type sources struct {
    Sources []source `json:"sources"`
}

type source struct {
    Url     string `json:"url"`
    Private bool   `json:"private"`
}

func (input *OpenAPISources) loadSchemas(ctx context.Context) ([]*openapi3.T, error) {
    var body []byte
    var err error

    private, _ := strconv.ParseBool(input.Private)

    if input.Path != "" {
        body, err = os.ReadFile(input.Path)
        if err != nil {
            return nil, err
        }
    } else if input.URL != "" {
        resp, err := http.Get(input.URL)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return nil, errors.New(resp.Status)
        }

        body, err = io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
    }

    var parsedSources *sources
    if err = json.Unmarshal(body, &parsedSources); err != nil {
        return nil, err
    }

    loader := openapi3.NewLoader()
    loader.Context = ctx
    loader.IsExternalRefsAllowed = true

    schemas := make([]*openapi3.T, 0)
    for _, s := range parsedSources.Sources {
        if input.shouldSkipURL(s.Url) {
            continue
        }
        if s.Private != private {
            continue
        }
        req, err := http.NewRequest(http.MethodGet, s.Url, nil)
        if err != nil {
            return nil, err
        }

        if ok, err := input.checkIfLoging(s.Url); err == nil && ok {
            req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv(input.AuthBearer)))
        }

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            if errors.Is(err, syscall.ECONNREFUSED) {
                continue // localhost ?
            }
            return nil, err
        }

        if resp.StatusCode != http.StatusOK {
            fmt.Println(s.Url)
            continue
        }

        body, err = io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }

        schema, err := loader.LoadFromData(body)
        if err != nil {
            return nil, err
        }

        schemas = append(schemas, schema)
    }

    return schemas, nil
}

func (input *OpenAPISources) packageName(schema *openapi3.T) string {
    if input.PackageRegex != "" && schema.Info != nil {
        reg, err := regexp.Compile(input.PackageRegex)
        if err != nil {
            fmt.Println(err)
            return schema.Info.Title
        }
        return reg.FindString(schema.Info.Title)
    }

    if schema.Info == nil {
        return ""
    }

    return schema.Info.Title
}

func (input *OpenAPISources) interpolateParameters(interpolator ParametersInterpolator) {
    input.InputBase.interpolateParameters(interpolator)
    input.Private = interpolator(input.Private)

    input.URL = interpolator(input.URL)
}

func (input *OpenAPISources) LoadSchemas(ctx context.Context) (ast.Schemas, error) {
    oapiSchemas, err := input.loadSchemas(ctx)
    if err != nil {
        return nil, err
    }

    var schemas ast.Schemas
    for _, oapiSchema := range oapiSchemas {
        schema, err := openapi.GenerateAST(ctx, oapiSchema, openapi.Config{
            Package:        input.packageName(oapiSchema),
            SchemaMetadata: input.schemaMetadata(),
            Validate:       input.Validate,
        })
        if err != nil {
            return nil, err
        }

        filteredSchemas, err := input.filterSchema(schema)
        if err != nil {
            return nil, err
        }

        schemas = append(schemas, filteredSchemas...)
    }

    return schemas, nil
}

func (input *OpenAPISources) checkIfLoging(url string) (bool, error) {
    client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    }}

    resp, err := client.Get(url)
    if err != nil {
        return false, err
    }

    if resp.StatusCode >= 300 && resp.StatusCode < 400 {
        return true, nil
    }

    return false, nil
}

func (input *OpenAPISources) shouldSkipURL(url string) bool {
    for _, skipURL := range input.SkipURLs {
        if skipURL == url {
            return true
        }
    }

    return false
}
