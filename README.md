# go-failure-sample

[morikuni/failure](https://github.com/morikuni/failure)を触ってみる。

## エラーの作り方

ざっくり`db`パッケージを用意する。

```go
package db

import (
	"github.com/morikuni/failure"
)

var NotFound failure.StringCode = "NotFound"

type Client struct {
	user string
}

func NewClient(user string) *Client {
	return &Client{
		user: user,
	}
}

func (c *Client) GetName(id int) (string, error) {
	return "", failure.New(NotFound)
}
```

上記実装の通り、`GetName()`を呼ぶと`NotFound`エラーを返す。
実際に動かしてみる。

```go
func main() {
	c := db.NewClient("user")
	userId := 0
	_, err := c.GetName(userId)
	if failure.Is(err, db.NotFound) {
		fmt.Println("error occurred: NotFound")
	}
}
```

実行結果

```text
❯ go run main.go
error occurred: NotFound
```




## References
- [morikuni/failure: failure is a utility package for handling application errors\.](https://github.com/morikuni/failure)