package feeds

type Feed struct {
	Id int64         `json:"id"`
	Title string     `json:"title"`
	Link string      `json:"link"`
	CreateTime int64 `json:"create_time"`
	version int32
}

func (fi *Feed) SetVersion(version int32) {
	fi.version = version
}

func (fi *Feed) GetVersion() int32 {
	return fi.version
}

type FeedItem struct {
	Id int64         `json:"id"`
	Title string     `json:"title"`
	Link string      `json:"link"`
	CreateTime int64 `json:"create_time"`
	FeedId int64 `json:"feed_id"`
	version int32
}

func (fi *FeedItem) SetVersion(version int32) {
	fi.version = version
}

func (fi *FeedItem) GetVersion() int32 {
	return fi.version
}


