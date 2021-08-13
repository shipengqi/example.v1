package dao

import (
	model2 "github.com/shipengqi/example.v1/apps/blog/model"
)

func (d *dao) GetUserRbac(userid uint) (*model2.UserRBAC, error) {

	// var gu model.GroupUser
	// rows, err := d.db.Model(gu).
	// 	Select(
	// 		"blog_group.name as group_name",
	// 		"blog_group.id as group_id",
	// 		"blog_role.name as role_name",
	// 		"blog_role.id as role_id",
	// 	).
	// 	Joins("JOIN blog_group ON blog_group.id = blog_group_user.group_id").
	// 	Joins("JOIN blog_group_role ON blog_group_role.group_id = blog_group.id").
	// 	Joins("JOIN blog_role ON blog_role.id = blog_group_role.role_id").
	// 	Where("blog_group_user.user_id = ?", userid).Rows()

	// Query Raw SQL
	rows, err := d.db.Raw("select blog_group.name as group_name, "+
		"blog_group.id as group_id, "+
		"blog_role.name as role_name, "+
		"blog_role.id as role_id from blog_group_user "+
		"join blog_group ON blog_group.id = blog_group_user.group_id "+
		"join blog_group_role ON blog_group_role.group_id = blog_group.id "+
		"join blog_role ON blog_role.id = blog_group_role.role_id "+
		"where blog_group_user.user_id = ? and " +
		"blog_group.locked = false and " +
		"blog_group.deleted = false and " +
		"blog_role.locked = false and " +
		"blog_role.deleted = false", userid).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ui := &model2.UserRBAC{
		Groups: []model2.Group{},
		Roles:  []model2.Role{},
	}
	for rows.Next() {
		var groupId, roleId uint
		var groupName, roleName string

		if err := rows.Scan(&groupName, &groupId, &roleName, &roleId); err != nil {
			return nil, err
		}
		ui.Roles = append(ui.Roles, model2.Role{
			Model: model2.Model{
				ID: roleId,
			},
			Name: roleName,
		})

		ui.Groups = append(ui.Groups, model2.Group{
			Model: model2.Model{
				ID: groupId,
			},
			Name: groupName,
		})
	}
	return ui, nil
}

func (d *dao) GetPermissionsWithRoles(roles []model2.Role) ([]model2.UserPermission, error) {
	var ups []model2.UserPermission
	for i := range roles {
		if roles[i].Deleted || roles[i].Locked {
			continue
		}
		var rp model2.RolePermission
		rows, err := d.db.Model(rp).
			Select(
				"blog_resource.url as resource_url",
				"blog_operation.name as operation_name",
			).
			Joins("JOIN blog_permission ON blog_permission.id = blog_role_permission.permission_id").
			Joins("JOIN blog_resource ON blog_resource.id = blog_permission.resource_id").
			Joins("JOIN blog_operation ON blog_operation.id = blog_permission.operation_id").
			Where("blog_role_permission.role_id = ?", roles[i].ID).Rows()
		if err != nil {
			return nil, err
		}


		for rows.Next() {
			var url string
			var method string
			if err := rows.Scan(&url, &method); err != nil {
				rows.Close()
				return nil, err
			}
			ups = append(ups, model2.UserPermission{URL: url, Method: method})
		}
		rows.Close()
	}
	return ups, nil
}

func (d *dao) GetUser(username string) (*model2.User, error) {
	u := &model2.User{}
	if err := d.db.Where("username = ?", username).Find(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (d *dao) AddUser(name, pass, phone, email string) error {
	u := &model2.User{
		Username: name,
		Password: pass,
		Phone:    phone,
		Email:    email,
	}
	if err := d.db.Create(u).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddGroup() error {
	g := &model2.Group{
		Name:        "",
		Description: "",
	}
	if err := d.db.Create(g).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddRole() error {
	r := &model2.Role{
		Name:        "",
		Description: "",
	}
	if err := d.db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddPermission() error {

	g := &model2.Permission{
		ResourceId:  0,
		OperationId: 0,
		Name:        "",
		Description: "",
	}

	if err := d.db.Create(g).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddResource() error {

	r := &model2.Resource{
		Name:        "",
		Description: "",
	}

	if err := d.db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddOperation() error {

	o := &model2.Operation{
		Name:        "",
		Description: "",
	}

	if err := d.db.Create(o).Error; err != nil {
		return err
	}

	return nil
}
