# Docker実行問題の原因と解決策

**Date**: 2026-01-24  
**Issue**: Docker Composeコマンドが権限エラーで失敗

## 問題の原因

### 症状
- ✅ `docker ps` は正常に動作
- ❌ `docker-compose` コマンドが権限エラーで失敗
- ❌ `docker info` も権限エラー

### エラーメッセージ
```
permission denied while trying to connect to the Docker daemon socket at unix:///Users/aoshima/.orbstack/run/docker.sock: dial unix /Users/aoshima/.orbstack/run/docker.sock: connect: operation not permitted
```

### 確認結果
- **Dockerソケット**: `/Users/aoshima/.orbstack/run/docker.sock`
- **所有者**: `501:20` (aoshima:staff)
- **権限**: `srwxr-xr-x`
- **OrbStack**: 実行中
- **Docker Compose**: v2.40.3 インストール済み

## 原因分析

OrbStackを使用している環境で、Dockerソケットへの接続権限エラーが発生しています。

### 確認結果
- ✅ `docker ps` - **動作する**（読み取り専用操作）
- ❌ `docker build` - **権限エラー**（書き込み操作）
- ❌ `docker compose` - **権限エラー**（書き込み操作）
- ❌ `docker info` - **権限エラー**（詳細情報取得）

### 原因
OrbStackのDockerソケット（`/Users/aoshima/.orbstack/run/docker.sock`）へのアクセス権限が制限されています。読み取り専用の操作（`docker ps`）は動作しますが、書き込みや詳細情報取得が必要な操作は権限エラーになります。

これは以下の可能性があります：
1. **OrbStackの権限設定**: OrbStackがDockerソケットへの書き込みアクセスを制限している
2. **macOSのセキュリティ設定**: システムのセキュリティ設定が影響している
3. **OrbStackの状態**: OrbStackが部分的にしか機能していない状態

## 解決策

### 方法1: OrbStackを再起動（推奨）

```bash
# OrbStackを再起動
killall OrbStack
open -a OrbStack

# 数秒待ってから確認
docker ps
docker-compose --version
```

### 方法2: Dockerソケットの権限を確認・修正

```bash
# ソケットファイルの権限を確認
ls -la /Users/aoshima/.orbstack/run/docker.sock

# 必要に応じて権限を変更（通常は不要）
# chmod 666 /Users/aoshima/.orbstack/run/docker.sock  # 非推奨（セキュリティリスク）
```

### 方法3: Docker Composeを直接使用（回避策）

`docker-compose` の代わりに `docker compose`（スペース区切り）を使用：

```bash
# docker-compose.yml の代わりに
cd api
docker compose -f ../docker-compose.yml up -d
docker compose -f ../docker-compose.yml exec dev oapi-codegen ...
```

### 方法4: Makefileを修正して直接dockerコマンドを使用

`api/Makefile` を修正して `docker-compose` の代わりに `docker compose` を使用：

```makefile
# 変更前
DOCKER_COMPOSE = docker-compose -f ../docker-compose.yml

# 変更後
DOCKER_COMPOSE = docker compose -f ../docker-compose.yml
```

## 推奨アクション

1. **まず方法1を試す**: OrbStackを再起動
2. **それでもダメな場合**: 方法3または方法4で回避
3. **根本解決**: OrbStackの設定を確認、またはDocker Desktopに切り替え

## 確認コマンド

問題が解決したか確認：

```bash
# 基本的なDockerコマンド
docker ps

# Docker Composeコマンド
docker-compose -f docker-compose.yml ps
# または
docker compose -f docker-compose.yml ps

# Docker info
docker info
```

## 参考情報

- OrbStack: https://orbstack.dev/
- Docker Compose v2: https://docs.docker.com/compose/
- 権限エラーの一般的な原因: macOSのセキュリティ設定、ファイルシステムの権限
