# 日本語訳の作業について

## ドキュメントのビルド

フルビルドするには以下のコマンドを実行します。
``` 
make doc
``` 

差分ビルドするには以下のコマンドを実行します。
```
make doc-incremental
```

ビルドが正常に完了したら、以下のコマンドを実行するとローカルでウェブサーバが起動します。これによりビルドしたドキュメントをブラウザで確認できます。

```
make doc-serve
```

なお、ドキュメントのビルド時にMakefile内のclientターゲットで一部のコマンドをビルドするようになってます。
そのためincusのソースファイルも消さずに残してあります。

## textlintによる日本語訳のチェック

textlintについては以下の記事の参照してください。

* [textlintで日本語の文章をチェックする | Web Scratch](https://efcl.info/2015/09/10/introduce-textlint/)
* [textlint + prhで表記ゆれを検出する | Web Scratch](https://efcl.info/2015/09/14/textlint-rule-prh/)

### 初回セットアップ

[Getting Started | Volta](https://docs.volta.sh/guide/getting-started)の手順でVoltaをインストールし
[pnpm Support | Volta](https://docs.volta.sh/advanced/pnpm)の設定をしてください。

レポジトリ内で以下のコマンドを実行してください。

```
volta pin node@18
volta pin pnpm@latest
```

その後以下のコマンドを実行して、textlintとルールのパッケージをインストールします。

```
pnpm install
```

### textlint実行方法

以下のコマンドでtextlintを実行してください。

```
npm run textlint
```

### このプロジェクトでのtextlintのルール設定

設定ファイルは `.textlintrc.json` にあります。

ルール一覧は[Collection of textlint rule · textlint/textlint Wiki](https://github.com/textlint/textlint/wiki/Collection-of-textlint-rule)を参照してください。

[prh/prh: proofreading helper](https://github.com/prh/prh)用のルールは[lxc-jp/textlint-prh-rules: A collection of prh rules](https://github.com/lxc-jp/textlint-prh-rules)をgit submoduleで`.textlint-prh-rules/`に配置しています。

