package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	ID                   = "id"
	MACHINE_INFO         = "machine_info"
	STATS                = "stats"
	STAT_PREFIX          = "Stat"
	STAT_REQUEST_PREFIX  = "StatRequests"
	STAT_ADDRESS_REQUEST = "StatAddressRequest"
	STAT_ADDRESS         = "StatAddress"
)

type StatDAL interface {
	SaveRequestStats(request *Request, stats Stat) (string, error)
	GetStatForRequest(requestId string) (Stat, error)
}

type StatDataAccess struct {
	pool DBPool
}

func NewStatDataAccess(pool DBPool) *StatDataAccess {
	return &StatDataAccess{pool}
}

type Stat map[string]string

type Request struct {
	RequestId   string `redis:request_id`
	MachineInfo string `redis:machine_info`
	Url         string `redis:url`
}

func getPrefixKey(prefix, requestId string) string {
	return fmt.Sprintf("%s:%s", prefix, requestId)
}

func (s *StatDataAccess) SaveRequestStats(request *Request, stats Stat) (string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	err := AddStat(conn, request.RequestId, stats)

	if err != nil {
		return timestamp, err
	}

	err = AddRequest(conn, request)

	if err != nil {
		return timestamp, err
	}

	// TODO: Get the address from the machine info.
	splitKey := strings.Split(request.MachineInfo, ".")
	address := splitKey[0]

	err = AddAddressRequest(conn, address, timestamp, request.RequestId)

	if err != nil {
		return timestamp, err
	}

	err = AddStatAddress(conn, address, timestamp)

	if err != nil {
		return timestamp, err
	}

	err = conn.Flush()

	if err != nil {
		return "", err
	}

	return timestamp, nil
}

func AddStat(conn redis.Conn, requestId string, stat Stat) error {
	return conn.Send("HMSET", redis.Args{}.Add(getPrefixKey(STAT_PREFIX, requestId)).AddFlat(stat)...)
}

func AddRequest(conn redis.Conn, request *Request) error {
	return conn.Send("HMSET", redis.Args{}.Add(
		getPrefixKey(STAT_REQUEST_PREFIX, request.RequestId)).AddFlat(request)...)
}

func AddAddressRequest(conn redis.Conn, address, timestamp, requestId string) error {
	return conn.Send("ZADD", getPrefixKey(STAT_ADDRESS_REQUEST, address), timestamp, requestId)
}

func AddStatAddress(conn redis.Conn, address, timestamp string) error {
	return conn.Send("ZADD", STAT_ADDRESS, timestamp, address)
}

func (s *StatDataAccess) GetStatForRequest(requestId string) (Stat, error) {
	conn := s.pool.Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", getPrefixKey(STAT_PREFIX, requestId)))

	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	return ScanMap(v)
}
