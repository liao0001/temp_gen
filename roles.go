package models

import (
	"github.com/liao0001/orm"
	"time"
)

type Role struct {
	ID         int       `json:"id" xorm:"id pk autoincr"`
	Name  string `json:"name" xorm:"name"` //角色名称,
	PermissionIds  []int `json:"permission_ids" xorm:"permission_ids json"` //权限数组,
	CreatedAt  time.Time `json:"created_at" xorm:"created_at created"`
	UpdatedAt  time.Time `json:"updated_at" xorm:"updated_at updated"`
	DeletedAt  time.Time `json:"deleted_at" xorm:"deleted_at deleted"`
}

func (*Role) TableName() string {
	return "roles"
}


type RoleRepo struct {
	*orm.Collection
}

func (r *RoleRepo) GetFields() []FieldInfo {
	return []FieldInfo{
	  {Name: "name", Label: "角色名称", FieldType: "text", IndexWidth: 140, IsDisplay: true, Editable: false},
	  {Name: "permission_ids", Label: "权限数组", FieldType: "text", IndexWidth: 240, IsDisplay: true, Editable: false},
	  {Name: "created_at", Label: "创建时间", FieldType: "datetime", IndexWidth: 140, Editable: true},
	  {Name: "updated_at", Label: "更新时间", FieldType: "datetime", IndexWidth: 140, Editable: true},
	}
}

func (r *RoleRepo) GetAllRoles() ([]Role, error) {
	var arr []Role
	err := r.Where().All(&arr)
	return arr, err
}

func (r *RoleRepo) GetRolesByPage(page int, pageSize int, asc string, desc string, query orm.Cond) ([]Role, error) {
	var arr []Role
	var err error
	if len(asc) > 0 {
		err = r.Where(query).Limit(pageSize).Offset(pageSize * (page - 1)).GroupBy(asc).All(&arr)
	} else if len(desc) > 0 {
		err = r.Where(query).Limit(pageSize).Offset(pageSize * (page - 1)).Desc(desc).All(&arr)
	} else {
		err = r.Where(query).Limit(pageSize).Offset(pageSize * (page - 1)).Desc("id").All(&arr)
	}
	return arr, err
}

func (r *RoleRepo) GetAllRolesByIds(ids []int) ([]Role, error) {
	var arr []Role
	err := r.Where(orm.Cond{"id in": ids}).All(&arr)
	return arr, err
}

func (r *RoleRepo) GetAllRoleMap() (map[int]Role, error) {
	var m map[int]Role
	err := r.Where().All(&m)
	return m, err
}

func (r *RoleRepo) GetRoleMapByIds(ids []int) (map[int]Role, error) {
	var m map[int]Role
	err := r.Where(orm.Cond{"id in": ids}).All(&m)
	return m, err
}

func (r *RoleRepo) GetRoleById(id int) (*Role, error) {
	var o Role
	err := r.ID(id).Get(&o)
	return &o, err
}

func (r *RoleRepo) InsertRole(role *Role) (int, error) {
	idInter, err := r.Insert(role)
	if err != nil {
		return 0, err
	}
	return (int)(idInter.(int64)), nil
}
