// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

// Entity identifies a single entity.
type Entity struct {
	Tag string
}

// Entities identifies multiple entities.
type Entities struct {
	Entities []Entity
}

// LifeResult holds the life status of a single entity, or an error
// indicating why it is not available.
type LifeResult struct {
	Life  Life
	Error *Error
}

// LifeResults holds the life or error status of multiple entities.
type LifeResults struct {
	Results []LifeResult
}

// MachineSetStatus holds a machine tag, status and extra info.
type MachineSetStatus struct {
	Tag    string
	Status Status
	Info   string
}

// MachinesSetStatus holds the parameters for making a Machiner.SetStatus call.
type MachinesSetStatus struct {
	Machines []MachineSetStatus
}

// MachineAgentGetMachinesResults holds the results of a
// machineagent.API.GetMachines call.
type MachineAgentGetMachinesResults struct {
	Machines []MachineAgentGetMachinesResult
}

// MachineAgentGetMachinesResult holds the results of a
// machineagent.API.GetMachines call for a single machine.
type MachineAgentGetMachinesResult struct {
	Life  Life
	Jobs  []MachineJob
	Error *Error
}

// AgentTools describes the tools for a given Agent. This is mostly a flattened
// state.Tools description, plus an agent Tag field.
type AgentTools struct {
	Tag    string
	Major  int
	Minor  int
	Patch  int
	Build  int
	Arch   string
	Series string
	URL    string
}

// AgentToolsResult holds the tools and possibly error for a given Agent request
type AgentToolsResult struct {
	AgentTools AgentTools
	Error      *Error
}

// AgentToolsResults is a list of tools for various requested agents.
type AgentToolsResults struct {
	Tools []AgentToolsResult
}

// Set what tools are being run for multiple agents
type SetAgentTools struct {
	AgentTools []AgentTools
}

// The result of setting the tools for one agent
type SetAgentToolsResult struct {
	Tag   string
	Error *Error
}

// The result of setting the tools for many agents
type SetAgentToolsResults struct {
	Results []SetAgentToolsResult
}

// PasswordChanges holds the parameters for making a SetPasswords call.
type PasswordChanges struct {
	Changes []PasswordChange
}

// PasswordChange specifies a password change for the entity
// with the given tag.
type PasswordChange struct {
	Tag      string
	Password string
}

// A NotifyWatcher will send events when something changes.
// It does not send content for those changes.
type NotifyWatcher interface {
	Changes() <-chan struct{}
	Stop() error
	Err() error
}

// NotifyWatchResult holds an NotifyWatcher id and an error (if any).
type NotifyWatchResult struct {
	NotifyWatcherId string
	Error           *Error
}

// NotifyWatchResults holds the results for any API call which ends up
// returning a list of NotifyWatchers
type NotifyWatchResults struct {
	Results []NotifyWatchResult
}

// LifecycleWatchResults holds the results of API calls
// that watch the lifecycle of a set of objects.
// It is used both for the initial Watch request
// and for subsequent Next requests.
type LifecycleWatchResults struct {
	// LifeCycleWatcherId holds the id of the newly
	// created watcher. It will be empty for a Next
	// request.
	LifecycleWatcherId string

	// Ids holds the list of entity ids.
	// For a Watch request, it holds all entity ids being
	// watched; for a Next request, it holds the ids of those
	// that have changed.
	Ids []string
}
