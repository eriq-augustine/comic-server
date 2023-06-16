package config;

// For the defaulted getters, the defualt will be returned on ANY error
// (even if the key exists, but is of the wrong type).
// This package is meant for read-only options.

import (
    _ "embed"
    "encoding/json"
    "fmt"
    "io"
    "math"
    "os"
    "strings"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

//go:embed configs/default_config.json
var DEFAULT_CONFIG string;

var options map[string]any = make(map[string]any);

func init() {
    err := LoadString(DEFAULT_CONFIG);
    if (err != nil) {
        log.Fatal().Err(err).Msg("Failed to load the default config.");
    }

    var rawLogLevel = GetString("log.level");
    level, err := zerolog.ParseLevel(rawLogLevel);
    if (err != nil) {
        log.Fatal().Err(err).Str("level", rawLogLevel).Msg("Failed to parse the logging level.");
    }
    zerolog.SetGlobalLevel(level);
}

// See LoadReader().
func LoadFile(path string) error {
    file, err := os.Open(path);
    if (err != nil) {
        return fmt.Errorf("Could not open config file (%s): %w.", path, err);
    }
    defer file.Close();

    err = LoadReader(file);
    if (err != nil) {
        return fmt.Errorf("Unable to decode config file (%s): %w.", path, err);
    }

    return nil;
}

// See LoadReader().
func LoadString(text string) error {
    err := LoadReader(strings.NewReader(text));
    if (err != nil) {
        return fmt.Errorf("Unable to decode config from string: %w.", err);
    }

    return nil;
}

// Load data into the configuration.
// This will not clear out an existing configuration (so can load multiple files).
// If there are any key conflicts, the data loaded last will win.
// If you want to clear the config, use Reset().
func LoadReader(reader io.Reader) error {
    decoder := json.NewDecoder(reader);

    var fileOptions map[string]any;

    err := decoder.Decode(&fileOptions);
    if (err != nil) {
        return err;
    }

    for key, val := range fileOptions {
        // encoding/json uses float64 as its default numeric type.
        // Check if it is actually an integer.
        floatVal, ok := val.(float64);
        if (ok) {
            if (math.Trunc(floatVal) == floatVal) {
                val = int(floatVal);
            }
        }

        options[key] = val;
    }

    return nil;
}

func Reset() {
    options = make(map[string]any);
}

func Has(key string) bool {
    _, present := options[key];
    return present;
}

func Get(key string) any {
    val, present := options[key];
    if (!present) {
        log.Fatal().Str("key", key).Msg("Config key does not exist.");
    }

    return val;
}

func GetDefault(key string, defaultVal any) any {
    if (!Has(key)) {
        return defaultVal;
    }

    val, _ := options[key];
    return val;
}

func GetString(key string) string {
    val := Get(key);

    stringVal, ok := val.(string);
    if (!ok) {
        log.Fatal().Str("key", key).Interface("value", val).Msg("Config option is not a string type.");
    }

    return stringVal;
}

func GetStringDefault(key string, defaultVal string) string {
    val := GetDefault(key, defaultVal);

    stringval, ok := val.(string);
    if (!ok) {
        return defaultVal;
    }

    return stringval;
}

func GetInt(key string) int {
    val := Get(key);

    intVal, ok := val.(int);
    if (!ok) {
        log.Fatal().Str("key", key).Interface("value", val).Msg("Config option is not an int type.");
    }

    return intVal;
}

func GetIntDefault(key string, defaultVal int) int {
    val := GetDefault(key, defaultVal);

    intVal, ok := val.(int);
    if (!ok) {
        return defaultVal;
    }

    return intVal;
}

func GetBool(key string) bool {
    val := Get(key);

    boolVal, ok := val.(bool);
    if (!ok) {
        log.Fatal().Str("key", key).Interface("value", val).Msg("Config option is not a bool type.");
    }

    return boolVal;
}

func GetBoolDefault(key string, defaultVal bool) bool {
    val := GetDefault(key, defaultVal);

    boolVal, ok := val.(bool);
    if (!ok) {
        return defaultVal;
    }

    return boolVal;
}
