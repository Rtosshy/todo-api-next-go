# 007 OpenAPI

- Status: Accepted
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

- REST APIを作成するにあたって、リクエスト/レスポンスのJSONとGo構造体の対応を定義する必要がある
- APIの入出力はフロントエンドとバックエンドの境界になるため、仕様のずれを早く検知できる形にしたい
- リクエストボディのバリデーションや型定義など、手書きすると単純だが繰り返しになりやすい実装を自動化したい
- 個人開発兼学習目的のプロジェクトであり、API仕様を明示しながらOpenAPI駆動開発の流れを経験したい

## Decision

- API仕様は`backend/api/openapi.yaml`にOpenAPI 3.0形式で記述する
- バックエンドではoapi-codegenを使用し、Goの型定義とGin向けのハンドラ定義を生成する
- フロントエンドではOrvalを使用し、TypeScriptのAPIクライアントとZodスキーマを生成する

## Rationale

- Go構造体、TypeScript型、APIクライアントを同じ`openapi.yaml`から決定的に生成でき、手書きによる型ずれや入力ミスを減らせる
- JSONのパース、必須項目、型、形式などの単純なリクエスト検証を仕様に寄せられるため、ハンドラ実装をユースケース呼び出しとレスポンス変換に集中させやすい
- フロントエンドでも同じ仕様からクライアントコードを生成できるため、API変更時に影響範囲をコンパイルエラーや生成差分として検知しやすい
- Swagger UIで仕様を確認でき、実装前後のAPI確認や手動疎通がしやすい
- APIの境界を先に明文化することで、バックエンド内部の設計とHTTPインターフェースを分けて考えやすくなる

## Consequences

- 良い影響
  - API仕様、Go型、TypeScriptクライアントの単一の参照元を`openapi.yaml`に寄せられる
  - ハンドラの型が生成されるため、API仕様と実装の不一致に気づきやすい
  - フロントエンドは手書きのfetch/axios呼び出しを減らせる
  - Swagger UIでAPI仕様をブラウザから確認できる
- 受け入れる制約
  - APIを変更するたびに`openapi.yaml`を先に更新し、バックエンドとフロントエンドの生成コマンドを実行する必要がある
  - 生成コードは直接編集せず、生成元の仕様または変換層を修正する運用にする
  - OpenAPIで表現しづらいドメイン制約は、生成されたバリデーションだけに任せずユースケース層やドメイン層で検証する
  - oapi-codegenやOrvalの生成結果に実装構造が影響を受けるため、ツール更新時は生成差分を確認する必要がある

## Alternatives Considered

- 手書きのGo構造体と手書きのフロントエンドAPIクライアント: 実装量は少なく始められるが、API変更時にバックエンドとフロントエンドの型ずれを起こしやすいため不採用
- Goコード/コメントからOpenAPI仕様を生成する方式: 実装を起点に仕様を生成できるためバックエンド単体では扱いやすいが、API仕様をフロントエンドとの契約として先に定義・確認する流れにしづらい。また、仕様の意図がGoの構造体やコメントに分散しやすいため不採用
- gRPC / Protocol Buffers: 型安全な契約には向いているが、今回のtodoアプリはブラウザから扱うREST APIで十分であり、学習・実装コストがスコープに対して大きいため不採用

## References

- OpenAPI Specification: https://spec.openapis.org/oas/latest.html
- oapi-codegen: https://github.com/oapi-codegen/oapi-codegen
- Orval: https://orval.dev/
- 関連ADR: [001-gin.md](./001-gin.md)
