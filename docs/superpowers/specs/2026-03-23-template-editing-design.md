# ProgramTemplate編集と自動バージョン管理

## 概要

ProgramTemplateにスマートなバージョン管理付きの編集機能を追加する。Programから参照されていないテンプレートはその場で更新（in-place update）。Programが参照している場合は新バージョンを作成し、旧テンプレートを自動的にarchiveする。

## 動機

現在、ProgramTemplateは作成後に変更できない。ユーザーがテンプレートを編集したいケース：

- 処方箋のタイポや間違いの修正
- 次のトレーニングサイクルに向けたセット数/レップ数/RPEの微調整
- 他人から受け取ったテンプレートのカスタマイズ（将来のシェア機能）

編集機能がないと、テンプレートを削除して再作成するしかなく、そこから生成されたProgramとのリンクが失われる。

## 設計

### データモデル

`program_templates` テーブルに1カラム追加：

```sql
ALTER TABLE program_templates
  ADD COLUMN source_template_id UUID REFERENCES program_templates(id) ON DELETE SET NULL;
```

- Nullable — 元祖テンプレートや手動作成のものは `NULL`
- `ON DELETE SET NULL` — 派生元が削除されてもチェーンが優雅に切れるだけ

ドメイン構造体の変更：

```go
type ProgramTemplate struct {
    // ...既存フィールド...
    SourceTemplateID *uuid.UUID `json:"source_template_id,omitempty"`
}
```

リポジトリインターフェースの追加：

```go
Update(ctx context.Context, tmpl *domain.ProgramTemplate) error
```

### APIエンドポイント

```
POST /program-templates/{id}/edit
```

リクエストボディ：`ProgramTemplateCreate` と同じ構造（name, description, notes, metadata, weeks, days_per_week, entries）。

| ケース | ステータス | ボディ |
|------|--------|------|
| 紐づくProgramなし → in-place update | 200 OK | 更新後のテンプレート |
| 紐づくProgramあり → 新バージョン + archive | 201 Created | 新テンプレート（新ID） |
| テンプレート未発見 | 404 | |
| アーカイブ済みテンプレート | 400 | |
| バリデーションエラー | 400 | |

バージョニング時に新しいリソース（新ID）を作成するため冪等ではなく、`PUT` ではなく `POST` を使用する。

### OpenAPIスキーマ変更

1. **新エンドポイント**: `POST /program-templates/{id}/edit`、リクエストボディは `ProgramTemplateCreate`
2. **レスポンススキーマ**: 既存の `ProgramTemplate` スキーマに `source_template_id`（nullable UUID）を追加
3. **リクエストボディ**: `ProgramTemplateCreate` スキーマを再利用（新しいリクエスト型は不要）

### バックエンド処理フロー

```
EditProgramTemplate(id, body):
  1. 既存テンプレート取得 (GetByID)
     → 404 未発見
     → 400 アーカイブ済み

  2. ボディのバリデーション (ValidateProgramTemplate)
     → 400 バリデーション失敗

  3. 紐づくProgram確認 (ExistsByProgramTemplateID)

  4a. 紐づくProgramなし:
      BEGIN TRANSACTION
      - DELETE FROM program_template_entries WHERE program_template_id = id
      - UPDATE program_templates SET name=..., updated_at=NOW() WHERE id = id
      - INSERT program_template_entries (新しいentries)
      COMMIT
      → return 200, 更新後テンプレート

  4b. 紐づくProgramあり:
      BEGIN TRANSACTION
      - 新テンプレート作成 (新UUID, source_template_id = 旧ID)
      - 新エントリーを新UUIDで作成
      - 旧テンプレートをarchive (archived_atを設定)
      COMMIT
      → return 201, 新テンプレート
```

### フロントエンドUIフロー

**Editボタン押下時：**

```
[Editボタン]
    │
    ├── GET /programs?program_template_id={id}&limit=1
    │
    ├── data.length == 0
    │     → エディタ画面へ遷移
    │
    └── data.length > 0
          → ダイアログ: 「このテンプレートは使用中のプログラムがあります。
             編集すると新しいバージョンとして保存されます。」
            [キャンセル] [続行]
          → 続行 → エディタ画面へ遷移
```

**保存時：**

```
[Save]
    │
    ├── POST /program-templates/{id}/edit
    │
    ├── 200 OK → そのまま詳細画面（IDは変わらない）
    │
    └── 201 Created → 新テンプレートのIDで詳細画面にリダイレクト
```

**エディタ:** 既存のテンプレート作成コンポーネント（SessionAccordion等）を再利用し、既存テンプレートの内容をプリフィルする。

**詳細画面:** `source_template_id` がある場合、「派生元: {元テンプレート名}」リンクを小さく表示する。

### 編集動作まとめ

| 紐づくProgram | 旧テンプレート | 新テンプレート | ユーザー体験 |
|-----------------|-------------|-------------|-----------------|
| なし | その場で更新 | N/A | 通常の編集 |
| あり（全ステータス） | アーカイブ | source_template_id付きで作成 | ダイアログ付き新バージョン |

### 紐づくProgramの範囲

`program_template_id` が一致する全Program（ステータス問わず：created, ongoing, completed, cancelled）を対象とする。

## スコープ外

- シェア用のExport/Import
- 整合性検証用のコンテンツハッシュ
- バージョン履歴一覧UI（「派生元」リンクのみ）
- バッチ編集
- テンプレートのマージ

## 参考

- [GitHub Issue #141](https://github.com/aoshimash/rx/issues/141)
- `docs/DOMAIN_MODEL.md` — 3層ライフサイクル
- `docs/PHILOSOPHY.md` — 「Dumb Backend」とAPI-Firstの原則
