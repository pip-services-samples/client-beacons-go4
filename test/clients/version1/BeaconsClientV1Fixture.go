package test_clients1

import (
	"context"
	"testing"

	clients1 "github.com/pip-services-samples/client-beacons-go/clients/version1"
	data1 "github.com/pip-services-samples/service-beacons-go/data/version1"
	cquery "github.com/pip-services4/pip-services4-go/pip-services4-data-go/query"
	"github.com/stretchr/testify/assert"
)

type BeaconsClientV1Fixture struct {
	BEACON1 *data1.BeaconV1
	BEACON2 *data1.BeaconV1
	client  clients1.IBeaconsClientV1
	ctx     context.Context
}

func NewBeaconsClientV1Fixture(client clients1.IBeaconsClientV1) *BeaconsClientV1Fixture {
	c := &BeaconsClientV1Fixture{}

	c.BEACON1 = &data1.BeaconV1{
		Id:     "1",
		Udi:    "00001",
		Type:   data1.AltBeacon,
		SiteId: "1",
		Label:  "TestBeacon1",
		Center: data1.GeoPointV1{Type: "Point", Coordinates: []float32{0.0, 0.0}},
		Radius: 50,
	}

	c.BEACON2 = &data1.BeaconV1{
		Id:     "2",
		Udi:    "00002",
		Type:   data1.IBeacon,
		SiteId: "1",
		Label:  "TestBeacon2",
		Center: data1.GeoPointV1{Type: "Point", Coordinates: []float32{2.0, 2.0}},
		Radius: 70,
	}

	c.client = client
	c.ctx = context.Background()

	return c
}

func (c *BeaconsClientV1Fixture) testCreateBeacons(t *testing.T) {
	// Create the first beacon
	beacon, err := c.client.CreateBeacon(c.ctx, *c.BEACON1)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, c.BEACON1.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON1.SiteId, beacon.SiteId)
	assert.Equal(t, c.BEACON1.Type, beacon.Type)
	assert.Equal(t, c.BEACON1.Label, beacon.Label)
	assert.NotNil(t, beacon.Center)

	// Create the second beacon
	beacon, err = c.client.CreateBeacon(c.ctx, *c.BEACON2)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, c.BEACON2.Udi, beacon.Udi)
	assert.Equal(t, c.BEACON2.SiteId, beacon.SiteId)
	assert.Equal(t, c.BEACON2.Type, beacon.Type)
	assert.Equal(t, c.BEACON2.Label, beacon.Label)
	assert.NotNil(t, beacon.Center)
}

func (c *BeaconsClientV1Fixture) TestCrudOperations(t *testing.T) {
	var beacon1 *data1.BeaconV1

	// Create items
	c.testCreateBeacons(t)

	// Get all beacons
	page, err := c.client.GetBeacons(c.ctx, *cquery.NewEmptyFilterParams(), *cquery.NewEmptyPagingParams())
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.Data, 2)
	beacon1 = &page.Data[0]

	// Update the beacon
	beacon1.Label = "ABC"
	beacon, err := c.client.UpdateBeacon(c.ctx, *beacon1)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)
	assert.Equal(t, "ABC", beacon.Label)

	// Get beacon by udi
	beacon, err = c.client.GetBeaconByUdi(c.ctx, beacon1.Udi)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Delete the beacon
	beacon, err = c.client.DeleteBeaconById(c.ctx, beacon1.Id)
	assert.Nil(t, err)
	assert.NotNil(t, beacon)
	assert.Equal(t, beacon1.Id, beacon.Id)

	// Try to get deleted beacon
	beacon, err = c.client.GetBeaconById(c.ctx, beacon1.Id)
	assert.Nil(t, err)
	assert.Empty(t, beacon)
}

func (c *BeaconsClientV1Fixture) TestCalculatePosition(t *testing.T) {
	// Create items
	c.testCreateBeacons(t)

	// Calculate position for one beacon
	position, err := c.client.CalculatePosition(c.ctx, "1", []string{"00001"})
	assert.Nil(t, err)
	assert.NotNil(t, position)
	assert.Equal(t, "Point", position.Type)
	assert.Equal(t, (float32)(0), position.Coordinates[0])
	assert.Equal(t, (float32)(0), position.Coordinates[1])
}
