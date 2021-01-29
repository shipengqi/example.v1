package dao

import "github.com/shipengqi/example.v1/blog/model"

func (d *dao) AddUser() error {
	u := &model.User{
		Username: "",
		Password: "",
		Phone:    "",
		Email:    "",
		Sex:      0,
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
		Role:        "",
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