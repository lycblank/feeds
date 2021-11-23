package users

type User struct {
	Uid int64            `json:"uid"`
	Profile *UserProfile `json:"profile"`
	Feeds []*UserFeed    `json:"feeds"`
}

func (u *User) AddFeed(feedId int64) {
	u.Feeds = append(u.Feeds, &UserFeed{
		FeedId: feedId,
		UnReadNum: 0,
	})
}

type UserProfile struct {
	Nickname string `json:"nickname"`
}

type UserFeed struct {
	FeedId int64
	UnReadNum int32
}
