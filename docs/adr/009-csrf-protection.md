# 009 CSRF Protection

- Status: Proposed
- Date: 2026-04-30
- Supersedes: None
- Superseded by: None

## Context

- このアプリケーションは Next.js frontend から Gin backend の API を呼び出すブラウザ向け todo アプリケーションである
- [008 Authentication](./008-authentication.md) で、認証情報は JWT を HttpOnly cookie に保存する方針にした
- HttpOnly cookie は JavaScript から JWT を直接読み取れないため XSS 発生時の token 漏えいリスクを下げられるが、ブラウザが cookie を自動送信するため CSRF リスクが発生する
- 認証が不要な `signup`、`login`、`logout` と、認証が必要な Task API では、意図しない外部サイトからの送信を拒否する必要がある
- frontend と backend は開発環境では別 origin で動作し、本番環境でも Vercel と Render のような別 origin 構成を取りうるため、CORS と credential 付きリクエストを前提にする必要がある
- OpenAPI によるリクエスト検証と JWT 認証 middleware を使っているため、CSRF 検証も Gin middleware として共通化し、各 handler に個別実装を持ち込まない構成にしたい
- CSRF token は frontend が状態変更リクエストに header として付与できる必要がある。一方で、サーバー側に CSRF token 用の session store を追加すると、この小規模アプリケーションでは運用負荷が増える

## Decision

- Double Submit Cookie パターンによる CSRF 対策を採用する
- backend は `GET /api/v1/csrf` で CSRF token を生成し、`_csrf` cookie と response body の両方で返す
- frontend は取得した token をメモリ上に保持し、API リクエスト時に `X-CSRF-Token` header として送信する
- backend は CSRF 保護対象の API で `_csrf` cookie と `X-CSRF-Token` header の値が一致することを Gin middleware で検証する
- middleware の適用順は、CSRF 検証、OpenAPI リクエスト検証、JWT 認証の順にする

## Rationale

- Cookie 認証を使う以上、ブラウザが cookie を自動送信する性質に対して明示的な防御が必要である
- 攻撃者サイトは被害者ブラウザに cookie を送信させることはできるが、同一生成元ポリシーと CORS により、この API の `GET /api/v1/csrf` response body から token を読み取り、任意の `X-CSRF-Token` header を正しく付与することはできない
- CSRF token を cookie と header の両方で検証することで、単に cookie が送信されたという事実だけでは状態変更 API を実行できない
- CSRF token 用の server-side session store を持たないため、Redis などの追加コンポーネントを導入せずに実装できる
- CSRF 検証を middleware に置くことで、`signup`、`login`、`logout`、Task API に対して同じルールを適用でき、handler の責務をアプリケーション処理に集中できる
- JWT 認証より先に CSRF 検証を行うことで、不正な cross-site request を軽い検証で早期に拒否できる
- frontend の Axios instance に `withCredentials: true` と `X-CSRF-Token` 付与処理を集約することで、各画面や hook が CSRF header を個別に意識しなくてよい

## Consequences

- 良い影響
  - HttpOnly cookie + JWT の構成を維持しながら、CSRF リスクを明示的に扱える
  - CSRF 検証を Gin middleware に集約でき、handler ごとの実装漏れを避けやすい
  - CSRF token 用の server-side session store を追加せずに済む
  - OpenAPI validation や JWT authentication と同じ routing 層で、API 横断のセキュリティ処理として扱える

- 受け入れる制約
  - frontend はアプリ起動時または状態変更 API 呼び出し前に `GET /api/v1/csrf` を呼び、token を保持する必要がある
  - CSRF token の有効期限は `_csrf` cookie の `MaxAge` に依存するため、期限切れ時には token を再取得する必要がある
  - credential 付き cross-origin request を使うため、backend の CORS 設定では許可 origin、`AllowCredentials`、`X-CSRF-Token` header を正しく管理する必要がある
  - cookie の `Secure`、`SameSite`、`Domain` は環境差分を持つため、開発環境と本番環境で設定を誤ると CSRF token の cookie が保存または送信されない
  - Vercel と Render のように eTLD+1 が異なる別サイト構成では、Safari ITP や Chrome のサードパーティ cookie 制限により、`SameSite=None; Secure` でも cookie が保存または送信されない場合がある
  - 現在の実装では token をメモリに保持するため、ページ reload 後は再取得が必要になる

## Alternatives Considered

- SameSite cookie のみに依存する: 実装は単純だが、frontend と backend が別 origin で動く構成では `SameSite=None` が必要になる場合があり、CSRF 対策として不十分になるため不採用
- Origin / Referer header のみを検証する: 追加 token が不要で実装も軽いが、proxy やブラウザ設定の影響を受ける可能性があり、token 検証ほど明示的にリクエスト意図を確認できないため不採用
- Synchronizer Token パターン: server-side session に token を保存して検証するため堅牢だが、session store の導入が必要になる。このアプリケーションでは JWT による stateless な認証方針と実装規模を優先するため不採用
- Authorization header に Bearer token を付与して cookie 認証をやめる: CSRF リスクは下げやすいが、frontend 側で認証 token を保存・管理する必要があり、[008 Authentication](./008-authentication.md) の HttpOnly cookie 方針と合わないため不採用
- 同一 origin 配信のみを前提にする: CSRF と cookie 設定は単純になるが、現時点の開発環境とデプロイ構成では frontend と backend を別 origin として扱う必要があるため不採用

## References

- 関連ADR: [008-authentication.md](./008-authentication.md), [005-backend-architecture.md](./005-backend-architecture.md), [007-openapi.md](./007-openapi.md)
- 実装: `backend/internal/infra/web/gin/middleware/csrf.go`, `backend/internal/infra/web/gin/router/router.go`, `backend/internal/infra/web/gin/middleware/cors.go`, `backend/pkg/cookie/helper.go`, `frontend/src/lib/csrf-store.ts`, `frontend/src/lib/api-client-instance.ts`, `backend/api/openapi.yaml`
- 調査メモ: `README.md` の「モバイルブラウザからのログイン失敗」
