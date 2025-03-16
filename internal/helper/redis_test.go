package helper

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisClient = db

	key := "test-key"
	expectedValue := map[string]string{"foo": "bar"}
	expectedJSON, _ := json.Marshal(expectedValue)

	mock.ExpectGet(key).SetVal(string(expectedJSON))

	var actualValue map[string]string
	err := GetCache(key, &actualValue)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, actualValue)
	mock.ExpectationsWereMet()
}

func TestGetCache_NotFound(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisClient = db

	key := "non-existent"

	mock.ExpectGet(key).RedisNil()

	var actualValue map[string]string
	err := GetCache(key, &actualValue)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, redis.Nil))
	mock.ExpectationsWereMet()
}

func TestDeleteCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	redisClient = db

	key := "test-key"

	mock.ExpectDel(key).SetVal(1)

	err := DeleteCache(key)

	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
