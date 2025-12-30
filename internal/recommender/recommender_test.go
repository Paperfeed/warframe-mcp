package recommender

import (
	"testing"
	"time"

	"github.com/paperfeed/warframe-mcp/internal/alecaframe"
	"github.com/paperfeed/warframe-mcp/internal/market"
	"github.com/paperfeed/warframe-mcp/internal/worldstate"
)

func TestNewRecommender(t *testing.T) {
	wsClient := worldstate.NewClient(5 * time.Minute)
	mClient := market.NewClient(10 * time.Minute)

	rec := NewRecommender(nil, wsClient, mClient)

	if rec == nil {
		t.Fatal("Expected recommender to be created, got nil")
	}

	if rec.worldStateClient == nil {
		t.Error("Expected worldStateClient to be set")
	}

	if rec.marketClient == nil {
		t.Error("Expected marketClient to be set")
	}
}

func TestRecommendPrimeParts(t *testing.T) {
	wsClient := worldstate.NewClient(5 * time.Minute)
	mClient := market.NewClient(10 * time.Minute)

	rec := NewRecommender(nil, wsClient, mClient)

	recommendation, err := rec.RecommendPrimeParts()
	if err != nil {
		t.Logf("Warning: RecommendPrimeParts failed (API may be down): %v", err)
		return
	}

	if recommendation == nil {
		t.Fatal("Expected recommendation to be returned, got nil")
	}

	t.Logf("Prime farming recommendation:")
	t.Logf("  Summary: %s", recommendation.Summary)
	t.Logf("  Currently unvaulted: %d items", len(recommendation.CurrentlyUnvaulted))
	t.Logf("  Best fissures: %d available", len(recommendation.BestFissures))

	if len(recommendation.BestFissures) > 0 {
		fissure := recommendation.BestFissures[0]
		t.Logf("  Top fissure: %s %s on %s", fissure.Tier, fissure.MissionType, fissure.Node)
	}
}

func TestRecommendStanding(t *testing.T) {
	wsClient := worldstate.NewClient(5 * time.Minute)
	mClient := market.NewClient(10 * time.Minute)

	rec := NewRecommender(nil, wsClient, mClient)

	recommendation, err := rec.RecommendStanding()
	if err != nil {
		t.Logf("Warning: RecommendStanding failed (API may be down): %v", err)
		return
	}

	if recommendation == nil {
		t.Fatal("Expected recommendation to be returned, got nil")
	}

	t.Logf("Standing farming recommendation:")
	t.Logf("  Summary: %s", recommendation.Summary)
	t.Logf("  Active syndicates: %d", len(recommendation.ActiveSyndicates))
	t.Logf("  Best missions: %d available", len(recommendation.BestMissions))
	t.Logf("  Current events: %d", len(recommendation.CurrentEvents))
}

func TestRecommendCredits(t *testing.T) {
	wsClient := worldstate.NewClient(5 * time.Minute)
	mClient := market.NewClient(10 * time.Minute)

	rec := NewRecommender(nil, wsClient, mClient)

	recommendation, err := rec.RecommendCredits()
	if err != nil {
		t.Logf("Warning: RecommendCredits failed (API may be down): %v", err)
		return
	}

	if recommendation == nil {
		t.Fatal("Expected recommendation to be returned, got nil")
	}

	t.Logf("Credit farming recommendation:")
	t.Logf("  Summary: %s", recommendation.Summary)
	t.Logf("  Best methods: %d available", len(recommendation.BestMethods))
	t.Logf("  Active bonuses: %d", len(recommendation.ActiveBonuses))

	if len(recommendation.BestMethods) > 0 {
		method := recommendation.BestMethods[0]
		t.Logf("  Top method: %s (%s)", method.Name, method.Efficiency)
	}
}

func TestRecommenderWithAlecaframe(t *testing.T) {
	// This test requires actual credentials, so we skip if not available
	t.Skip("Skipping Alecaframe integration test (requires credentials)")

	aClient := alecaframe.NewClient("test-hash", "test-token")
	wsClient := worldstate.NewClient(5 * time.Minute)
	mClient := market.NewClient(10 * time.Minute)

	rec := NewRecommender(aClient, wsClient, mClient)

	if rec.alecaframeClient == nil {
		t.Error("Expected alecaframeClient to be set")
	}

	recommendation, err := rec.RecommendPrimeParts()
	if err != nil {
		t.Fatalf("RecommendPrimeParts failed: %v", err)
	}

	// With Alecaframe, we should have user relics
	if len(recommendation.UserRelics) == 0 {
		t.Log("Warning: Expected user relics to be populated with Alecaframe client")
	}
}

func TestBuildPrimePartSummary(t *testing.T) {
	rec := &Recommender{}

	unvaulted := []string{"Mesa Prime", "Ash Prime"}
	fissures := []FissureRecommendation{
		{
			Node:        "Mot (Void)",
			Tier:        "Axi",
			MissionType: "Survival",
			TimeLeft:    "45m",
		},
	}
	relics := []string{"Axi A5 (10)", "Neo N3 (5)"}

	summary := rec.buildPrimePartSummary(unvaulted, fissures, relics)

	if summary == "" {
		t.Error("Expected summary to be generated, got empty string")
	}

	t.Logf("Generated summary:\n%s", summary)
}

func TestBuildStandingSummary(t *testing.T) {
	rec := &Recommender{}

	syndicates := []SyndicateInfo{
		{
			Name:     "Steel Meridian",
			Missions: []string{"Node 1", "Node 2"},
		},
	}
	events := []string{"Cetus: Day (2h remaining)"}

	summary := rec.buildStandingSummary(syndicates, events)

	if summary == "" {
		t.Error("Expected summary to be generated, got empty string")
	}

	t.Logf("Generated summary:\n%s", summary)
}

func TestBuildCreditSummary(t *testing.T) {
	rec := &Recommender{}

	methods := []CreditMethod{
		{
			Name:       "Index (High Risk)",
			Credits:    "75,000-200,000 per round",
			Efficiency: "Excellent",
		},
	}
	bonuses := []string{"Event credit bonus active!"}

	summary := rec.buildCreditSummary(methods, bonuses)

	if summary == "" {
		t.Error("Expected summary to be generated, got empty string")
	}

	t.Logf("Generated summary:\n%s", summary)
}
