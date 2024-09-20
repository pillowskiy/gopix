package postgres

import (
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/batch"
)

var imageLikesBatchAgg = batch.NewKGAggregator[likeBatchItem]()

type imageLikesAnalytics struct {
	ImageID       domain.ID `db:"image_id"`
	InsertedCount int       `db:"inserted_count"`
	RemovedCount  int       `db:"removed_count"`
}

type likeBatchItem struct {
	ImageID domain.ID `db:"image_id"`
	UserID  domain.ID `db:"user_id"`
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
