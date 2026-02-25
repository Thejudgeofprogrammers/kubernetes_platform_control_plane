package domain

type ClientStatus string

const (
	ClientStatusCreated    ClientStatus = "created"
	ClientStatusRunning    ClientStatus = "running"
	ClientStatusStopping   ClientStatus = "stopping"
	ClientStatusStopped    ClientStatus = "stopped"
	ClientStatusRestarting ClientStatus = "restarting"
	ClientStatusDeploying  ClientStatus = "deploying"
	ClientStatusDeleting   ClientStatus = "deleting"
	ClientStatusDisabled   ClientStatus = "disabled"
)

type ActionType string

const (
	ActionStart   ActionType = "start"
	ActionStop    ActionType = "stop"
	ActionRestart ActionType = "restart"
	ActionDelete  ActionType = "delete"
	ActionCreate  ActionType = "create"
	ActionDeploy  ActionType = "deploy"
	ActionUpdate  ActionType = "update"
)

type AccessRole string

const (
	RoleOwner    AccessRole = "owner"
	RoleOperator AccessRole = "operator"
	RoleViewer   AccessRole = "viewer"
)
