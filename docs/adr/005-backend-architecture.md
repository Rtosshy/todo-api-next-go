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
- 学習目的として、依存方向を意識したバックエンド設計を経験したい

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- クリーンアーキテクチャの考え方を一部取り入れた、レイヤードアーキテクチャを採用する
- 依存方向は外側から内側に向ける
  - `infra/web` はHTTPフレームワークやOpenAPI生成コードを扱い、`usecase` を呼び出す
  - `infra/db` はGORMやDAOを扱い、`domain/repository` のインターフェースを実装する
  - `usecase` はアプリケーションロジックとトランザクション境界を扱い、`domain` と repository interface に依存する
  - `domain` はエンティティ、値オブジェクト、ドメイン不変条件を扱い、外部ライブラリに依存しない
- 厳密なClean Architectureの層名や構成にはこだわらず、このプロジェクトの規模に合わせて簡略化する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- HandlerにHTTPの入出力変換、Usecaseにアプリケーションロジック、Repositoryに永続化処理を分けることで、責務の境界を読み取りやすくなる
- `domain` / `usecase` がGinやGORMに直接依存しないため、ユースケース単位のテストを書きやすい
- Task作成時のStatus解決のように複数Repositoryをまたぐ処理を、Usecase層に集約できる
- トランザクション境界をUsecase層に置く方針（[004-transaction-boundary.md](./004-transaction-boundary.md)）と整合する
- Gin、GORM、OpenAPI生成コードの影響範囲を `infra` 層に寄せることで、外部ライブラリ由来の変更がドメインロジックへ波及しにくくなる
- 小規模アプリケーションに対してはやや重い構成だが、学習目的と今後の拡張しやすさを優先する

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

- 良い影響
  - ユースケース単位でテストしやすい
  - HTTP、DB、ドメインロジックの責務が分かれ、変更時の影響範囲を追いやすい
  - GinやGORMの型がUsecase/Domain層に漏れにくい
  - 複数Repositoryをまたぐ処理をUsecase層で表現しやすい
- 受け入れる制約
  - 小規模なCRUDに対してはファイル数と実装量が増える
  - レイヤー分割による認知負荷が増える
  - 単純な処理でもHandler、Usecase、Repository、Domainのどこに書くべきか判断が必要になる
  - インターフェースや変換処理が増え、実装初期の速度は落ちる

## Follow-up Decisions

<!--
このADRでは決めきらず、次のステップで追加する意思決定を書く。
-->

- 各層の責務と、どの処理をどの層に置くかの判断基準
- DTO / Input / Domain Entity / DAO の変換責務
- Repository interface を `domain/repository` に置く方針の理由と代替案
- DIの組み立て場所と粒度
- パッケージ構成を技術単位にするか、機能単位にするか

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- HandlerからGORMを直接呼ぶシンプルな構成: 実装量は最も少ないが、HTTP、永続化、ドメインロジックが密結合になりやすく、ユースケース単位のテストが難しくなるため不採用
- 一般的な3層構成（Handler / Service / Repository）: このプロジェクトでも成立するが、Domainの不変条件やRepository interfaceの依存方向を明示しづらく、外部ライブラリへの依存を内側に持ち込みやすいため不採用
- 厳密なClean Architecture: 依存方向を明確にできるが、Interface AdapterやPresenterなどの層を厳密に分けると、このプロジェクトの規模に対して抽象が重くなるため不採用

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- 関連ADR: [001-gin.md](./001-gin.md), [002-gorm.md](./002-gorm.md), [004-transaction-boundary.md](./004-transaction-boundary.md), [007-openapi.md](./007-openapi.md)
- 実装: `backend/internal/domain`, `backend/internal/usecase`, `backend/internal/infra/db`, `backend/internal/infra/web/gin`
