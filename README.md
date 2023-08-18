# simplist

PostgreSQLの配列型を使用した、List機能のみを持ったWeb APIです。

設定管理を想定しています。

## 使い方

起動
```
docker run --env KEYS="foo,bar" --env DB_URL="<YOUR_DB_URL>" -p 8080:8080 takumi3488/simplist
```

全て取得
```
$ curl localhost:8080/lists
[{"key":"foo","items":[]},{"key":"bar","items":[]}]
```

1つ取得
```
$ curl localhost:8080/lists/foo
{"key":"foo","items":[]}
```

更新
```
$ curl -X PUT \
    -H "Content-Type:application/json" \
    -d '{"items":["a","b","c"]}' \
    localhost:8080/lists/foo
$ curl localhost:8080/lists/foo
{"key":"foo","items":["a","b","c"]}
```

## 環境変数

- KEYS: カンマ区切りでリスト名を設定します
- DB_URL: PostgreSQLのDSNを指定します

## 特徴

キーと配列のみでデータを管理するため、以下のメリットがあります。

- 最低限のプロパティで操作しやすい
- 配列要素の並べ替えが容易

その他のメリットとして以下が挙げられます。

- 起動時にテーブル作成やキーの挿入を行うため、環境変数とDBさえあればすぐに起動できる

## 注意

- 認証機能の実装予定はありません。外部公開せずに使うか、前段にNginx等を置いて認証機能を追加するのがおすすめです。
- 再起動せずにキーを追加することはできません。停止させたくない場合にはローリングアップデート等が使える環境で使用してください。
