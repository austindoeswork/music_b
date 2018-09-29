package main
import(
 "fmt"
 "net/http"
 "io/ioutil"
)

func add(x int, y int) int {
  return x + y
}


func main() {
    resp, _ := http.Get("http://example.com")
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
}

