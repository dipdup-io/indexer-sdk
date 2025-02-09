// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go
//
// Generated by this command:
//
//	mockgen -source=storage.go -destination=mock/storage.go -package=mock -typed
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/dipdup-net/indexer-sdk/pkg/storage"
	bun "github.com/uptrace/bun"
	gomock "go.uber.org/mock/gomock"
)

// MockTable is a mock of Table interface.
type MockTable[M storage.Model] struct {
	ctrl     *gomock.Controller
	recorder *MockTableMockRecorder[M]
}

// MockTableMockRecorder is the mock recorder for MockTable.
type MockTableMockRecorder[M storage.Model] struct {
	mock *MockTable[M]
}

// NewMockTable creates a new mock instance.
func NewMockTable[M storage.Model](ctrl *gomock.Controller) *MockTable[M] {
	mock := &MockTable[M]{ctrl: ctrl}
	mock.recorder = &MockTableMockRecorder[M]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTable[M]) EXPECT() *MockTableMockRecorder[M] {
	return m.recorder
}

// CursorList mocks base method.
func (m *MockTable[M]) CursorList(ctx context.Context, id, limit uint64, order storage.SortOrder, cmp storage.Comparator) ([]M, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CursorList", ctx, id, limit, order, cmp)
	ret0, _ := ret[0].([]M)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CursorList indicates an expected call of CursorList.
func (mr *MockTableMockRecorder[M]) CursorList(ctx, id, limit, order, cmp any) *MockTableCursorListCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CursorList", reflect.TypeOf((*MockTable[M])(nil).CursorList), ctx, id, limit, order, cmp)
	return &MockTableCursorListCall[M]{Call: call}
}

// MockTableCursorListCall wrap *gomock.Call
type MockTableCursorListCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableCursorListCall[M]) Return(arg0 []M, arg1 error) *MockTableCursorListCall[M] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableCursorListCall[M]) Do(f func(context.Context, uint64, uint64, storage.SortOrder, storage.Comparator) ([]M, error)) *MockTableCursorListCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableCursorListCall[M]) DoAndReturn(f func(context.Context, uint64, uint64, storage.SortOrder, storage.Comparator) ([]M, error)) *MockTableCursorListCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetByID mocks base method.
func (m *MockTable[M]) GetByID(ctx context.Context, id uint64) (M, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(M)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTableMockRecorder[M]) GetByID(ctx, id any) *MockTableGetByIDCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTable[M])(nil).GetByID), ctx, id)
	return &MockTableGetByIDCall[M]{Call: call}
}

// MockTableGetByIDCall wrap *gomock.Call
type MockTableGetByIDCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableGetByIDCall[M]) Return(arg0 M, arg1 error) *MockTableGetByIDCall[M] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableGetByIDCall[M]) Do(f func(context.Context, uint64) (M, error)) *MockTableGetByIDCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableGetByIDCall[M]) DoAndReturn(f func(context.Context, uint64) (M, error)) *MockTableGetByIDCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// IsNoRows mocks base method.
func (m *MockTable[M]) IsNoRows(err error) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsNoRows", err)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsNoRows indicates an expected call of IsNoRows.
func (mr *MockTableMockRecorder[M]) IsNoRows(err any) *MockTableIsNoRowsCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsNoRows", reflect.TypeOf((*MockTable[M])(nil).IsNoRows), err)
	return &MockTableIsNoRowsCall[M]{Call: call}
}

// MockTableIsNoRowsCall wrap *gomock.Call
type MockTableIsNoRowsCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableIsNoRowsCall[M]) Return(arg0 bool) *MockTableIsNoRowsCall[M] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableIsNoRowsCall[M]) Do(f func(error) bool) *MockTableIsNoRowsCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableIsNoRowsCall[M]) DoAndReturn(f func(error) bool) *MockTableIsNoRowsCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LastID mocks base method.
func (m *MockTable[M]) LastID(ctx context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastID", ctx)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastID indicates an expected call of LastID.
func (mr *MockTableMockRecorder[M]) LastID(ctx any) *MockTableLastIDCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastID", reflect.TypeOf((*MockTable[M])(nil).LastID), ctx)
	return &MockTableLastIDCall[M]{Call: call}
}

// MockTableLastIDCall wrap *gomock.Call
type MockTableLastIDCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableLastIDCall[M]) Return(arg0 uint64, arg1 error) *MockTableLastIDCall[M] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableLastIDCall[M]) Do(f func(context.Context) (uint64, error)) *MockTableLastIDCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableLastIDCall[M]) DoAndReturn(f func(context.Context) (uint64, error)) *MockTableLastIDCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockTable[M]) List(ctx context.Context, limit, offset uint64, order storage.SortOrder) ([]M, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, limit, offset, order)
	ret0, _ := ret[0].([]M)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTableMockRecorder[M]) List(ctx, limit, offset, order any) *MockTableListCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTable[M])(nil).List), ctx, limit, offset, order)
	return &MockTableListCall[M]{Call: call}
}

// MockTableListCall wrap *gomock.Call
type MockTableListCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableListCall[M]) Return(arg0 []M, arg1 error) *MockTableListCall[M] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableListCall[M]) Do(f func(context.Context, uint64, uint64, storage.SortOrder) ([]M, error)) *MockTableListCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableListCall[M]) DoAndReturn(f func(context.Context, uint64, uint64, storage.SortOrder) ([]M, error)) *MockTableListCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Save mocks base method.
func (m_2 *MockTable[M]) Save(ctx context.Context, m M) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Save", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTableMockRecorder[M]) Save(ctx, m any) *MockTableSaveCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTable[M])(nil).Save), ctx, m)
	return &MockTableSaveCall[M]{Call: call}
}

// MockTableSaveCall wrap *gomock.Call
type MockTableSaveCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableSaveCall[M]) Return(arg0 error) *MockTableSaveCall[M] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableSaveCall[M]) Do(f func(context.Context, M) error) *MockTableSaveCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableSaveCall[M]) DoAndReturn(f func(context.Context, M) error) *MockTableSaveCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m_2 *MockTable[M]) Update(ctx context.Context, m M) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTableMockRecorder[M]) Update(ctx, m any) *MockTableUpdateCall[M] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTable[M])(nil).Update), ctx, m)
	return &MockTableUpdateCall[M]{Call: call}
}

// MockTableUpdateCall wrap *gomock.Call
type MockTableUpdateCall[M storage.Model] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTableUpdateCall[M]) Return(arg0 error) *MockTableUpdateCall[M] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTableUpdateCall[M]) Do(f func(context.Context, M) error) *MockTableUpdateCall[M] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTableUpdateCall[M]) DoAndReturn(f func(context.Context, M) error) *MockTableUpdateCall[M] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockTransactable is a mock of Transactable interface.
type MockTransactable struct {
	ctrl     *gomock.Controller
	recorder *MockTransactableMockRecorder
}

// MockTransactableMockRecorder is the mock recorder for MockTransactable.
type MockTransactableMockRecorder struct {
	mock *MockTransactable
}

// NewMockTransactable creates a new mock instance.
func NewMockTransactable(ctrl *gomock.Controller) *MockTransactable {
	mock := &MockTransactable{ctrl: ctrl}
	mock.recorder = &MockTransactableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactable) EXPECT() *MockTransactableMockRecorder {
	return m.recorder
}

// BeginTransaction mocks base method.
func (m *MockTransactable) BeginTransaction(ctx context.Context) (storage.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction", ctx)
	ret0, _ := ret[0].(storage.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockTransactableMockRecorder) BeginTransaction(ctx any) *MockTransactableBeginTransactionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockTransactable)(nil).BeginTransaction), ctx)
	return &MockTransactableBeginTransactionCall{Call: call}
}

// MockTransactableBeginTransactionCall wrap *gomock.Call
type MockTransactableBeginTransactionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactableBeginTransactionCall) Return(arg0 storage.Transaction, arg1 error) *MockTransactableBeginTransactionCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactableBeginTransactionCall) Do(f func(context.Context) (storage.Transaction, error)) *MockTransactableBeginTransactionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactableBeginTransactionCall) DoAndReturn(f func(context.Context) (storage.Transaction, error)) *MockTransactableBeginTransactionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockTransaction) Add(ctx context.Context, model any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockTransactionMockRecorder) Add(ctx, model any) *MockTransactionAddCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockTransaction)(nil).Add), ctx, model)
	return &MockTransactionAddCall{Call: call}
}

// MockTransactionAddCall wrap *gomock.Call
type MockTransactionAddCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionAddCall) Return(arg0 error) *MockTransactionAddCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionAddCall) Do(f func(context.Context, any) error) *MockTransactionAddCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionAddCall) DoAndReturn(f func(context.Context, any) error) *MockTransactionAddCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// BulkSave mocks base method.
func (m *MockTransaction) BulkSave(ctx context.Context, models []any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkSave", ctx, models)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkSave indicates an expected call of BulkSave.
func (mr *MockTransactionMockRecorder) BulkSave(ctx, models any) *MockTransactionBulkSaveCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkSave", reflect.TypeOf((*MockTransaction)(nil).BulkSave), ctx, models)
	return &MockTransactionBulkSaveCall{Call: call}
}

// MockTransactionBulkSaveCall wrap *gomock.Call
type MockTransactionBulkSaveCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionBulkSaveCall) Return(arg0 error) *MockTransactionBulkSaveCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionBulkSaveCall) Do(f func(context.Context, []any) error) *MockTransactionBulkSaveCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionBulkSaveCall) DoAndReturn(f func(context.Context, []any) error) *MockTransactionBulkSaveCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Close mocks base method.
func (m *MockTransaction) Close(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTransactionMockRecorder) Close(ctx any) *MockTransactionCloseCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTransaction)(nil).Close), ctx)
	return &MockTransactionCloseCall{Call: call}
}

// MockTransactionCloseCall wrap *gomock.Call
type MockTransactionCloseCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionCloseCall) Return(arg0 error) *MockTransactionCloseCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionCloseCall) Do(f func(context.Context) error) *MockTransactionCloseCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionCloseCall) DoAndReturn(f func(context.Context) error) *MockTransactionCloseCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// CopyFrom mocks base method.
func (m *MockTransaction) CopyFrom(ctx context.Context, tableName string, data []storage.Copiable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyFrom", ctx, tableName, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// CopyFrom indicates an expected call of CopyFrom.
func (mr *MockTransactionMockRecorder) CopyFrom(ctx, tableName, data any) *MockTransactionCopyFromCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFrom", reflect.TypeOf((*MockTransaction)(nil).CopyFrom), ctx, tableName, data)
	return &MockTransactionCopyFromCall{Call: call}
}

// MockTransactionCopyFromCall wrap *gomock.Call
type MockTransactionCopyFromCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionCopyFromCall) Return(arg0 error) *MockTransactionCopyFromCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionCopyFromCall) Do(f func(context.Context, string, []storage.Copiable) error) *MockTransactionCopyFromCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionCopyFromCall) DoAndReturn(f func(context.Context, string, []storage.Copiable) error) *MockTransactionCopyFromCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Exec mocks base method.
func (m *MockTransaction) Exec(ctx context.Context, query string, params ...any) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range params {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockTransactionMockRecorder) Exec(ctx, query any, params ...any) *MockTransactionExecCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, params...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTransaction)(nil).Exec), varargs...)
	return &MockTransactionExecCall{Call: call}
}

// MockTransactionExecCall wrap *gomock.Call
type MockTransactionExecCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionExecCall) Return(arg0 int64, arg1 error) *MockTransactionExecCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionExecCall) Do(f func(context.Context, string, ...any) (int64, error)) *MockTransactionExecCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionExecCall) DoAndReturn(f func(context.Context, string, ...any) (int64, error)) *MockTransactionExecCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Flush mocks base method.
func (m *MockTransaction) Flush(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockTransactionMockRecorder) Flush(ctx any) *MockTransactionFlushCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockTransaction)(nil).Flush), ctx)
	return &MockTransactionFlushCall{Call: call}
}

// MockTransactionFlushCall wrap *gomock.Call
type MockTransactionFlushCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionFlushCall) Return(arg0 error) *MockTransactionFlushCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionFlushCall) Do(f func(context.Context) error) *MockTransactionFlushCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionFlushCall) DoAndReturn(f func(context.Context) error) *MockTransactionFlushCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HandleError mocks base method.
func (m *MockTransaction) HandleError(ctx context.Context, err error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleError", ctx, err)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleError indicates an expected call of HandleError.
func (mr *MockTransactionMockRecorder) HandleError(ctx, err any) *MockTransactionHandleErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleError", reflect.TypeOf((*MockTransaction)(nil).HandleError), ctx, err)
	return &MockTransactionHandleErrorCall{Call: call}
}

// MockTransactionHandleErrorCall wrap *gomock.Call
type MockTransactionHandleErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionHandleErrorCall) Return(arg0 error) *MockTransactionHandleErrorCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionHandleErrorCall) Do(f func(context.Context, error) error) *MockTransactionHandleErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionHandleErrorCall) DoAndReturn(f func(context.Context, error) error) *MockTransactionHandleErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Rollback mocks base method.
func (m *MockTransaction) Rollback(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTransactionMockRecorder) Rollback(ctx any) *MockTransactionRollbackCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTransaction)(nil).Rollback), ctx)
	return &MockTransactionRollbackCall{Call: call}
}

// MockTransactionRollbackCall wrap *gomock.Call
type MockTransactionRollbackCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionRollbackCall) Return(arg0 error) *MockTransactionRollbackCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionRollbackCall) Do(f func(context.Context) error) *MockTransactionRollbackCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionRollbackCall) DoAndReturn(f func(context.Context) error) *MockTransactionRollbackCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Tx mocks base method.
func (m *MockTransaction) Tx() *bun.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tx")
	ret0, _ := ret[0].(*bun.Tx)
	return ret0
}

// Tx indicates an expected call of Tx.
func (mr *MockTransactionMockRecorder) Tx() *MockTransactionTxCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tx", reflect.TypeOf((*MockTransaction)(nil).Tx))
	return &MockTransactionTxCall{Call: call}
}

// MockTransactionTxCall wrap *gomock.Call
type MockTransactionTxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionTxCall) Return(arg0 *bun.Tx) *MockTransactionTxCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionTxCall) Do(f func() *bun.Tx) *MockTransactionTxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionTxCall) DoAndReturn(f func() *bun.Tx) *MockTransactionTxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Update mocks base method.
func (m *MockTransaction) Update(ctx context.Context, model any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, model)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTransactionMockRecorder) Update(ctx, model any) *MockTransactionUpdateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTransaction)(nil).Update), ctx, model)
	return &MockTransactionUpdateCall{Call: call}
}

// MockTransactionUpdateCall wrap *gomock.Call
type MockTransactionUpdateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockTransactionUpdateCall) Return(arg0 error) *MockTransactionUpdateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockTransactionUpdateCall) Do(f func(context.Context, any) error) *MockTransactionUpdateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockTransactionUpdateCall) DoAndReturn(f func(context.Context, any) error) *MockTransactionUpdateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockModel is a mock of Model interface.
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel.
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance.
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// TableName mocks base method.
func (m *MockModel) TableName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableName")
	ret0, _ := ret[0].(string)
	return ret0
}

// TableName indicates an expected call of TableName.
func (mr *MockModelMockRecorder) TableName() *MockModelTableNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableName", reflect.TypeOf((*MockModel)(nil).TableName))
	return &MockModelTableNameCall{Call: call}
}

// MockModelTableNameCall wrap *gomock.Call
type MockModelTableNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockModelTableNameCall) Return(arg0 string) *MockModelTableNameCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockModelTableNameCall) Do(f func() string) *MockModelTableNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockModelTableNameCall) DoAndReturn(f func() string) *MockModelTableNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockCopiable is a mock of Copiable interface.
type MockCopiable struct {
	ctrl     *gomock.Controller
	recorder *MockCopiableMockRecorder
}

// MockCopiableMockRecorder is the mock recorder for MockCopiable.
type MockCopiableMockRecorder struct {
	mock *MockCopiable
}

// NewMockCopiable creates a new mock instance.
func NewMockCopiable(ctrl *gomock.Controller) *MockCopiable {
	mock := &MockCopiable{ctrl: ctrl}
	mock.recorder = &MockCopiableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCopiable) EXPECT() *MockCopiableMockRecorder {
	return m.recorder
}

// Columns mocks base method.
func (m *MockCopiable) Columns() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Columns")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Columns indicates an expected call of Columns.
func (mr *MockCopiableMockRecorder) Columns() *MockCopiableColumnsCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Columns", reflect.TypeOf((*MockCopiable)(nil).Columns))
	return &MockCopiableColumnsCall{Call: call}
}

// MockCopiableColumnsCall wrap *gomock.Call
type MockCopiableColumnsCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCopiableColumnsCall) Return(arg0 []string) *MockCopiableColumnsCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCopiableColumnsCall) Do(f func() []string) *MockCopiableColumnsCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCopiableColumnsCall) DoAndReturn(f func() []string) *MockCopiableColumnsCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Flat mocks base method.
func (m *MockCopiable) Flat() []any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flat")
	ret0, _ := ret[0].([]any)
	return ret0
}

// Flat indicates an expected call of Flat.
func (mr *MockCopiableMockRecorder) Flat() *MockCopiableFlatCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flat", reflect.TypeOf((*MockCopiable)(nil).Flat))
	return &MockCopiableFlatCall{Call: call}
}

// MockCopiableFlatCall wrap *gomock.Call
type MockCopiableFlatCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockCopiableFlatCall) Return(arg0 []any) *MockCopiableFlatCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockCopiableFlatCall) Do(f func() []any) *MockCopiableFlatCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockCopiableFlatCall) DoAndReturn(f func() []any) *MockCopiableFlatCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
