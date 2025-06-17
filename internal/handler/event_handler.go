package handler

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gin-gonic/gin"

	alpacaApi "omnichart-server/internal/alpaca"
	"omnichart-server/internal/models"
	"omnichart-server/internal/supabase"
)

func parseSummaryText(rawText string) (string, []string) {
	lines := strings.Split(strings.TrimSpace(rawText), "\n")

	if len(lines) == 0 {
		return "", nil
	}

	title := strings.TrimSpace(lines[0])
	bullets := []string{}
	for _, line := range lines[1:] {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			bullets = append(bullets, trimmed)
		}
	}

	return title, bullets
}

// GetEventsHandler godoc
// @Summary Gets events for a given ticker and timeframe
// @Tags event
// @Param ticker path string true "Ticker symbol"
// @Param from query string true "timeframe of events to fetch in days"
// @Success 200 {array} models.Event
// @Failure 400 {object} map[string]interface{}
// @Router /events/{event_id} [get]
func GetEventsHandler(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))
	timeframe, _ := strconv.Atoi(c.Query("timeframe"))
	log.Println(timeframe)
	events, err := supabase.GetEvents(ticker)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Fetching from database"})
		return
	}
	if len(events) != 0 {
		c.JSON(http.StatusOK, events)
		return
	}

	res, err := alpacaApi.MarketData.GetCorporateActions(marketdata.GetCorporateActionsRequest{Symbols: []string{ticker}, Start: civil.DateOf(time.Now().AddDate(0, 0, -timeframe*3)), TotalLimit: 3})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)

	res2, err2 := alpacaApi.MarketData.GetNews(marketdata.GetNewsRequest{Symbols: []string{ticker}, Start: time.Now().AddDate(0, 0, -timeframe*3), TotalLimit: 3})
	if err2 != nil {
		log.Println(err2)
	}
	// log.Println(res2[0].URL)
	recentEvent := res2[0]

	client := anthropic.NewClient()

	content := anthropic.ContentBlockParamUnion{
		OfText: &anthropic.TextBlockParam{
			Type: "text",
			Text: "Summarize the following article from a url into max 3 key summary sentences and a summarized title. " +
				"Do not include any redundant system phrases. the headline and summary sentences should not be questions" +
				"Have the title and each bullet point separated by one newline and nothing else: " + recentEvent.Headline + " " + recentEvent.URL,
		},
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			{
				Role:    "user",
				Content: []anthropic.ContentBlockParamUnion{content},
			},
		},
		Model: anthropic.ModelClaude3_7SonnetLatest,
	})
	if err != nil {
		panic(err.Error())
	}

	rawText := message.Content[0].Text
	title, bullets := parseSummaryText(rawText)

	summarizedEvent := models.Event{
		ID:           uuid.New(),
		Timestamp:    recentEvent.CreatedAt,
		Title:        title,
		SourceUrl:    recentEvent.URL,
		Content:      strings.Join(bullets, "\n"), // or use JSON if you prefer bullet structure
		EventTypesID: 19, // AI_generated Events
	}

	if _, _, err := supabase.Client.From("events").Insert(summarizedEvent, false, "representation", "", "").Execute(); err != nil {
		log.Println("Error inserting event:", err)
		return
	}

	summarizedTickerEvent := models.TickerEvent{
		ID:         uuid.New(),
		Ticker:     ticker,
		EventId:    summarizedEvent.ID,
		StartIndex: 0,
		EndIndex:   0,
	}

	if _, _, err := supabase.Client.From("ticker_event").Insert(summarizedTickerEvent, false, "representation", "", "").Execute(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting ticker_event"})
		return
	}
	
	summarizedTickerEvent.Event = &summarizedEvent

	c.JSON(http.StatusOK, []models.TickerEvent{summarizedTickerEvent})
}
