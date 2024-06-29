package avtechservice

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sensor-data-collection-service/internal/devices"
	"sensor-data-collection-service/internal/sensordb"
	"sensor-data-collection-service/internal/sensordb/sqlc"
	"sensor-data-collection-service/internal/utils"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AvtechService struct {
	Ctx           context.Context
	Wg            *sync.WaitGroup
	Logger        *slog.Logger
	PgPool        *pgxpool.Pool
	ApiUrl        string
	DeviceListMap map[string]sqlc.GetDevicesRow
}

func NewAvtechService(
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	pgpool *pgxpool.Pool,
	apiUrl string,
	deviceListMap map[string]sqlc.GetDevicesRow,
) *AvtechService {
	return &AvtechService{
		Ctx:           ctx,
		Wg:            wg,
		Logger:        logger,
		PgPool:        pgpool,
		ApiUrl:        apiUrl,
		DeviceListMap: deviceListMap,
	}
}

func (a *AvtechService) Run() {
	for {
		select {
		case <-a.Ctx.Done():
			a.Logger.Info("avtech station service context done signal detected")
			a.Wg.Done()
			return
		default:
			a.Logger.Info("Getting avtech station data...")
			resp, err := a.GetAvtechData()
			if err != nil {
				a.Logger.Error("Error getting avtech station data", "error", err)
				continue
			}

			a.Logger.Debug("avtech station response", "response", resp)

			a.ProcessAvtechResponse(resp)

			waitDuration := 1 * time.Minute
			timer := time.NewTimer(waitDuration)

			select {
			case <-a.Ctx.Done():
				timer.Stop()
				a.Logger.Info("avtech station service context done signal detected while waiting")
				a.Wg.Done()
				return
			case <-timer.C:
			}
		}
	}
}

func (a *AvtechService) GetAvtechData() (devices.AvtechResponse, error) {
	var response devices.AvtechResponse

	resp, err := http.Get(a.ApiUrl)
	if err != nil {
		return devices.AvtechResponse{}, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return devices.AvtechResponse{}, err
	}

	return response, nil
}

func (a *AvtechService) ProcessAvtechResponse(resp devices.AvtechResponse) {
	a.Logger.Debug("Processing avtech sensor response")

	tempF := utils.ConvertStrToFloat32(resp.Sensor[0].TempF)
	tempC := utils.ConvertStrToFloat32(resp.Sensor[0].TempC)
	timestamp := utils.GetUTCTimestamp()
	humidity := float32(0)
	pressure := float32(0)
	deviceID := a.DeviceListMap["avtech_basement_rack"].DeviceID
	deviceTypeID := a.DeviceListMap["avtech_basement_rack"].DeviceTypeID

	avtechDataParams := sqlc.InsertReadingParams{
		Timestamp:        timestamp,
		TempF:            &tempF,
		TempC:            &tempC,
		Humidity:         &humidity,
		AbsolutePressure: &pressure,
		DeviceTypeID:     deviceTypeID,
		DeviceID:         &deviceID,
	}

	conn := utils.AcquirePoolConn(a.PgPool)
	defer conn.Release()

	db := sensordb.NewSensorDataDB(conn)
	err := db.InsertReading(a.Ctx, avtechDataParams)
	if err != nil {
		a.Logger.Error("Error writing avtech shared data", "error", err)
	} else {
		a.Logger.Info("Successful avtech shared data write")
	}
}
