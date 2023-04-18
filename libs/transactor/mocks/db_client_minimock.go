package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i route256/libs/transactor.DbClient -o ./mocks/db_client_minimock.go -n DbClientMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DbClientMock implements transactor.DbClient
type DbClientMock struct {
	t minimock.Tester

	funcBeginTx          func(ctx context.Context, t1 pgx.TxOptions) (t2 pgx.Tx, err error)
	inspectFuncBeginTx   func(ctx context.Context, t1 pgx.TxOptions)
	afterBeginTxCounter  uint64
	beforeBeginTxCounter uint64
	BeginTxMock          mDbClientMockBeginTx

	funcGetPool          func() (pp1 *pgxpool.Pool)
	inspectFuncGetPool   func()
	afterGetPoolCounter  uint64
	beforeGetPoolCounter uint64
	GetPoolMock          mDbClientMockGetPool
}

// NewDbClientMock returns a mock for transactor.DbClient
func NewDbClientMock(t minimock.Tester) *DbClientMock {
	m := &DbClientMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.BeginTxMock = mDbClientMockBeginTx{mock: m}
	m.BeginTxMock.callArgs = []*DbClientMockBeginTxParams{}

	m.GetPoolMock = mDbClientMockGetPool{mock: m}

	return m
}

type mDbClientMockBeginTx struct {
	mock               *DbClientMock
	defaultExpectation *DbClientMockBeginTxExpectation
	expectations       []*DbClientMockBeginTxExpectation

	callArgs []*DbClientMockBeginTxParams
	mutex    sync.RWMutex
}

// DbClientMockBeginTxExpectation specifies expectation struct of the DbClient.BeginTx
type DbClientMockBeginTxExpectation struct {
	mock    *DbClientMock
	params  *DbClientMockBeginTxParams
	results *DbClientMockBeginTxResults
	Counter uint64
}

// DbClientMockBeginTxParams contains parameters of the DbClient.BeginTx
type DbClientMockBeginTxParams struct {
	ctx context.Context
	t1  pgx.TxOptions
}

// DbClientMockBeginTxResults contains results of the DbClient.BeginTx
type DbClientMockBeginTxResults struct {
	t2  pgx.Tx
	err error
}

// Expect sets up expected params for DbClient.BeginTx
func (mmBeginTx *mDbClientMockBeginTx) Expect(ctx context.Context, t1 pgx.TxOptions) *mDbClientMockBeginTx {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DbClientMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &DbClientMockBeginTxExpectation{}
	}

	mmBeginTx.defaultExpectation.params = &DbClientMockBeginTxParams{ctx, t1}
	for _, e := range mmBeginTx.expectations {
		if minimock.Equal(e.params, mmBeginTx.defaultExpectation.params) {
			mmBeginTx.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmBeginTx.defaultExpectation.params)
		}
	}

	return mmBeginTx
}

// Inspect accepts an inspector function that has same arguments as the DbClient.BeginTx
func (mmBeginTx *mDbClientMockBeginTx) Inspect(f func(ctx context.Context, t1 pgx.TxOptions)) *mDbClientMockBeginTx {
	if mmBeginTx.mock.inspectFuncBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("Inspect function is already set for DbClientMock.BeginTx")
	}

	mmBeginTx.mock.inspectFuncBeginTx = f

	return mmBeginTx
}

// Return sets up results that will be returned by DbClient.BeginTx
func (mmBeginTx *mDbClientMockBeginTx) Return(t2 pgx.Tx, err error) *DbClientMock {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DbClientMock.BeginTx mock is already set by Set")
	}

	if mmBeginTx.defaultExpectation == nil {
		mmBeginTx.defaultExpectation = &DbClientMockBeginTxExpectation{mock: mmBeginTx.mock}
	}
	mmBeginTx.defaultExpectation.results = &DbClientMockBeginTxResults{t2, err}
	return mmBeginTx.mock
}

// Set uses given function f to mock the DbClient.BeginTx method
func (mmBeginTx *mDbClientMockBeginTx) Set(f func(ctx context.Context, t1 pgx.TxOptions) (t2 pgx.Tx, err error)) *DbClientMock {
	if mmBeginTx.defaultExpectation != nil {
		mmBeginTx.mock.t.Fatalf("Default expectation is already set for the DbClient.BeginTx method")
	}

	if len(mmBeginTx.expectations) > 0 {
		mmBeginTx.mock.t.Fatalf("Some expectations are already set for the DbClient.BeginTx method")
	}

	mmBeginTx.mock.funcBeginTx = f
	return mmBeginTx.mock
}

// When sets expectation for the DbClient.BeginTx which will trigger the result defined by the following
// Then helper
func (mmBeginTx *mDbClientMockBeginTx) When(ctx context.Context, t1 pgx.TxOptions) *DbClientMockBeginTxExpectation {
	if mmBeginTx.mock.funcBeginTx != nil {
		mmBeginTx.mock.t.Fatalf("DbClientMock.BeginTx mock is already set by Set")
	}

	expectation := &DbClientMockBeginTxExpectation{
		mock:   mmBeginTx.mock,
		params: &DbClientMockBeginTxParams{ctx, t1},
	}
	mmBeginTx.expectations = append(mmBeginTx.expectations, expectation)
	return expectation
}

// Then sets up DbClient.BeginTx return parameters for the expectation previously defined by the When method
func (e *DbClientMockBeginTxExpectation) Then(t2 pgx.Tx, err error) *DbClientMock {
	e.results = &DbClientMockBeginTxResults{t2, err}
	return e.mock
}

// BeginTx implements transactor.DbClient
func (mmBeginTx *DbClientMock) BeginTx(ctx context.Context, t1 pgx.TxOptions) (t2 pgx.Tx, err error) {
	mm_atomic.AddUint64(&mmBeginTx.beforeBeginTxCounter, 1)
	defer mm_atomic.AddUint64(&mmBeginTx.afterBeginTxCounter, 1)

	if mmBeginTx.inspectFuncBeginTx != nil {
		mmBeginTx.inspectFuncBeginTx(ctx, t1)
	}

	mm_params := &DbClientMockBeginTxParams{ctx, t1}

	// Record call args
	mmBeginTx.BeginTxMock.mutex.Lock()
	mmBeginTx.BeginTxMock.callArgs = append(mmBeginTx.BeginTxMock.callArgs, mm_params)
	mmBeginTx.BeginTxMock.mutex.Unlock()

	for _, e := range mmBeginTx.BeginTxMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.t2, e.results.err
		}
	}

	if mmBeginTx.BeginTxMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmBeginTx.BeginTxMock.defaultExpectation.Counter, 1)
		mm_want := mmBeginTx.BeginTxMock.defaultExpectation.params
		mm_got := DbClientMockBeginTxParams{ctx, t1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmBeginTx.t.Errorf("DbClientMock.BeginTx got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmBeginTx.BeginTxMock.defaultExpectation.results
		if mm_results == nil {
			mmBeginTx.t.Fatal("No results are set for the DbClientMock.BeginTx")
		}
		return (*mm_results).t2, (*mm_results).err
	}
	if mmBeginTx.funcBeginTx != nil {
		return mmBeginTx.funcBeginTx(ctx, t1)
	}
	mmBeginTx.t.Fatalf("Unexpected call to DbClientMock.BeginTx. %v %v", ctx, t1)
	return
}

// BeginTxAfterCounter returns a count of finished DbClientMock.BeginTx invocations
func (mmBeginTx *DbClientMock) BeginTxAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.afterBeginTxCounter)
}

// BeginTxBeforeCounter returns a count of DbClientMock.BeginTx invocations
func (mmBeginTx *DbClientMock) BeginTxBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmBeginTx.beforeBeginTxCounter)
}

// Calls returns a list of arguments used in each call to DbClientMock.BeginTx.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmBeginTx *mDbClientMockBeginTx) Calls() []*DbClientMockBeginTxParams {
	mmBeginTx.mutex.RLock()

	argCopy := make([]*DbClientMockBeginTxParams, len(mmBeginTx.callArgs))
	copy(argCopy, mmBeginTx.callArgs)

	mmBeginTx.mutex.RUnlock()

	return argCopy
}

// MinimockBeginTxDone returns true if the count of the BeginTx invocations corresponds
// the number of defined expectations
func (m *DbClientMock) MinimockBeginTxDone() bool {
	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginTxMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBeginTx != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		return false
	}
	return true
}

// MinimockBeginTxInspect logs each unmet expectation
func (m *DbClientMock) MinimockBeginTxInspect() {
	for _, e := range m.BeginTxMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to DbClientMock.BeginTx with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.BeginTxMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		if m.BeginTxMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to DbClientMock.BeginTx")
		} else {
			m.t.Errorf("Expected call to DbClientMock.BeginTx with params: %#v", *m.BeginTxMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcBeginTx != nil && mm_atomic.LoadUint64(&m.afterBeginTxCounter) < 1 {
		m.t.Error("Expected call to DbClientMock.BeginTx")
	}
}

type mDbClientMockGetPool struct {
	mock               *DbClientMock
	defaultExpectation *DbClientMockGetPoolExpectation
	expectations       []*DbClientMockGetPoolExpectation
}

// DbClientMockGetPoolExpectation specifies expectation struct of the DbClient.GetPool
type DbClientMockGetPoolExpectation struct {
	mock *DbClientMock

	results *DbClientMockGetPoolResults
	Counter uint64
}

// DbClientMockGetPoolResults contains results of the DbClient.GetPool
type DbClientMockGetPoolResults struct {
	pp1 *pgxpool.Pool
}

// Expect sets up expected params for DbClient.GetPool
func (mmGetPool *mDbClientMockGetPool) Expect() *mDbClientMockGetPool {
	if mmGetPool.mock.funcGetPool != nil {
		mmGetPool.mock.t.Fatalf("DbClientMock.GetPool mock is already set by Set")
	}

	if mmGetPool.defaultExpectation == nil {
		mmGetPool.defaultExpectation = &DbClientMockGetPoolExpectation{}
	}

	return mmGetPool
}

// Inspect accepts an inspector function that has same arguments as the DbClient.GetPool
func (mmGetPool *mDbClientMockGetPool) Inspect(f func()) *mDbClientMockGetPool {
	if mmGetPool.mock.inspectFuncGetPool != nil {
		mmGetPool.mock.t.Fatalf("Inspect function is already set for DbClientMock.GetPool")
	}

	mmGetPool.mock.inspectFuncGetPool = f

	return mmGetPool
}

// Return sets up results that will be returned by DbClient.GetPool
func (mmGetPool *mDbClientMockGetPool) Return(pp1 *pgxpool.Pool) *DbClientMock {
	if mmGetPool.mock.funcGetPool != nil {
		mmGetPool.mock.t.Fatalf("DbClientMock.GetPool mock is already set by Set")
	}

	if mmGetPool.defaultExpectation == nil {
		mmGetPool.defaultExpectation = &DbClientMockGetPoolExpectation{mock: mmGetPool.mock}
	}
	mmGetPool.defaultExpectation.results = &DbClientMockGetPoolResults{pp1}
	return mmGetPool.mock
}

// Set uses given function f to mock the DbClient.GetPool method
func (mmGetPool *mDbClientMockGetPool) Set(f func() (pp1 *pgxpool.Pool)) *DbClientMock {
	if mmGetPool.defaultExpectation != nil {
		mmGetPool.mock.t.Fatalf("Default expectation is already set for the DbClient.GetPool method")
	}

	if len(mmGetPool.expectations) > 0 {
		mmGetPool.mock.t.Fatalf("Some expectations are already set for the DbClient.GetPool method")
	}

	mmGetPool.mock.funcGetPool = f
	return mmGetPool.mock
}

// GetPool implements transactor.DbClient
func (mmGetPool *DbClientMock) GetPool() (pp1 *pgxpool.Pool) {
	mm_atomic.AddUint64(&mmGetPool.beforeGetPoolCounter, 1)
	defer mm_atomic.AddUint64(&mmGetPool.afterGetPoolCounter, 1)

	if mmGetPool.inspectFuncGetPool != nil {
		mmGetPool.inspectFuncGetPool()
	}

	if mmGetPool.GetPoolMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetPool.GetPoolMock.defaultExpectation.Counter, 1)

		mm_results := mmGetPool.GetPoolMock.defaultExpectation.results
		if mm_results == nil {
			mmGetPool.t.Fatal("No results are set for the DbClientMock.GetPool")
		}
		return (*mm_results).pp1
	}
	if mmGetPool.funcGetPool != nil {
		return mmGetPool.funcGetPool()
	}
	mmGetPool.t.Fatalf("Unexpected call to DbClientMock.GetPool.")
	return
}

// GetPoolAfterCounter returns a count of finished DbClientMock.GetPool invocations
func (mmGetPool *DbClientMock) GetPoolAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPool.afterGetPoolCounter)
}

// GetPoolBeforeCounter returns a count of DbClientMock.GetPool invocations
func (mmGetPool *DbClientMock) GetPoolBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetPool.beforeGetPoolCounter)
}

// MinimockGetPoolDone returns true if the count of the GetPool invocations corresponds
// the number of defined expectations
func (m *DbClientMock) MinimockGetPoolDone() bool {
	for _, e := range m.GetPoolMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetPoolMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetPoolCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetPool != nil && mm_atomic.LoadUint64(&m.afterGetPoolCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetPoolInspect logs each unmet expectation
func (m *DbClientMock) MinimockGetPoolInspect() {
	for _, e := range m.GetPoolMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to DbClientMock.GetPool")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetPoolMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetPoolCounter) < 1 {
		m.t.Error("Expected call to DbClientMock.GetPool")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetPool != nil && mm_atomic.LoadUint64(&m.afterGetPoolCounter) < 1 {
		m.t.Error("Expected call to DbClientMock.GetPool")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *DbClientMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockBeginTxInspect()

		m.MinimockGetPoolInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *DbClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *DbClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockBeginTxDone() &&
		m.MinimockGetPoolDone()
}