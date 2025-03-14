package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"t800/internal/common"
	"t800/internal/monitoring"
)

// DecisionMaker handles AI-based decision making
type DecisionMaker struct {
	baseURL string
	logger  *monitoring.Logger
	model   string
}

// CombatDecision represents the AI's decision for combat
type CombatDecision struct {
	Action      string  `json:"action"`       // "move", "attack", "defend", "retreat"
	Target      string  `json:"target"`       // Target ID if applicable
	Weapon      string  `json:"weapon"`       // Weapon to use if attacking
	Priority    int     `json:"priority"`     // Priority level (1-10)
	Confidence  float64 `json:"confidence"`   // Confidence in the decision (0-1)
	Explanation string  `json:"explanation"`  // Explanation of the decision
}

// EngagementDecision represents the AI's decision for threat engagement
type EngagementDecision struct {
	ShouldEngage bool    `json:"should_engage"`
	Confidence   float64 `json:"confidence"`
	Explanation  string  `json:"explanation"`
}

// NewDecisionMaker creates a new AI decision maker
func NewDecisionMaker(logger *monitoring.Logger) (*DecisionMaker, error) {
	baseURL := os.Getenv("OLLAMA_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2" // Default model
	}

	return &DecisionMaker{
		baseURL: baseURL,
		logger:  logger,
		model:   model,
	}, nil
}

// callOllama makes a request to the Ollama API and returns the cleaned response
func (d *DecisionMaker) callOllama(ctx context.Context, prompt string, result interface{}) error {
	// Prepare the request body
	requestBody := map[string]interface{}{
		"model":  d.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create and send the request
	req, err := http.NewRequestWithContext(ctx, "POST", d.baseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get AI decision: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	var ollamaResponse struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &ollamaResponse); err != nil {
		return fmt.Errorf("failed to parse Ollama response: %v", err)
	}

	// Clean and parse the response
	cleanResponse := d.cleanResponse(ollamaResponse.Response)
	if err := json.Unmarshal([]byte(cleanResponse), result); err != nil {
		return fmt.Errorf("failed to parse AI decision: %v (response: %s)", err, cleanResponse)
	}

	return nil
}

// cleanResponse removes markdown code block markers and trims whitespace
func (d *DecisionMaker) cleanResponse(response string) string {
	clean := response
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	return strings.TrimSpace(clean)
}

// MakeCombatDecision makes a decision based on current state and threats
func (d *DecisionMaker) MakeCombatDecision(
	ctx context.Context,
	currentLoc common.Location,
	activeThreat *common.Threat,
	healthStatus map[string]float64,
	availableWeapons []string,
) (*CombatDecision, error) {
	prompt := fmt.Sprintf(`You are the AI core of a T800 combat robot. Analyze the following situation and make a tactical decision.

Current Location: (%.2f, %.2f, %.2f)
Active Threat: %s (Severity: %d, Location: (%.2f, %.2f, %.2f))
Health Status: %v
Available Weapons: %v

Make a tactical decision considering:
1. Distance to threat
2. Threat severity
3. Current health status
4. Available weapons
5. Strategic advantage

IMPORTANT: Respond with ONLY a valid JSON object in the following format:
{
    "action": "move", "attack", "defend", or "retreat",
    "target": "target ID if applicable",
    "weapon": "weapon to use if attacking",
    "priority": number between 1-10,
    "confidence": number between 0-1,
    "explanation": "brief explanation of the decision"
}

Do not include any text before or after the JSON object.`,
		currentLoc.X, currentLoc.Y, currentLoc.Z,
		activeThreat.ID, activeThreat.Severity,
		activeThreat.Location.X, activeThreat.Location.Y, activeThreat.Location.Z,
		healthStatus,
		availableWeapons)

	var decision CombatDecision
	if err := d.callOllama(ctx, prompt, &decision); err != nil {
		return nil, err
	}

	d.logger.Info(fmt.Sprintf("AI Decision: %s (Confidence: %.2f) - %s",
		decision.Action, decision.Confidence, decision.Explanation))

	return &decision, nil
}

// ShouldEngageProactively determines if the robot should engage a threat proactively
func (d *DecisionMaker) ShouldEngageProactively(
	ctx context.Context,
	threat common.Threat,
	currentLoc common.Location,
	healthStatus map[string]float64,
) (bool, error) {
	prompt := fmt.Sprintf(`Analyze if the T800 should proactively engage this threat.

Threat: %s (Severity: %d, Location: (%.2f, %.2f, %.2f))
Current Location: (%.2f, %.2f, %.2f)
Health Status: %v

Consider:
1. Threat severity
2. Distance
3. Current health status
4. Strategic advantage

IMPORTANT: Respond with ONLY a valid JSON object in the following format:
{
    "should_engage": true or false,
    "confidence": number between 0-1,
    "explanation": "brief explanation"
}

Do not include any text before or after the JSON object.`,
		threat.ID, threat.Severity,
		threat.Location.X, threat.Location.Y, threat.Location.Z,
		currentLoc.X, currentLoc.Y, currentLoc.Z,
		healthStatus)

	var decision EngagementDecision
	if err := d.callOllama(ctx, prompt, &decision); err != nil {
		return false, err
	}

	d.logger.Info(fmt.Sprintf("AI Engagement Decision: %v (Confidence: %.2f) - %s",
		decision.ShouldEngage, decision.Confidence, decision.Explanation))

	return decision.ShouldEngage, nil
} 