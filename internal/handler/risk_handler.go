package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// Struct for the Finnhub metric response
type finnhubMetricResponse struct {
    Metric struct {
        Beta float64 `json:"beta"`
    } `json:"metric"`
}

// classifyBetaRisk categorizes the risk level based on the beta value
func classifyBetaRisk(beta float64) string {
    switch {
    case beta < 0:
        return "Inverse Market Risk (Negative Beta)"
    case beta == 0:
        return "No Market Risk"
    case beta > 0 && beta < 1:
        return "Low Risk / Defensive Stock"
    case beta == 1:
        return "Average Market Risk"
    case beta > 1 && beta <= 2:
        return "High Risk / Growth Stock"
    case beta > 2:
        return "Very High Risk / Speculative Stock"
    default:
        return "Unknown"
    }
}

// GET /beta/:ticker
func GetBetaHandler(c *gin.Context) {
    ticker := strings.ToUpper(c.Param("ticker"))
    apiKey := os.Getenv("FINNHUB_API_KEY")

    if apiKey == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Finnhub API key not configured"})
        return
    }

    url := fmt.Sprintf("https://finnhub.io/api/v1/stock/metric?symbol=%s&metric=all&token=%s", ticker, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Finnhub API"})
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Finnhub API returned status %d", resp.StatusCode)})
        return
    }

    var metricResp finnhubMetricResponse
    if err := json.NewDecoder(resp.Body).Decode(&metricResp); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Finnhub response"})
        return
    }

    beta := metricResp.Metric.Beta
    if beta == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Beta not found or is zero for this ticker"})
        return
    }

    classification := classifyBetaRisk(beta)

    c.JSON(http.StatusOK, gin.H{
        "ticker":        ticker,
        "beta":          beta,
        "riskCategory":  classification,
    })
}
