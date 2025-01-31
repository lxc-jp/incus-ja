# ドキュメントへのコントリビュート

私たちは Incus をできるだけ簡単に使えるようにしたいと考えています。
ですので、Incus を使用するユーザーが必要とする（よくある使い方すべてをカバーし、典型的な質問に答えるような）情報を含んだドキュメントを提供を目指しています。

いろいろな方法でドキュメントにコントリビュートできます。
あなたのコントリビュートを感謝します！

コントリビュートする典型的な方法は以下のとおりです。

- コードにあなたがコントリビュートする新機能や機能改善についてのドキュメントを追加あるいは更新します。私たちはドキュメントの更新をレビューし、あなたのコードとともにマージします。
- Incus を使っているときに感じた疑問点を明確にするようなドキュメントを追加あるいは更新します。それらの修正は Pull Request またはフォーラムの[Tutorials](https://discuss.linuxcontainers.org/c/tutorials/16)セクションへの投稿で送ってください。新しいチュートリアルはドキュメントへ含める（リンク経由で参照あるいは実際のコンテンツをインクルード）ことを検討します。
- ドキュメントの修正を依頼するために、[GitHub](https://github.com/canonical/incus/issues)でドキュメントに関するイシューを作成します。
- 質問や提案を[フォーラム](https://discuss.linuxcontainers.org)に投稿してください。
  私たちはイシューを評価し、適切にドキュメントを更新します。
- [IRC](https://web.libera.chat/#lxc)の`#lxc`チャンネルで質問したり提案してください。IRC の動的な性質のため、IRC の投稿に回答や反応することを保証はできませんが、チャンネルを注視して受け取ったフィードバックにもとづいてドキュメントを改善するつもりです。

% Include content from [../README.md](../README.md)
```{include} ../README.md
    :start-after: <!-- Include start docs -->
```

Pull Request をオープンすると、ドキュメントのプレビュー出力が自動的にビルドされます。

## ドキュメントの自動チェック

GitHub はドキュメントのスペル、リンク先が正しいか、Markdown ファイルのフォーマットが正しいか、差別用語を使用していないか、を自動的にチェックします。

ローカルでも以下のコマンドでこれらのチェックができます（してください！）。

- スペルをチェックする：`make doc-spellcheck`
- リンクの有効性をチェックする：`make doc-linkcheck`
- Markdown のフォーマットをチェックする：`make doc-lint`
- 差別用語を使っていないかをチェックする：`make doc-woke`

上記を実行するためには、以下のものが必要です:

- Python 3.8以降
- `venv` pythonパッケージ
- スペルチェック用に`aspell`ツール
- `mdl` markdown lintツール

## ドキュメントの設定オプション

```{note}
現在ドキュメントの設定オプションをコード内のコメントに移行中です。
現時点では、一部の設定オプションはこのアプローチに沿っていません。
```

ドキュメントの設定オプションは Go コード内のコメントから抽出されます。
コード内の`gendoc:generate`で始まるコメントを参照してください。

設定オプションを追加または変更する際は、それに必要なドキュメントのコメントを含めるようにしてください。

次に`make generate-config`を実行して`doc/config_options.txt`ファイルを再生成してください。更新されたファイルはチェックインしてください。

ドキュメントでは、設定オプションのグループを表示するために、`doc/config_options.txt`の一部をインクルードしています。
たとえば、コアのサーバーオプションをインクルードするには以下のようにします。

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
