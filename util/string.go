package util

import (
    "strings"
)

// A hyphen (dash).
const UNSAFE_JOIN_REPLACE_CHAR = "-";

// A figure dash ("‒").
const UNSAFE_JOIN_DELIM = "\u2012";

// Join together strings by first replacing any instances of UNSAFE_JOIN_DELIM with UNSAFE_JOIN_REPLACE_CHAR,
// and then joining the strings together with UNSAFE_JOIN_DELIM.
// This could mean that UnsafeSplit(UnsafeJoin(values)) != values.
// However, we will use a delimiter that is rare and typically represented visually the same as the replacement character.
func UnsafeJoin(values []string) string {
    newValues := make([]string, len(values));
    for i, value := range values {
        newValues[i] = strings.ReplaceAll(value, UNSAFE_JOIN_DELIM, UNSAFE_JOIN_REPLACE_CHAR);
    }

    return strings.Join(newValues, UNSAFE_JOIN_DELIM);
}

func UnsafeSplit(text string) []string {
    return strings.Split(text, UNSAFE_JOIN_DELIM);
}
