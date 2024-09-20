package postgres

import (
	"fmt"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/pkg/batch"

	nanoid "github.com/matoous/go-nanoid/v2"
)

var imageViewsBatchAgg = batch.NewKGAggregator[viewBatchItem]()

type viewBatchItem struct {
	ImageID domain.ID  `db:"image_id"`
	UserID  *domain.ID `db:"user_id"`
}

func (i viewBatchItem) Group() string {
	return imageGroupKey(i.ImageID)
}

func (i viewBatchItem) Key() string {
	if i.UserID == nil {
		return fmt.Sprintf("%v:%s", i.ImageID, nanoid.Must(8))
	}
	return imageWithUserKey(i.ImageID, *i.UserID)
}

func (i viewBatchItem) Count() int {
	return 1
}
