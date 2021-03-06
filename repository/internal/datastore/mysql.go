package datastore

import (
	"context"
	"database/sql"
	"go-api-samp/application"
	"go-api-samp/model/entity"
	"go-api-samp/model/errors"
	"go-api-samp/util/log"
)

/*
 * Repository(datastore)の実装
 */

type MySQLClient struct {
	Db *sql.DB
}

func (c *MySQLClient) AddWeather(ctx context.Context, locationId, weather int, date, comment string) error {
	logger := log.GetLogger()

	sql := "INSERT INTO WEATHER(dat, weather, location_id, comment) VALUES(?, ?, ?, ?)"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(date, weather, locationId, comment)
	if err != nil {
		logger.Error(ctx, "failed to execute add weather query.", err)
		return errors.DataStoreSystemError(err)
	}

	return nil
}

func (c *MySQLClient) UpdateWeather(ctx context.Context, locationId, weather int, date, comment string) error {
	logger := log.GetLogger()

	sql := "UPDATE WEATHER SET weather = ?, comment = ? WHERE location_id = ? AND dat = ?"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	_, err = stmt.Query(weather, comment, locationId, date)
	if err != nil {
		logger.Error(ctx, "failed to execute update weather query.", err)
		return errors.DataStoreSystemError(err)
	}

	return nil
}

func (c *MySQLClient) GetWeather(ctx context.Context, locationId int, date string) (*entity.Weather, error) {
	logger := log.GetLogger()

	sql := "SELECT dat, weather, location_id, city, comment FROM WEATHER as w INNER JOIN LOCATION as l ON w.location_id = l.id WHERE l.id = ? AND dat = ?"
	stmt, err := c.Db.Prepare(sql)
	if err != nil {
		logger.Error(ctx, "failed to prepare statement.", err)
		return nil, errors.DataStoreSystemError(err)
	}
	defer stmt.Close()

	row, err := stmt.Query(locationId, date)
	if err != nil {
		logger.Error(ctx, "failed to execute get location query.", err)
		return nil, errors.DataStoreSystemError(err)
	}

	if !row.Next() {
		m := "weather not found"
		logger.Info(ctx, m)
		return nil, errors.DataStoreValueNotFoundSystemError(err)
	}

	w := &entity.Weather{
		Location: &entity.Location{},
	}
	if err := row.Scan(&w.Dat, &w.Weather, &w.Location.Id, &w.Location.City, &w.Comment); err != nil {
		logger.Error(ctx, "failed to scan.", err)
		return nil, errors.DataStoreSystemError(err)
	}

	return w, nil
}

func (c *MySQLClient) FindLocation(ctx context.Context, locationId int) error {
	logger := log.GetLogger()

	m := application.GetLocationsMap()
	if _, ok := m[locationId]; !ok {
		logger.Warn(ctx, "location id is not mapped.")
		return errors.DataStoreValueNotFoundSystemError(nil)
	}

	return nil
}
