# id_generator

## example
```
import "github.com/HaroldHoo/id_generator"
```
```
// machineID, like IP, 8bit
id_generator.DefaultInstanceId = 154

// Init()
id_gen := id_generator.Init()

// dataID, like Mysql table ID, 10bit
var dataID = uint64(256)

// generat
id,_ = id_gen.NextId(dataID)
```

