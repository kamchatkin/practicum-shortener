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
	alias, err := Get(context.Background(), getM(dbRef, m), config.DefaultAlias)
	assert.NoError(t, err)
	assert.Equal(t, config.DefaultSource, alias.Source)
}

func getM(dbRef storage.Storage, m *mocks.MockStorage) *storage.Storage {
	dbRef = m

	return &dbRef
}

func TestSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockStorage(ctrl)

	m.EXPECT().Set(context.Background(), config.DefaultAlias, config.DefaultSource, int64(1)).Return(nil)

	var dbRef storage.Storage
	err := Set(context.Background(), getM(dbRef, m), config.DefaultAlias, config.DefaultSource, int64(1))
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
	m.EXPECT().SetBatch(context.Background(), batch, int64(1)).Return(nil)
	err := SetBatch(context.Background(), getM(dbRef, m), batch, int64(1))
	assert.NoError(t, err)

	alias := models.Alias{Alias: k1, Source: v1}

	m.EXPECT().Get(context.Background(), k1).Return(alias, nil)
	alias, err = Get(context.Background(), getM(dbRef, m), k1)

	assert.NoError(t, err)
	assert.True(t, alias.Found())
}
