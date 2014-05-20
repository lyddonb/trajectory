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

func (p *MockTaskDataAccess) GetRequests(address string) ([]string, error) {
	args := p.Mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (p *MockTaskDataAccess) GetRequestTaskKeys(requestId string) (
	[]string, error) {
	args := p.Mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (p *MockTaskDataAccess) GetAddresses() ([]string, error) {
	args := p.Mock.Called()
	return args.Get(0).([]string), args.Error(1)
}

func TestListTaskRequestsWithDuplicates(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

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

	result, _ := api.ListRequests("")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, result["parentreqeustid"], timestamp2)
	assert.Equal(t, result["parentreqeustid2"], timestamp)
}

func TestListTaskRequestsWithNoItems(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	requestResult := []string{}

	dal.On("GetRequests").Return(requestResult, nil)

	result, _ := api.ListRequests("")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 0)
}

func TestListTaskRequestsWithSingleItem(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	timestamp := time.Now().UTC().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	requestResult := []string{
		"parentreqeustid",
		timestampStr,
	}

	dal.On("GetRequests").Return(requestResult, nil)

	result, _ := api.ListRequests("")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["parentreqeustid"], timestamp)
}

func TestListTaskAddressesWithDuplicates(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	timestamp := time.Now().UTC().Unix()
	timestamp2 := time.Now().UTC().Unix() + 100
	timestampStr := strconv.FormatInt(timestamp, 10)
	timestampStr2 := strconv.FormatInt(timestamp2, 10)

	addressResult := []string{
		"addressid",
		timestampStr,
		"addressid2",
		timestampStr,
		"addressid",
		timestampStr2,
	}

	dal.On("GetAddresses").Return(addressResult, nil)

	result, _ := api.ListAddresses()

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, result["addressid"], timestamp2)
	assert.Equal(t, result["addressid2"], timestamp)
}

func TestListTaskAddressesWithNoItems(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	addressResult := []string{}

	dal.On("GetAddresses").Return(addressResult, nil)

	result, _ := api.ListAddresses()

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 0)
}

func TestListTaskAddressesWithSingleItem(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	timestamp := time.Now().UTC().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	addressResult := []string{
		"addressid",
		timestampStr,
	}

	dal.On("GetAddresses").Return(addressResult, nil)

	result, _ := api.ListAddresses()

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["addressid"], timestamp)
}

func TestListRequestTaskKeysWithDuplicates(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	timestamp := time.Now().UTC().Unix()
	timestamp2 := time.Now().UTC().Unix() + 100
	timestampStr := strconv.FormatInt(timestamp, 10)
	timestampStr2 := strconv.FormatInt(timestamp2, 10)

	taskKeyResult := []string{
		"taskkeyid",
		timestampStr,
		"taskkeyid2",
		timestampStr,
		"taskkeyid",
		timestampStr2,
	}

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)

	result, _ := api.ListRequestTaskKeys("requestid")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, result["taskkeyid"], timestamp2)
	assert.Equal(t, result["taskkeyid2"], timestamp)
}

func TestListRequestTaskKeysWithNoItems(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	taskKeyResult := []string{}

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)

	result, _ := api.ListRequestTaskKeys("requestid")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 0)
}

func TestListRequestTaskKeysWithSingleItem(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	timestamp := time.Now().UTC().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)

	taskKeyResult := []string{
		"taskkeyid",
		timestampStr,
	}

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)

	result, _ := api.ListRequestTaskKeys("requestid")

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result["taskkeyid"], timestamp)
}
