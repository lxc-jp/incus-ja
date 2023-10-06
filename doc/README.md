# Incusドキュメント

Incusの日本語ドキュメントは、<https://incus-ja.readthedocs.io/ja/latest/>(原文のドキュメントは<https://linuxcontainers.org/incus/docs/latest/>)で閲覧できます。

GitHubでもドキュメントの基本的なレンダリングを提供していますが、includeやクリッカブルリンクなどの重要な機能が欠落しています。そのため、[公開ドキュメント](https://incus-ja.readthedocs.io/ja/latest/)を読むことをお勧めします。

## どのように動作するか

<!-- Include start docs -->

### ドキュメントのフレームワーク

Incusのドキュメントは[Sphinx](https://www.sphinx-doc.org/en/master/index.html)でビルドされます。

ドキュメントは[Markdown](https://commonmark.org/)と[MyST](https://myst-parser.readthedocs.io/)の拡張で書かれています。
構文のヘルプやガイドラインについては、[ドキュメントチートシート](https://incus-ja.readthedocs.io/ja/latest/doc-cheat-sheet/) ([ソース](https://raw.githubusercontent.com/lxc-jp/incus-ja/main/doc/doc-cheat-sheet.md))を参照してください。

構成に関しては、このドキュメントでは[Diátaxis](https://diataxis.fr/)アプローチを採用しています。

### ドキュメントのビルド

ドキュメントをビルドするには、リポジトリのルートディレクトリから`make doc`を実行します。このコマンドは必要なツールをインストールして、出力を`doc/html/`ディレクトリにレンダリングします。
変更されたファイルのみを対象に（ツールを再インストールすることなく）ドキュメントを更新するには、`make doc-incremental`を実行します。

Pull Requestをオープンする前に、ドキュメントが警告なしでビルドできることを確認してください（警告はエラーとして扱われます）。
ドキュメントをローカルでプレビューするには、`make doc-serve`を実行し[`http://localhost:8001`](http://localhost:8001)をブラウザで開いてください。
