package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	parcel.Number = p
	require.NoError(t, err)
	require.NotEmpty(t, p)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, parcel, gp)

	err = store.Delete(p)
	require.NoError(t, err)

	_, err = store.Get(p)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	newAddress := "new test address"
	err = store.SetAddress(p, newAddress)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, newAddress, gp.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	err = store.SetStatus(p, ParcelStatusSent)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, gp.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		p, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, p)
		parcels[i].Number = p
		parcelMap[p] = parcels[i]
	}

	var storedParcels []Parcel
	storedParcels, err = store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
