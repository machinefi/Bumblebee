package types

import (
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrOverTimeLimit      = errors.New("over the time limit")
	ErrOverMaxWorkerID    = errors.New("over the worker id limit")
	ErrInvalidSystemClock = errors.New("invalid system clock")
)

// NewSnowflakeFactory |1|
func NewSnowflakeFactory(bitsWorkerID, bitsSequence, gap uint, base time.Time) *SnowflakeFactory {
	sf := &SnowflakeFactory{
		bitsWorkerID:  bitsWorkerID,
		bitsSequence:  bitsSequence,
		bitsTimestamp: 63 - bitsWorkerID - bitsSequence,
		unit:          time.Duration(gap) * time.Millisecond,
		base:          base,
	}
	sf.init()
	return sf
}

type SnowflakeFactory struct {
	bitsWorkerID  uint
	bitsSequence  uint
	bitsTimestamp uint
	unit          time.Duration
	base          time.Time

	maxWorkerID   uint32
	maxSequence   uint32
	maxTime       time.Time
	baseTimestamp uint64
}

func (f *SnowflakeFactory) init() {
	f.maxSequence = 1<<f.bitsSequence - 1
	f.maxWorkerID = 1<<f.bitsWorkerID - 1
	f.baseTimestamp = f.SnowFlakeTimestamp(f.base)
	maxTimestamp := uint64(1<<f.bitsTimestamp - 1)
	f.maxTime = time.Unix(
		int64(time.Duration(f.baseTimestamp+maxTimestamp)*f.unit/time.Second),
		0,
	)
}

func (f *SnowflakeFactory) MaxWorkerID() uint32 { return f.maxWorkerID }

func (f *SnowflakeFactory) MaxSequence() uint32 { return f.maxSequence }

func (f *SnowflakeFactory) MaxTime() time.Time { return f.maxTime }

func (f *SnowflakeFactory) Sleep(d time.Duration) time.Duration {
	return d*f.unit - time.Duration(time.Now().UnixNano())%f.unit*time.Nanosecond
}

func (f *SnowflakeFactory) Elapsed() uint64 {
	return f.SnowFlakeTimestamp(time.Now()) - f.baseTimestamp
}

func (f *SnowflakeFactory) BuildID(worker, seq uint32, elapsed uint64) (uint64, error) {
	if elapsed >= 1<<f.bitsTimestamp {
		return 0, ErrOverTimeLimit
	}
	return elapsed<<(f.bitsSequence+f.bitsWorkerID) | uint64(seq)<<f.bitsWorkerID | uint64(worker), nil
}

func (f *SnowflakeFactory) SnowFlakeTimestamp(t time.Time) uint64 {
	return uint64(t.UnixNano() / int64(f.unit))
}

func (f *SnowflakeFactory) MaskSequence(sequence uint32) uint32 {
	return sequence & f.maxSequence
}

func (f *SnowflakeFactory) NewSnowflake(workerID uint32) (*Snowflake, error) {
	if workerID > f.maxWorkerID {
		return nil, ErrOverMaxWorkerID
	}
	return &Snowflake{
		f:      f,
		worker: workerID,
		mtx:    &sync.Mutex{},
	}, nil
}

func NewSnowflake(worker uint32) (*Snowflake, error) {
	start, _ := time.Parse(time.RFC3339, "2010-10-24T07:30:06Z07:00")
	return NewSnowflakeFactory(10, 12, 1, start).NewSnowflake(worker)
}

type Snowflake struct {
	f        *SnowflakeFactory
	worker   uint32
	elapsed  uint64
	sequence uint32
	mtx      *sync.Mutex
}

func (s *Snowflake) WorkerID() uint32 { return s.worker }

func (s *Snowflake) ID() (uint64, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	elapsed := s.f.Elapsed()
	if s.elapsed < elapsed {
		s.elapsed = elapsed
		s.sequence = genRandomSequence(9)
		return s.f.BuildID(s.worker, s.sequence, s.elapsed)
	}

	if s.elapsed > elapsed {
		elapsed = s.f.Elapsed()
		if s.elapsed > elapsed {
			return 0, ErrInvalidSystemClock
		}
	}

	s.sequence = s.f.MaskSequence(s.sequence + 1)
	if s.sequence == 0 {
		s.elapsed = s.elapsed + 1
		time.Sleep(s.f.Sleep(time.Duration(s.elapsed - elapsed)))
	}

	return s.f.BuildID(s.worker, s.sequence, s.elapsed)
}

func genRandomSequence(n int32) uint32 {
	return uint32(rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(n))
}

func WorkerIDFromIP(ipv4 net.IP) uint32 {
	if ipv4 == nil {
		return 0
	}
	ip := ipv4.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}

func WorkerIDFromLocalIP() uint32 {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}

	var ipv4 net.IP
	addrs, _ := net.LookupIP(hostname)
	for _, addr := range addrs {
		if ipv4 = addr.To4(); ipv4 != nil {
			break
		}
	}
	return WorkerIDFromIP(ipv4)
}
