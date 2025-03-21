// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// UnifiedRepositoryMock is an autogenerated mock type for the UnifiedRepository type
type UnifiedRepositoryMock struct {
	mock.Mock
}

type UnifiedRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *UnifiedRepositoryMock) EXPECT() *UnifiedRepositoryMock_Expecter {
	return &UnifiedRepositoryMock_Expecter{mock: &_m.Mock}
}

// AddLink provides a mock function with given fields: ctx, chatID, link
func (_m *UnifiedRepositoryMock) AddLink(ctx context.Context, chatID int64, link *models.Link) error {
	ret := _m.Called(ctx, chatID, link)

	if len(ret) == 0 {
		panic("no return value specified for AddLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *models.Link) error); ok {
		r0 = rf(ctx, chatID, link)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnifiedRepositoryMock_AddLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddLink'
type UnifiedRepositoryMock_AddLink_Call struct {
	*mock.Call
}

// AddLink is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
//   - link *models.Link
func (_e *UnifiedRepositoryMock_Expecter) AddLink(ctx interface{}, chatID interface{}, link interface{}) *UnifiedRepositoryMock_AddLink_Call {
	return &UnifiedRepositoryMock_AddLink_Call{Call: _e.mock.On("AddLink", ctx, chatID, link)}
}

func (_c *UnifiedRepositoryMock_AddLink_Call) Run(run func(ctx context.Context, chatID int64, link *models.Link)) *UnifiedRepositoryMock_AddLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(*models.Link))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_AddLink_Call) Return(_a0 error) *UnifiedRepositoryMock_AddLink_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UnifiedRepositoryMock_AddLink_Call) RunAndReturn(run func(context.Context, int64, *models.Link) error) *UnifiedRepositoryMock_AddLink_Call {
	_c.Call.Return(run)
	return _c
}

// AddTagToLink provides a mock function with given fields: ctx, chatID, url, tag
func (_m *UnifiedRepositoryMock) AddTagToLink(ctx context.Context, chatID int64, url string, tag string) error {
	ret := _m.Called(ctx, chatID, url, tag)

	if len(ret) == 0 {
		panic("no return value specified for AddTagToLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) error); ok {
		r0 = rf(ctx, chatID, url, tag)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnifiedRepositoryMock_AddTagToLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTagToLink'
type UnifiedRepositoryMock_AddTagToLink_Call struct {
	*mock.Call
}

// AddTagToLink is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
//   - url string
//   - tag string
func (_e *UnifiedRepositoryMock_Expecter) AddTagToLink(ctx interface{}, chatID interface{}, url interface{}, tag interface{}) *UnifiedRepositoryMock_AddTagToLink_Call {
	return &UnifiedRepositoryMock_AddTagToLink_Call{Call: _e.mock.On("AddTagToLink", ctx, chatID, url, tag)}
}

func (_c *UnifiedRepositoryMock_AddTagToLink_Call) Run(run func(ctx context.Context, chatID int64, url string, tag string)) *UnifiedRepositoryMock_AddTagToLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_AddTagToLink_Call) Return(_a0 error) *UnifiedRepositoryMock_AddTagToLink_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UnifiedRepositoryMock_AddTagToLink_Call) RunAndReturn(run func(context.Context, int64, string, string) error) *UnifiedRepositoryMock_AddTagToLink_Call {
	_c.Call.Return(run)
	return _c
}

// CreateUser provides a mock function with given fields: ctx, chatID
func (_m *UnifiedRepositoryMock) CreateUser(ctx context.Context, chatID int64) error {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, chatID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnifiedRepositoryMock_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type UnifiedRepositoryMock_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *UnifiedRepositoryMock_Expecter) CreateUser(ctx interface{}, chatID interface{}) *UnifiedRepositoryMock_CreateUser_Call {
	return &UnifiedRepositoryMock_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, chatID)}
}

func (_c *UnifiedRepositoryMock_CreateUser_Call) Run(run func(ctx context.Context, chatID int64)) *UnifiedRepositoryMock_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_CreateUser_Call) Return(_a0 error) *UnifiedRepositoryMock_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UnifiedRepositoryMock_CreateUser_Call) RunAndReturn(run func(context.Context, int64) error) *UnifiedRepositoryMock_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteLink provides a mock function with given fields: ctx, chatID, url
func (_m *UnifiedRepositoryMock) DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, error) {
	ret := _m.Called(ctx, chatID, url)

	if len(ret) == 0 {
		panic("no return value specified for DeleteLink")
	}

	var r0 *models.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) (*models.Link, error)); ok {
		return rf(ctx, chatID, url)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) *models.Link); ok {
		r0 = rf(ctx, chatID, url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, string) error); ok {
		r1 = rf(ctx, chatID, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnifiedRepositoryMock_DeleteLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteLink'
type UnifiedRepositoryMock_DeleteLink_Call struct {
	*mock.Call
}

// DeleteLink is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
//   - url string
func (_e *UnifiedRepositoryMock_Expecter) DeleteLink(ctx interface{}, chatID interface{}, url interface{}) *UnifiedRepositoryMock_DeleteLink_Call {
	return &UnifiedRepositoryMock_DeleteLink_Call{Call: _e.mock.On("DeleteLink", ctx, chatID, url)}
}

func (_c *UnifiedRepositoryMock_DeleteLink_Call) Run(run func(ctx context.Context, chatID int64, url string)) *UnifiedRepositoryMock_DeleteLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteLink_Call) Return(_a0 *models.Link, _a1 error) *UnifiedRepositoryMock_DeleteLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteLink_Call) RunAndReturn(run func(context.Context, int64, string) (*models.Link, error)) *UnifiedRepositoryMock_DeleteLink_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTagFromLink provides a mock function with given fields: ctx, chatID, url, tag
func (_m *UnifiedRepositoryMock) DeleteTagFromLink(ctx context.Context, chatID int64, url string, tag string) error {
	ret := _m.Called(ctx, chatID, url, tag)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTagFromLink")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, string) error); ok {
		r0 = rf(ctx, chatID, url, tag)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnifiedRepositoryMock_DeleteTagFromLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTagFromLink'
type UnifiedRepositoryMock_DeleteTagFromLink_Call struct {
	*mock.Call
}

// DeleteTagFromLink is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
//   - url string
//   - tag string
func (_e *UnifiedRepositoryMock_Expecter) DeleteTagFromLink(ctx interface{}, chatID interface{}, url interface{}, tag interface{}) *UnifiedRepositoryMock_DeleteTagFromLink_Call {
	return &UnifiedRepositoryMock_DeleteTagFromLink_Call{Call: _e.mock.On("DeleteTagFromLink", ctx, chatID, url, tag)}
}

func (_c *UnifiedRepositoryMock_DeleteTagFromLink_Call) Run(run func(ctx context.Context, chatID int64, url string, tag string)) *UnifiedRepositoryMock_DeleteTagFromLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteTagFromLink_Call) Return(_a0 error) *UnifiedRepositoryMock_DeleteTagFromLink_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteTagFromLink_Call) RunAndReturn(run func(context.Context, int64, string, string) error) *UnifiedRepositoryMock_DeleteTagFromLink_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteUser provides a mock function with given fields: ctx, chatID
func (_m *UnifiedRepositoryMock) DeleteUser(ctx context.Context, chatID int64) (*models.User, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.User, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.User); ok {
		r0 = rf(ctx, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnifiedRepositoryMock_DeleteUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteUser'
type UnifiedRepositoryMock_DeleteUser_Call struct {
	*mock.Call
}

// DeleteUser is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *UnifiedRepositoryMock_Expecter) DeleteUser(ctx interface{}, chatID interface{}) *UnifiedRepositoryMock_DeleteUser_Call {
	return &UnifiedRepositoryMock_DeleteUser_Call{Call: _e.mock.On("DeleteUser", ctx, chatID)}
}

func (_c *UnifiedRepositoryMock_DeleteUser_Call) Run(run func(ctx context.Context, chatID int64)) *UnifiedRepositoryMock_DeleteUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteUser_Call) Return(_a0 *models.User, _a1 error) *UnifiedRepositoryMock_DeleteUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UnifiedRepositoryMock_DeleteUser_Call) RunAndReturn(run func(context.Context, int64) (*models.User, error)) *UnifiedRepositoryMock_DeleteUser_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllActiveLinks provides a mock function with given fields: ctx
func (_m *UnifiedRepositoryMock) GetAllActiveLinks(ctx context.Context) ([]models.Link, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllActiveLinks")
	}

	var r0 []models.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Link, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Link); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnifiedRepositoryMock_GetAllActiveLinks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllActiveLinks'
type UnifiedRepositoryMock_GetAllActiveLinks_Call struct {
	*mock.Call
}

// GetAllActiveLinks is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UnifiedRepositoryMock_Expecter) GetAllActiveLinks(ctx interface{}) *UnifiedRepositoryMock_GetAllActiveLinks_Call {
	return &UnifiedRepositoryMock_GetAllActiveLinks_Call{Call: _e.mock.On("GetAllActiveLinks", ctx)}
}

func (_c *UnifiedRepositoryMock_GetAllActiveLinks_Call) Run(run func(ctx context.Context)) *UnifiedRepositoryMock_GetAllActiveLinks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_GetAllActiveLinks_Call) Return(_a0 []models.Link, _a1 error) *UnifiedRepositoryMock_GetAllActiveLinks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UnifiedRepositoryMock_GetAllActiveLinks_Call) RunAndReturn(run func(context.Context) ([]models.Link, error)) *UnifiedRepositoryMock_GetAllActiveLinks_Call {
	_c.Call.Return(run)
	return _c
}

// GetLinks provides a mock function with given fields: ctx, chatID
func (_m *UnifiedRepositoryMock) GetLinks(ctx context.Context, chatID int64) ([]models.Link, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for GetLinks")
	}

	var r0 []models.Link
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]models.Link, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []models.Link); ok {
		r0 = rf(ctx, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Link)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnifiedRepositoryMock_GetLinks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLinks'
type UnifiedRepositoryMock_GetLinks_Call struct {
	*mock.Call
}

// GetLinks is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *UnifiedRepositoryMock_Expecter) GetLinks(ctx interface{}, chatID interface{}) *UnifiedRepositoryMock_GetLinks_Call {
	return &UnifiedRepositoryMock_GetLinks_Call{Call: _e.mock.On("GetLinks", ctx, chatID)}
}

func (_c *UnifiedRepositoryMock_GetLinks_Call) Run(run func(ctx context.Context, chatID int64)) *UnifiedRepositoryMock_GetLinks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_GetLinks_Call) Return(_a0 []models.Link, _a1 error) *UnifiedRepositoryMock_GetLinks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UnifiedRepositoryMock_GetLinks_Call) RunAndReturn(run func(context.Context, int64) ([]models.Link, error)) *UnifiedRepositoryMock_GetLinks_Call {
	_c.Call.Return(run)
	return _c
}

// GetUser provides a mock function with given fields: ctx, chatID
func (_m *UnifiedRepositoryMock) GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.User, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.User); ok {
		r0 = rf(ctx, chatID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnifiedRepositoryMock_GetUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUser'
type UnifiedRepositoryMock_GetUser_Call struct {
	*mock.Call
}

// GetUser is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
func (_e *UnifiedRepositoryMock_Expecter) GetUser(ctx interface{}, chatID interface{}) *UnifiedRepositoryMock_GetUser_Call {
	return &UnifiedRepositoryMock_GetUser_Call{Call: _e.mock.On("GetUser", ctx, chatID)}
}

func (_c *UnifiedRepositoryMock_GetUser_Call) Run(run func(ctx context.Context, chatID int64)) *UnifiedRepositoryMock_GetUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_GetUser_Call) Return(_a0 *models.User, _a1 error) *UnifiedRepositoryMock_GetUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UnifiedRepositoryMock_GetUser_Call) RunAndReturn(run func(context.Context, int64) (*models.User, error)) *UnifiedRepositoryMock_GetUser_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateLinkFilters provides a mock function with given fields: ctx, chatID, url, filters
func (_m *UnifiedRepositoryMock) UpdateLinkFilters(ctx context.Context, chatID int64, url string, filters map[string][]string) error {
	ret := _m.Called(ctx, chatID, url, filters)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLinkFilters")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, map[string][]string) error); ok {
		r0 = rf(ctx, chatID, url, filters)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnifiedRepositoryMock_UpdateLinkFilters_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateLinkFilters'
type UnifiedRepositoryMock_UpdateLinkFilters_Call struct {
	*mock.Call
}

// UpdateLinkFilters is a helper method to define mock.On call
//   - ctx context.Context
//   - chatID int64
//   - url string
//   - filters map[string][]string
func (_e *UnifiedRepositoryMock_Expecter) UpdateLinkFilters(ctx interface{}, chatID interface{}, url interface{}, filters interface{}) *UnifiedRepositoryMock_UpdateLinkFilters_Call {
	return &UnifiedRepositoryMock_UpdateLinkFilters_Call{Call: _e.mock.On("UpdateLinkFilters", ctx, chatID, url, filters)}
}

func (_c *UnifiedRepositoryMock_UpdateLinkFilters_Call) Run(run func(ctx context.Context, chatID int64, url string, filters map[string][]string)) *UnifiedRepositoryMock_UpdateLinkFilters_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string), args[3].(map[string][]string))
	})
	return _c
}

func (_c *UnifiedRepositoryMock_UpdateLinkFilters_Call) Return(_a0 error) *UnifiedRepositoryMock_UpdateLinkFilters_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *UnifiedRepositoryMock_UpdateLinkFilters_Call) RunAndReturn(run func(context.Context, int64, string, map[string][]string) error) *UnifiedRepositoryMock_UpdateLinkFilters_Call {
	_c.Call.Return(run)
	return _c
}

// NewUnifiedRepositoryMock creates a new instance of UnifiedRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUnifiedRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *UnifiedRepositoryMock {
	mock := &UnifiedRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
