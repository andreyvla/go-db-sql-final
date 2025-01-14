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

func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")

	require.NoError(t, err)

	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)

	assert.Equal(t, parcel, stored)

	err = store.Delete(parcel.Number)
	require.NoError(t, err)
	_, err = store.Get(parcel.Number)

	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	newAddress := "new test address"

	err = store.SetAddress(parcel.Number, newAddress)
	require.NoError(t, err)

	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)

	assert.Equal(t, newAddress, stored.Address)
}

func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()
	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	NewStatus := ParcelStatusSent
	store.SetStatus(parcel.Number, NewStatus)
	require.NoError(t, err)

	stored, err := store.Get(parcel.Number)
	require.NoError(t, err)

	assert.Equal(t, NewStatus, stored.Status)
}

func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

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
		id, err := store.Add(parcel)

		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(parcelMap))
	for _, parcel := range storedParcels {
		_, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, parcelMap[parcel.Number], parcel)
	}

}
