package models

type TickerType int

func (TickerType) Stock() TickerType {
	return 1
}

func (TickerType) Index() TickerType {
	return 2
}

func (TickerType) Forex() TickerType {
	return 3
}

func (TickerType) Future() TickerType {
	return 4
}

func (TickerType) Metal() TickerType {
	return 5
}

func (TickerType) Other() TickerType {
	return 99
}
