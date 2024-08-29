package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pillowskiy/gopix/pkg/batch"

	nanoid "github.com/matoous/go-nanoid/v2"
)

var imageBatchConfig = batch.BatchConfig{Retries: 3, MaxSize: 10000}
var batchingCtxTimeout = time.Second * 5

func imageWithUserKey(imageID int, userID int) string {
	return fmt.Sprintf("%v:%v", imageID, userID)
}

func imageGroupKey(imageID int) string {
	return strconv.Itoa(imageID)
}

var imageViewsBatchAgg = batch.NewKGAggregator[viewBatchItem]()

type imageAnalyticsAgg struct {
	ImageID int `db:"image_id"`
	Count   int `db:"count"`
}

type viewBatchItem struct {
	ImageID int  `db:"image_id"`
	UserID  *int `db:"user_id"`
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
