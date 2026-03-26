package entitlements

import (
	"strings"

	"api-testing-kit/server/internal/auth"
)

type PlanSource string

const (
	PlanSourceGuest         PlanSource = "guest"
	PlanSourceSystemDefault PlanSource = "system_default"
	PlanSourceManual        PlanSource = "manual"
)

type CapabilityKey string

const (
	CapabilityCustomURLExecution CapabilityKey = "custom_url_execution"
	CapabilityHistoryDepth       CapabilityKey = "history_depth"
	CapabilityEnvironmentVars    CapabilityKey = "environment_variables"
	CapabilitySharedLinks        CapabilityKey = "shared_links"
)

type AccessPlan struct {
	Code   string     `json:"code"`
	Name   string     `json:"name"`
	Source PlanSource `json:"source"`
}

type AccessCapability struct {
	Key         CapabilityKey `json:"key"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Enabled     bool          `json:"enabled"`
	Limit       *int          `json:"limit,omitempty"`
	LimitLabel  string        `json:"limitLabel,omitempty"`
	Scope       string        `json:"scope"`
	Reason      string        `json:"reason,omitempty"`
}

type State struct {
	Plan         AccessPlan         `json:"plan"`
	Capabilities []AccessCapability `json:"capabilities"`
}

func Resolve(user auth.UserRecord) State {
	if user.ID == "" || !isActiveUser(user) {
		return guestState()
	}

	state := authenticatedState()
	if isPrivilegedRole(user.Role) {
		enableAllCapabilities(&state)
		state.Plan = AccessPlan{
			Code:   "pro",
			Name:   "Pro",
			Source: PlanSourceSystemDefault,
		}
	}

	return state
}

func Guest() State {
	return guestState()
}

func (s State) Capability(key CapabilityKey) (AccessCapability, bool) {
	for _, capability := range s.Capabilities {
		if capability.Key == key {
			return capability, true
		}
	}

	return AccessCapability{}, false
}

func (s State) Enabled(key CapabilityKey) bool {
	capability, ok := s.Capability(key)
	return ok && capability.Enabled
}

func (s State) HistoryDepthLimit() int {
	capability, ok := s.Capability(CapabilityHistoryDepth)
	if !ok || capability.Limit == nil {
		return 0
	}

	return *capability.Limit
}

func (s State) CanUseCustomURLs() bool {
	return s.Enabled(CapabilityCustomURLExecution)
}

func (s State) CanUseEnvironmentVariables() bool {
	return s.Enabled(CapabilityEnvironmentVars)
}

func (s State) CanCreateSharedLinks() bool {
	return s.Enabled(CapabilitySharedLinks)
}

func isActiveUser(user auth.UserRecord) bool {
	return strings.EqualFold(user.Status, "active") || user.Status == ""
}

func isPrivilegedRole(role string) bool {
	return strings.EqualFold(role, "admin") || strings.EqualFold(role, "owner")
}

func guestState() State {
	return State{
		Plan: AccessPlan{
			Code:   "guest",
			Name:   "Guest Preview",
			Source: PlanSourceGuest,
		},
		Capabilities: []AccessCapability{
			{
				Key:         CapabilityCustomURLExecution,
				Label:       "Custom URL execution",
				Description: "Guests stay on allowlisted templates and cannot replace the target URL.",
				Enabled:     false,
				Scope:       "guest",
				Reason:      "Sign in to send custom outbound requests.",
			},
			{
				Key:         CapabilityHistoryDepth,
				Label:       "History depth",
				Description: "Guest history is visible as a preview, but it is not persisted.",
				Enabled:     false,
				Scope:       "guest",
				Reason:      "Sign in to retain request history.",
			},
			{
				Key:         CapabilityEnvironmentVars,
				Label:       "Environment variables",
				Description: "Environment variables stay locked until the session is authenticated.",
				Enabled:     false,
				Scope:       "guest",
				Reason:      "Sign in to use request variables.",
			},
			{
				Key:         CapabilitySharedLinks,
				Label:       "Shared links",
				Description: "Guest sessions can preview sharing surfaces without creating durable links.",
				Enabled:     false,
				Scope:       "guest",
				Reason:      "Sign in to share saved requests.",
			},
		},
	}
}

func authenticatedState() State {
	historyDepth := 50

	return State{
		Plan: AccessPlan{
			Code:   "starter",
			Name:   "Starter",
			Source: PlanSourceSystemDefault,
		},
		Capabilities: []AccessCapability{
			{
				Key:         CapabilityCustomURLExecution,
				Label:       "Custom URL execution",
				Description: "Authenticated sessions can submit custom URLs after destination validation.",
				Enabled:     true,
				Scope:       "authenticated",
			},
			{
				Key:         CapabilityHistoryDepth,
				Label:       "History depth",
				Description: "Authenticated users retain a bounded history window for replay and inspection.",
				Enabled:     true,
				Limit:       &historyDepth,
				LimitLabel:  "50 runs",
				Scope:       "authenticated",
			},
			{
				Key:         CapabilityEnvironmentVars,
				Label:       "Environment variables",
				Description: "Environment variables remain locked on the base plan until a later tier unlocks them.",
				Enabled:     false,
				Scope:       "plan",
				Reason:      "Upgrade to unlock environment variable storage.",
			},
			{
				Key:         CapabilitySharedLinks,
				Label:       "Shared links",
				Description: "Readonly sharing remains part of a later tier.",
				Enabled:     false,
				Scope:       "plan",
				Reason:      "Upgrade to publish shared links.",
			},
		},
	}
}

func enableAllCapabilities(state *State) {
	if state == nil {
		return
	}

	historyDepth := 200
	for index := range state.Capabilities {
		state.Capabilities[index].Enabled = true
		state.Capabilities[index].Scope = "plan"
		state.Capabilities[index].Reason = ""
		if state.Capabilities[index].Key == CapabilityHistoryDepth {
			state.Capabilities[index].Limit = &historyDepth
			state.Capabilities[index].LimitLabel = "200 runs"
		}
	}
}
