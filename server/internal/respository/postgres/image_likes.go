package postgres

import "github.com/pillowskiy/gopix/pkg/batch"

var imageLikesBatchAgg = batch.NewKGAggregator[likeBatchItem]()

type imageLikesAnalyticsAgg struct {
	ImageID       int `db:"image_id"`
	LikesCount    int `db:"likes_count"`
	DislikesCount int `db:"dislikes_count"`
}

type likeBatchItem struct {
	ImageID int `db:"image_id"`
	UserID  int `db:"user_id"`
	Liked   bool
}

func (i likeBatchItem) Group() string {
	return imageGroupKey(i.ImageID)
}

func (i likeBatchItem) Key() string {
	return imageWithUserKey(i.ImageID, i.UserID)
}

func (i likeBatchItem) Count() int {
	if i.Liked {
		return 1
	}

	return -1
}
