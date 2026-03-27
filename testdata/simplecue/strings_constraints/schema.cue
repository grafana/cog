import "strings"

container: {
    minLengthConstraints: string & strings.MinRunes(1)
    maxLengthConstraints: string & strings.MaxRunes(64)
    minMaxLengthConstraints: string & strings.MinRunes(2) & strings.MaxRunes(8)
    regex: string & =~ "^[a-zA-Z0-9_-]+$"
    negativeRegex: string & !~ "^[a-zA-Z0-9_-]+$"
}
