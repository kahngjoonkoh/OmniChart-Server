package handler

import (
	"net/http"
	"strings"
	"time"
	"omnichart-server/internal/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
	"os"
)

type Split struct {
	EffectiveDate string  `json:"effective_date"` // e.g. "2023-10-20"
	Ratio         float64 `json:"ratio"`          // 2.0 = 2-for-1 split
}

type ActionSplit struct {
    Symbol      string  `json:"symbol"`
    ExDate      string  `json:"ex_date"`
    NewRate     float64 `json:"new_rate"`
    OldRate     float64 `json:"old_rate"`
    ProcessDate string  `json:"process_date"`
}

type corporateActionResponse struct {
    CorporateActions struct {
        ForwardSplits []ActionSplit `json:"forward_splits"`
    } `json:"corporate_actions"`
    NextPageToken *string `json:"next_page_token"`
}

func GetSplits(ticker string, start, end time.Time) ([]Split, error) {
    url := fmt.Sprintf(
        "https://data.alpaca.markets/v1/corporate-actions?symbols=%s&types=%s&start=%s&end=%s",
        ticker, "forward_split", start.Format("2006-01-02"), end.Format("2006-01-02"),
    )

    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("Apca-Api-Key-Id", os.Getenv("APCA_API_KEY_ID"))
    req.Header.Set("Apca-Api-Secret-Key", os.Getenv("APCA_API_SECRET_KEY"))

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var raw corporateActionResponse
    if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
        return nil, err
    }

    var splits []Split
    for _, fs := range raw.CorporateActions.ForwardSplits {
        ratio := fs.NewRate / fs.OldRate
        splits = append(splits, Split{
            EffectiveDate: fs.ExDate,
            Ratio:         ratio,
        })
    }

    return splits, nil
}

func GetHistoricalDataHandler(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))
	startStr := c.Query("start")
	endStr := c.Query("end")

	// Input validation (same as before)...

	start, _ := time.Parse(time.RFC3339, startStr)
	end, _ := time.Parse(time.RFC3339, endStr)

	bars, err := alpacaApi.MarketData.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     start,
		End:       end,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching historical bars", "details": err.Error()})
		return
	}

	// Get stock split data
	splits, err := GetSplits(ticker, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching stock splits", "details": err.Error()})
		return
	}

	fmt.Printf("✅ Got %d splits\n", len(splits))


	fmt.Println("==== SPLIT DEBUG START ====")
	for _, split := range splits {
		fmt.Printf("Effective Date: %s | Ratio: %.2f\n", split.EffectiveDate, split.Ratio)
	}
	fmt.Println("==== SPLIT DEBUG END ====")


	// Apply split adjustments
	for i := range bars {
		barTime := bars[i].Timestamp
		adjustment := 1.0

		for _, split := range splits {
			effDate, _ := time.Parse("2006-01-02", split.EffectiveDate)
			if barTime.Before(effDate) {
				adjustment *= split.Ratio
			}
		}

		if adjustment != 1.0 {
			bars[i].Open /= adjustment
			bars[i].High /= adjustment
			bars[i].Low /= adjustment
			bars[i].Close /= adjustment
			bars[i].Volume = uint64(float64(bars[i].Volume) * adjustment)
		}
	}

	c.JSON(http.StatusOK, bars)
}
