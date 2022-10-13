package id

import (
	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/base/types/snowflake_id"
)

func NewSFIDGenerator() (SFIDGenerator, error) {
	return NewSFIDGeneratorWithWorkerID(snowflake_id.WorkerIDFromLocalIP())
}

func NewSFIDGeneratorWithWorkerID(wid uint32) (SFIDGenerator, error) {
	sf, err := sff.NewSnowflake(wid)
	if err != nil {
		return nil, err
	}
	return &SFIDGeneratorImpl{sf}, nil
}

func MustNewSFIDGenerator() SFIDGenerator {
	g, err := NewSFIDGenerator()
	if err != nil {
		panic(err)
	}
	return g
}

func MustNewSFIDGeneratorWithWorkerID(wid uint32) SFIDGenerator {
	g, err := NewSFIDGeneratorWithWorkerID(wid)
	if err != nil {
		panic(err)
	}
	return g
}

type SFIDGeneratorImpl struct{ *snowflake_id.Snowflake }

func (sfg *SFIDGeneratorImpl) NewSFID() (types.SFID, error) {
	id, err := sfg.ID()
	if err != nil {
		return 0, err
	}
	return types.SFID(id), nil
}

func (sfg *SFIDGeneratorImpl) NewSFIDs(n int) (types.SFIDs, error) {
	var ids = make(types.SFIDs, 0, n)
	for i := 0; i < n; i++ {
		id, err := sfg.ID()
		if err != nil {
			return nil, err
		}
		ids = append(ids, types.SFID(id))
	}
	return ids, nil
}
