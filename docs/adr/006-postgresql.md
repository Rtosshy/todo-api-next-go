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

- todoアプリケーションでは User / Task / Status などのデータを永続化する必要がある
- Task と Status の関連や、ユーザーごとのデータ分離を扱うため、データの整合性を保ちやすいデータストアを選ぶ必要がある
- GORMを採用する方針（[002-gorm.md](./002-gorm.md)）と整合し、ローカル開発・本番環境で同じDB種別を使える構成にしたい
- 個人開発兼学習目的のプロジェクトであり、将来的な機能追加に耐えられる一方で、運用コストを抑えられる選択肢が望ましい
- 学生として利用しやすい無料枠やホスティングの選択肢があるDBを優先したい

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- アプリケーションの永続化用DBとしてPostgreSQLを採用する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- PostgreSQLはSupabase、Neon、Render、Railway、Heroku、AWS RDS、Google Cloud SQLなど、主要なホスティングサービスで広くサポートされており、個人開発から本番運用まで移行先の選択肢を確保しやすい
- PostgreSQLは個人開発で利用しやすいマネージドホスティングの選択肢が多く、学生でも運用コストを抑えて本番環境を用意しやすい
- User / Task / Status のような関連を持つデータを扱うため、外部キー制約やトランザクションを備えたRDBが適している
- GORMがPostgreSQLをサポートしており、ローカル開発と本番環境で同じDB種別を使いやすい
- todoアプリとしては小規模だが、将来的な機能追加や検索条件の追加にも対応しやすい

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

- 良い影響
  - 主要なホスティングサービスで広くサポートされており、個人開発から本番運用まで移行先の選択肢を確保しやすい
  - 関連データを外部キーやトランザクションで扱いやすく、User / Task / Status 間の整合性を保ちやすい
  - ローカル開発と本番環境で同じDB種別を使いやすく、環境差分を減らせる
  - 将来的な検索条件追加、集計、並び替えなどに対応しやすい
- 受け入れる制約
  - PostgreSQL固有機能を使いすぎると、他のRDBへ移行しづらくなる
  - SQLiteと比べるとローカル開発環境の準備が重くなる
  - スキーマ設計、マイグレーション、インデックス設計を継続的に管理する必要がある
  - ホスティングサービスごとの接続数、無料枠、バックアップ、リージョンなどの制約を確認する必要がある

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- MySQL: PostgreSQLと同じRDBであり、GORM対応や運用実績も十分にある。ただし、今回の要件ではPostgreSQLに対する明確な優位が少なく、PostgreSQL系のホスティング選択肢や学習価値を優先したため不採用。
- SQLite: ファイルベースで軽量であり、ローカル開発の導入は最も簡単。ただし、本番環境で複数ユーザー・複数接続を扱う構成では制約が出やすく、ローカルと本番でDB種別が分かれる可能性があるため不採用。
- Redis: 高速なインメモリデータストアであり、キャッシュやセッション管理には適している。ただし、User / Task / Status の関連、外部キー制約、トランザクションを主DBとして扱う用途には向かないため不採用。
- NoSQL: スキーマ変更に柔軟で、ドキュメント指向のデータやアクセスパターンが固定された高スケール用途には有力。ただし、今回のtodoアプリでは関連データの整合性をRDBで扱う方がシンプルなため不採用。

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- Supabase Database: https://supabase.com/docs/guides/database
- Neon: https://neon.com/
- Render Postgres: https://render.com/docs/postgresql
- Railway Databases: https://docs.railway.com/databases
- Vercel Storage: https://vercel.com/docs/storage
- AWS RDS: https://docs.aws.amazon.com/ja_jp/AmazonRDS/latest/UserGuide/Welcome.html
- Google Cloud SQL: https://docs.cloud.google.com/sql/docs
