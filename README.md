# id_generator

## example
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

