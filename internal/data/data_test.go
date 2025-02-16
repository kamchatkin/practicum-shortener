package data

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/mocks"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	value := models.Alias{Alias: config.DefaultAlias, Source: config.DefaultSource}
	m.EXPECT().Get(context.Background(), config.DefaultAlias).Return(value, nil)

	var dbRef storage.Storage
	dbRef = m
	alias, err := Get(context.Background(), &dbRef, config.DefaultAlias)
	assert.NoError(t, err)
	assert.Equal(t, config.DefaultSource, alias.Source)
}

func TestSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	m.EXPECT().Set(context.Background(), config.DefaultAlias, config.DefaultSource).Return(nil)
	var dbRef storage.Storage
	dbRef = m

	err := Set(context.Background(), &dbRef, config.DefaultAlias, config.DefaultSource)
	assert.NoError(t, err)
}

func TestSetBatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	k1 := "asdfgh"
	v1 := "https://yandex.ru/"

	batch := map[string]string{
		k1:       v1,
		"zxcvbn": "https://yandex.ru2/",
	}

	var dbRef storage.Storage
	dbRef = m

	m.EXPECT().SetBatch(context.Background(), batch).Return(nil)
	err := SetBatch(context.Background(), &dbRef, batch)
	assert.NoError(t, err)

	alias := models.Alias{Alias: k1, Source: v1}

	m.EXPECT().Get(context.Background(), k1).Return(alias, nil)
	alias, err = Get(context.Background(), &dbRef, k1)

	assert.NoError(t, err)
	assert.True(t, alias.Found())
}
