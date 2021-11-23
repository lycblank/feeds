package feeds

import (
	"context"
)

type FeedService struct {
	feedRepository FeedRepository
}

func (fs *FeedService) FetchFeedItems(ctx context.Context, fid int64, offset,limit int) ([]*FeedItem, int64, error){
	return fs.feedRepository.FetchFeedItems(ctx, fid, offset, limit)
}

func (fs *FeedService) FetchFeedItemsByIds(ctx context.Context, itemIds []int64) ([]*FeedItem, error) {
	return fs.feedRepository.GetFeedItems(ctx, itemIds)
}

func (fs *FeedService) AddFeed(ctx context.Context, feed *Feed) error {
	return nil
}

func (fs *FeedService) FetchFeedList(ctx context.Context, fids []int64) ([]*Feed, error) {
	return nil, nil
}

func (fs *FeedService) DelFeed(ctx context.Context, fid []int64) error {
	return nil
}

