package io

// there should be 3 types:
// LocalIO, DistributedIO, and Cache

type LRUCache struct {
	lruCache   map[string]*SUBCache
	regionSize int32
}

func (c *LRUCache) Init(baseIO IO, fileToRecordSize *map[string]int32, regionSize int32) {
	c.regionSize = regionSize
	c.lruCache = make(map[string]*SUBCache)

	for k, v := range *fileToRecordSize {
		c.lruCache[k] = &SUBCache{}
		c.lruCache[k].init(k, c.regionSize, v, baseIO)
	}
	a := c.lruCache["nodes.store"]
	print(a.maxUse)
}

func (c *LRUCache) ReadBytes(file string, offset, count int32) ([]byte, bool) {
	return c.lruCache[file].ReadBytes(file, offset, count)
}

func (c *LRUCache) WriteBytes(file string, offset int32, bytes *[]byte) bool {
	return c.lruCache[file].WriteBytes(file, offset, bytes)
}

func (c *LRUCache) CreateFile(file string) bool {
	return c.lruCache[file].CreateFile(file)
}

func (c *LRUCache) FileExists(file string) bool {
	return c.lruCache[file].FileExists(file)
}

func (c *LRUCache) DeleteFile(file string) bool {
	return c.lruCache[file].FileExists(file)
}

type SUBCache struct {
	baseIO       IO
	cache        map[int32]*[]byte
	cacheUsage   map[int32]int32
	maxCacheSize int32
	regionSize   int32
	maxUse       int32
	recordSize   int32
}

func (c *SUBCache) init(file string, numOfRecordsInRegion int32, recordSize int32, baseIO IO) {
	c.baseIO = baseIO
	c.cache = make(map[int32]*[]byte)
	c.cacheUsage = make(map[int32]int32)
	c.maxUse = numOfRecordsInRegion
	c.recordSize = recordSize
	c.regionSize = numOfRecordsInRegion * c.recordSize
}

func (c *SUBCache) ReadBytes(file string, offset, count int32) ([]byte, bool) {
	recordId := offset / c.recordSize
	regionId, ok := c.isInCache(recordId)
	if ok {
		return c.getFromCache(regionId, offset), true
	} else {
		//region offset
		regionId = offset / c.regionSize
		ofst := regionId * c.regionSize
		data, isOk := c.baseIO.ReadBytes(file, ofst, c.regionSize)
		if isOk {
			c.addToCache(regionId, data)
			return c.getFromCache(regionId, offset), true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (c *SUBCache) WriteBytes(file string, offset int32, bytes *[]byte) bool {
	ok := c.baseIO.WriteBytes(file, offset, bytes)
	if ok {
		//добавлять регион
		recordId := offset / c.recordSize
		regionId, ok := c.isInCache(recordId)
		if !ok {
			//region offset
			regionId = offset / c.regionSize
			ofst := regionId * c.regionSize
			data, isOk := c.baseIO.ReadBytes(file, ofst, c.regionSize)
			if isOk {
				c.addToCache(regionId, data)
			}
		}
		return true
	}
	return false
}

func (c *SUBCache) CreateFile(file string) bool {
	c.baseIO.CreateFile(file)
	return true
}

func (c *SUBCache) FileExists(file string) bool {
	return c.baseIO.FileExists(file)
}

func (c *SUBCache) DeleteFile(file string) bool {
	return c.baseIO.DeleteFile(file)
}

func (c *SUBCache) findMin() int32 {
	min := c.maxUse
	for k, v := range c.cacheUsage {
		if v < min {
			min = k
		}
	}
	return min
}

func (c *SUBCache) gc() {
	delete(c.cache, c.findMin())
}

func (c *SUBCache) getRegionId(id int32) int32 {
	return id / c.regionSize
}

func (c *SUBCache) isInCache(id int32) (int32, bool) {
	if _, ok := c.cache[id]; ok {
		return c.getRegionId(id), ok
	}
	return -1, false
}

func (c *SUBCache) getFromCache(regionId int32, offset int32) []byte {
	c.cacheUsage[regionId]++
	region := c.cache[regionId]
	position := offset % c.regionSize
	return (*region)[position : position+c.recordSize]
}

func (c *SUBCache) addToCache(regionId int32, data []byte) {
	if int32(len(c.cache)) < c.maxCacheSize {
		c.cache[regionId] = &data
		c.cacheUsage[regionId] = c.maxUse
	} else {
		//c.gc()
		c.cache[regionId] = &data
		c.cacheUsage[regionId] = c.maxUse
	}
	c.decreaseVals()
}

func (c *SUBCache) decreaseVals() {
	for _, v := range c.cacheUsage {
		v--
	}
}
