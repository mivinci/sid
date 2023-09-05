# sid

A bit.ly-like tiny id generator porting from [code.activestate.com/recipes/576918](http://code.activestate.com/recipes/576918) which uses a bit-shuffling approach to avoid generating consecutive, predictable URLs.  However, 
the algorithm is deterministic and guarantees that no collisions will occur.

The intended use is that incrementing, consecutive integers will be used as 
keys to generate the short URLs.  For example, when creating a new URL, the 
unique integer ID assigned by a database could be used to generate the URL 
by using this module.  Or a simple counter may be used.  As long as the same 
integer is not used twice, the same short URL will not be generated twice.


## Example

```go
import (
    "fmt"
    "github.com/mivinci/sid"
)

func main() {
    sid.Encode(12)      // yNvD
    sid.Decode("yNvD")  // 12
}
```

## LICENSE

sid is MIT licensed.