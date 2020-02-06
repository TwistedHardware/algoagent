package models

type DataSourceDuration int

func (DataSourceDuration) Daily() DataSourceDuration {
	return 86400
}
