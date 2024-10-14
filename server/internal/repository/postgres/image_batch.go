package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/repository/postgres/pgutils"
	"github.com/pillowskiy/gopix/pkg/batch"
	"github.com/pkg/errors"
)

var imageBatchConfig = batch.BatchConfig{Retries: 3, MaxSize: 1000}
var imageBatchTickDuration = time.Minute
var batchingCtxTimeout = time.Second * 5

type imageAnalyticsAgg struct {
	ImageID int `db:"image_id"`
	Count   int `db:"count"`
}

func imageWithUserKey(imageID domain.ID, userID domain.ID) string {
	return fmt.Sprintf("%v:%v", imageID, userID)
}

func imageGroupKey(imageID domain.ID) string {
	return imageID.String()
}

func (r *imageRepository) processViewsBatch(views []viewBatchItem) error {
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

	aggValues, aggParams, err := pgutils.BulkUpdateValues(aggResult, "image_id::bigint,count::int")
	if err != nil {
		return errors.Wrap(err, "imageRepository.batchViews.BulkUpdateValues")
	}

	aggQuery := fmt.Sprintf(`
  UPDATE images_analytics AS ia
  SET views_count = ia.views_count + p.count
  FROM (%s) AS p(image_id, count)
  WHERE ia.image_id = p.image_id;
  `, aggValues)

	if _, err := tx.ExecContext(ctx, aggQuery, aggParams...); err != nil {
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.BeginTx")
	}
	defer tx.Rollback()

	var insLikesBunch, delLikesBunch []likeBatchItem
	for _, l := range likes {
		if l.Liked {
			insLikesBunch = append(insLikesBunch, l)
		} else {
			delLikesBunch = append(delLikesBunch, l)
		}
	}

	likesAnalytics := make([]imageLikesAnalytics, 0, len(likes))

	insLikesAnalytics, err := r.queryBulkWriteLikes(ctx, tx, insLikesBunch)
	if err != nil {
		log.Printf("Error in queryBulkWriteLikes: %v", err)
	} else {
		likesAnalytics = append(likesAnalytics, insLikesAnalytics...)
	}

	delLikesAnalytics, err := r.queryBulkDeleteLikes(ctx, tx, delLikesBunch)
	if err != nil {
		log.Printf("Error in queryBulkWriteDislikes: %v", err)
	} else {
		likesAnalytics = append(likesAnalytics, delLikesAnalytics...)
	}

	aggValues, aggParams, err := pgutils.BulkUpdateValues(
		likesAnalytics,
		"image_id::bigint,inserted_count::int,removed_count::int",
	)
	if err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.BulkUpdateValues")
	}

	aggQuery := fmt.Sprintf(`
  UPDATE images_analytics AS ia
  SET likes_count = ia.likes_count + p.inserted_count - p.removed_count
  FROM (%s) AS p(image_id, inserted_count, removed_count)
  WHERE ia.image_id = p.image_id;
  `, aggValues)

	if _, err := tx.ExecContext(ctx, aggQuery, aggParams...); err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.AnalyticsExecContext")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "imageRepository.processLikesBatch.Commit")
	}

	return nil
}

func (r *imageRepository) queryBulkWriteLikes(
	ctx context.Context,
	tx *sqlx.Tx,
	likes []likeBatchItem,
) ([]imageLikesAnalytics, error) {
	if len(likes) == 0 {
		return make([]imageLikesAnalytics, 0), nil
	}

	call, err := r.likesBulkWriteInTx(likes)
	if err != nil {
		return nil, err
	}

	return r.queryBulkLikesCall(ctx, tx, call)
}

func (r *imageRepository) queryBulkDeleteLikes(
	ctx context.Context,
	tx *sqlx.Tx,
	likes []likeBatchItem,
) ([]imageLikesAnalytics, error) {
	if len(likes) == 0 {
		return make([]imageLikesAnalytics, 0), nil
	}

	call, err := r.likesBulkDeleteInTx(likes)
	if err != nil {
		return nil, err
	}

	return r.queryBulkLikesCall(ctx, tx, call)
}

func (r *imageRepository) queryBulkLikesCall(
	ctx context.Context,
	tx *sqlx.Tx,
	call InTxQueryCallContext,
) ([]imageLikesAnalytics, error) {
	rowx, err := call(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.processLikesBatch.InTxQueryCall")
	}

	agg, err := scanToStructSliceOf[imageLikesAnalytics](rowx)
	if err != nil {
		return nil, errors.Wrap(err, "imageRepository.processLikesBatch.SliceScan")
	}

	return agg, nil
}

func (r *imageRepository) likesBulkWriteInTx(likes []likeBatchItem) (InTxQueryCallContext, error) {
	query := `
    WITH inserted AS (
      INSERT INTO images_to_likes (image_id, user_id)
      VALUES(:image_id, :user_id)
      ON CONFLICT DO NOTHING
      RETURNING image_id
    )
    SELECT image_id, COUNT(*) AS inserted_count, 0 AS removed_count 
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
    SELECT image_id, 0 AS inserted_count, COUNT(*) AS removed_count
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
