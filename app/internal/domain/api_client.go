package domain

import (
	"time"
)

type APIClient struct {
	ID             string
	APIServiceID   string
	Name           string
	Description    string
	status         ClientStatus
	ActiveConfigID *string
	CreatedAt      time.Time
}

func NewAPIClient(id, name, description, apiServiceID string) *APIClient {
	return &APIClient{
		ID:           id,
		Name:         name,
		Description:  description,
		APIServiceID: apiServiceID,
		status:       ClientStatusCreated,
		CreatedAt:    time.Now(),
	}
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
		if to == ClientStatusRunning || to == ClientStatusDisabled {
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
		if to == ClientStatusRunning ||
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
	return c.status == ClientStatusCreated || c.status == ClientStatusStopped
}

func (c *APIClient) CanStop() bool {
	return c.status == ClientStatusRunning
}

func (c *APIClient) GetStatus() ClientStatus {
	return c.status
}
