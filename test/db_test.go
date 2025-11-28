package test

import (
	"prisma/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Config = config.NewViper()
var Log = config.Log(Config)

func TestMongoConnect_NotPanics(t *testing.T) {

	assert.NotPanics(t, func() {
		db := config.MongoConnect(Config, Log)
		assert.NotNil(t, db)
	})
}

func TestPostgreConnect_NotPanics(t *testing.T) {

	assert.NotPanics(t, func() {
		db := config.PostgresConnect(Config, Log)
		assert.NotNil(t, db)
	})
}

func TestRedisConnect_NotPanics(t *testing.T) {
	assert.NotPanics(t, func() {
		db := config.NewRedisClient(Config, Log)
		assert.NotNil(t, db)
	})
}
