package repository

import (
	"context"
	"encoding/json"
	"log"
	"time"

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

	old, err := r.C.HGet(ctx, key, "timestamp").Result()
	if err != nil && err != redis.Nil {
		log.Printf("[Redis] HGet failed (%s): %v", key, err)
		return err
	}

	if old == b.DataHora {
		if err := r.C.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
			log.Printf("[Redis] Expire failed (%s): %v", key, err)
			return err
		}
		return nil
	}

	raw, err := json.Marshal(b)
	if err != nil {
		log.Printf("[Redis] JSON marshal failed (%s): %v", key, err)
		return err
	}

	if err := r.C.HSet(ctx, key, map[string]interface{}{
		"latitude":  b.Latitude,
		"longitude": b.Longitude,
		"linha":     b.Linha,
		"timestamp": b.DataHora,
		"raw":       string(raw),
	}).Err(); err != nil {
		log.Printf("[Redis] HSet failed (%s): %v", key, err)
		return err
	}

	if err := r.C.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
		log.Printf("[Redis] Expire failed (%s): %v", key, err)
		return err
	}

	lon := util.Convert(b.Longitude)
	lat := util.Convert(b.Latitude)

	if err := r.C.GeoAdd(ctx, "buses", &redis.GeoLocation{
		Name:      b.Ordem,
		Longitude: lon,
		Latitude:  lat,
	}).Err(); err != nil {
		log.Printf("[Redis] GeoAdd failed (%s): %v", b.Ordem, err)
		return err
	}

	if err := r.C.SAdd(ctx, "line:"+b.Linha, b.Ordem).Err(); err != nil {
		log.Printf("[Redis] SAdd failed (line:%s): %v", b.Linha, err)
		return err
	}

	return nil
}

func (r *RedisRepo) GetBus(ctx context.Context, ordem string) (*model.BusPosition, error) {
	key := "bus:" + ordem

	raw, err := r.C.HGet(ctx, key, "raw").Result()
	if err != nil {
		log.Printf("[Redis] GetBus failed (%s): %v", key, err)
		return nil, err
	}

	var bus model.BusPosition
	if err := json.Unmarshal([]byte(raw), &bus); err != nil {
		log.Printf("[Redis] Unmarshal failed (%s): %v", key, err)
		return nil, err
	}

	return &bus, nil
}

func (r *RedisRepo) GetLine(ctx context.Context, line string) ([]model.BusPosition, error) {
	key := "line:" + line

	ids, err := r.C.SMembers(ctx, key).Result()
	if err != nil {
		log.Printf("[Redis] SMembers failed (%s): %v", key, err)
		return nil, err
	}

	result := []model.BusPosition{}

	for _, id := range ids {
		bus, err := r.GetBus(ctx, id)
		if err != nil {
			log.Printf("[Redis] Missing bus %s in line %s", id, line)
			continue
		}
		result = append(result, *bus)
	}

	log.Printf("[Redis] GetLine %s -> %d buses", line, len(result))
	return result, nil
}
