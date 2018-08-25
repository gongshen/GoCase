package cmap

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Segment interface {
	Put(p Pair) (bool, error)
	Get(key string) Pair
	GetWithHash(key string, hash uint64) Pair
	Delete(key string) bool
	//获取段内总键值对数量
	Size() uint64
}

type segment struct {
	buckets           []Bucket
	bucketsLen        int
	pairTotal         uint64
	pairRedistributor PairRedistributor
	lock              sync.Mutex
}

func newSegment(bucketNumber int, pairRedistribute PairRedistributor) Segment {
	if bucketNumber <= 0 {
		bucketNumber = DEFAULT_BUCKET_NUMBER
	}
	if pairRedistribute == nil {
		pairRedistribute = newDefaultPairRedistributor(DEFAULT_BUCKET_LOAD_FACTOR, bucketNumber)
	}
	buckets := make([]Bucket, bucketNumber)
	for i := 0; i < bucketNumber; i++ {
		buckets[i] = newBucket()
	}
	return &segment{
		buckets:           buckets,
		bucketsLen:        bucketNumber,
		pairRedistributor: pairRedistribute,
	}
}

func (se *segment) Put(p Pair) (bool, error) {
	se.lock.Lock()
	b := se.buckets[int(p.Hash())%se.bucketsLen]
	ok, err := b.Put(p, nil)
	if ok {
		newTotal := atomic.AddUint64(&se.pairTotal, 1)
		se.redistribute(newTotal, b.Size())
	}
	se.lock.Unlock()
	return ok, err
}

func (se *segment) Get(key string) Pair {
	return se.GetWithHash(key, hash(key))
}

func (se *segment) Delete(key string) bool {
	se.lock.Lock()
	b := se.buckets[int(hash(key))%se.bucketsLen]
	ok := b.Delete(key, nil)
	if ok {
		newTotal := atomic.AddUint64(&se.pairTotal, ^uint64(0))
		se.redistribute(newTotal, b.Size())
	}
	se.lock.Unlock()
	return ok
}

func (se *segment) GetWithHash(key string, keyHash uint64) Pair {
	se.lock.Lock()
	b := se.buckets[int(hash(key))%se.bucketsLen]
	pair := b.Get(key)
	se.lock.Unlock()
	return pair
}

func (se *segment) Size() uint64 {
	return atomic.LoadUint64(&se.pairTotal)
}

func (se *segment) redistribute(pairTotal uint64, bucketSize uint64) (err error) {
	defer func() {
		if p := recover(); p != nil {
			if np, ok := p.(error); ok {
				err = newIllegalRedistributorError(np.Error())
			} else {
				err = newIllegalRedistributorError(fmt.Sprintf("%s", p))
			}
		}
	}()
	se.pairRedistributor.UpdateThreshold(pairTotal, se.bucketsLen)
	bucketStatus := se.pairRedistributor.CheckBucketStatus(pairTotal, bucketSize)
	newBuckets, changed := se.pairRedistributor.Redistribute(bucketStatus, se.buckets)
	if changed {
		se.buckets = newBuckets
		se.bucketsLen = len(se.buckets)
	}
	return nil
}
