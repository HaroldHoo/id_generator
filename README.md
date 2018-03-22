# id_generator

## Layout
```
         |<----timestamp---->|<-----ExtraID----->|
layout:  [-------------------|----|------|-------]
         0                   32                  64


         |<InstID>|<-DataID->|<---NextID--->|
ExtraID: [--------|----------|--------------]
         0        8          18             32

InstID:  Machine Instance ID, like IP, 8bit.
DataID:  Data ID, like Mysql table ID Or UserID, 10bit.
NextId:  A Increment Number in the same second, 14bit.
```

## Example
```
import "github.com/HaroldHoo/id_generator"
```
```
// Machine Instance ID, like IP, 8bit
id_generator.DefaultInstanceId = 154

// Data ID, like Mysql table ID, 10bit
var dataID = uint64(256)

// generate
id,_ := id_generator.NextId(dataID)
```
