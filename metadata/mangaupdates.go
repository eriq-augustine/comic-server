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

func crawlMangaUdates(query string, year string, series *model.Series) ([]*model.MetadataCrawl, error) {
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

func mangaupdatesSearch(query string, year string) ([]string, error) {
    values := neturl.Values{};
    values.Set("search", query);
    searchURL := BASE_SEARCH_URL + "?" + values.Encode();

    page, err := util.GetWithCache(searchURL);
    if (err != nil) {
        return nil, err;
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(page));
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

    page, err := util.GetWithCache(url);
    if (err != nil) {
        return nil, err;
    }

    doc, err := goquery.NewDocumentFromReader(strings.NewReader(page));
    if (err != nil) {
        return nil, err;
    }

    crawl := model.EmptyCrawl();

    source := SOURCE_MANGA_UPDATES;
    crawl.MetadataSource = &source;
    crawl.MetadataSourceID = &id;

    crawl.Name = doc.Find(`span.releasestitle`).First().Text();

    // Parse out all the metadata blocks.
    metadataBlocks := make(map[string]*goquery.Selection);
    doc.Find(`div.sCat`).Each(func(id int, ele *goquery.Selection) {
        metadataBlocks[ele.Text()] = ele.Next().Clone();
    });

    _, exists := metadataBlocks["Year"];
    if (exists) {
        year, err := strconv.Atoi(strings.TrimSpace(metadataBlocks["Year"].Text()));
        if (err != nil) {
            log.Warn().Err(err).Str("source", SOURCE_MANGA_UPDATES).Str("id", id).Msg("Failed to parse year.");
        } else {
            crawl.Year = &year;
        }
    }

    _, exists = metadataBlocks["Author(s)"];
    if (exists) {
        author := metadataBlocks["Author(s)"].Text();
        crawl.Author = &author;
    }

    return crawl, nil;
}
