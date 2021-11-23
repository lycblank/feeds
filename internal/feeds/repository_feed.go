package feeds

import "context"

type FeedRepository interface {
	GetFeedItems(ctx context.Context, itemId []int64) ([]*FeedItem, error)
	FetchFeedItems(ctx context.Context, fid int64, offset, limit int) ([]*FeedItem, int64, error)
	GetFeed(ctx context.Context, fid int64) (*Feed, error)

	GetAllFeed(ctx context.Context) ([]*Feed, error)
	SaveFeedItems(ctx context.Context, items []*FeedItem) error
}