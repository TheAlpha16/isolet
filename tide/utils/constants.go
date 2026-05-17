package utils

const (
	// Reconciliation Reasons
	ReasonReconciliationFailed = "ReconciliationFailed"

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

	// Label Values
	LabelValInstance       = "instance"
	LabelValTideController = "tide-controller"
	LabelValDeployment     = "deployment"
	TeamDynamic            = "dynamic"

	// Traefik
	EntryPointWeb       = "web"
	EntryPointWebsecure = "websecure"
	RouteKindRule       = "Rule"

	// Condition Types
	ConditionDeploymentReady = "DeploymentReady"
	ConditionServiceReady    = "ServiceReady"
	ConditionIngressReady    = "IngressReady"
	ConditionReady           = "Ready"

	FieldOwner        = "tide-controller"
	InstanceFinalizer = "challenges.isolet.dev/instance-finalizer"
)
