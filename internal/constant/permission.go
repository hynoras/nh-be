package constant

type ManageAction string

const (
	Create ManageAction = "create"
	Update ManageAction = "update"
	Delete ManageAction = "delete"
)

const (
	ViewUser              = "user:view"
	ManageUser            = "user:manage"
	ViewPermissionGroup   = "permission_group:view"
	ManagePermissionGroup = "permission_group:manage"
	ViewExperiment        = "experiment:view"
	ManageExperiment      = "experiment:manage"
)
