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

- トランザクション用のリポジトリメソッドを用意する: ユースケースごとに専用メソッドが必要になり、ユースケース数に比例してリポジトリが肥大化する。リポジトリがアプリケーションロジックを抱えてしまい、ユースケース層が薄くなる
  ```go
  // リポジトリ側に「Task作成 + Status解決」をまとめて閉じ込める形
  func (r *taskRepository) CreateWithStatus(ctx context.Context, task *domain.Task, status *domain.Status) (*domain.Task, error) {
      var created *domain.Task
      err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
          // ここに Status の GetOrCreate と Task の Create が同居する
          return nil
      })
      return created, err
  }
  ```

- Unit of Work（RepositoryManagerを介して同一トランザクションを共有するリポジトリ群を一括で扱う）: トランザクション境界とリポジトリ群を一括で表現できるが、Manager / Factory / ジェネリクスによる抽象が必要で、本リポジトリの規模に対して構成が重い（参考: GREE Tech Blog）
  ```go
  // application層: UoWとRepositoryManagerのインターフェース
  type TaskCreateRepoManager interface {
      TaskRepository() repository.TaskRepository
      StatusRepository() repository.StatusRepository
  }

  type UnitOfWork[T any] interface {
      DoInTx(ctx context.Context, fn func(ctx context.Context, repoManager T) error) error
  }

  // infra層: db.Transactionの中でtxを共有したRepositoryManagerをfactoryで組み立てる
  type unitOfWork[T any] struct {
      db      *gorm.DB
      factory func(db *gorm.DB) T
  }

  func (u *unitOfWork[T]) DoInTx(ctx context.Context, fn func(ctx context.Context, repoManager T) error) error {
      return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
          return fn(ctx, u.factory(tx))
      })
  }

  // ユースケース側
  func (tu *taskUsecase) Create(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
      return tu.uow.DoInTx(ctx, func(ctx context.Context, rm TaskCreateRepoManager) error {
          status, err := rm.StatusRepository().GetOrCreate(ctx, ...)
          if err != nil {
              return err
          }
          _, err = rm.TaskRepository().Create(ctx, ...)
          return err
      })
  }
  ```

- ユースケースで `db.Transaction` を直接呼び、得られたtxをリポジトリに渡すDIパターン: 関数シグネチャからトランザクション参加が明示される利点はあるが、ユースケースが `*gorm.DB` に直接依存し、リポジトリのインターフェースにも `tx *gorm.DB` が現れることでusecase/domainの抽象がGORMに汚染される。トランザクション制御の差し替えも難しくなる
  ```go
  // リポジトリのインターフェースに tx が露出する
  type TaskRepository interface {
      Create(ctx context.Context, tx *gorm.DB, task *domain.Task) (*domain.Task, error)
  }

  // ユースケースが *gorm.DB を直接知り、tx を引数として配り回す
  func (tu *taskUsecase) Create(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
      return tu.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
          status, err := tu.sr.GetOrCreate(ctx, tx, ...)
          if err != nil {
              return err
          }
          _, err = tu.tr.Create(ctx, tx, ...)
          return err
      })
  }
  ```

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- Goのトランザクション管理パターン (GREE Tech): https://tech.gree-x.com/golang-transaction-pattern/
- レイヤードアーキテクチャにおけるトランザクション (karamaru-alpha): https://karamaru-alpha.com/posts/layered-tx/
- GORM Transactions: https://gorm.io/docs/transactions.html
- 関連ADR: [002-gorm.md](./002-gorm.md), [005-backend-architecture.md](./005-backend-architecture.md)
- 実装: `backend/internal/usecase/transaction_manager.go`, `backend/internal/infra/db/transaction_manager.go`, `backend/internal/infra/db/context.go`, `backend/internal/infra/db/base_repository.go`
- 利用例: `backend/internal/usecase/task.go`, `backend/internal/usecase/user.go`
