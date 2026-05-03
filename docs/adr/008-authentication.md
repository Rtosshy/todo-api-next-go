# 008 Authentication

- Status: Proposed
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

<!--
なぜこの決定が必要になったのか。
前提、制約、課題を書く。
-->

- このアプリケーションは、ユーザーごとに Task を作成・取得・更新・削除する todo アプリケーションである
- Task API では、リクエスト元のユーザーを識別し、そのユーザーに紐づくデータだけを操作できるようにする必要がある
- frontend は Next.js、backend は Gin で実装しており、ブラウザから API を呼び出す構成である
- ブラウザ向け Web アプリケーションでは、認証情報の保持方法として HttpOnly cookieを使う方式と、Authorization header に Bearer token を付与する方式が広く使われている
- このプロジェクトでは、JavaScript から認証情報を直接扱わない構成にするか、frontend側で token を明示的に管理する構成にするかを決める必要があった
- cookie を使う場合、ブラウザが認証情報を自動送信する性質により CSRF リスクも考慮する必要がある

## Decision

<!--
何を決めたのか。
できるだけ短く、明確に書く。
-->

- HttpOnly cookie + JWT を採用する

## Rationale

<!--
なぜその決定にしたのか。
判断基準、重視したことを書く。
-->

- Next.js frontend から Gin backend の API を呼び出すブラウザアプリケーションであり、認証状態を cookie に保持することで、ブラウザの標準的な送信機構を利用できる
- HttpOnly cookie を使うことで、frontend の JavaScript から JWT を直接読み取れないようにし、localStorage などに token を保存する方式と比べて、XSS 発生時の token 漏えいリスクを下げられる
- frontend 側で access token の保存、取得、Authorization header への付与、logout 時の削除を個別に実装しなくてよく、認証情報の管理責務を backend 側に寄せられる
- JWT を使うことで、サーバー側に session store を持たずに、署名検証と claim の確認で認証済みユーザーを判定できる
- このアプリケーションは小規模な todo アプリケーションであり、server-side session store を導入するよりも、JWT の検証を middleware に集約する構成の方が実装と運用の負荷を抑えられる
- cookie はブラウザにより自動送信されるため CSRF リスクを持つが、CSRF token の検証と組み合わせることで、そのリスクを明示的に扱える

## Consequences

<!--
この決定によって起きる影響を書く。
良い影響だけでなく、受け入れる制約や運用負荷も書く。
-->

- 良い影響
  - frontend 側で JWT を保存・取得・削除する処理を持たなくてよくなり、認証情報の扱いを backend 側に集約できる
  - 認証必須 API の認証処理を Gin middleware に集約でき、各 handler で個別に token 検証を実装しなくてよい
  - middleware で検証した `user_id` をリクエストコンテキストに載せることで、後続の処理が認証済みユーザーを共通の方法で参照できる
  - server-side session store を持たないため、Redis などの追加コンポーネントを運用しなくてよい

- 受け入れる制約
  - cookie はブラウザにより自動送信されるため、CSRF 対策を必ず組み合わせる必要がある
  - JWT の署名に使う `SECRET` を安全に管理する必要がある
  - `SECRET` を変更すると既存 JWT は検証できなくなるため、鍵ローテーションを行う場合は再ログインの発生を受け入れる必要がある
  - 発行済み JWT は有効期限まで利用できるため、サーバー側で即時に無効化する仕組みは現時点では持たない
  - logout は `token` cookie を削除する処理であり、発行済み JWT 自体を失効させるものではない
  - cookie の `Secure`、`SameSite`、`Domain` は環境により適切な設定が異なるため、開発環境と本番環境の差分を管理する必要がある
  - JWT の有効期限が切れた場合、ユーザーは再ログインする必要がある

## Alternatives Considered

<!--
検討したが採用しなかった選択肢を書く。
それぞれを不採用にした理由も書く。
-->

- Authorization header に Bearer token を付与する方式: API クライアントやモバイルアプリでは扱いやすいが、ブラウザアプリケーションでは frontend 側で token の保存、取得、削除、リクエストへの付与を管理する必要がある。localStorage などに token を保存すると XSS 発生時の漏えいリスクが大きくなるため不採用
- Server-side session + Cookie: サーバー側で session を管理するため、強制ログアウトや即時失効を実装しやすい。一方で Redis や DB などの session store が必要になり、この小規模 todo アプリケーションでは運用コストが増えるため不採用
- Basic 認証: HTTP 標準の認証方式で実装は単純だが、リクエストごとにユーザー名とパスワード相当の認証情報を送る必要がある。一般的なブラウザアプリケーションのログイン、ログアウト、認証状態の維持を扱う方式としては柔軟性が低いため不採用

## References

<!--
参考資料、Issue、PR、公式ドキュメントなど。
-->

- 関連ADR: [005-backend-architecture.md](./005-backend-architecture.md), [007-openapi.md](./007-openapi.md), [009-csrf-protection.md](./009-csrf-protection.md)
- 実装: `backend/internal/usecase/user.go`, `backend/internal/infra/web/gin/handler/user.go`, `backend/internal/infra/web/gin/middleware/jwt.go`, `backend/internal/infra/web/gin/router/router.go`, `backend/pkg/cookie/helper.go`, `backend/api/openapi.yaml`
