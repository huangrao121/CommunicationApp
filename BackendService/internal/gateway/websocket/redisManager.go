package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/huangrao121/CommunicationApp/BackendService/config"
	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	redisClusterClient *redis.ClusterClient
	nodeID             string
}

func NewRedisManager(cfg *config.Config, nodeID string) *RedisManager {
	redisOpts := &redis.ClusterOptions{
		Addrs:          cfg.Redis.ClusterAddrs,
		Password:       cfg.Redis.Password,
		PoolSize:       cfg.Redis.PoolSize,
		MinIdleConns:   cfg.Redis.MinIdleConns,
		MaxIdleConns:   cfg.Redis.MaxIdleConns,
		MaxActiveConns: cfg.Redis.MaxActiveConns,
		PoolTimeout:    cfg.Redis.IdleTimeout,
	}
	return &RedisManager{
		redisClusterClient: redis.NewClusterClient(redisOpts),
		nodeID:             nodeID,
	}
}

func (ulm *RedisManager) SetUserLocation(ctx context.Context, userID string, location string) error {
	return ulm.redisClusterClient.Set(context.Background(),
		fmt.Sprintf("user_location:%s", userID),
		ulm.nodeID,
		0).Err()
}

func (ulm *RedisManager) GetUserLocation(ctx context.Context, userID string) (string, error) {
	return ulm.redisClusterClient.Get(context.Background(),
		fmt.Sprintf("user_location:%s", userID)).Result()
}

func (ulm *RedisManager) UnregisterUser(ctx context.Context, userID string) error {
	return ulm.redisClusterClient.Del(ctx, fmt.Sprintf("user_location:%s", userID)).Err()
}

func (ulm *RedisManager) GetNodeID() string {
	return ulm.nodeID
}

func (ulm *RedisManager) SetGroupMemberByID(ctx context.Context, groupID string, members []uuid.UUID) error {
	jsonData, err := json.Marshal(members)
	if err != nil {
		log.Panic("Error marshalling group members: %w", err)
		return err
	}
	return ulm.redisClusterClient.Set(ctx,
		fmt.Sprintf("group_member_by_id:%s", groupID),
		jsonData,
		0).Err()
}

func (ulm *RedisManager) GetGroupMemberByID(ctx context.Context, groupID string) ([]uuid.UUID, error) {
	var members []uuid.UUID
	membersData, err := ulm.redisClusterClient.Get(ctx, fmt.Sprintf("group_member_by_id:%s", groupID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal group members: %w", err)
	}
	json.Unmarshal([]byte(membersData), &members)
	return members, nil
}
