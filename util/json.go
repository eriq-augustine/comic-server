package util

import (
   "encoding/json"

   "github.com/rs/zerolog/log"
)

func ToJSON(data any) (string, error) {
   bytes, err := json.Marshal(data);
   if (err != nil) {
      log.Error().Err(err).Interface("object", data).Msg("Could not marshal.");
      return "", err;
   }

   return string(bytes), nil;
}
