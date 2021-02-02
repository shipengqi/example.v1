package dao

import "github.com/shipengqi/example.v1/blog/model"

func (d *dao) GetUserRbac(userid uint) (*model.UserRBAC, error) {

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
		"where blog_group_user.user_id = ?", userid).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ui := &model.UserRBAC{
		Groups: []model.Group{},
		Roles:  []model.Role{},
	}
	for rows.Next() {
		var groupId, roleId uint
		var groupName, roleName string

		if err := rows.Scan(&groupName, &groupId, &roleName, &roleId); err != nil {
			return nil, err
		}
		ui.Roles = append(ui.Roles, model.Role{
			Model: model.Model{
				ID: roleId,
			},
			Name: roleName,
		})

		ui.Groups = append(ui.Groups, model.Group{
			Model: model.Model{
				ID: groupId,
			},
			Name: groupName,
		})
	}
	return ui, nil
}

func (d *dao) GetPermissionsWithRoles(roles []model.Role) ([]model.UserPermission, error) {
	var ups []model.UserPermission
	for i := range roles {
		if roles[i].Deleted || roles[i].Locked {
			continue
		}
		var rp model.RolePermission
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
			ups = append(ups, model.UserPermission{URL: url, Method: method})
		}
		rows.Close()
	}
	return ups, nil
}

func (d *dao) GetUser(username string) (*model.User, error) {
	u := &model.User{}
	if err := d.db.Where("username = ?", username).Find(u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (d *dao) AddUser(name, pass, phone, email string) error {
	u := &model.User{
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
	g := &model.Group{
		Name:        "",
		Description: "",
	}
	if err := d.db.Create(g).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddRole() error {
	r := &model.Role{
		Name:        "",
		Description: "",
	}
	if err := d.db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddPermission() error {

	g := &model.Permission{
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

	r := &model.Resource{
		Name:        "",
		Description: "",
	}

	if err := d.db.Create(r).Error; err != nil {
		return err
	}

	return nil
}

func (d *dao) AddOperation() error {

	o := &model.Operation{
		Name:        "",
		Description: "",
	}

	if err := d.db.Create(o).Error; err != nil {
		return err
	}

	return nil
}
