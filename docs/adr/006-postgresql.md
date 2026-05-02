# 006 PostgreSQL

- Status: Accepted
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

<!--
なぜこの決定が必要になったのか。
前提、制約、課題を書く。
-->

- todo appのデータを永続化する必要がある
- DBが必要になる
- 今回のユースケースにあったデータストアが必要
- 学生の身分で無料ホスティングを備えているDBだから

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- DBにPostgreSQLを使用する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- 将来の拡張性も踏まえて高機能なデータストアである必要がある

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

-

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- MySQL
- SQLite
- Redis
- NoSQL

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

-
