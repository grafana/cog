import "strings"

#Pattern: =~"^[a-zA-Z0-9_-]+$"
#NegativePattern: !~"^[a-zA-Z0-9_-]+$"

container: {
    minLengthConstraints: string & strings.MinRunes(1)
    maxLengthConstraints: string & strings.MaxRunes(64)
    minMaxLengthConstraints: string & strings.MinRunes(2) & strings.MaxRunes(8)
    regex: string & =~ "^[a-zA-Z0-9_-]+$"
    negativeRegex: string & !~ "^[a-zA-Z0-9_-]+$"
    regexRef: string & #Pattern
    negativeRegexRef: string & #NegativePattern
}
