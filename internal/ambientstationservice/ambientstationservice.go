package ambientstationservice

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

type AmbientStationService struct {
	Ctx           context.Context
	Wg            *sync.WaitGroup
	Logger        *slog.Logger
	PgPool        *pgxpool.Pool
	ApiUrl        string
	DeviceListMap map[string]devices.Devices
}

func NewAmbientStationService(
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	pgpool *pgxpool.Pool,
	apiUrl string,
	deviceListMap map[string]devices.Devices,
) *AmbientStationService {
	return &AmbientStationService{
		Ctx:           ctx,
		Wg:            wg,
		Logger:        logger,
		PgPool:        pgpool,
		ApiUrl:        apiUrl,
		DeviceListMap: deviceListMap,
	}
}

func (a *AmbientStationService) Run() {
	for {
		select {
		case <-a.Ctx.Done():
			a.Logger.Info("Ambient station service context done signal detected")
			a.Wg.Done()
			return
		default:
			a.Logger.Info("Getting ambient station data...")
			resp, err := a.GetStationData()
			if err != nil {
				a.Logger.Error("Error getting ambient station data", "error", err)
				continue
			}

			a.Logger.Debug("Ambient station response", "response", resp)

			a.ProcessAmbientStationResponse(resp)

			waitDuration := 1 * time.Minute
			timer := time.NewTimer(waitDuration)

			select {
			case <-a.Ctx.Done():
				timer.Stop()
				a.Logger.Info("Ambient station service context done signal detected while waiting")
				a.Wg.Done()
				return
			case <-timer.C:
			}
		}
	}
}

func (a *AmbientStationService) GetStationData() (devices.AmbientStationResponse, error) {
	resp, err := http.Get(a.ApiUrl)
	if err != nil {
		return nil, err
	}
	// log.Println("Ambient Weather Station API call successful")

	defer resp.Body.Close()

	var response devices.AmbientStationResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *AmbientStationService) ProcessAmbientStationResponse(resp devices.AmbientStationResponse) {
	a.Logger.Debug("Processing ambient station response")

	timestamp := utils.GetUTCTimestamp()
	dateUTC := int32(resp[0].LastData.DateUTC)
	insideTempF := float32(resp[0].LastData.InsideTempF)
	insideTempC := utils.FahrenheitToCelsius(insideTempF)
	insideFeelsLikeTempF := float32(resp[0].LastData.FeelsLikeInside)
	outsideTempF := float32(resp[0].LastData.OutsideTempF)
	outsideFeelsLikeTempF := float32(resp[0].LastData.FeelsLikeOutside)
	insideHumidity := float32(resp[0].LastData.InsideHumidity)
	outsideHumidity := float32(resp[0].LastData.OutsideHumidity)
	insideDewPoint := float32(resp[0].LastData.DewPointInside)
	outsideDewPoint := float32(resp[0].LastData.DewPointOutside)
	baroRelative := float32(resp[0].LastData.BarometricPressureRelIn)
	baroAbsolute := float32(resp[0].LastData.BarometricPressureAbsIn)
	windDirection := float32(resp[0].LastData.WindDirection)
	windSpeedMph := float32(resp[0].LastData.WindSpeedMPH)
	windSpeedGustMph := float32(resp[0].LastData.WindGustMPH)
	maxDailyGust := float32(resp[0].LastData.MaxDailyGust)
	hourlyRainInches := float32(resp[0].LastData.HourlyRainIn)
	eventRainInches := float32(resp[0].LastData.EventRainIn)
	dailyRainInches := float32(resp[0].LastData.DailyRainIn)
	weeklyRainInches := float32(resp[0].LastData.WeeklyRainIn)
	monthlyRainInches := float32(resp[0].LastData.MonthlyRainIn)
	totalRainInches := float32(resp[0].LastData.TotalRainIn)
	uvIndex := float32(resp[0].LastData.UVIndex)
	solarRadiation := float32(resp[0].LastData.SolarRadiation)
	outsideBattStatus := int32(resp[0].LastData.OutsideBattStatus)
	co2BattStatus := int32(resp[0].LastData.BattCO2)
	deviceID := a.DeviceListMap["ambient_wx_station"].DeviceID
	deviceTypeID := a.DeviceListMap["ambient_wx_station"].DeviceTypeID

	ambientStationUniqueParams := sqlc.InsertAmbientStationDataParams{
		Timestamp:             timestamp,
		Date:                  &resp[0].LastData.Date,
		Timezone:              &resp[0].LastData.TZ,
		DateUtc:               &dateUTC,
		InsideTempFeelsLikeF:  &insideFeelsLikeTempF,
		OutsideTempF:          &outsideTempF,
		OutsideTempFeelsLikeF: &outsideFeelsLikeTempF,
		OutsideHumidity:       &outsideHumidity,
		InsideDewPoint:        &insideDewPoint,
		OutsideDewPoint:       &outsideDewPoint,
		RelativePressure:      &baroRelative,
		WindDirection:         &windDirection,
		WindSpeedMph:          &windSpeedMph,
		WindSpeedGustMph:      &windSpeedGustMph,
		MaxDailyGustMph:       &maxDailyGust,
		HourlyRainInches:      &hourlyRainInches,
		EventRainInches:       &eventRainInches,
		DailyRainInches:       &dailyRainInches,
		WeeklyRainInches:      &weeklyRainInches,
		MonthlyRainInches:     &monthlyRainInches,
		TotalRainInches:       &totalRainInches,
		LastRain:              &resp[0].LastData.LastRain,
		UvIndex:               &uvIndex,
		SolarRadiation:        &solarRadiation,
		OutsideBattStatus:     &outsideBattStatus,
		Co2BattStatus:         &co2BattStatus,
		DeviceID:              &deviceID,
	}

	sharedParams := sqlc.InsertReadingParams{
		Timestamp:        timestamp,
		TempF:            &insideTempF,
		TempC:            &insideTempC,
		Humidity:         &insideHumidity,
		AbsolutePressure: &baroAbsolute,
		DeviceTypeID:     &deviceTypeID,
		DeviceID:         &deviceID,
	}

	db := sensordb.NewDbWrapper(a.PgPool)
	defer db.Conn.Release()

	err := db.DB.InsertAmbientStationData(a.Ctx, ambientStationUniqueParams)
	if err != nil {
		a.Logger.Error("Error writing ambient unique data", "error", err)
	} else {
		a.Logger.Info("Successful ambient unique data write")
	}

	err = db.DB.InsertReading(a.Ctx, sharedParams)
	if err != nil {
		a.Logger.Error("Error writing ambient shared data", "error", err)
	} else {
		a.Logger.Info("Successful ambient shared data write")
	}
}
