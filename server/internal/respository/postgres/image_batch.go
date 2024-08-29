package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	goErrors "errors"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/pkg/batch"
	"github.com/pkg/errors"
)

var imageBatchConfig = batch.BatchConfig{Retries: 3, MaxSize: 10000}
var batchingCtxTimeout = time.Second * 5

type imageAnalyticsAgg struct {
	ImageID int `db:"image_id"`
	Count   int `db:"count"`
}

func imageWithUserKey(imageID int, userID int) string {
	return fmt.Sprintf("%v:%v", imageID, userID)
}

func imageGroupKey(imageID int) string {
	return strconv.Itoa(imageID)
}

func (r *imageRepository) processViewsBatch(views []viewBatchItem) error {
	if len(views) == 0 {
		return nil
	}

	ctx, close := context.WithTimeout(context.Background(), batchingCtxTimeout)
	defer close()

	start := time.Now()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.BeginTx")
	}
	defer tx.Rollback()

	// This additionally loads the database, but this way we get a correct result of batch
	query := `
    WITH inserted AS (
      INSERT INTO images_to_views (image_id, user_id)
      VALUES(:image_id, :user_id)
      ON CONFLICT DO NOTHING
      RETURNING image_id
    )
    SELECT image_id, COUNT(*) AS count
    FROM inserted
    GROUP BY image_id;
  `

	query, params, err := tx.BindNamed(query, views)

	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.Named")
	}

	rows, err := tx.QueryxContext(ctx, query, params...)
	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.QueryxContext")
	}

	aggResult, err := scanToStructSliceOf[imageAnalyticsAgg](rows)
	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.SliceScan")
	}

	aggQuery := `
    UPDATE images_analytics AS ia
    SET views_count = ia.views_count + p.count
    FROM (VALUES (CAST(:image_id AS int), CAST(:count AS int))) AS p(image_id, count)
    WHERE ia.image_id = p.image_id
  `
	if _, err := tx.NamedExecContext(ctx, aggQuery, aggResult); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.AnalyticsExecContext")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.Commit")
	}

	log.Printf("Batching views took %s", time.Since(start))

	return nil
}

func (r *imageRepository) processLikesBatch(likes []likeBatchItem) error {
	ctx, close := context.WithTimeout(context.Background(), batchingCtxTimeout)
	defer close()

	var aggLikes, aggDislikes []likeBatchItem
	for _, l := range likes {
		if l.Liked {
			aggLikes = append(aggLikes, l)
		} else {
			aggDislikes = append(aggDislikes, l)
		}
	}

	writeLikes, writeErr := r.likesBulkWriteInTx(aggLikes)
	delLikes, delErr := r.likesBulkDeleteInTx(aggDislikes)
	if err := goErrors.Join(writeErr, delErr); err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.BeginTx")
	}
	defer tx.Rollback()

	writeLikesRowx, writeErr := writeLikes(ctx, tx)
	delLikesRowx, delErr := delLikes(ctx, tx)
	if err := goErrors.Join(writeErr, delErr); err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.InTxQueryCall")
	}

	writeLikesRes, writeErr := scanToStructSliceOf[imageAnalyticsAgg](writeLikesRowx)
	delLikesRes, delErr := scanToStructSliceOf[imageAnalyticsAgg](delLikesRowx)
	if err := goErrors.Join(writeErr, delErr); err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.SliceScan")
	}

	var likesAggRes []imageLikesAnalyticsAgg
	likesAggIndexes := make(map[int]int)

	for i, agg := range writeLikesRes {
		likesAggRes = append(likesAggRes, imageLikesAnalyticsAgg{
			ImageID:    agg.ImageID,
			LikesCount: agg.Count,
		})
		likesAggIndexes[agg.ImageID] = i
	}

	for _, agg := range delLikesRes {
		aggIndex, exists := likesAggIndexes[agg.ImageID]

		if exists {
			likesAggRes[aggIndex].DislikesCount = agg.Count
		} else {
			likesAggRes = append(likesAggRes, imageLikesAnalyticsAgg{
				ImageID:       agg.ImageID,
				DislikesCount: agg.Count,
			})
		}
	}

	aggQuery := `
    UPDATE images_analytics AS ia
    SET views_count = ia.views_count + p.count
    FROM (
      VALUES (CAST(:image_id AS int), CAST(:likes_count AS int), CAST(:dislikes_count AS int))
    ) AS p(image_id, count)
    WHERE ia.image_id = p.image_id
  `
	if _, err := tx.NamedExecContext(ctx, aggQuery, likesAggRes); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.AnalyticsExecContext")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.Commit")
	}

	return nil
}

func (r *imageRepository) likesBulkWriteInTx(likes []likeBatchItem) (InTxQueryCallContext, error) {
	query := `
    WITH inserted AS (
      INSERT INTO images_to_likes (image_id, user_id)
      VALUES(:image_id, :user_id)
      ON CONFLICT DO NOTHING
      RETURNING image_id
    )
    SELECT image_id, COUNT(*) AS count
    FROM inserted
    GROUP BY image_id;
  `

	query, params, err := r.db.BindNamed(query, likes)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.likesBulkWriteInTx.Named")
	}

	return func(ctx context.Context, tx *sqlx.Tx) (*sqlx.Rows, error) {
		return tx.QueryxContext(ctx, query, params...)
	}, nil
}

func (r *imageRepository) likesBulkDeleteInTx(likes []likeBatchItem) (InTxQueryCallContext, error) {
	query := `
    WITH deleted AS (
      DELETE FROM images_to_likes (image_id, user_id)
      WHERE image_id = :image_id AND user_id = :user_id
      RETURNING image_id
    )
    SELECT image_id, COUNT(*) AS count
    FROM inserted
    GROUP BY image_id;
  `

	query, params, err := r.db.BindNamed(query, likes)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.likesBulkDeleteInTx.Named")
	}

	return func(ctx context.Context, tx *sqlx.Tx) (*sqlx.Rows, error) {
		return tx.QueryxContext(ctx, query, params...)
	}, nil
}
