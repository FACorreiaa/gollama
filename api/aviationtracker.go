package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/FACorreiaa/go-ollama/api/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func handleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05-07:00")
}

/*Airline Migration function */

func MigrateAirlineAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airline").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAirlineData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirlineData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airlines")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.AirlineApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	// Insert data from the JSON
	if _, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"airline"},
		[]string{"fleet_average_age", "airline_id", "callsign", "hub_code", "iata_code", "icao_code", "country_iso2",
			"date_founded", "iata_prefix_accounting", "airline_name", "country_name", "fleet_size", "status", "type",
			"created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].FleetAverageAge,
				res.Data[i].AirlineId,
				res.Data[i].Callsign,
				res.Data[i].HubCode,
				res.Data[i].IataCode,
				res.Data[i].IcaoCode,
				res.Data[i].CountryIso2,
				res.Data[i].DateFounded,
				res.Data[i].IataPrefixAccounting,
				res.Data[i].AirlineName,
				res.Data[i].CountryName,
				res.Data[i].FleetSize,
				res.Data[i].Status,
				res.Data[i].Type,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airline table")
		return err
	}

	slog.Info("Data inserted into the airline table")
	return nil
}

/*Aircraft Migration function */

func MigrateAircraftAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM aircraft").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAircraftData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAircraftData(conn *pgxpool.Pool) error {
	res := new(structs.AircraftApiData)
	data, err := fetchAviationStackData("aircraft_types")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"aircraft"},
		[]string{"iata_code", "aircraft_name", "plane_type_id", "created_at"},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].IataCode,
				res.Data[i].AircraftName,
				res.Data[i].PlaneTypeId,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into aircraft table")
		return err
	}

	slog.Info("Data inserted into the aircraft table")

	return nil
}

/*Tax Migration function */

func MigrateTaxAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM tax").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertTaxData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertTaxData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("taxes")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.TaxApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"tax"},
		[]string{"tax_id", "tax_name", "iata_code", "created_at"},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].TaxId, res.Data[i].TaxName, res.Data[i].IataCode,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into tax table")
		return err
	}

	slog.Info("Data inserted into the aircraft table")
	return nil
}

/* Airplane */

func MigrateAirplaneAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airplane").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airline table, fetch from the external API
		if err := fetchDataAndInsertAirplaneData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirplaneData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("airplanes")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	res := new(structs.AirplaneApiData)
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airplane"},
		[]string{"iata_type", "airplane_id", "airline_iata_code", "iata_code_long", "iata_code_short",
			"airline_icao_code", "construction_number", "delivery_date", "engines_count", "engines_type",
			"first_flight_date", "icao_code_hex", "line_number", "model_code", "registration_number",
			"test_registration_number", "plane_age", "plane_class", "model_name", "plane_owner", "plane_series",
			"plane_status", "production_line", "registration_date", "rollout_date", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].IataType,
				res.Data[i].AirplaneId,
				res.Data[i].AirlineIataCode,
				res.Data[i].IataCodeLong,
				res.Data[i].IataCodeShort,
				res.Data[i].AirlineIcaoCode,
				res.Data[i].ConstructionNumber,
				res.Data[i].DeliveryDate.Time,
				res.Data[i].EnginesCount,
				res.Data[i].EnginesType,
				res.Data[i].FirstFlightDate.Time,
				res.Data[i].IcaoCodeHex,
				res.Data[i].LineNumber,
				res.Data[i].ModelCode,
				res.Data[i].RegistrationNumber,
				res.Data[i].TestRegistrationNumber,
				res.Data[i].PlaneAge,
				res.Data[i].PlaneClass,
				res.Data[i].ModelName,
				res.Data[i].PlaneOwner,
				res.Data[i].PlaneSeries,
				res.Data[i].PlaneStatus,
				res.Data[i].ProductionLine,
				res.Data[i].RegistrationDate.Time,
				res.Data[i].RolloutDate.Time,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airplane table")
		return err
	}

	slog.Info("Data inserted into the airplane table")

	return nil
}

/* Airports */

func MigrateAirportAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM airport").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := fetchDataAndInsertAirportData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertAirportData(conn *pgxpool.Pool) error {
	res := new(structs.AirportApiData)
	data, err := fetchAviationStackData("airports")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"airport"},
		[]string{"gmt", "airport_id", "iata_code", "city_iata_code", "icao_code",
			"country_iso2", "geoname_id", "latitude", "longitude", "airport_name",
			"country_name", "phone_number", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].GMT, res.Data[i].AirportId, res.Data[i].IataCode,
				res.Data[i].CityIataCode, res.Data[i].IcaoCode, res.Data[i].CountryIso2,
				res.Data[i].GeonameId, res.Data[i].Latitude, res.Data[i].Longitude,
				res.Data[i].AirportName, res.Data[i].CountryName, res.Data[i].PhoneNumber,
				res.Data[i].Timezone, formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into airports table")
		return err
	}

	slog.Info("Data inserted into the airport table")
	return nil
}

/* Countries */

func MigrateCountryAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM country").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the country table, fetch from the external API
		if err := fetchDataAndInsertCountryData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertCountryData(conn *pgxpool.Pool) error {
	res := new(structs.CountryApiData)
	data, err := fetchAviationStackData("countries")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"country"},
		[]string{"country_name", "country_iso2", "country_iso3", "country_iso_numeric", "population",
			"capital", "continent", "currency_name", "currency_code", "fips_code",
			"phone_prefix", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].CountryName,
				res.Data[i].CountryIso2,
				res.Data[i].CountryIso3,
				res.Data[i].CountryIsoNumeric,
				res.Data[i].Population,
				res.Data[i].Capital,
				res.Data[i].Continent,
				res.Data[i].CurrencyName,
				res.Data[i].CurrencyCode,
				res.Data[i].FipsCode,
				res.Data[i].PhonePrefix,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into country table")
		return err
	}

	slog.Info("Data inserted into the country table")

	return nil
}

/* Cities */

func MigrateCityAPIData(conn *pgxpool.Pool) error {
	slog.Info("Running API check")
	ctx := context.Background()
	slog.Info("checking for data on the DB")

	var count int
	if err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM city").Scan(&count); err != nil {
		handleError(err, "Error querying the table")
		return err
	}

	if count == 0 {
		// No data in the airport table, fetch from the external API
		if err := fetchDataAndInsertCityData(conn); err != nil {
			handleError(err, "Error inserting data")
			return err
		}
	}
	slog.Info("Migrations finished")
	return nil
}

func fetchDataAndInsertCityData(conn *pgxpool.Pool) error {
	data, err := fetchAviationStackData("cities")
	if err != nil {
		handleError(err, "error fetching data")
		return err
	}

	res := new(structs.CityApiData)

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&res); err != nil {
		handleError(err, "error unmarshaling API response")
		return err
	}

	//Insert data from the json
	if _, err := conn.CopyFrom(

		context.Background(),
		pgx.Identifier{"city"},
		[]string{"gmt", "city_id", "iata_code", "country_iso2", "geoname_id",
			"latitude", "longitude", "city_name", "timezone", "created_at",
		},
		pgx.CopyFromSlice(len(res.Data), func(i int) ([]interface{}, error) {
			return []interface{}{
				res.Data[i].GMT,
				res.Data[i].CityId,
				res.Data[i].IataCode,
				res.Data[i].CountryIso2,
				res.Data[i].GeonameId,
				res.Data[i].Latitude,
				res.Data[i].Longitude,
				res.Data[i].CityName,
				res.Data[i].Timezone,
				formatTime(time.Now()),
			}, nil
		}),
	); err != nil {
		handleError(err, "error inserting data into cities table")
		return err
	}

	slog.Info("Data inserted into the city table")

	return nil
}

/* Fetch data from endpoint */

// REFACTOR FUNCTION TO DELETE fetchData...

func fetchAviationStackData(endpoint string, queryParams ...string) ([]byte, error) {
	accessKey := os.Getenv("AVIATION_STACK_API_KEY")
	if accessKey == "" {
		return nil, fmt.Errorf("missing API access key")
	}

	baseURL := "http://api.aviationstack.com/v1/"

	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	// Set the endpoint path
	parsedURL.Path += endpoint

	// Create a new query parameters object
	query := parsedURL.Query()

	// Add the access key parameter
	query.Set("access_key", accessKey)

	// Add additional query parameters
	if len(queryParams) > 0 {
		for _, param := range queryParams {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				query.Set(parts[0], parts[1])
			}
		}
	}

	parsedURL.RawQuery = query.Encode()

	finalURL := parsedURL.String()

	response, err := http.Get(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("something is not ok")
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	defer response.Body.Close()

	return body, nil
}
