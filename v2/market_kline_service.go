package binance

import (
	"context"
	"fmt"
	"net/http"
)

// KlinesService list klines
type MarketKlinesService struct {
	c         *Client
	symbol    string
	interval  string
	limit     *int
	startTime *int64
	endTime   *int64
}

// Symbol set symbol
func (s *MarketKlinesService) Symbol(symbol string) *MarketKlinesService {
	s.symbol = symbol
	return s
}

// Interval set interval
func (s *MarketKlinesService) Interval(interval string) *MarketKlinesService {
	s.interval = interval
	return s
}

// Limit set limit
func (s *MarketKlinesService) Limit(limit int) *MarketKlinesService {
	s.limit = &limit
	return s
}

// StartTime set startTime
func (s *MarketKlinesService) StartTime(startTime int64) *MarketKlinesService {
	s.startTime = &startTime
	return s
}

// EndTime set endTime
func (s *MarketKlinesService) EndTime(endTime int64) *MarketKlinesService {
	s.endTime = &endTime
	return s
}

// Do send request
func (s *MarketKlinesService) Do(ctx context.Context, opts ...RequestOption) (res []*MarketKline, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/marketKlines",
	}
	r.setParam("symbol", s.symbol)
	r.setParam("interval", s.interval)
	if s.limit != nil {
		r.setParam("limit", *s.limit)
	}
	if s.startTime != nil {
		r.setParam("startTime", *s.startTime)
	}
	if s.endTime != nil {
		r.setParam("endTime", *s.endTime)
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*MarketKline{}, err
	}
	j, err := newJSON(data)
	if err != nil {
		return []*MarketKline{}, err
	}
	num := len(j.MustArray())
	res = make([]*MarketKline, num)
	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		if len(item.MustArray()) < 11 {
			err = fmt.Errorf("invalid kline response")
			return []*MarketKline{}, err
		}
		res[i] = &MarketKline{
			OpenTime:         item.GetIndex(0).MustInt64(),
			Open:             item.GetIndex(1).MustString(),
			High:             item.GetIndex(2).MustString(),
			Low:              item.GetIndex(3).MustString(),
			Close:            item.GetIndex(4).MustString(),
			Volume:           item.GetIndex(5).MustString(),
			CloseTime:        item.GetIndex(6).MustInt64(),
			QuoteAssetVolume: item.GetIndex(7).MustString(),
			TradeNum:         item.GetIndex(8).MustInt64(),
		}
	}
	return res, nil
}

// MarketKline define kline info
type MarketKline struct {
	OpenTime         int64  `json:"openTime"`
	Open             string `json:"open"`
	High             string `json:"high"`
	Low              string `json:"low"`
	Close            string `json:"close"`
	Volume           string `json:"volume"`
	CloseTime        int64  `json:"closeTime"`
	QuoteAssetVolume string `json:"quoteAssetVolume"`
	TradeNum         int64  `json:"tradeNum"`
}
