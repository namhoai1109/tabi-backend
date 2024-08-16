package rbac

import (
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	// Add permission for representative
	r.AddPolicy(model.RoleRepresentative, model.ObjectCompany, model.ActionAny)
	r.AddPolicy(model.RoleRepresentative, model.ObjectBank, model.ActionAny)
	r.AddPolicy(model.RoleRepresentative, model.ObjectBooking, model.ActionCreate)
	r.AddPolicy(model.RoleRepresentative, model.ObjectBranch, model.ActionAny)
	r.AddPolicy(model.RoleRepresentative, model.ObjectBranchManager, model.ActionAny)
	r.AddPolicy(model.RoleRepresentative, model.ObjectGeneralType, model.ActionViewAll)
	r.AddPolicy(model.RoleRepresentative, model.ObjectRoom, model.ActionView)
	r.AddPolicy(model.RoleRepresentative, model.ObjectRoom, model.ActionViewAll)
	r.AddPolicy(model.RoleRepresentative, model.ObjectRoom, model.ActionUpdate)

	// Add permission for host
	r.AddPolicy(model.RoleHost, model.ObjectBank, model.ActionAny)
	r.AddPolicy(model.RoleHost, model.ObjectBooking, model.ActionViewAll)
	r.AddPolicy(model.RoleHost, model.ObjectBooking, model.ActionUpdateAll)
	r.AddPolicy(model.RoleHost, model.ObjectBranch, model.ActionAny)
	r.AddPolicy(model.RoleHost, model.ObjectBranchManager, model.ActionAny)
	r.AddPolicy(model.RoleHost, model.ObjectGeneralType, model.ActionViewAll)
	r.AddPolicy(model.RoleHost, model.ObjectRoom, model.ActionAny)
	r.AddPolicy(model.RoleHost, model.ObjectRoomType, model.ActionAny)

	// Add permission for branch manager
	r.AddPolicy(model.RoleBranchManager, model.ObjectBank, model.ActionViewAll)
	r.AddPolicy(model.RoleBranchManager, model.ObjectBooking, model.ActionViewAll)
	r.AddPolicy(model.RoleBranchManager, model.ObjectBooking, model.ActionUpdateAll)
	r.AddPolicy(model.RoleBranchManager, model.ObjectBranch, model.ActionView)
	r.AddPolicy(model.RoleBranchManager, model.ObjectBranch, model.ActionUpdate)
	r.AddPolicy(model.RoleBranchManager, model.ObjectFacility, model.ActionViewAll)
	r.AddPolicy(model.RoleBranchManager, model.ObjectGeneralType, model.ActionViewAll)
	r.AddPolicy(model.RoleBranchManager, model.ObjectRoom, model.ActionAny)
	r.AddPolicy(model.RoleBranchManager, model.ObjectRoomType, model.ActionAny)

	// Add role permission for admin
	r.AddPolicy(model.RoleAdmin, model.ObjectAny, model.ActionAny)

	r.GetModel().PrintPolicy()

	return r
}
