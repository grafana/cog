package envvars

import "os"

// VarUpdateGolden is the name of the env var to trigger updating golden test files.
const VarUpdateGolden = "COG_UPDATE_GOLDEN"

// UpdateGoldenFiles determines whether tests should update txtar
// archives on failures.
// It is controlled by setting COG_UPDATE_GOLDEN to a non-empty string like "true".
var UpdateGoldenFiles = os.Getenv(VarUpdateGolden) != "" //nolint: gochecknoglobals
