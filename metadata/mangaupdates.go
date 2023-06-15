package metadata

import (
    neturl "net/url"
    "regexp"
    "strconv"
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/model"
    "github.com/eriq-augustine/comic-server/util"
)
const SOURCE_MANGA_UPDATES = "MangaUpdates";

const BASE_SEARCH_URL = "https://www.mangaupdates.com/series.html";
const BASE_SERIES_URL = "https://www.mangaupdates.com/series/";

func init() {
    metadataSources[SOURCE_MANGA_UPDATES] = crawlMangaUdates;
}

func crawlMangaUdates(query string, year int, series *model.Series) ([]*model.MetadataCrawl, error) {
    ids, err := mangaupdatesSearch(query, year);
    if (err != nil) {
        return nil, err;
    }

    crawls := make([]*model.MetadataCrawl, 0);

    for _, id := range ids {
        crawl, err := managaupdatesFetchSeries(id);
        if (err != nil) {
            return nil, err;
        }

        if (crawl != nil) {
            crawl.SourceSeries = series;
            crawls = append(crawls, crawl);
        }
    }

    return crawls, nil;
}

func mangaupdatesSearch(query string, year int) ([]string, error) {
    values := neturl.Values{};
    values.Set("search", query);
    searchURL := BASE_SEARCH_URL + "?" + values.Encode();

    pageBytes, err := util.GetWithCache(searchURL);
    if (err != nil) {
        return nil, err;
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(pageBytes)));
    if (err != nil) {
        return nil, err;
    }

    var ids = make([]string, 0);

    doc.Find(`div.text > a[alt="Series Info"]`).Each(func(id int, ele *goquery.Selection) {
        url, exists := ele.Attr("href");
        if (!exists) {
            return;
        }

        match := regexp.MustCompile(`^.*www.mangaupdates.com/series/(\w+)/.*$`).FindStringSubmatch(url);
        if (match == nil) {
            return;
        }

        ids = append(ids, match[1]);
    });

    return ids, nil;
}

func managaupdatesFetchSeries(id string) (*model.MetadataCrawl, error) {
    url := BASE_SERIES_URL + neturl.PathEscape(id);

    pageBytes, err := util.GetWithCache(url);
    if (err != nil) {
        return nil, err;
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(pageBytes)));
    if (err != nil) {
        return nil, err;
    }

    crawl := model.EmptyCrawl(SOURCE_MANGA_UPDATES, id);

    crawl.Name = cleanHTMLText(doc.Find(`span.releasestitle`).First());
    crawl.URL = &url;

    // Parse out all the metadata blocks.
    metadataBlocks := make(map[string]*goquery.Selection);
    doc.Find(`div.sCat`).Each(func(id int, ele *goquery.Selection) {
        metadataBlocks[ele.Text()] = ele.Next().Clone();
    });

    node, exists := metadataBlocks["Year"];
    if (exists) {
        year, err := strconv.Atoi(strings.TrimSpace(node.Text()));
        if (err != nil) {
            log.Warn().Err(err).Str("source", SOURCE_MANGA_UPDATES).Str("id", id).Msg("Failed to parse year.");
        } else {
            crawl.Year = &year;
        }
    }

    node, exists = metadataBlocks["Author(s)"];
    if (exists) {
        author := cleanHTMLText(node);

        author = strings.ReplaceAll(author, "[Add]", "");
        author = regexp.MustCompile(`\s*\n\s*`).ReplaceAllString(author, ", ");
        author = regexp.MustCompile(`\s+`).ReplaceAllString(author, " ");

        if (author != "") {
            crawl.Author = &author;
        }
    }

    node, exists = metadataBlocks["Description"];
    if (exists) {
        var description string;

        longDescription := node.Find(`div#div_desc_more`);
        if (longDescription.Length() > 0) {
            longDescription.Find(`a`).Remove();
            description = cleanHTMLText(longDescription);
        } else {
            description = cleanHTMLText(node);
        }

        if (description != "") {
            crawl.Description = &description;
        }
    }

    node, exists = metadataBlocks["Image [Report Inappropriate Content]"];
    if (exists) {
        imageURL, exists := node.Find(`img`).Attr("src");
        if (exists) {
            path, err := util.FetchImage(imageURL);
            if (err != nil) {
                log.Warn().Err(err).Str("source", SOURCE_MANGA_UPDATES).Str("id", id).Str("url", imageURL).Msg("Failed to fetch image.");
            } else {
                crawl.CoverImagePath = &path;
            }
        }
    }

    return crawl, nil;
}

func cleanHTMLText(node *goquery.Selection) string {
    node = node.Clone();
    html, err := node.Html();
    if (err != nil) {
        log.Warn().Err(err).Msg("Could not get html for text cleaning.");
        return "";
    }

    node.SetHtml(strings.ReplaceAll(html, `<br/>`, "\n"));
    text := node.Text();

    text = strings.TrimSpace(text);

    if (text == "N/A") {
        return "";
    }

    return text;
}
