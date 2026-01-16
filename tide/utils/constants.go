package utils

const (
	// Event Reasons
	EventReasonCreated                = "Created"
	EventReasonStatusUpdateFailed     = "StatusUpdateFailed"
	EventReasonDeleting               = "Deleting"
	EventReasonExpired                = "Expired"
	EventReasonInstanceStaged         = "InstanceStaged"
	EventReasonInstanceRunning        = "InstanceRunning"
	EventReasonInstanceFailed         = "InstanceFailed"
	EventReasonDeploymentCreated      = "DeploymentCreated"
	EventReasonDeploymentUpdated      = "DeploymentUpdated"
	EventReasonDeploymentFailed       = "DeploymentCreationFailed"
	EventReasonDeploymentUpdateFailed = "DeploymentUpdateFailed"
	EventReasonServiceCreated         = "ServiceCreated"
	EventReasonServiceUpdated         = "ServiceUpdated"
	EventReasonServiceFailed          = "ServiceCreationFailed"
	EventReasonServiceUpdateFailed    = "ServiceUpdateFailed"
	EventReasonServiceDeleted         = "ServiceDeleted"
	EventReasonServiceDeletionFailed  = "ServiceDeletionFailed"
	EventReasonIngressCreated         = "IngressRouteCreated"
	EventReasonIngressUpdated         = "IngressRouteUpdated"
	EventReasonIngressFailed          = "IngressRouteCreationFailed"
	EventReasonIngressUpdateFailed    = "IngressRouteUpdateFailed"
	EventReasonEndpointsResolved      = "EndpointsResolved"
	EventReasonResolutionFailed       = "EndpointResolutionFailed"

	// Label Keys
	LabelAppName      = "app.kubernetes.io/name"
	LabelAppPartOf    = "app.kubernetes.io/part-of"
	LabelAppManagedBy = "app.kubernetes.io/managed-by"
	LabelAppComponent = "app.kubernetes.io/component"

	LabelChallengeID   = "challenges.isolet.dev/id"
	LabelChallengeSlug = "challenges.isolet.dev/challenge"
	LabelChallengeCID  = "challenges.isolet.dev/challenge-id"
	LabelChallengeType = "challenges.isolet.dev/type"
	LabelTeamID        = "challenges.isolet.dev/team"

	// Condition Types
	ConditionDeploymentReady = "DeploymentReady"
	ConditionServiceReady    = "ServiceReady"
	ConditionIngressReady    = "IngressReady"
	ConditionReady           = "Ready"

	// Field Owner for SSA
	FieldOwner = "tide-controller"
)
