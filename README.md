# id_generator

## Layout
```
         |<----timestamp---->|<-----ExtraID----->|
layout:  [-------------------|----|------|-------]
         64                  32                  0


         |<InstID>|<-DataID->|<---NextID--->|
ExtraID: [--------|----------|--------------]
         32       24         14             0

InstID:  Machine Instance ID, like IP, 8bit.
DataID:  Data ID, like Mysql table ID Or UserID, 10bit.
NextId:  A Increment Number in the same second, 14bit.
```

## Example
```
import "github.com/HaroldHoo/id_generator"
```
```
// Set the cacheFile
id_generator.SetDefaultCacheFile("/dev/shm/id_generator")

// Machine Instance ID, like IP, 8bit
id_generator.DefaultInstanceId = 154

// Data ID, like Mysql table ID, 10bit
var dataID = uint64(256)

// generate
id,_ := id_generator.NextId(dataID)
```

