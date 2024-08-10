(incus-alias)=
# コマンドエイリアスを管理するには

Incus コマンドラインクライアントでは良く使うコマンドのエイリアスを追加できます。
長いコマンドのショートカットとして、あるいは既存のコマンドに自動的にフラグを追加するために、エイリアスを使用できます。

エイリアスの管理には[`incus alias`](incus_alias.md)コマンドを使います。

[`incus alias`](incus_alias.md)コマンドでは、以下のサブコマンドが使えます:

- 新しいコマンドエイリアスの追加は`incus alias add`
- コマンドエイリアスの一覧表示は`incus alias list`
- コマンドエイリアスの削除は`incus alias remove`
- コマンドエイリアスの改名は`incus alias rename`

すべての利用可能なサブコマンドとパラメーターを見るには[`incus alias --help`](incus_alias.md)を実行してください。

## コマンドエイリアスを追加するには

インスタンスを削除する際に必ず確認を求めるようにするには[`incus delete`](incus_delete.md)に常に`incus delete --interactive`を実行するようにエイリアスを作成します。

以下のコマンドは`delete`という名前でコマンドエイリアスを_追加_し同じIncusのコマンドを`--interactive`フラグつきで実行します。

    incus alias add delete "delete --interactive"

`myinstance`と呼ばれるインスタンスを削除するために`incus delete mycontainer`を実行した際にIncusのコマンドラインクライアントは`incus delete`を`incus delete --interactive`に置き換えて`incus delete --interactive myinstance`を実行することに注意してください。

コマンドエイリアスをIncusコマンドと同じ名前で登録すると、コマンドエイリアスはIncusコマンドを隠します。

文字通りに同じ名前のIncusコマンドを実行したい場合は、まずコマンドエイリアスを削除する必要があります。
さらに、パラメータ（上の例ではコンテナの名前）つきのコマンドエイリアスを使う場合、`@ARGS`という文字列でパラメータを別の場所に手動で置かない限り、Incusのコマンドラインクライアントはパラメータをエイリアスされたコマンドの最後に置きます。

最後に、コマンドエイリアス内のコマンドはクォートで囲む必要があります。

## すべてのコマンドエイリアスを一覧表示するには

設定されたすべてのエイリアスを見るには、[`incus alias list`](incus_alias_list.md)を実行します。

## コマンドエイリアスを削除するには

既存のコマンドエイリアスを削除するには[`incus alias remove`](incus_alias_remove.md)にコマンドエイリアスの名前を追加して実行します。

## コマンドエイリアスを改名するには

既存のコマンドエイリアスを改名するには、[`incus alias rename`](incus_alias_rename.md)に既存のコマンドエイリアスの名前と新しいコマンドエイリアスの名前を指定して実行します。

## ビルトインの`shell`エイリアス

Incusは`shell`というビルトインのコマンドエイリアスがあります。このエイリアスは[`incus exec`](incus_exec.md)コマンドをベースにしており、`exec @ARGS@ -- su -l`を実行します。

`incus shell myinstance`を実行すると、このコマンドは`incus exec mycontainer -- su -l`に展開されます。

```
$ incus alias list
+-----------+----------------------+
|   ALIAS   |        TARGET        |
+-----------+----------------------+
| shell     | exec @ARGS@ -- su -l |
+-----------+----------------------+
```

このコマンドエイリアスのターゲットは`exec @ARGS@ -- su -l`です。

`--`は`-l`のようなパラメータを処理しないように指示するコマンドラインの約束事です。`--`がないと、展開された`incus exec mycontainer su -l`というコマンドは失敗します。これは`-l`を処理しないのではなくIncusコマンドラインクライアントが解釈しようとしてしまうからです。

`su -l`コマンドは`su -`や`su --login`と動議です。
ログインシェルを`root`ユーザーでインスタンス内に起動します。
コマンドは`root`ユーザーでログインシェルを起動するために必要な設定ファイルを読みます。

`shell`エイリアスはIncusサーバーにビルトインされています。そのため、Incusクライアントでは削除できません。
削除しようとすると、エイリアスが存在しないというエラーになります。

```
$ incus alias remove shell
Error: Alias shell doesn't exist
$
```

`shell`という名前でエイリアスを登録すると、ビルトインのコマンドエイリアスを隠すことになります。
つまり、Incusコマンドラインクライアントは新しく追加されたエイリアスを使い、代わりにビルトインのコマンドエイリアスは隠されます。追加した`shell`エイリアスを削除すると、ビルトインのエイリアスが再び現れます。

## インスタンス内で非rootのシェルを起動するコマンドエイリアスを使うには

いくつかのIncusイメージは以下の表に示すように非rootのユーザー名を作成するように設定されています。

| ディストリビューション          | ユーザー         | イメージ |
| :----------- | :--------------: | :----------- |
| Alpine | `alpine` | `images:alpine/edge/cloud` |
| Debian | `debian` | `images:debian/12/cloud` |
| Fedora | `fedora` | `images:fedora/40/cloud` |
| Ubuntu | `ubuntu` | `images:ubuntu/24.04/cloud` |

以下のコマンドで非rootのユーザー名でインスタンス内のシェルを起動できます。

```
$ incus launch images:debian/12/cloud mycontainer
Launching mycontainer
$ incus exec mycontainer -- su -l debian
debian@mycontainer:~$
```

Incusコマンドエイリアスを使うことで、そのインスタンスへのシェルを起動するコマンドエイリアスも作れます。
次のコマンドエイリアスでは、`debian`というユーザー名に`su -l`するように指定しています。

```
$ incus alias add debian 'exec @ARGS@ -- su -l debian'
$
```

これで、以下の便利なコマンドでインスタンス内にシェルを起動できます:

```
$ incus debian mycontainer
debian@mycontainer:~$
```

## その他

_コマンドエイリアス_は_イメージエイリアス_とは違うことに注意してください。
イメージエイリアスはイメージの別名で、通常はより短いな目やそのイメージの別の一般的なニーモニックです。

イメージエイリアスはIncus APIの一部でサーバーサイドの概念ですが、コマンドエイリアスは純粋にコマンドラインツールの設定です。
