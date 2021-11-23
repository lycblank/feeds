package feeds

import (
	"context"
	"errors"
	"github.com/lycblank/feeds/internal/infrastructure/conf"
	"github.com/mmcdole/gofeed"
	"time"
)

type FeedTaskService struct {
	feedRepository FeedRepository
	fetchInterval time.Duration
	exit chan struct{}
}

func NewFeedTaskService(config *conf.Config, feedRepository FeedRepository) *FeedTaskService {
	return &FeedTaskService{
		feedRepository:feedRepository,
		fetchInterval: time.Duration(config.Feed.FetchInterval) * time.Second,
		exit: make(chan struct{}),
	}
}

func (fts *FeedTaskService) Run() {
	fts.runFetch()

	tm := time.NewTimer(fts.fetchInterval)
	for {
		select {
		case <-tm.C:
			fts.runFetch()
			tm.Reset(fts.fetchInterval)
		case <-fts.exit:
			return
		}
	}
}

func (fts *FeedTaskService) runFetch() {
	feeds, err := fts.getAllFeed()
	if err != nil {
		return
	}
	for _, feed := range feeds {
		items, err := fts.fetch(feed)
		if err != nil {
			continue
		}
		_ = fts.saveFeedItems(items)
	}
}

func (fts *FeedTaskService) fetch(feed *Feed) ([]*FeedItem, error) {
	if feed == nil || feed.Link == "" {
		return nil, errors.New("feed is invalidate")
	}
	fp := gofeed.NewParser()
	rawFeed, err := fp.ParseURL(feed.Link)
	if err != nil {
		return nil, err
	}
	nowTime := time.Now().Unix()
	items := make([]*FeedItem, 0, len(rawFeed.Items))
	for _, item := range rawFeed.Items {
		items = append(items, &FeedItem{
			Title:item.Title,
			Link:item.Link,
			CreateTime:nowTime,
			FeedId:feed.Id,
		})
	}
	return items, nil
}

func (fts *FeedTaskService) getAllFeed() ([]*Feed, error) {
	return fts.feedRepository.GetAllFeed(context.Background())
}

func (fts *FeedTaskService) saveFeedItems(items []*FeedItem) error {
	return fts.feedRepository.SaveFeedItems(context.Background(), items)
}




