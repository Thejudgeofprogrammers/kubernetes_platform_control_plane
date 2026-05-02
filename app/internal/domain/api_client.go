package domain

import (
	"strings"
	"time"
)

type APIClient struct {
	ID             string
	Slug           string
	APIServiceID   string
	Name           string
	Description    string
	status         ClientStatus
	ActiveConfigID *string
	CreatedAt      time.Time
}

// CREATED    → DEPLOYING
// DEPLOYING  → RUNNING
// RUNNING    → STOPPING / RESTARTING / DEPLOYING / DELETING
// STOPPING   → STOPPED
// STOPPED    → DEPLOYING / DELETING
// RESTARTING → RUNNING
// DELETING   → DISABLED

//            ┌──────────────┐
//            │   CREATED    │
//            └──────┬───────┘
//                   ↓
//            ┌──────────────┐
//            │  DEPLOYING   │
//            └──────┬───────┘
//                   ↓
//            ┌──────────────┐
//            │   RUNNING    │
//            └──────┬───────┘
//         ┌─────────┴─────────┐
//         ↓                   ↓
//  ┌──────────────┐    ┌──────────────┐
//  │  STOPPING    │    │ RESTARTING   │
//  └──────┬───────┘    └──────┬───────┘
//         ↓                   ↓
//  ┌──────────────┐    ┌──────────────┐
//  │   STOPPED    │    │   RUNNING    │
//  └──────┬───────┘
//         ↓
//  ┌──────────────┐
//  │  DELETING    │
//  └──────┬───────┘
//         ↓
//  ┌──────────────┐
//  │  DISABLED    │
//  └──────────────┘

func NewAPIClient(id, name, description, apiServiceID string) *APIClient {
	return &APIClient{
		ID:           id,
		Slug:         generateSlug(name),
		Name:         name,
		Description:  description,
		APIServiceID: apiServiceID,
		status:       ClientStatusCreated,
		CreatedAt:    time.Now(),
	}
}

func generateSlug(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func (c *APIClient) ActivateConfig(configID string) {
	c.ActiveConfigID = &configID
}

func (c *APIClient) Transition(to ClientStatus) error {
	if c.status == to {
		return nil
	}

	switch c.status {

	case ClientStatusCreated:
		if to == ClientStatusDeploying || to == ClientStatusDisabled {
			c.status = to
			return nil
		}

	case ClientStatusStopping:
		if to == ClientStatusStopped {
			c.status = to
			return nil
		}

	case ClientStatusRunning:
		if to == ClientStatusRestarting ||
			to == ClientStatusStopping ||
			to == ClientStatusDeleting ||
			to == ClientStatusDeploying {
			c.status = to
			return nil
		}

	case ClientStatusRestarting:
		if to == ClientStatusRunning {
			c.status = to
			return nil
		}

	case ClientStatusStopped:
		if to == ClientStatusDeploying ||
			to == ClientStatusDeleting ||
			to == ClientStatusDisabled {
			c.status = to
			return nil
		}

	case ClientStatusDeleting:
		if to == ClientStatusDisabled {
			c.status = to
			return nil
		}

	case ClientStatusDisabled:
		return ErrInvalidStateTransition

	case ClientStatusDeploying:
		if to == ClientStatusRunning {
			c.status = to
			return nil
		}

	}

	return ErrInvalidStateTransition
}

func (c *APIClient) CanStart() bool {
	return c.status == ClientStatusCreated ||
		c.status == ClientStatusStopped
}

func (c *APIClient) CanStop() bool {
	return c.status == ClientStatusRunning
}

func (c *APIClient) CanRestart() bool {
	return c.status == ClientStatusRunning
}

func (c *APIClient) CanDelete() bool {
	return c.status != ClientStatusDeleting &&
		c.status != ClientStatusDisabled
}

func (c *APIClient) CanDeploy() bool {
	return c.status == ClientStatusRunning ||
		c.status == ClientStatusStopped
}

func (c *APIClient) GetStatus() ClientStatus {
	return c.status
}
