# Backend

## 概要

Go + Clean Architecture による Todo API サーバー

- **レイヤー分離**: Controller / Usecase / Gateway で責務を明確化
- **型安全**: OpenAPI 3.0 からコード自動生成
- **セキュリティ**: CSRF + JWT 多層認証
- **テスタブル**: インターフェース駆動で疎結合な設計

## アーキテクチャ

### ディレクトリ構造

```bash
backend/
  ├── entity/          # ドメインモデル（Task, User, Status）
  ├── usecase/         # ビジネスロジック（認証、検証）
  ├── adapter/
  │   ├── gateway/     # リポジトリ実装（DBアクセス）
  │   └── controller/  # HTTPハンドラ、ミドルウェア
  └── infrastructure/  # 外部依存（GORM, Ginサーバー）
```

## ローカル PC でのセットアップ

### 必要環境

- Go 1.25+
- PostgreSQL
- Docker
- make

#### 1. リポジトリのクローン

```bash
git clone https://github.com/Rtosshy/todo-app-next-go.git
cd todo-app-next-go
cd backend
```

#### 2. PostgreSQL の起動

```bash
# Docker で PostgreSQL を起動
make external-up
```

#### 3. 環境変数の設定

`.env.development`を用意

```env
DB_USER=app
DB_PASSWORD=password
DB_NAME=api_database
DB_HOST=localhost
DB_PORT=5432
DB_SSL_MODE=disable
```

カスタマイズが必要な場合は編集してください（デフォルト値は `infrastructure/database/config.go` 参照）

#### 4. サーバーの起動

```bash
make run
```

サーバーが `http://localhost:8080` で起動します。

---

### Docker を使った起動（推奨）

すべてを Docker で起動する場合：

```bash
# Backend + PostgreSQL + Swagger UI をまとめて起動
make docker-compose-up

# 停止
make docker-compose-down
```

起動後：

- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8001`

## 開発方法

### よく使うコマンド

```bash
# サーバー起動
make run

# OpenAPI 仕様からコード生成
make gen

# テスト実行
go test ./...

# コード品質チェック
make prettier  # lint + vet + fmt + imports

# 利用可能なコマンド一覧
make help
```

### OpenAPI 駆動開発のワークフロー

1. `api/openapi.yaml` を編集
2. `make gen` でコード生成
3. 生成された `adapter/controller/presenter/api.go` を基に実装

## AWS EC2 へのデプロイ

### 前提条件

- AWS EC2 インスタンス作成済み (Amazon Linux, t3.micro)
- RDS PostgreSQL 作成済み

#### 1. Go のインストール

```bash
# システムのアップデート
sudo yum update -y

# Goのダウンロード
cd /tmp
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz

# Goの展開
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz

# パスの設定
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# インストール確認
go version
# 出力: go version go1.25.5 linux/amd64
```

#### 2. Git のインストール

```bash
sudo yum install -y git
```

#### 3. アプリケーションのクローン

```bash
cd ~
git clone https://github.com/Rtosshy/todo-app-next-go.git
cd todo-app-next-go/backend
```

#### 4. 環境変数の設定

```bash
vi .env
```

以下の内容を入力

```env
APP_ENV=production
DB_USER=postgres
DB_PASSWORD=<RDS-master-password>
DB_NAME=todo_db
DB_HOST=<RDS-endpoint>
DB_PORT=5432
DB_DRIVER=postgres
DB_SSL_MODE=required
SECRET=<random-secret-key>
API_DOMAIN=
WEB_HOST=0.0.0.0
WEB_PORT=8080
WEB_CORS_ALLOW_ORIGINS=http://<frontend-server-dns>
```

#### 5. アプリケーションのビルド

ローカルでクロスコンパイル(推奨)

```bash
# ローカルマシンで
make build-linux

# EC2にアップロード
scp -i <key-file>.pem bin_linux ec2-user@<backend-ip>:~/todo-app-next-go/backend
```

EC2 で直接ビルド(時間がかかる)

```bash
# EC2で
cd ~/todo-app-next-go/backend
go mod download
go build -o bin_linux ./cmd/server/main.go
```

#### 6. 動作確認(テスト起動)

```bash
# テスト起動
APP_ENV=production ./bin_linux

# 別のターミナルでテスト
curl http://localhost:8080/api
```

#### 7. systemd サービスの作成

```bash
sudo vi /etc/systemd/system/todo-backend.service
```

以下の内容を入力

```ini
[Unit]
Description=Todo Backend API
After=network.target

[Service]
Type=simple
User=ec2-user
WorkingDirectory=/home/ec2-user/todo-app-next-go/backend
Environment="APP_ENV=production"
ExecStart=/home/ec2-user/todo-app-next-go/backend/bin_linux
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

保存後

```bash
# サービスの有効化と起動
sudo systemctl daemon-reload
sudo systemctl enable todo-backend
sudo systemctl start todo-backend

# 状態確認
sudo systemctl status todo-backend

# ログ確認
sudo journalctl -u todo-backend -f
```

## その他

### API 仕様

- OpenAPI 仕様: `api/openapi.yaml`
- Swagger UI: `http://localhost:8001` (docker-compose 起動時)

### Curl での API テストコマンド

```bash
# /api/health
curl http://localhost:8080/api/health
```

```bash
# /api/v1/csrf
curl http://localhost:8080/api/v1/csrf
```

```bash
# /api/v1/signup
curl -c cookies.txt -X POST http://localhost:8080/api/v1/signup \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Cookie: _csrf=<csrf-token>" \
  -d '{
    "user": {
        "email": "email@email.com",
        "password": "password"
    }
}'
```

```bash
# /api/v1/login
curl -c cookies.txt -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Cookie: _csrf=<csrf-token>" \
  -d '{
    "user": {
        "email": "email@email.com",
        "password": "email"
    }
}'
```

```bash
# /api/v1/tasks
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Cookie: _csrf=<csrf-token>; token=<jwt-token>" \
  -d '{
  "kind": "task",
  "name": "test",
  "status": {
    "name": "todo"
  }
}'
```

```bash
# /api/v1/tasks
curl -X GET http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: <csrf-token>" \
  -H "Cookie: _csrf=<csrf-token>; token=<jwt-token>"
```
