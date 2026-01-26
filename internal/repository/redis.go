package repository

import (
	"context"
	"encoding/json"

	"github.com/gomes800/bus-api-go/internal/model"
	"github.com/gomes800/bus-api-go/internal/util"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	C *redis.Client
}

func NewRedisRepo(c *redis.Client) *RedisRepo {
	return &RedisRepo{C: c}
}

func (r *RedisRepo) SaveBus(ctx context.Context, b model.BusPosition) error {
	key := "bus:" + b.Ordem

	old, _ := r.C.HGet(ctx, key, "timestamp").Result()
	if old == b.DataHora {
		return nil
	}

	raw, _ := json.Marshal(b)

	_, err := r.C.HSet(ctx, key, map[string]interface{}{
		"latitude":  b.Latitude,
		"longitude": b.Longitude,
		"linha":     b.Linha,
		"timestamp": b.DataHora,
		"raw":       string(raw),
	}).Result()
	if err != nil {
		return err
	}

	lon := util.Convert(b.Longitude)
	lat := util.Convert(b.Latitude)

	_, err = r.C.GeoAdd(ctx, "buses", &redis.GeoLocation{
		Name:      b.Ordem,
		Longitude: lon,
		Latitude:  lat,
	}).Result()
	if err != nil {
		return err
	}

	_, err = r.C.SAdd(ctx, "line:"+b.Linha, b.Ordem).Result()

	return err
}

func (r *RedisRepo) GetBus(ctx context.Context, ordem string) (*model.BusPosition, error) {
	raw, err := r.C.HGet(ctx, "bus:"+ordem, "raw").Result()
	if err != nil {
		return nil, err
	}
	var bus model.BusPosition
	json.Unmarshal([]byte(raw), &bus)
	return &bus, nil
}

func (r *RedisRepo) GetLine(ctx context.Context, line string) ([]model.BusPosition, error) {
	ids, err := r.C.SMembers(ctx, "line:"+line).Result()
	if err != nil {
		return nil, err
	}

	result := []model.BusPosition{}
	for _, id := range ids {
		bus, err := r.GetBus(ctx, id)
		if err == nil {
			result = append(result, *bus)
		}
	}
	return result, nil
}
