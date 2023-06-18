package web

import (
    "net/http"

    "github.com/eriq-augustine/comic-server/util"
)

func handleImageBlob(matches []string, response http.ResponseWriter, request *http.Request) error {
    imageRelPath := matches[1];
    imagePath := util.GetImagePath(imageRelPath);

    if (!util.PathExists(imagePath)) {
        http.NotFound(response, request);
        return nil;
    }

    http.ServeFile(response, request, imagePath);
    return nil;
}
