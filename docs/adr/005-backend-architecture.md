# 005 Backend Architecture

- Status: Accepted
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

<!--
なぜこの決定が必要になったのか。
前提、制約、課題を書く。
-->

- todoアプリケーションとしては小規模だが、User / Task / Status など複数のドメイン概念を扱う
- ハンドラ、DBアクセス、ドメインロジックが密結合になると、機能追加時に変更箇所が広がりやすくなる
- GinやGORMなどの外部ライブラリの都合を、ドメインロジックやユースケースに直接持ち込みたくない
- ユースケース単位でテストしやすい設計にし、実装変更時の品質を担保したい
- 将来的に外部API連携や通知などの外部I/Oを伴う機能追加を想定している
- 学習目的として、依存方向、責務分離、テストしやすい境界設計を経験したい

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- クリーンアーキテクチャの考え方を一部取り入れた、レイヤードアーキテクチャを採用する
- 依存方向は外側から内側に向ける
  - `infra/web` はGin、OpenAPI生成コード、HTTPの入出力変換を扱い、`usecase` を呼び出す
  - `infra/db` はGORM、DAO、Domain/DAO変換を扱い、`domain/repository` のインターフェースを実装する
  - `usecase` はアプリケーションロジック、Input、トランザクション境界を扱い、`domain` と repository interface に依存する
  - `domain` はエンティティ、値オブジェクト、ドメイン不変条件を扱い、外部ライブラリに依存しない
- Repository interface は `domain/repository` に置く
- DTO / Input / Domain Entity / DAO の変換責務は境界ごとに置く
  - HTTP DTO と OpenAPI生成型の変換は `infra/web`
  - ユースケース入力の組み立てとアプリケーションロジックは `usecase`
  - Domain Entity と DAO の変換は `infra/db`
- DIは `infra/web` 起点で、server/router/handler生成時にRepository、TxManager、Usecaseを組み立てる
- パッケージ構成は `domain` / `usecase` / `infra/db` / `infra/web` の技術レイヤー単位を採用する
- 将来の外部API連携や通知は、`infra/external`、`infra/notification` のような用途別の `infra` 配下へ実装を追加する
- 外部API連携や通知の抽象インターフェースは、ユースケースが必要とするPortとして `usecase` 側に置き、`infra` 側が実装する
- 厳密なClean Architectureの層名や構成にはこだわらず、このプロジェクトの規模に合わせて簡略化する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- HandlerにHTTPの入出力変換、Usecaseにアプリケーションロジック、Repositoryに永続化処理を分けることで、責務の境界を読み取りやすくなる
- `domain` / `usecase` がGinやGORMに直接依存しないため、ユースケース単位のテストを書きやすい
- 外部API連携や通知が追加されても、Usecaseが処理の流れを制御し、具体的な外部I/Oの実装を `infra` 側に閉じ込められる
- Task作成時のStatus解決のように複数Repositoryをまたぐ処理を、Usecase層に集約できる
- トランザクション境界をUsecase層に置く方針（[004-transaction-boundary.md](./004-transaction-boundary.md)）と整合する
- Gin、GORM、OpenAPI生成コードの影響範囲を `infra` 層に寄せることで、外部ライブラリ由来の変更がドメインロジックへ波及しにくくなる
- Repository interface を `domain/repository` に置くことで、Usecaseは永続化の実装ではなくドメイン側の抽象に依存できる
- 技術レイヤー単位のパッケージ構成にすることで、現時点の小規模な機能数でも依存方向を把握しやすい
- 小規模アプリケーションに対してはやや重い構成だが、学習と設計理解、将来の拡張性、テスト容易性を優先する
- ファイル数や変換処理が増えるコストは、変更時の保守性と品質担保のために受け入れる

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

- 良い影響
  - ユースケース単位でテストしやすい
  - HTTP、DB、ドメインロジックの責務が分かれ、変更時の影響範囲を追いやすい
  - GinやGORMの型がUsecase/Domain層に漏れにくい
  - OpenAPI生成型の影響範囲をHTTP境界に寄せやすい
  - 複数Repositoryをまたぐ処理をUsecase層で表現しやすい
  - 外部API連携や通知の実装を差し替えやすい
- 受け入れる制約
  - 小規模なCRUDに対してはファイル数と実装量が増える
  - レイヤー分割による認知負荷が増える
  - 単純な処理でもHandler、Usecase、Repository、Domainのどこに書くべきか判断が必要になる
  - インターフェースや変換処理が増え、実装初期の速度は落ちる
  - 通知方式、非同期ジョブ、リトライ、外部APIごとの詳細設計は、必要になった時点で別ADRとして決める

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- HandlerからGORMを直接呼ぶシンプルな構成: 実装量は最も少ないが、HTTP、永続化、ドメインロジックが密結合になりやすく、ユースケース単位のテストが難しくなるため不採用
- 一般的な3層構成（Handler / Service / Repository）: このプロジェクトでも成立するが、Domainの不変条件やRepository interfaceの依存方向を明示しづらく、外部ライブラリへの依存を内側に持ち込みやすいため不採用
- 機能単位のパッケージ構成: user/task/status のように機能ごとの見通しは良くなるが、現時点では機能数が少なく、学習目的としてレイヤー間の依存方向を明示する方を優先するため不採用
- `cmd/server/main.go` でDIを組み立てる構成: エントリポイントに依存組み立てを集約できるが、現行構成ではGinのrouter/handler生成と密接に関わるため、`infra/web` 起点で組み立てる方針を採用する
- 外部API連携や通知の実装をUsecaseから直接呼ぶ構成: 抽象が少なく実装は速いが、外部サービスの変更やテスト時の差し替えが難しくなるため不採用
- 厳密なClean Architecture: 依存方向を明確にできるが、Interface AdapterやPresenterなどの層を厳密に分けると、このプロジェクトの規模に対して抽象が重くなるため不採用

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- 関連ADR: [001-gin.md](./001-gin.md), [002-gorm.md](./002-gorm.md), [004-transaction-boundary.md](./004-transaction-boundary.md), [007-openapi.md](./007-openapi.md)
- 実装: `backend/internal/domain`, `backend/internal/usecase`, `backend/internal/infra/db`, `backend/internal/infra/web/gin`
- 将来の実装候補: `backend/internal/infra/external`, `backend/internal/infra/notification`
