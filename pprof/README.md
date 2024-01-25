## pprof
ビルド後に下記で集計

`go tool pprof -http=":8081" http://localhost:8080/debug/pprof/profile?seconds=5`


## 参考文献
- [ドキュメント](https://github.com/google/pprof/tree/main)
- [Goのpprofの使い方【基礎編】](https://christina04.hatenablog.com/entry/golang-pprof-basic)