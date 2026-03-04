package constant

import "errors"

type ErrorDetail error
type ErrorMessage string

const (
	ErrInvalidSortOrder = "Invalid sort order"

	//auth
	ErrAuthorizationFailed  = "Authorization failed"
	ErrVerifyTokenFailed    = "Failed to verify token"
	ErrSignUpFailed         = "Failed to sign up"
	ErrLoginFailed          = "Failed to login"
	ErrLogoutFailed         = "Failed to logout"
	ErrChangePasswordFailed = "Failed to change password"

	//permission
	ErrGetAllPermissionFailed         = "Failed to get permissions"
	ErrGetPermissionDetailFailed      = "Failed to get permission detail"
	ErrCreatePermissionFailed         = "Failed to create permission"
	ErrUpdatePermissionFailed         = "Failed to update permission"
	ErrDeletePermissionFailed         = "Failed to delete permission"
	ErrGetAllPermissionGroupFailed    = "Failed to get permission groups"
	ErrGetPermissionGroupDetailFailed = "Failed to get permission group detail"
	ErrCreatePermissionGroupFailed    = "Failed to create permission group"
	ErrUpdatePermissionGroupFailed    = "Failed to update permission group"
	ErrDeletePermissionGroupFailed    = "Failed to delete permission group"

	//user
	ErrGetAllUsersFailed   = "Failed to get users"
	ErrGetUserDetailFailed = "Failed to get user detail"
	ErrCreateUserFailed    = "Failed to create user"
	ErrUpdateUserFailed    = "Failed to update user"
	ErrDeleteUserFailed    = "Failed to delete user"

	//experiment
	ErrGetAllExperimentFailed            = "Failed to get experiments"
	ErrGetExperimentDetailFailed         = "Failed to get experiment detail"
	ErrCreateExperimentFailed            = "Failed to create experiment"
	ErrUpdateExperimentFailed            = "Failed to update experiment"
	ErrDeleteExperimentFailed            = "Failed to delete experiment"
	ErrInvalidStatusTransition           = "Invalid status transition"
	ErrAssignProcedureToExperimentFailed = "Failed to assign procedure to experiment"

	//experiment result
	ErrGetExperimentResultDetailFailed = "Failed to get experiment result detail"
	ErrCreateExperimentResultFailed    = "Failed to create experiment result"
	ErrUpdateExperimentResultFailed    = "Failed to update experiment result"
	ErrDeleteExperimentResultFailed    = "Failed to delete experiment result"

	//procedure
	ErrGetAllProceduresFailed    = "Failed to get procedures"
	ErrGetProcedureDetailFailed  = "Failed to get procedure detail"
	ErrCreateProcedureFailed     = "Failed to create procedure"
	ErrUpdateProcedureFailed     = "Failed to update procedure"
	ErrUpdateProcedureStepFailed = "Failed to update procedure step"
	ErrDeleteProcedureFailed     = "Failed to delete procedure"

	//observation
	ErrGetAllObservationsFailed = "Failed to get observations"
	ErrGetProcedureStepsFailed  = "Failed to get procedure steps"
	ErrCreateObservationFailed  = "Failed to create observation"
)

var (
	ErrSortOrderShouldBeASCOrDESC ErrorDetail = errors.New("sort order should be ASC or DESC")
	ErrInvalidIDFormat            ErrorDetail = errors.New("invalid id format")

	// auth
	ErrInvalidCredentials        ErrorDetail = errors.New("invalid credentials")
	ErrVerificationTokenNotFound ErrorDetail = errors.New("token not found")
	ErrUnauthenticated           ErrorDetail = errors.New("unauthenticated")
	ErrEmailAlreadyExists        ErrorDetail = errors.New("email already exists")
	ErrVerificationTokenExpired  ErrorDetail = errors.New("token expired")
	ErrInvalidVerificationToken  ErrorDetail = errors.New("invalid token, please check token type")
	ErrSessionNotFound           ErrorDetail = errors.New("session not found")

	//permission
	ErrPermissionNotFound               ErrorDetail = errors.New("permission not found")
	ErrPermissionGroupNotFound          ErrorDetail = errors.New("permission group not found")
	ErrNotNullPermissions               ErrorDetail = errors.New("permissions can not be null")
	ErrCannotDeleteSuperAdmin           ErrorDetail = errors.New("can not delete Super Admin. At least one must exist")
	ErrForbidViewPermissions            ErrorDetail = errors.New("you are not allowed to view permissions")
	ErrPermissionGroupNameAlreadyExists ErrorDetail = errors.New("permission group name already exists")
	ErrForbidViewPermissionGroups       ErrorDetail = errors.New("you are not allowed to view permission groups")
	ErrForbidViewPermissionGroup        ErrorDetail = errors.New("you are not allowed to view permission group")
	ErrForbidCreatePermissionGroup      ErrorDetail = errors.New("you are not allowed to create permission group")
	ErrForbidUpdatePermissionGroup      ErrorDetail = errors.New("you are not allowed to update permission group")
	ErrForbidDeletePermissionGroup      ErrorDetail = errors.New("you are not allowed to delete permission group")

	//user
	ErrUserNotFound                      ErrorDetail = errors.New("user not found")
	ErrDuplicateUsername                 ErrorDetail = errors.New("username already exists")
	ErrDuplicateEmail                    ErrorDetail = errors.New("email already exists")
	ErrForbidCreateUser                  ErrorDetail = errors.New("you do not have permission to create users")
	ErrForbidViewUsers                   ErrorDetail = errors.New("you do not have permission to view users")
	ErrForbidViewUser                    ErrorDetail = errors.New("you do not have permission to view this user")
	ErrForbidUpdateUser                  ErrorDetail = errors.New("you do not have permission to update this user")
	ErrForbidDeleteUser                  ErrorDetail = errors.New("you do not have permission to delete this user")
	ErrInvalidUsernameLength             ErrorDetail = errors.New("username must be between 3 and 30 characters")
	ErrInvalidUsernameChars              ErrorDetail = errors.New("username can only contain letters, numbers, dots, and underscores")
	ErrUsernameMustStartWithLetter       ErrorDetail = errors.New("username must start with a letter")
	ErrUsernameMustEndWithLetterOrNumber ErrorDetail = errors.New("username must end with a letter or number")
	ErrUsernameNoConsecutiveSpecialChars ErrorDetail = errors.New("username cannot contain consecutive dots or underscores")
	ErrUsernameNoAdjacentSpecialChars    ErrorDetail = errors.New("username cannot contain adjacent dots and underscores")
	ErrReservedUsername                  ErrorDetail = errors.New("username is reserved and cannot be used")

	// experiment
	ErrExperimentNotFound                              ErrorDetail = errors.New("experiment not found")
	ErrForbidCreateExperiment                          ErrorDetail = errors.New("you do not have permission to create experiments")
	ErrForbidViewExperiments                           ErrorDetail = errors.New("you do not have permission to view experiments")
	ErrForbidViewExperiment                            ErrorDetail = errors.New("you do not have permission to view this experiment")
	ErrForbidUpdateExperiment                          ErrorDetail = errors.New("you do not have permission to update this experiment")
	ErrForbidDeleteExperiment                          ErrorDetail = errors.New("you do not have permission to delete this experiment")
	ErrForbidAssignProcedureToExperiment               ErrorDetail = errors.New("you do not have permission to assign procedure to experiment")
	ErrStatusTransitionFromDraftToPlanning             ErrorDetail = errors.New("Invalid status transition, only draft can be transition to planning")
	ErrStatusTransitionFromPlanningToRunning           ErrorDetail = errors.New("Invalid status transition, only planning can be transition to running")
	ErrStatusTransitionFromRunningToCompletedOrAborted ErrorDetail = errors.New("Invalid status transition, only running can be transition to completed or aborted")
	ErrExperimentConflict                              ErrorDetail = errors.New("the experiment was modified by another request, please retry")
	ErrExperimentAlreadyInTargetState                  ErrorDetail = errors.New("experiment is already in target state")
	ErrDuplicateProcedureAssignment                    ErrorDetail = errors.New("This procedure is already assigned to the experiment")

	//experiment result
	ErrExperimentResultNotFound      ErrorDetail = errors.New("experiment result not found")
	ErrExperimentResultAlreadyExists ErrorDetail = errors.New("experiment result already exists for this experiment")
	ErrForbidCreateExperimentResult  ErrorDetail = errors.New("you do not have permission to create experiment results")
	ErrForbidViewExperimentResult    ErrorDetail = errors.New("you do not have permission to view this experiment result")
	ErrForbidUpdateExperimentResult  ErrorDetail = errors.New("you do not have permission to update this experiment result")
	ErrInvalidOutcome                ErrorDetail = errors.New("invalid outcome value")
	ErrInvalidConfidenceLevel        ErrorDetail = errors.New("invalid confidence level value")
	ErrExperimentResultConflict      ErrorDetail = errors.New("the experiment result was modified by another request, please retry")

	//procedure
	ErrProcedureNotFound      ErrorDetail = errors.New("procedure not found")
	ErrProcedureAlreadyExists ErrorDetail = errors.New("procedure already exists")
	ErrProcedureConflict      ErrorDetail = errors.New("This procedure is modified by another request, please retry")
	ErrForbidViewProcedure    ErrorDetail = errors.New("you do not have permission to view this procedure")
	ErrForbidCreateProcedure  ErrorDetail = errors.New("you do not have permission to create this procedure")
	ErrForbidUpdateProcedure  ErrorDetail = errors.New("you do not have permission to update this procedure")
	ErrForbidDeleteProcedure  ErrorDetail = errors.New("you do not have permission to delete this procedure")
	ErrProcedureStepNotFound  ErrorDetail = errors.New("procedure step not found")
	ErrProcedureStepConflict  ErrorDetail = errors.New("This procedure step is modified by another request, please retry")

	//observation
	ErrObservationNotFound     ErrorDetail = errors.New("observation not found")
	ErrForbidViewObservation   ErrorDetail = errors.New("you do not have permission to view observation")
	ErrForbidCreateObservation ErrorDetail = errors.New("you do not have permission to create this observation")
)
