package postgres

import (
	"strconv"
	"time"

	"github.com/pillowskiy/gopix/pkg/batch"
)

var imageBatchConfig = batch.BatchConfig{Retries: 3, MaxSize: 10000}
var batchingCtxTimeout = time.Second * 5

var imageViewsBatchAgg = batch.NewMapAggregator[viewBatchItem]()

type imageAnalyticsAgg struct {
	ImageID int `db:"image_id"`
	Count   int `db:"count"`
}

type viewBatchItem struct {
	ImageID int  `db:"image_id"`
	UserID  *int `db:"user_id"`
}

func (i viewBatchItem) Group() string {
	return strconv.Itoa(i.ImageID)
}
