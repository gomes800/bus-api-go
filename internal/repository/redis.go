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

func (r *RedisRepo) CleanupLine(ctx context.Context, line string) error {
	lineKey := "line:" + line

	ids, err := r.C.SMembers(ctx, lineKey).Result()
	if err != nil {
		return err
	}

	pipe := r.C.Pipeline()
	existsCmds := make([]*redis.IntCmd, 0, len(ids))
	for _, id := range ids {
		existsCmds = append(existsCmds, pipe.Exists(ctx, "bus:"+id))
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	dead := make([]string, 0)
	for i, id := range ids {
		if existsCmds[i].Val() == 0 {
			dead = append(dead, id)
		}
	}

	if len(dead) > 0 {
		if err := r.C.SRem(ctx, lineKey, dead).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisRepo) CleanupGeo(ctx context.Context, max int) error {
	var cursor uint64
	checked := 0

	for checked < max {
		items, nextCursor, err := r.C.ZScan(ctx, "buses", cursor, "*", 500).Result()
		if err != nil {
			return err
		}

		members := make([]string, 0, len(items)/2)
		for i := 0; i < len(items); i += 2 {
			members = append(members, items[i])
		}

		if len(members) > 0 {
			pipe := r.C.Pipeline()
			existsCmds := make([]*redis.IntCmd, 0, len(members))
			for _, m := range members {
				existsCmds = append(existsCmds, pipe.Exists(ctx, "bus:"+m))
			}
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}

			dead := make([]string, 0)
			for i, m := range members {
				if existsCmds[i].Val() == 0 {
					dead = append(dead, m)
				}
			}

			if len(dead) > 0 {
				if err := r.C.ZRem(ctx, "buses", dead).Err(); err != nil {
					return err
				}
				log.Printf("[Cleanup] Geo removed %d dead buses", len(dead))
			}

			checked += len(members)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
