package persistence

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/jinzhu/copier"
	"github.com/lycblank/feeds/internal/feeds"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hash/fnv"
	"io"
)

type MysqlFeedRepository struct {
	gdb *gorm.DB
}

func NewMysqlFeedRepository(gdb *gorm.DB) feeds.FeedRepository {
	gdb.AutoMigrate(&FeedGorm{},&FeedItemGorm{})
	return &MysqlFeedRepository{
		gdb:gdb,
	}
}

func (fr *MysqlFeedRepository) GetFeedItems(ctx context.Context, itemId []int64) ([]*feeds.FeedItem, error) {
	var items []*FeedItemGorm
	err := fr.gdb.Model(FeedItemGorm{}).Where("id in ?", itemId).Find(&items).Error
	if err != nil {
		return nil, err
	}

	feedItems := make([]*feeds.FeedItem, 0, len(items))
	for _, item := range items {
		feedItems = append(feedItems, item.ToFeedItem())
	}

	return feedItems, nil
}
func (fr *MysqlFeedRepository) FetchFeedItems(ctx context.Context, fid int64, offset, limit int) ([]*feeds.FeedItem, int64, error) {
	var items []*FeedItemGorm
	err := fr.gdb.Model(FeedItemGorm{}).Where("feed_id = ?", fid).Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	feedItems := make([]*feeds.FeedItem, 0, len(items))
	for _, item := range items {
		feedItems = append(feedItems, item.ToFeedItem())
	}

	var totalCount int64
	fr.gdb.Model(FeedItemGorm{}).Where("feed_id = ?", fid).Count(&totalCount)

	return feedItems, totalCount, nil
}

func (fr *MysqlFeedRepository) GetFeed(ctx context.Context, fid int64) (*feeds.Feed, error) {
	feed := &FeedGorm{}
	err := fr.gdb.Model(FeedGorm{}).Where("id = ?", fid).Find(feed).Error
	if err != nil {
		return nil, err
	}

	return feed.ToFeed(), nil
}
func (fr *MysqlFeedRepository) GetAllFeed(ctx context.Context) ([]*feeds.Feed, error) {
	var feedList []*FeedGorm
	err := fr.gdb.Model(FeedGorm{}).Find(&feedList).Error
	if err != nil {
		return nil, err
	}

	fs := make([]*feeds.Feed, 0, len(feedList))
	for _, feed := range feedList {
		fs = append(fs, feed.ToFeed())
	}

	return fs, nil
}
func (fr *MysqlFeedRepository) SaveFeedItems(ctx context.Context, items []*feeds.FeedItem) error {
	feedItems := make([]*FeedItemGorm, 0, len(items))
	for _, item := range items {
		feedItem := &FeedItemGorm{}
		feedItem.TransFeedItem(item)
		feedItems = append(feedItems, feedItem)
	}
	err := fr.gdb.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(feedItems, len(feedItems)).Error
	if err != nil {
		return err
	}
	for i:=0;i<len(items);i++{
		items[i].SetVersion(feedItems[i].Version)
	}
	return nil
}

type FeedItemGorm struct {
	Id int64         `gorm:"column:id;primaryKey;autoIncrement"`
	Title string     `gorm:"column:title"`
	Link string      `gorm:"column:link;size:1024"`
	CreateTime int64 `gorm:"column:create_time"`
	FeedId int64 	 `gorm:"column:feed_id"`
	Version int32 	 `gorm:"column:version"`
	Hash string `gorm:"column:hash;unique_index:hc"`
	Hash64 uint64 `gorm:"column:hash64;;unique_index:hc"`
}
func (FeedItemGorm) TableName() string {
	return "t_feed_item"
}

func (fi *FeedItemGorm) ToFeedItem() *feeds.FeedItem {
	feedItem := &feeds.FeedItem{}
	_ = copier.Copy(feedItem, fi)
	feedItem.SetVersion(fi.Version)

	return feedItem
}

func (fi *FeedItemGorm) TransFeedItem(feedItem *feeds.FeedItem) {
	_ = copier.Copy(fi, feedItem)
	fi.Version = feedItem.GetVersion()

	h := md5.New()
	io.WriteString(h, feedItem.Link)
	fi.Hash = hex.EncodeToString(h.Sum(nil))

	hash64 := fnv.New64()
	io.WriteString(hash64, feedItem.Link)
	fi.Hash64 = hash64.Sum64()
}


type FeedGorm struct {
	Id int64         `gorm:"column:id;primaryKey;autoIncrement"`
	Title string     `gorm:"column:title"`
	Link string      `gorm:"column:link;size:1024"`
	Version int32 	 `gorm:"column:version"`
	Hash string `gorm:"column:hash;unique_index:hc"`
	Hash64 uint64 `gorm:"column:hash64;;unique_index:hc"`
}
func (FeedGorm) TableName() string {
	return "t_feed"
}

func (f *FeedGorm) ToFeed() *feeds.Feed {
	feed := &feeds.Feed{}
	_ = copier.Copy(feed, f)
	feed.SetVersion(f.Version)

	return feed
}




