package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	resp, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	response, err := client.Do(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}

	response.Header.Set("User-Agent", "gator")

	defer response.Body.Close()

	feed, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	rss := &RSSFeed{}
	xml.Unmarshal(feed, &rss)

	html.UnescapeString(rss.Channel.Title)
	html.UnescapeString(rss.Channel.Description)

	for i := range rss.Channel.Item {
		html.UnescapeString(rss.Channel.Item[i].Title)
		html.UnescapeString(rss.Channel.Item[i].Description)
	}

	return rss, nil
}
