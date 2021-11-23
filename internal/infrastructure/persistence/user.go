package persistence

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/lycblank/feeds/internal/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlUserRepository struct {
	gdb *gorm.DB
}

func NewMysqlUserRepository(gdb *gorm.DB) users.UserRepository {
	return &MysqlUserRepository{
		gdb:gdb,
	}
}

func (ur *MysqlUserRepository) GetUser(ctx context.Context, uid int64) (*users.User, error) {
	userProfile := &UserProfileGorm{}
	if err := ur.gdb.Model(UserProfileGorm{}).Where("uid = ?", uid).Find(&userProfile).Error; err != nil {
		return nil, err
	}

	var userFeeds []*UserFeedGorm
	ur.gdb.Model(UserFeedGorm{}).Where("uid = ?", uid).Find(&userFeeds)

	user := &users.User{
		Uid: uid,
		Profile: userProfile.ToUserProfile(),
		Feeds: make([]*users.UserFeed, 0, len(userFeeds)),
	}

	for _, userFeed := range userFeeds {
		var totalCount int64
		feed := userFeed.ToUserFeed()
		ur.gdb.Model(UserFeedItemGorm{}).Where("feed_id = ? and is_new = ?", userFeed.Id, true).Count(&totalCount)
		feed.UnReadNum = int32(totalCount)
		user.Feeds = append(user.Feeds, feed)
	}

	return user, nil
}

func (ur *MysqlUserRepository) SaveUser(ctx context.Context, user *users.User) error {
	err := ur.gdb.Transaction(func(tx *gorm.DB) error {
		userProfile := &UserProfileGorm{}
		userProfile.TransUserProfile(user.Uid, user.Profile)
		if err := tx.Save(userProfile).Error; err != nil {
			return err
		}
		userFeeds := make([]*UserFeedGorm, 0, len(user.Feeds))
		for _, feed := range user.Feeds {
			userFeed := &UserFeedGorm{}
			userFeed.TransUserFeed(user.Uid, feed)
			userFeeds = append(userFeeds, userFeed)
		}
		return tx.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(userFeeds, len(userFeeds)).Error
	})
	return err
}


type UserProfileGorm struct {
	Uid int64       `gorm:"column:uid;primaryKey"`
	Nickname string `gorm:"column:nickname"`
}

func (UserProfileGorm) TableName() string {
	return "t_user_profile"
}

func (up *UserProfileGorm) ToUserProfile() *users.UserProfile {
	userProfile := &users.UserProfile{}
	_ = copier.Copy(userProfile, up)
	return userProfile
}

func (up *UserProfileGorm) TransUserProfile(uid int64, userProfile *users.UserProfile) {
	_ = copier.Copy(up, userProfile)
	up.Uid = uid
}

type UserFeedGorm struct {
	Id int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Uid int64    `gorm:"column:uid;unique_index:uf"`
	FeedId int64 `gorm:"column:feed_id;unique_index:uf"`
}

func (UserFeedGorm) TableName() string {
	return "t_user_feed"
}

func (uf *UserFeedGorm) ToUserFeed() *users.UserFeed {
	userFeed := &users.UserFeed{}
	_ = copier.Copy(userFeed, uf)
	return userFeed
}

func (uf *UserFeedGorm) TransUserFeed(uid int64, userFeed *users.UserFeed) {
	_ = copier.Copy(uf, userFeed)
	uf.Uid = uid
}

type UserFeedItemGorm struct {
	Id int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Uid int64    `gorm:"column:uid"`
	FeedId int64 `gorm:"column:feed_id"`
	ItemId int64 `gorm:"column:item_id"`
	IsNew bool 	 `gorm:"column:is_new"`
}

func (UserFeedItemGorm) TableName() string {
	return "t_user_feed_item"
}
