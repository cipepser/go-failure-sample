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

`whitelist.txt`の`Oepn`に失敗したり、ホワイトリストに存在しないアドレスの場合、`Forbidden`エラーを返す関数として`CheckPermitted`を実装する。
許可されたアドレスにマッチした場合のみ、エラーが`nil`になる。[^1]

[^1]: 実環境で使う場合には毎回ファイルを開き直すのはパフォーマンスが悪い。今回はエラーハンドリングしたいだけなので、こんな実装になった。

```go
func (c *Client) CheckPermitted(address string) error {
	f, err := os.Open(WHITELIST)
	if err != nil {
		return failure.Wrap(err,
			failure.Context{"package": "os"},
			failure.Messagef("failed to open %s", WHITELIST),
		)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		permittedAddress := sc.Text()
		if permittedAddress == address {
			return nil
		}
	}

	return failure.New(Forbidden)
}
```

実行するmain関数。

```go
package main

import (
	"fmt"
	"github.com/cipepser/go-failure-sample/db"
	"github.com/morikuni/failure"
)

func main() {
	c := db.NewClient("user")
	if err := c.CheckPermitted("mallory@example.com"); err != nil {
		fmt.Println("============ Error ============")
		fmt.Printf("Error = %v\n", err)

		msg, _ := failure.MessageOf(err)
		fmt.Printf("Message = %v\n", msg)

		cs, _ := failure.CallStackOf(err)
		fmt.Printf("CallStack = %v\n", cs)

		fmt.Printf("Cause = %v\n", failure.CauseOf(err))

		fmt.Println()
		fmt.Println("============ Detail ============")
		fmt.Printf("%+v\n", err)
	}
}

```

### ホワイトリストにアドレスが存在しないパターン

`mallory@example.com`という、`whitelist.txt`に存在しないアドレスを入力に`CheckPermitted`を実行する。

試したかったポイントとしては、以下。

- [morikuni/failureのExample](https://github.com/morikuni/failure#example)がどのような出力となるのか
- スタックトレースが取れるか

実行結果

```text
❯ go run main.go
============ Error ============
Error = db.(*Client).CheckPermitted: code(Forbidden)
Code = Forbidden
Message =
CallStack = db.(*Client).CheckPermitted: main.main: runtime.main: goexit
Cause = code(Forbidden)

============ Detail ============
[db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:87
    code(Forbidden)
[CallStack]
    [db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:87
    [main.main] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/main.go:16
    [runtime.main] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/proc.go:203
    [runtime.goexit] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/asm_amd64.s:1357
```

### os.Openに失敗するパターン

`whitelist.txt`を一時的に削除し、`os.Open`に失敗する状況にする。

試したかったポイントとしては、以下。

- errの`Wrap`
- `Wrap`してもスタックトレースが取れるか
- `failure.Context`の挙動
- `failure.Messagef`の挙動

実行結果

```text
❯ go run main.go
============ Error ============
Error = db.(*Client).CheckPermitted: package=os: failed to open whitelist.txt: open whitelist.txt: no such file or directory
Code = <nil>
Message = failed to open whitelist.txt
CallStack = db.(*Client).CheckPermitted: main.main: runtime.main: goexit
Cause = no such file or directory

============ Detail ============
[db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:71
    package = os
    message("failed to open whitelist.txt")
    *os.PathError("open whitelist.txt: no such file or directory")
    syscall.Errno("no such file or directory")
[CallStack]
    [db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:71
    [main.main] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/main.go:16
    [runtime.main] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/proc.go:203
    [runtime.goexit] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/asm_amd64.s:1357
```

`Code = <nil>`になった。
Detailのほうに書いてあるとおり、`*os.PathError`が実際のエラーなのだが、
自前で実装している範囲で単純に`Wrap`するだけだと`nil`になってしまう。
そこで、以下を`failure.Translate`を使って書き換える。

変更前

```go
		return failure.Wrap(err,
			failure.Context{"package": "os"},
			failure.Messagef("failed to open %s", WHITELIST),
		)
```

変更後

```go
		return failure.Translate(err, Forbidden,
			failure.Context{"package": "os"},
			failure.Messagef("failed to open %s", WHITELIST),
		)
```

実行結果

```text
❯ go run main.go
============ Error ============
Error = db.(*Client).CheckPermitted: package=os: failed to open whitelist.txt: code(Forbidden): open whitelist.txt: no such file or directory
Code = Forbidden
Message = failed to open whitelist.txt
CallStack = db.(*Client).CheckPermitted: main.main: runtime.main: goexit
Cause = no such file or directory

============ Detail ============
[db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:71
    package = os
    message("failed to open whitelist.txt")
    code(Forbidden)
    *os.PathError("open whitelist.txt: no such file or directory")
    syscall.Errno("no such file or directory")
[CallStack]
    [db.(*Client).CheckPermitted] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/db/db.go:71
    [main.main] /Users/cipepser/.go/src/github.com/cipepser/go-failure-sample/main.go:16
    [runtime.main] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/proc.go:203
    [runtime.goexit] /usr/local/Cellar/go/1.13.5/libexec/src/runtime/asm_amd64.s:1357
```

`Code = Forbidden`になった。しかもDetailのほうは`*os.PathError`のままなので、原因を追うことができる。  
最終的にどのエラーを返すのかは、ドメインに応じて設計する必要があるが、設計さえすればやりたいことが実現できる。


// TODO: db_testに書く


## unwrap

`whitelist.txt`が存在しない状態とする。

```go
func main() {
	c := db.NewClient("user")
	if err := c.CheckPermitted("mallory@example.com"); err != nil {
		fmt.Println("============ Unwrap ============")
		iter := failure.NewIterator(err)
		for iter.Next() {
			fmt.Println("***********")
			err := iter.Error()
			fmt.Printf("Error = %v\n", err)
			code, ok := failure.CodeOf(err)
			if !ok {
				fmt.Println("not ok")
				continue
			}
			switch code {
			case db.Forbidden:
				fmt.Println(db.Forbidden)
			case db.NotFound:
				fmt.Println(db.NotFound)
			default:
				fmt.Println("default")
			}
		}
	}
}
```

実行結果

```text
============ Unwrap ============
***********
Error = db.(*Client).CheckPermitted: package=os: failed to open whitelist.txt: code(Forbidden): open whitelist.txt: no such file or directory
Forbidden
***********
Error = db.(*Client).CheckPermitted: package=os: failed to open whitelist.txt: code(Forbidden): open whitelist.txt: no such file or directory
Forbidden
***********
Error = package=os: failed to open whitelist.txt: code(Forbidden): open whitelist.txt: no such file or directory
Forbidden
***********
Error = failed to open whitelist.txt: code(Forbidden): open whitelist.txt: no such file or directory
Forbidden
***********
Error = code(Forbidden): open whitelist.txt: no such file or directory
Forbidden
***********
Error = open whitelist.txt: no such file or directory
not ok
***********
Error = no such file or directory
not ok
```

- `Unwrap`メソッドが明示的にあるわけではないので、`failure.Iterator`を用いる
→どこでWrapされているのかわからず全部出力するしかない？
→TODO: そもそもどんなハンドリングをしたい？
- `CodeOf`を使うことで、switch文を用いたエラーハンドリングが可能
- `not ok`となるのは標準の`os` packageから。

// TODO: うーん。たぶんstructにエラーに必要な情報を格納してreturnしたいんだよな



## References
- [morikuni/failure: failure is a utility package for handling application errors\.](https://github.com/morikuni/failure)