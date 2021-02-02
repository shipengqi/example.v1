package dao

import "github.com/shipengqi/example.v1/blog/model"

type userRow struct {
	GroupId   uint   `json:"group_id"`
	RoleId    uint   `json:"role_id"`
	GroupName string `json:"group_name"`
	RoleName  string `json:"role_name"`
}

func (d *dao) GetUserRbac(userid uint) (*model.UserRBAC, error) {

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
		r := userRow{}
		if err := rows.Scan(&r.GroupName, &r.GroupId, &r.RoleName, &r.RoleId); err != nil {
			return nil, err
		}
		ui.Roles = append(ui.Roles, model.Role{
			Model: model.Model{
				ID: r.RoleId,
			},
			Name: r.RoleName,
		})

		ui.Groups = append(ui.Groups, model.Group{
			Model: model.Model{
				ID: r.GroupId,
			},
			Name: r.GroupName,
		})
	}
	return ui, nil
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