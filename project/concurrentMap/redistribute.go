package cmap

import (
	"fmt"
	"sync/atomic"
)

type BucketStatus uint8

const (
	BUCKET_STATUS_NORMAL     BucketStatus = 0
	BUCKET_STATUS_OVERWEIGHT BucketStatus = 1
	BUCKET_STATUS_UDERWEIGHT BucketStatus = 2
)

type PairRedistributor interface {
	UpdateThreshold(pairTotal uint64, bucketNumber int)
	CheckBucketStatus(pairTatal uint64, bucketSize uint64) (bucketStatus BucketStatus)
	Redistribute(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool)
}

type pairRedistributor struct {
	loadFactor           float64
	upperThreshold       uint64
	overwightBucketCount uint64
	emptyBucketCount     uint64
}

func newDefaultPairRedistributor(loadFactor float64, bucketNumber int) PairRedistributor {
	if loadFactor <= 0 {
		loadFactor = DEFAULT_BUCKET_LOAD_FACTOR
	}
	pr := &pairRedistributor{}
	pr.loadFactor = loadFactor
	pr.UpdateThreshold(0, bucketNumber)
	return pr
}

var bucketCountTemplate = `Bucket count:
	pairTotal:%d
	bucketNumber:%d
	average:%f
	upperThreshold:%d
`

func (pr *pairRedistributor) UpdateThreshold(pairTotal uint64, bucketNumber int) {
	var average float64
	average = float64(pairTotal / uint64(bucketNumber))
	if average < 100 {
		average = 100
	}
	defer func() {
		fmt.Sprintf(bucketCountTemplate,
			pairTotal,
			bucketNumber,
			average,
			atomic.LoadUint64(&pr.upperThreshold))
	}()
	atomic.StoreUint64(&pr.upperThreshold, uint64(average*pr.loadFactor))
}

var bucketStatusTemplate = `Bucket status:
	pairTotal:%d
	bucketSize:%d
	upperthreshold:%d
	overweightBucketCount:%d
	emptyBucketCount:%d
	bucketStatus:%d
`

func (pr *pairRedistributor) CheckBucketStatus(pairTotal uint64, bucketSize uint64) (bucketStatus BucketStatus) {
	defer func() {
		fmt.Sprintf(bucketStatusTemplate,
			pairTotal,
			bucketSize,
			atomic.LoadUint64(&pr.upperThreshold),
			atomic.LoadUint64(&pr.overwightBucketCount),
			atomic.LoadUint64(&pr.emptyBucketCount),
			bucketStatus)
	}()
	if bucketSize > DEFAULT_BUCKET_MAX_SIZE || bucketSize >= atomic.LoadUint64(&pr.upperThreshold) {
		atomic.AddUint64(&pr.overwightBucketCount, 1)
		bucketStatus = BUCKET_STATUS_OVERWEIGHT
		return
	}
	if bucketSize == 0 {
		atomic.AddUint64(&pr.emptyBucketCount, 1)
		bucketStatus = BUCKET_STATUS_UDERWEIGHT
		return
	}
	return
}

var redistributionTemplate = `Redistribution:
	bucketStatus:%d
	currentNumber:%d
	newNumber:%d
`

func (pr *pairRedistributor) Redistribute(bucketStatus BucketStatus, buckets []Bucket) (newBuckets []Bucket, changed bool) {
	currentNumber := uint64(len(buckets))
	newNumber := currentNumber
	switch bucketStatus {
	case BUCKET_STATUS_OVERWEIGHT:
		if atomic.LoadUint64(&pr.overwightBucketCount)*4 < currentNumber {
			return nil, false
		}
		newNumber = currentNumber << 1
	case BUCKET_STATUS_UDERWEIGHT:
		if atomic.LoadUint64(&pr.emptyBucketCount)*4 < currentNumber || currentNumber < 100 {
			return nil, false
		}
		newNumber = currentNumber >> 1
		if newNumber < 2 {
			newNumber = 2
		}
	default:
		return nil, false
	}
	if currentNumber == newNumber {
		atomic.StoreUint64(&pr.emptyBucketCount, 0)
		atomic.StoreUint64(&pr.overwightBucketCount, 0)
		return nil, false
	}
	var pairs []Pair
	for _, b := range buckets {
		for p := b.GetFirstPair(); p != nil; p = p.Next() {
			pairs = append(pairs, p)
		}
	}
	if newNumber > currentNumber {
		for i := uint64(0); i < currentNumber; i++ {
			buckets[i].Clear(nil)
		}
		for j := newNumber - currentNumber; j > 0; j-- {
			buckets = append(buckets, newBucket())
		}
	} else {
		buckets := make([]Bucket, newNumber)
		for i := uint64(0); i < newNumber; i++ {
			buckets[i] = newBucket()
		}
	}
	for _, p := range pairs {
		index := int(p.Hash() % newNumber)
		b := buckets[index]
		b.Put(p, nil)
	}
	atomic.StoreUint64(&pr.overwightBucketCount, 0)
	atomic.StoreUint64(&pr.emptyBucketCount, 0)
	return buckets, true
}
