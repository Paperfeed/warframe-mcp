package recommender

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/paperfeed/warframe-mcp/internal/alecaframe"
	"github.com/paperfeed/warframe-mcp/internal/market"
	"github.com/paperfeed/warframe-mcp/internal/worldstate"
)

// Recommender provides intelligent recommendations based on account and world state
type Recommender struct {
	alecaframeClient *alecaframe.Client
	worldStateClient *worldstate.Client
	marketClient     *market.Client
}

// NewRecommender creates a new recommender
func NewRecommender(
	alecaframeClient *alecaframe.Client,
	worldStateClient *worldstate.Client,
	marketClient *market.Client,
) *Recommender {
	return &Recommender{
		alecaframeClient: alecaframeClient,
		worldStateClient: worldStateClient,
		marketClient:     marketClient,
	}
}

// PrimePartRecommendation represents a recommendation for farming prime parts
type PrimePartRecommendation struct {
	Summary        string
	CurrentlyUnvaulted []string
	BestFissures   []FissureRecommendation
	MarketValue    []MarketRecommendation
	UserRelics     []string
}

// FissureRecommendation represents a fissure to run
type FissureRecommendation struct {
	Node        string
	Tier        string
	MissionType string
	Faction     string
	TimeLeft    string
	IsVoidStorm bool
}

// MarketRecommendation represents valuable items to farm for trading
type MarketRecommendation struct {
	ItemName     string
	AvgPrice     float64
	Demand       int
	CurrentlyAvailable bool
}

// StandingRecommendation represents syndicate standing recommendations
type StandingRecommendation struct {
	Summary            string
	ActiveSyndicates   []SyndicateInfo
	BestMissions       []string
	CurrentEvents      []string
}

// SyndicateInfo contains syndicate information
type SyndicateInfo struct {
	Name     string
	Missions []string
	Rewards  []string
}

// CreditRecommendation represents credit farming recommendations
type CreditRecommendation struct {
	Summary       string
	BestMethods   []CreditMethod
	ActiveBonuses []string
	CurrentEvents []string
}

// CreditMethod represents a credit farming method
type CreditMethod struct {
	Name        string
	Credits     string
	Requirements string
	Duration    string
	Efficiency  string
}

// RecommendPrimeParts provides recommendations for farming prime parts
func (r *Recommender) RecommendPrimeParts() (*PrimePartRecommendation, error) {
	ws, err := r.worldStateClient.GetWorldState()
	if err != nil {
		return nil, fmt.Errorf("getting world state: %w", err)
	}

	vaultTrader, err := r.worldStateClient.GetVaultTrader()
	if err != nil {
		return nil, fmt.Errorf("getting vault trader: %w", err)
	}

	// Get user's relics if available
	var userRelics []string
	if r.alecaframeClient != nil {
		inventory, err := r.alecaframeClient.GetRelicInventory()
		if err == nil && inventory != nil {
			for _, relic := range inventory.Relics {
				if relic.Count > 0 {
					userRelics = append(userRelics, fmt.Sprintf("%s %s (%s) x%d", relic.Type, relic.Name, relic.Refinement, relic.Count))
				}
			}
		}
	}

	// Get currently unvaulted primes from Varzia
	var unvaulted []string
	if vaultTrader != nil && vaultTrader.Active {
		for _, schedule := range vaultTrader.Schedule {
			// Check if schedule is currently active
			start, _ := time.Parse(time.RFC3339, schedule.Start)
			end, _ := time.Parse(time.RFC3339, schedule.End)
			now := time.Now()
			if now.After(start) && now.Before(end) {
				unvaulted = append(unvaulted, schedule.Featured...)
			}
		}
	}

	// Get active fissures
	var fissures []FissureRecommendation
	for _, f := range ws.Fissures {
		if f.Active {
			timeLeft := "Unknown"
			if f.Expiry != "" {
				expiry, err := time.Parse(time.RFC3339, f.Expiry)
				if err == nil {
					timeLeft = time.Until(expiry).Round(time.Minute).String()
				}
			}

			fissures = append(fissures, FissureRecommendation{
				Node:        f.Node,
				Tier:        f.Tier,
				MissionType: f.MissionType,
				Faction:     f.Enemy,
				TimeLeft:    timeLeft,
				IsVoidStorm: f.IsStorm,
			})
		}
	}

	// Sort fissures by tier (higher tiers first)
	sort.Slice(fissures, func(i, j int) bool {
		tierOrder := map[string]int{"Requiem": 5, "Axi": 4, "Neo": 3, "Meso": 2, "Lith": 1}
		return tierOrder[fissures[i].Tier] > tierOrder[fissures[j].Tier]
	})

	summary := r.buildPrimePartSummary(unvaulted, fissures, userRelics)

	return &PrimePartRecommendation{
		Summary:            summary,
		CurrentlyUnvaulted: unvaulted,
		BestFissures:       fissures[:min(5, len(fissures))], // Top 5 fissures
		UserRelics:         userRelics,
	}, nil
}

// RecommendStanding provides syndicate standing recommendations
func (r *Recommender) RecommendStanding() (*StandingRecommendation, error) {
	ws, err := r.worldStateClient.GetWorldState()
	if err != nil {
		return nil, fmt.Errorf("getting world state: %w", err)
	}

	var syndicates []SyndicateInfo
	var bestMissions []string

	// Process syndicate missions
	for _, sm := range ws.SyndicateMissions {
		info := SyndicateInfo{
			Name:     sm.Syndicate,
			Missions: sm.Nodes,
		}
		syndicates = append(syndicates, info)

		for _, node := range sm.Nodes {
			bestMissions = append(bestMissions, fmt.Sprintf("%s - %s", sm.Syndicate, node))
		}
	}

	// Check for events that provide standing
	var currentEvents []string
	for _, event := range ws.Events {
		if event.Active && strings.Contains(strings.ToLower(event.Description), "standing") {
			currentEvents = append(currentEvents, event.Description)
		}
	}

	// Add open world cycles for bounty optimization
	if ws.CetusCycle != nil {
		state := "Day"
		if !ws.CetusCycle.IsDay {
			state = "Night"
		}
		currentEvents = append(currentEvents, fmt.Sprintf("Cetus: %s (resets in %s)", state, ws.CetusCycle.TimeLeft))
	}

	if ws.VallisEarth != nil {
		state := "Warm"
		if !ws.VallisEarth.IsWarm {
			state = "Cold"
		}
		currentEvents = append(currentEvents, fmt.Sprintf("Orb Vallis: %s (resets in %s)", state, ws.VallisEarth.TimeLeft))
	}

	if ws.CambionDrift != nil {
		currentEvents = append(currentEvents, fmt.Sprintf("Cambion Drift: %s", ws.CambionDrift.State))
	}

	summary := r.buildStandingSummary(syndicates, currentEvents)

	return &StandingRecommendation{
		Summary:          summary,
		ActiveSyndicates: syndicates,
		BestMissions:     bestMissions,
		CurrentEvents:    currentEvents,
	}, nil
}

// RecommendCredits provides credit farming recommendations
func (r *Recommender) RecommendCredits() (*CreditRecommendation, error) {
	ws, err := r.worldStateClient.GetWorldState()
	if err != nil {
		return nil, fmt.Errorf("getting world state: %w", err)
	}

	methods := []CreditMethod{
		{
			Name:        "Index (High Risk)",
			Credits:     "75,000-200,000 per round",
			Requirements: "Neptune - Warframe with decent survivability",
			Duration:    "5-10 minutes",
			Efficiency:  "Excellent",
		},
		{
			Name:        "Profit-Taker Orb",
			Credits:     "125,000 + 1 Credit Booster",
			Requirements: "Max rank Solaris United, Archwing, good gear",
			Duration:    "5-8 minutes",
			Efficiency:  "Excellent (with booster)",
		},
		{
			Name:        "Railjack Missions",
			Credits:     "100,000-400,000 per mission",
			Requirements: "Railjack, Veil Proxima access",
			Duration:    "10-15 minutes",
			Efficiency:  "Very Good",
		},
	}

	var activeBonuses []string
	var currentEvents []string

	// Check for credit boosting events
	for _, event := range ws.Events {
		if event.Active {
			desc := strings.ToLower(event.Description)
			if strings.Contains(desc, "credit") || strings.Contains(desc, "credits") {
				currentEvents = append(currentEvents, event.Description)
				activeBonuses = append(activeBonuses, "Event credit bonus active!")
			}
		}
	}

	// Check Arbitration for credit boost
	if ws.Arbitration != nil && ws.Arbitration.Active {
		currentEvents = append(currentEvents, fmt.Sprintf("Arbitration: %s (%s)", ws.Arbitration.Node, ws.Arbitration.Type))
		methods = append(methods, CreditMethod{
			Name:        "Arbitration",
			Credits:     "Varies (with credit boosters from shop)",
			Requirements: "Completed Star Chart",
			Duration:    "Endless",
			Efficiency:  "Good",
		})
	}

	// Check for Dark Sector missions (passive credit farming)
	activeBonuses = append(activeBonuses, "Dark Sector missions provide credit bonuses")

	summary := r.buildCreditSummary(methods, activeBonuses)

	return &CreditRecommendation{
		Summary:       summary,
		BestMethods:   methods,
		ActiveBonuses: activeBonuses,
		CurrentEvents: currentEvents,
	}, nil
}

// Helper methods to build summaries
func (r *Recommender) buildPrimePartSummary(unvaulted []string, fissures []FissureRecommendation, relics []string) string {
	var parts []string

	if len(unvaulted) > 0 {
		parts = append(parts, fmt.Sprintf("Currently unvaulted: %s", strings.Join(unvaulted, ", ")))
	}

	if len(fissures) > 0 {
		parts = append(parts, fmt.Sprintf("\nBest fissures to run:\n"))
		for i, f := range fissures {
			if i >= 3 {
				break
			}
			parts = append(parts, fmt.Sprintf("  • %s %s on %s (%s left)", f.Tier, f.MissionType, f.Node, f.TimeLeft))
		}
	}

	if len(relics) > 0 {
		parts = append(parts, fmt.Sprintf("\n\nYou have %d unique relics in your inventory", len(relics)))
	}

	if len(parts) == 0 {
		return "No specific prime farming recommendations at this time. Check back when fissures are active!"
	}

	return strings.Join(parts, "\n")
}

func (r *Recommender) buildStandingSummary(syndicates []SyndicateInfo, events []string) string {
	var parts []string

	if len(syndicates) > 0 {
		parts = append(parts, fmt.Sprintf("Active syndicates: %d", len(syndicates)))
		for _, syn := range syndicates {
			if len(syn.Missions) > 0 {
				parts = append(parts, fmt.Sprintf("\n%s: %d missions available", syn.Name, len(syn.Missions)))
			}
		}
	}

	if len(events) > 0 {
		parts = append(parts, "\n\nCurrent world state:")
		for _, event := range events {
			parts = append(parts, fmt.Sprintf("  • %s", event))
		}
	}

	if len(parts) == 0 {
		return "Complete syndicate daily missions and open-world bounties for standing."
	}

	return strings.Join(parts, "\n")
}

func (r *Recommender) buildCreditSummary(methods []CreditMethod, bonuses []string) string {
	var parts []string

	parts = append(parts, "Best credit farming methods:")
	for i, method := range methods {
		if i >= 3 {
			break
		}
		parts = append(parts, fmt.Sprintf("\n%d. %s: %s (%s)", i+1, method.Name, method.Credits, method.Efficiency))
	}

	if len(bonuses) > 0 {
		parts = append(parts, "\n\nActive bonuses:")
		for _, bonus := range bonuses {
			parts = append(parts, fmt.Sprintf("  • %s", bonus))
		}
	}

	return strings.Join(parts, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
