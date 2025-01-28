# コードへのコントリビュート

開発環境をセットアップし、Incus の新機能に取り組みを開始するには以下の手順に従ってください。

## 依存ライブラリのビルド

依存ライブラリをビルドするには{ref}`installing_from_source`の手順に従ってください。

## あなたのforkのremoteを追加

依存ライブラリをビルドし終わったら、GitHub の fork を remote として追加できます。

    git remote add myfork git@github.com:<your_username>/incus.git
    git remote update

次にこちらに切り替えます。

    git checkout myfork/main

## Incusのビルド

最後にリポジトリ内で`make`を実行すればこのプロジェクトのあなたの fork をビルドできます。

この時点で、あなたが最も行いたいであろうことはあなたの fork 上にあなたの変更のための新しいブランチを作ることです。

```bash
git checkout -b [name_of_your_new_branch]
git push myfork [name_of_your_new_branch]
```

## Incusの新しいコントリビュータのための重要な注意事項

- 永続データは`INCUS_DIR`ディレクトリーに保管されます。これは`incus admin init`で作成されます。
  `INCUS_DIR`のデフォルトは`/var/lib/incus`です。
- 開発中はバージョン衝突を避けるため、Incus のあなたの fork 用に`INCUS_DIR`の値を変更すると良いでしょう。
- あなたのソースからコンパイルされる実行ファイルはデフォルトでは`$(go env GOPATH)/bin`に生成されます。
   - あなたの変更をテストするときはこれらの実行ファイル（インストール済みかもしれないグローバルの`incus admin`ではなく）を明示的に起動する必要があります。
   - これらの実行ファイルを適切なオプションを指定してもっと便利に呼び出せるよう`~/.bashrc`にエイリアスを作るという選択も良いでしょう。
- 既存のインストール済み Incus のデーモンを実行するための`systemd`サービスが設定されている場合はバージョン衝突を避けるためにサービスを無効にすると良いでしょう。
