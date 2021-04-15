package dao

import (
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/cuigh/swirl/dao/bolt"
	"github.com/cuigh/swirl/dao/mongo"
	"github.com/cuigh/swirl/misc"
	"github.com/cuigh/swirl/model"
)

var (
	value = lazy.Value{New: create}
)

// Interface is the interface that wraps all dao methods.
type Interface interface {
	Init()
	Close()

	RoleGet(id string) (*model.Role, error)
	RoleList() (roles []*model.Role, err error)
	RoleCreate(role *model.Role) error
	RoleUpdate(role *model.Role) error
	RoleDelete(id string) error

	UserCreate(user *model.User) error
	UserUpdate(user *model.User) error
	UserList(args *model.UserListArgs) (users []*model.User, count int, err error)
	UserCount() (int, error)
	UserGetByID(id string) (*model.User, error)
	UserGetByName(loginName string) (*model.User, error)
	UserBlock(id string, blocked bool) error
	UserDelete(id string) error

	ProfileUpdateInfo(user *model.User) error
	ProfileUpdatePassword(id, pwd, salt string) error

	SessionUpdate(session *model.Session) error
	SessionGet(token string) (*model.Session, error)

	RegistryCreate(registry *model.Registry) error
	RegistryUpdate(registry *model.Registry) error
	RegistryGet(id string) (*model.Registry, error)
	RegistryList() (registries []*model.Registry, err error)
	RegistryDelete(id string) error

	StackList() (stacks []*model.Stack, err error)
	StackGet(name string) (*model.Stack, error)
	StackCreate(stack *model.Stack) error
	StackUpdate(stack *model.Stack) error
	StackDelete(name string) error
	// StackMigrate migrates stacks from old archive collection. This method will be removed after v0.8.
	StackMigrate()

	TemplateList(args *model.TemplateListArgs) (tpls []*model.Template, count int, err error)
	TemplateGet(id string) (*model.Template, error)
	TemplateCreate(tpl *model.Template) error
	TemplateUpdate(tpl *model.Template) error
	TemplateDelete(id string) error

	EventCreate(event *model.Event) error
	EventList(args *model.EventListArgs) (events []*model.Event, count int, err error)

	PermGet(resType, resID string) (*model.Perm, error)
	PermUpdate(perm *model.Perm) error
	PermDelete(resType, resID string) error

	SettingGet() (setting *model.Setting, err error)
	SettingUpdate(setting *model.Setting) error

	ChartGet(name string) (*model.Chart, error)
	ChartBatch(names ...string) ([]*model.Chart, error)
	ChartList() (charts []*model.Chart, err error)
	ChartCreate(chart *model.Chart) error
	ChartUpdate(chart *model.Chart) error
	ChartDelete(name string) error

	DashboardGet(name, key string) (dashboard *model.ChartDashboard, err error)
	DashboardUpdate(dashboard *model.ChartDashboard) error
}

// Get return a dao instance according to DB_TYPE.
func Get() (Interface, error) {
	v, err := value.Get()
	if err != nil {
		return nil, err
	}
	return v.(Interface), nil
}

func create() (d interface{}, err error) {
	var i Interface
	switch misc.Options.DBType {
	case "", "mongo":
		i, err = mongo.New(misc.Options.DBAddress)
	case "bolt":
		i, err = bolt.New(misc.Options.DBAddress)
	default:
		err = errors.New("Unknown database type: " + misc.Options.DBType)
	}

	if err == nil {
		i.Init()
	}
	return i, err
}
