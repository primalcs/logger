# logger

golang logger based on linux syslog

Create logger with options: 
```
package main

import (
	"context"

	"github.com/primalcs/logger"
)

func main() {
	ctx := context.Background()
	lg, err := logger.NewLogger(ctx,
		logger.WithNSQWriter("127.0.0.1:4151", "new_topic"),
		logger.WithLogLevel(logger.DEBUG),
		logger.WithDelimiter(logger.DefaultDelimiter),
	)
	if err != nil {
		panic(err)
	}
	lg.Log(logger.DEBUG, "prefix_tag", "message", "key1", "value1", "key2", "value2")
}

```

For more Options - options.go
