package test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"prisma/Database"
	"prisma/config"
)

var Config = config.NewViper()

func TestMongoConnect_NotPanics(t *testing.T) {

	assert.NotPanics(t, func() {
		db := database.MongoConnect(Config)
		assert.NotNil(t, db) 
	})
}

func TestPostgreConnect_NotPanics(t *testing.T) {

	assert.NotPanics(t, func() {
		db := database.PostgreConnect(Config)
		assert.NotNil(t, db) 
	})
}


