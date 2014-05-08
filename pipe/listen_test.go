package pipe

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListener struct {
	mock.Mock
}

func (l *MockListener) Accept() (net.Conn, error) {
	args := l.Mock.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(net.Conn), nil
}

func (l *MockListener) Close() error {
	args := l.Mock.Called()
	return args.Error(0)
}

func (l *MockListener) Addr() net.Addr {
	return nil
}

type MockPipeline struct {
	mock.Mock
	stopped bool
	handled bool
	errored bool
}

func (p *MockPipeline) Handler(connection net.Conn) {
	p.Mock.Called()
	p.stopped = true
	p.handled = true
}

func (p *MockPipeline) Open() bool {
	return !p.stopped
}

func (p *MockPipeline) Error(error) {
	p.Mock.Called()
	p.stopped = true
	p.errored = true
}

func (p *MockPipeline) Parse([]byte) {
	p.stopped = true
	p.handled = true
}

type MockConnection struct {
	mock.Mock
	net.Conn
}

func (c *MockConnection) Close() error {
	args := c.Mock.Called()
	return args.Error(0)
}

func (c *MockConnection) Read(b []byte) (n int, err error) {
	args := c.Mock.Called()
	return args.Int(0), args.Error(1)
}

func TestFailedListen(t *testing.T) {
	listener := new(MockListener)
	pipeline := new(MockPipeline)
	pipeline.stopped = false

	listener.On("Close").Return(nil)
	listener.On("Accept").Return(nil, errors.New("Failure"))
	listener.Mock.AssertNotCalled(t, "Addr")

	pipeline.Mock.AssertNotCalled(t, "Handler")
	pipeline.On("Error").Return()

	Listen(listener, pipeline)

	listener.Mock.AssertExpectations(t)
	pipeline.Mock.AssertExpectations(t)

	assert.True(t, pipeline.errored)
	assert.False(t, pipeline.handled)
}

func TestSuccessfulListen(t *testing.T) {
	listener := new(MockListener)
	pipeline := new(MockPipeline)
	connection := new(MockConnection)
	pipeline.stopped = false

	listener.On("Close").Return(nil)
	listener.On("Accept").Return(connection, nil)
	listener.Mock.AssertNotCalled(t, "Addr")

	pipeline.On("Handler").Return()

	Listen(listener, pipeline)

	listener.Mock.AssertExpectations(t)
	pipeline.Mock.AssertExpectations(t)
	connection.Mock.AssertExpectations(t)

	assert.False(t, pipeline.errored)
	assert.True(t, pipeline.handled)
}

func TestHandleSuccessfullClient(t *testing.T) {
	pipeline := new(MockPipeline)
	pipeline.stopped = false
	connection := new(MockConnection)

	connection.On("Close").Return(nil)
	connection.On("Read").Return(1, nil)

	// Figure out how to make this work.
	//pipeline.On("Parse").Return().Twice()

	handleClient(connection, pipeline)

	connection.Mock.AssertExpectations(t)
	pipeline.Mock.AssertExpectations(t)
}
