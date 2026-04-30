# 004 Transaction Boundary

- Status: Accepted
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

<!--
なぜこの決定が必要になったのか。
前提、制約、課題を書く。
-->

- 1リクエスト中に複数集約をまたぐ永続化が発生する（例: Task作成時に Status の `GetOrCreate` と Task の `Create` を一貫性をもって書き込む必要がある）
- ドメインの不変条件を満たす単位でコミット/ロールバックを制御したく、その単位を明示的に置く層を決める必要があった
- ハンドラやリポジトリに `*gorm.DB` を直接持たせると層間の依存が漏れ、テスト容易性も損なわれる
- GORMを採用した（[002-gorm.md](./002-gorm.md)）うえで、トランザクション境界の引き方と `*gorm.DB` の伝播方法を決める必要があった

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- トランザクション境界はユースケース層に置く
- ユースケース層に `TxManager` インターフェースを定義し、ユースケースから `RunInTx(ctx, fn)` を呼んで境界を起動する
- 実装は `infra/db` に置き、GORMの `db.Transaction(func(tx *gorm.DB) error)` クロージャに委譲する（Begin/Commit/Rollback はGORMに任せる）
- トランザクション中の `*gorm.DB` は `context.Context` に格納して下位層に伝播する
- リポジトリは `baseRepository.db(ctx)` を経由してDBを取得し、ctxにtxがあれば自動で参加、なければ素のDBで動作する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- ドメインの整合性単位＝ユースケースなので、境界の置き場所として自然
- `TxManager` をusecaseパッケージのインターフェースとして定義することで、ユースケース層がGORMに直接依存せずDIで切り替えられる
- contextでの伝播により、リポジトリのシグネチャに `tx *gorm.DB` のような実装詳細を出さずに済む
- GORMのクロージャAPIは Commit/Rollback 漏れを構造的に防げる
- 単一集約のRead系（`Get`, `GetAll`）は `RunInTx` で包まず素直にリポジトリ呼び出しにする → 不要なオーバーヘッドを避ける

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

- 良い影響
  - ユースケースを読むだけでトランザクション範囲が把握できる
  - リポジトリとユースケースのテストがそれぞれ独立しやすい
  - GORM固有の型がユースケース層に漏れない
- 受け入れる制約
  - context経由でtxを渡すため、契約がコンパイル時に見えづらい（`GetTx` はランタイム判定）
  - ネストした `RunInTx` 呼び出しの挙動はGORMのSavepoint仕様に従う点を理解しておく必要がある
  - Read-onlyな処理で誤って `RunInTx` を使うとオーバーヘッドが乗る

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- トランザクション用のリポジトリメソッドを用意する: ユースケースの数だけリポジトリメソッドが肥大化
  - コード例
- unit of work: 今回だと複雑すぎる
  - コード例
- リポジトリ層で境界を引く: 複数リポジトリにまたがる整合性を表現できず不採用
  - コード例
- ユースケース層でトランザクション処理を直接呼び出す: usecaseにdbの処理を書く必要がある
  - コード例
- `*gorm.DB` を引数で引き回す: usecase/repositoryのシグネチャがGORMに汚染されるため不採用

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- トランザクションの書き方の記事二つ
- からまるさんのやつ
- greeeのやつ
- GORM Transactions: https://gorm.io/docs/transactions.html
- 関連ADR: [002-gorm.md](./002-gorm.md), [005-backend-architecture.md](./005-backend-architecture.md)
- 実装: `backend/internal/usecase/transaction_manager.go`, `backend/internal/infra/db/transaction_manager.go`, `backend/internal/infra/db/context.go`, `backend/internal/infra/db/base_repository.go`
- 利用例: `backend/internal/usecase/task.go`, `backend/internal/usecase/user.go`
