package fmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func test(t *testing.T, src, dst string) {
	module := "github.com/PinkD/gofimpt"
	assert.Equal(t, dst, FormatCode(module, src))
}

func TestFormatWithoutImport(t *testing.T) {
	code1 := `
package main 

func main() {
	println("Hello World!")
}
`
	test(t, code1, code1)
}

func TestFormatCode(t *testing.T) {
	code1 := `
package main 
import (
	"fmt"
	"github.com/PinkD/gofimpt/errors"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
)

func main() {
	println("Hello World!")
}
`
	code2 := `
package main 
import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/PinkD/gofimpt/errors"
)

func main() {
	println("Hello World!")
}
`
	test(t, code1, code2)
}

func TestFormatWithComment(t *testing.T) {
	code1 := `
package main 
import (
	"fmt"
	// comment 1
	// comment 2
	"github.com/PinkD/gofimpt/errors"
	"github.com/pkg/errors"
	"os"
	"sort"
	"strings"
)

func main() {
	println("Hello World!")
}
`
	code2 := `
package main 
import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"

	// comment 1
	// comment 2
	"github.com/PinkD/gofimpt/errors"
)

func main() {
	println("Hello World!")
}
`
	test(t, code1, code2)
}

func TestFormatWithoutBuiltin(t *testing.T) {
	code1 := `
package main 
import (
	"github.com/PinkD/gofimpt/errors"
	"github.com/pkg/errors"
)

func main() {
	println("Hello World!")
}
`
	code2 := `
package main 
import (
	"github.com/pkg/errors"

	"github.com/PinkD/gofimpt/errors"
)

func main() {
	println("Hello World!")
}
`
	test(t, code1, code2)
}

func TestFormatWithoutBuiltinAndThirdParty(t *testing.T) {
	code1 := `
package main 
import (
	// comment 1
	// comment 2
	"github.com/PinkD/gofimpt/errors"
)

func main() {
	println("Hello World!")
}
`
	code2 := `
package main 
import (
	// comment 1
	// comment 2
	"github.com/PinkD/gofimpt/errors"
)

func main() {
	println("Hello World!")
}
`
	test(t, code1, code2)
}

func TestFormatMultiImport(t *testing.T) {
	code1 := `
package main 
import (
	"fmt"
	"github.com/PinkD/gofimpt/errors"
	"sort"
	"strings"
)

import (
	"github.com/pkg/errors"
	"os"
)

func main() {
	println("Hello World!")
}
`
	code2 := `
package main 
import (
	"fmt"
	"sort"
	"strings"

	"github.com/PinkD/gofimpt/errors"
)

import (
	"os"

	"github.com/pkg/errors"
)

func main() {
	println("Hello World!")
}
`
	test(t, code1, code2)
}
