package api

import (
	"strconv"
	"testing"
	"time"

	"github.com/lyddonb/trajectory/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskDataAccess struct {
	mock.Mock
	db.TaskDataAccess
	pool db.DBPool
}

func (p *MockTaskDataAccess) GetRequests(machine string) ([]string, error) {
	args := p.Mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

func TestListTaskRequestsWithDuplicates(t *testing.T) {
	dal := new(MockTaskDataAccess)

	timestamp := time.Now().UTC().Unix()
	timestamp2 := time.Now().UTC().Unix() + 100
	timestampStr := strconv.FormatInt(timestamp, 10)
	timestampStr2 := strconv.FormatInt(timestamp2, 10)

	requestResult := []string{
		"parentreqeustid",
		timestampStr,
		"parentreqeustid2",
		timestampStr,
		"parentreqeustid",
		timestampStr2,
	}

	dal.On("GetRequests").Return(requestResult, nil)

	result, _ := ListTaskRequests(dal, "")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, result["parentreqeustid"], timestamp2)
	assert.Equal(t, result["parentreqeustid2"], timestamp)
}

func TestListTaskRequestsWithNoItems(t *testing.T) {
	dal := new(MockTaskDataAccess)

	requestResult := []string{}

	dal.On("GetRequests").Return(requestResult, nil)

	result, _ := ListTaskRequests(dal, "")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 0)
}

func TestListTaskRequestsWithSingleItem(t *testing.T) {
	dal := new(MockTaskDataAccess)

	timestamp := time.Now().UTC().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	requestResult := []string{
		"parentreqeustid",
		timestampStr,
	}

	dal.On("GetRequests").Return(requestResult, nil)

	result, _ := ListTaskRequests(dal, "")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["parentreqeustid"], timestamp)
}
