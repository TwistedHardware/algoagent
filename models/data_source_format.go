package models

type DataSourceFormat int

func (DataSourceFormat) JSON() DataSourceFormat {
	return 1
}

func (DataSourceFormat) CSVWithHeader() DataSourceFormat {
	return 2
}
