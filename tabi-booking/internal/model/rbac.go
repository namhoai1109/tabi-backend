package model

// RBAC roles
const (
	RoleAdmin          = "ADM"
	RoleRepresentative = "REP"
	RoleBranchManager  = "BMA"
	RoleClient         = "CLI"
	RoleHost           = "HST"
)

// RBAC objects
const (
	ObjectAny           = "*"
	ObjectBank          = "BANK"
	ObjectBooking       = "BOOKING"
	ObjectBranch        = "BRANCH"
	ObjectBranchManager = "BRANCH_MANAGER"
	ObjectCompany       = "COMPANY"
	ObjectFacility      = "FACILITY"
	ObjectGeneralType   = "GENERAL_TYPE"
	ObjectRoom          = "ROOM"
	ObjectRoomType      = "ROOM_TYPE"
)

// RBAC actions
const (
	ActionAny       = "*"
	ActionViewAll   = "view_all"
	ActionView      = "view"
	ActionCreateAll = "create_all"
	ActionCreate    = "create"
	ActionUpdateAll = "update_all"
	ActionUpdate    = "update"
	ActionDeleteAll = "delete_all"
	ActionDelete    = "delete"
	ActionApprove   = "approve"
)
