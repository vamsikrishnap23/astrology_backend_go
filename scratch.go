package main
import (
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"time"
)

func main() {
	go func() {
		// Run server in background
	}()
	time.Sleep(1 * time.Second)
	// We'll just run it directly without spinning up the server since it's already tested mostly.
}
