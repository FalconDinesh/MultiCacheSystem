// Code generated by MockGen. DO NOT EDIT.
// Source: Internal/cache/interface.go

package cache

// import (
//     reflect "reflect"
//     time "time"

//     gomock "github.com/golang/mock/gomock"
// )

// // MockCache is a mock of Cache interface
// type MockCache struct {
//     ctrl     *gomock.Controller
//     recorder *MockCacheMockRecorder
// }

// // MockCacheMockRecorder is the mock recorder for MockCache
// type MockCacheMockRecorder struct {
//     mock *MockCache
// }

// // NewMockCache creates a new mock instance
// func NewMockCache(ctrl *gomock.Controller) *MockCache {
//     mock := &MockCache{ctrl: ctrl}
//     mock.recorder = &MockCacheMockRecorder{mock}
//     return mock
// }

// // EXPECT returns an object that allows the caller to indicate expected use
// func (m *MockCache) EXPECT() *MockCacheMockRecorder {
//     return m.recorder
// }

// // Get mocks base method
// func (m *MockCache) Get(key string) (interface{}, error) {
//     m.ctrl.T.Helper()
//     ret := m.ctrl.Call(m, "Get", key)
//     ret0, _ := ret[0].(interface{})
//     ret1, _ := ret[1].(error)
//     return ret0, ret1
// }

// // Get indicates an expected call of Get
// func (mr *MockCacheMockRecorder) Get(key interface{}) *gomock.Call {
//     mr.mock.ctrl.T.Helper()
//     return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), key)
// }

// // GetWithTTL mocks base method
// func (m *MockCache) GetWithTTL(key string) (interface{}, time.Duration, error) {
//     m.ctrl.T.Helper()
//     ret := m.ctrl.Call(m, "GetWithTTL", key)
//     ret0, _ := ret[0].(interface{})
//     ret1, _ := ret[1].(time.Duration)
//     ret2, _ := ret[2].(error)
//     return ret0, ret1, ret2
// }

// // GetWithTTL indicates an expected call of GetWithTTL
// func (mr *MockCacheMockRecorder) GetWithTTL(key interface{}) *gomock.Call {
//     mr.mock.ctrl.T.Helper()
//     return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithTTL", reflect.TypeOf((*MockCache)(nil).GetWithTTL), key)
// }

// // Set mocks base method
// func (m *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
//     m.ctrl.T.Helper()
//     ret := m.ctrl.Call(m, "Set", key, value, ttl)
//     ret0, _ := ret[0].(error)
//     return ret0
// }

// // Set indicates an expected call of Set
// func (mr *MockCacheMockRecorder) Set(key, value, ttl interface{}) *gomock.Call {
//     mr.mock.ctrl.T.Helper()
//     return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache)(nil).Set), key, value, ttl)
// }

// // Delete mocks base method
// func (m *MockCache) Delete(key string) error {
//     m.ctrl.T.Helper()
//     ret := m.ctrl.Call(m, "Delete", key)
//     ret0, _ := ret[0].(error)
//     return ret0
// }

// // Delete indicates an expected call of Delete
// func (mr *MockCacheMockRecorder) Delete(key interface{}) *gomock.Call {
//     mr.mock.ctrl.T.Helper()
//     return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCache)(nil).Delete), key)
// }

// // ClearAll mocks base method
// // func (m *MockCache) ClearAll() error {
// //     m.ctrl.T.Helper()
// //     ret := m.ctrl.Call(m, "ClearAll")
// //     ret0, _ := ret[0].(error)
// //     return ret0
// // }

// // // ClearAll indicates an expected call of ClearAll
// // func (mr *MockCacheMockRecorder) ClearAll() *gomock.Call {
// //     mr.mock.ctrl.T.Helper()
// //     return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearAll", reflect.TypeOf((*MockCache)(nil).ClearAll))
// // }
