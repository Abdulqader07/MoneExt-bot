package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/thevolcanomanishere/gud-price/generated/go/ethereum"
	"github.com/thevolcanomanishere/gud-price/generated/go/rpc"
)

// PriceResult is the clean output every asset function returns.
type PriceResult struct {
	Name      string
	Price     float64
	UpdatedAt string // human-readable, e.g. "2026-10-09 04:36:23"
}

// fetchPrice is the shared helper — no asset-specific logic here.
func fetchPrice(name, feed string) (PriceResult, error) {
	data, err := rpc.ReadLatestPrice(feed)
	if err != nil {
		return PriceResult{}, fmt.Errorf("%s: %w", name, err)
	}

	price, err := strconv.ParseFloat(data.Answer, 64)
	if err != nil {
		return PriceResult{}, fmt.Errorf("%s parse: %w", name, err)
	}

	return PriceResult{
		Name:      name,
		Price:     price,
		UpdatedAt: data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func GetGold() (PriceResult, error) {
	return fetchPrice("Gold (XAU/USD)", ethereum.XAU_USD)
}

func GetSilver() (PriceResult, error) {
	return fetchPrice("Silver (XAG/USD)", ethereum.XAG_USD)
}

func GetBitcoin() (PriceResult, error) {
	return fetchPrice("Bitcoin (BTC/USD)", ethereum.BTC_USD)
}

func GetEthereum() (PriceResult, error) {
	return fetchPrice("Ethereum (ETH/USD)", ethereum.ETH_USD)
}

func GetSolana() (PriceResult, error) {
	return fetchPrice("Solana (SOL/USD)", ethereum.SOL_USD)
}

// GetAllAssets fetches every asset concurrently.
func GetAllAssets() []PriceResult {
	type job struct {
		fn func() (PriceResult, error)
	}

	jobs := []job{
		{GetGold}, {GetSilver}, {GetBitcoin}, {GetEthereum}, {GetSolana},
	}

	results := make([]PriceResult, len(jobs))
	var wg sync.WaitGroup

	for i, j := range jobs {
		wg.Add(1)
		go func(idx int, f func() (PriceResult, error)) {
			defer wg.Done()
			res, err := f()
			if err != nil {
				results[idx] = PriceResult{Name: "ERROR", UpdatedAt: err.Error()}
				return
			}
			results[idx] = res
		}(i, j.fn)
	}

	wg.Wait()
	return results
}

// formatHumanTime converts a timestamp into something like "just now", "3m ago".
func formatHumanTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
