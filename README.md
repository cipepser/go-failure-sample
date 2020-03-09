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

## テストを書いてみる

以下の`GetName`をテストすることを考える。

```go
func (c *Client) GetName(id int) (string, error) {
	for _, c := range customers {
		if c.id == id {
			return c.name, nil
		}
 	}
	return "", failure.New(NotFound)
}
```

以下のようにすればPASSする。

気になったところは以下。

- `wantErr`の型が`failure.StringCode`になったが、正しい？それ以外の型が出てこないのかまだわかっていない
- 正常系は、`wantErr`を`nil`にしたい（1つ目のテストにあるように`""`になってしまった）

でも上記くらいでほとんど実用では困らなそう。  
標準エラーもwrapするなどして、`failure`に統一してしまえば、
`if err != nil && !failure.Is(err, tt.wantErr) {`
のチェックしか出てこない（はず）なので、きれいだと思う。

```go
func TestClient_GetName(t *testing.T) {
	type fields struct {
		user string
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr failure.StringCode
	}{
		{
			name: "Alice",
			args: args{
				id: 0,
			},
			want:    "Alice",
			wantErr: "",
		},
		{
			name: "",
			args: args{
				id: -1,
			},
			want:    "",
			wantErr: NotFound,
		},
	}

	_ = NewCustomer("Alice", "alice@example.com")
	_ = NewCustomer("Bob", "bob@example.com")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				user: tt.fields.user,
			}
			got, err := c.GetName(tt.args.id)
			if err != nil && !failure.Is(err, tt.wantErr) {
				t.Errorf("GetName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetName() got = %v, want %v", got, tt.want)
			}
		})
	}
	cleanUpCustomers()
}
```

## エラーをwrapする

TODO: unimplemented

例題として、以下を考える。

- `whitelist.txt`を個別ファイルとして用意
- リストに含まれる`address`のみをアクセス許可

ホワイトリスト形式にしたため、ファイル自体が開けない場合もアクセス遮断する。
要はファイルを開けない場合のアクセス遮断で、wrapを実現する。

`whitelist.txt`

```text
alice@example.com
bob@example.com
```


## References
- [morikuni/failure: failure is a utility package for handling application errors\.](https://github.com/morikuni/failure)