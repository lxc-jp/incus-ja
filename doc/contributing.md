# Incusにコントリビュートするには

% Include content from [../CONTRIBUTING.md](../CONTRIBUTING.md)
```{include} ../CONTRIBUTING.md
    :start-after: <!-- Include start contributing -->
    :end-before: <!-- Include end contributing -->
```

## 開発を始める

開発環境をセットアップし、Incusの新機能に取り組みを開始するには以下の手順に従ってください。

### 依存ライブラリのビルド

依存ライブラリをビルドするには{ref}`installing_from_source`の手順に従ってください。

### あなたのforkのremoteを追加

依存ライブラリをビルドし終わったら、GitHubのforkをremoteとして追加できます。

    git remote add myfork git@github.com:<your_username>/incus.git
    git remote update

次にこちらに切り替えます。

    git checkout myfork/main

### Incusのビルド

最後にリポジトリ内で`make`を実行すればこのプロジェクトのあなたのforkをビルドできます。

この時点で、あなたが最も行いたいであろうことはあなたのfork上にあなたの変更のための新しいブランチを作ることです。

```bash
git checkout -b [name_of_your_new_branch]
git push myfork [name_of_your_new_branch]
```

### Incusの新しいコントリビュータのための重要な注意事項

- 永続データは`INCUS_DIR`ディレクトリに保管されます。これは`incus admin init`で作成されます。
  `INCUS_DIR`のデフォルトは`/var/lib/incus`です。
- 開発中はバージョン衝突を避けるため、Incusのあなたのfork用に`INCUS_DIR`の値を変更すると良いでしょう。
- あなたのソースからコンパイルされる実行ファイルはデフォルトでは`$(go env GOPATH)/bin`に生成されます。
   - あなたの変更をテストするときはこれらの実行ファイル（インストール済みかもしれないグローバルの`incus admin`ではなく）を明示的に起動する必要があります。
   - これらの実行ファイルを適切なオプションを指定してもっと便利に呼び出せるよう`~/.bashrc`にエイリアスを作るという選択も良いでしょう。
- 既存のインストール済みIncusのデーモンを実行するための`systemd`サービスが設定されている場合はバージョン衝突を避けるためにサービスを無効にすると良いでしょう。

## ドキュメントへのコントリビュート

私たちはIncusをできるだけ簡単に使えるようにしたいと考えています。
ですので、Incusを使用するユーザーが必要とする（よくある使い方すべてをカバーし、典型的な質問に答えるような）情報を含んだドキュメントを提供を目指しています。

いろいろな方法でドキュメントにコントリビュートできます。
あなたのコントリビュートを感謝します！

コントリビュートする典型的な方法は以下のとおりです。

- コードにあなたがコントリビュートする新機能や機能改善についてのドキュメントを追加あるいは更新します。私たちはドキュメントの更新をレビューし、あなたのコードとともにマージします。
- Incusを使っているときに感じた疑問点を明確にするようなドキュメントを追加あるいは更新します。それらの修正はPull Requestまたはフォーラムの[Tutorials](https://discuss.linuxcontainers.org/c/tutorials/16)セクションへの投稿で送ってください。新しいチュートリアルはドキュメントへ含める（リンク経由で参照あるいは実際のコンテンツをインクルード）ことを検討します。
- ドキュメントの修正を依頼するために、[GitHub](https://github.com/canonical/incus/issues)でドキュメントに関するイシューを作成します。
- 質問や提案を[フォーラム](https://discuss.linuxcontainers.org)に投稿してください。
  私たちはイシューを評価し、適切にドキュメントを更新します。
- [IRC](https://web.libera.chat/#lxc)の`#lxc`チャンネルで質問したり提案してください。IRCの動的な性質のため、IRCの投稿に回答や反応することを保証はできませんが、チャンネルを注視して受け取ったフィードバックに基づいてドキュメントを改善するつもりです。

% Include content from [README.md](README.md)
```{include} README.md
    :start-after: <!-- Include start docs -->
```

Pull Requestをオープンすると、ドキュメントのプレビュー出力が自動的にビルドされます。

### ドキュメントの自動チェック

GitHubはドキュメントのスペル、リンク先が正しいか、Markdownファイルのフォーマットが正しいか、差別用語を使用していないか、を自動的にチェックします。

ローカルでも以下のコマンドでこれらのチェックができます（してください！）。

- スペルをチェックする：`make doc-spellcheck`
- リンクの有効性をチェックする：`make doc-linkcheck`
- Markdownのフォーマットをチェックする：`make doc-lint`
- 差別用語を使っていないかをチェックする：`make doc-woke`

### ドキュメントの設定オプション

```{note}
現在ドキュメントの設定オプションをコード内のコメントに移行中です。
現時点では、一部の設定オプションはこのアプローチに沿っていません。
```

ドキュメントの設定オプションはGoコード内のコメントから抽出されます。
コード内の`gendoc:generate`で始まるコメントを参照してください。

設定オプションを追加または変更する際は、それに必要なドキュメントのコメントを含めるようにしてください。

次に`make generate-config`を実行して`doc/config_options.txt`ファイルを再生成してください。更新されたファイルはチェックインしてください。

ドキュメントでは、設定オプションのグループを表示するために、`doc/config_options.txt`の一部をインクルードしています。
たとえば、コアのサーバオプションをインクルードするには以下のようにします。

````
% Include content from [config_options.txt](config_options.txt)
```{include} config_options.txt
    :start-after: <!-- config group server-core start -->
    :end-before: <!-- config group server-core end -->
```
````

既存のグループに設定を追加した場合は、ドキュメントファイルを更新する必要はありません。
新しいオプションは自動的に取り込まれます。
新しいグループを追加した場合のみ、ドキュメントファイルにインクルードを追加する必要があります。
