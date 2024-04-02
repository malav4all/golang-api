package alert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/malav4all/golang-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func alertIDKey(id uint64) string {
	return fmt.Sprintf("alert:%d", id)
}

func (r *RedisRepo) InsertAlert(ctx context.Context, order model.Alert) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := alertIDKey(order.AlertID)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "alerts", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to alerts set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("alert does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Alert, error) {
	key := alertIDKey(id)
	fmt.Println(key)
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Alert{}, ErrNotExist
	} else if err != nil {
		return model.Alert{}, fmt.Errorf("get order: %w", err)
	}

	var alert model.Alert
	err = json.Unmarshal([]byte(value), &alert)
	if err != nil {
		return model.Alert{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	return alert, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := alertIDKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get order: %w", err)
	}

	if err := txn.SRem(ctx, "alerts", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from alerts set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Alert) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := alertIDKey(order.AlertID)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("set order: %w", err)
	}

	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Alerts []model.Alert
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "alerts", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Alerts: []model.Alert{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	response := make([]model.Alert, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Alert

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		response[i] = order
	}

	return FindResult{
		Alerts: response,
		Cursor: cursor,
	}, nil
}
