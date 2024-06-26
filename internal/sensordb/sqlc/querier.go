// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"context"
)

type Querier interface {
	GetDeviceIdByName(ctx context.Context, deviceName string) (int32, error)
	GetDeviceTypeIDByDeviceType(ctx context.Context, deviceType string) (int32, error)
	GetDeviceTypeIdByName(ctx context.Context, deviceName string) (*int32, error)
	GetDevices(ctx context.Context) ([]*GetDevicesRow, error)
	GetMqttTopicData(ctx context.Context) ([]*GetMqttTopicDataRow, error)
	GetUniqueMqttTopics(ctx context.Context) ([]*string, error)
	InsertAmbientStationData(ctx context.Context, arg InsertAmbientStationDataParams) error
	InsertAqaraUniqueData(ctx context.Context, arg InsertAqaraUniqueDataParams) error
	InsertReading(ctx context.Context, arg InsertReadingParams) error
	InsertSonoffSmartPlugReading(ctx context.Context, arg InsertSonoffSmartPlugReadingParams) error
}

var _ Querier = (*Queries)(nil)
