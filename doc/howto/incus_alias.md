(incus-alias)=
# コマンドエイリアスを管理するには

Incus コマンドラインクライアント`incus`では良く使うコマンドのエイリアスを追加できます。
長いコマンドのショートカットとして、あるいは既存のコマンドに自動的にフラグを追加するために、エイリアスを使用できます。

エイリアスの管理には[`incus alias`](incus_alias.md)コマンドを使います。

[`incus alias`](incus_alias.md)コマンドでは、以下のサブコマンドが使えます:

- 新しいコマンドエイリアスの追加は`incus alias add`
- コマンドエイリアスの一覧表示は`incus alias list`
- コマンドエイリアスの削除は`incus alias remove`
- コマンドエイリアスのリネームは`incus alias rename`

すべての利用可能なサブコマンドとパラメーターを見るには[`incus alias --help`](incus_alias.md)を実行してください。

```{note}
_コマンドエイリアス_ は _{ref}`イメージエイリアス <images>`_ とは異なります。
イメージエイリアスはイメージの別名で、通常はそのイメージのより短い名前や別の覚えやすい名前です。

イメージエイリアスはサーバーサイドの概念でIncus APIの一部ですが、コマンドエイリアスは純粋にコマンドラインツールの設定の一部です。
```

## コマンドエイリアスを追加するには

エイリアスを作成するには、[`incus alias add`](incus_alias_add.md)を実行し、エイリアス名と（クォートで囲んだ）エイリアスコマンドを指定します。

```
$ incus alias add my-alias "image list"
$ incus my-alias
┌───────┬──────────────┬────────┬────────────────────────────────────────┬──────────────┬───────────┬───────────┬──────────────────────┐
│ ALIAS │ FINGERPRINT  │ PUBLIC │              DESCRIPTION               │ ARCHITECTURE │   TYPE    │   SIZE    │     UPLOAD DATE      │
├───────┼──────────────┼────────┼────────────────────────────────────────┼──────────────┼───────────┼───────────┼──────────────────────┤
│       │ 3b3bd7f47fca │ no     │ Debian bookworm amd64 (20260608_05:24) │ x86_64       │ CONTAINER │ 106.22MiB │ 2026/06/08 22:01 -03 │
└───────┴──────────────┴────────┴────────────────────────────────────────┴──────────────┴───────────┴───────────┴──────────────────────┘
```

コマンドエイリアスがIncusコマンドと同じ名前の場合、コマンドエイリアスはIncusコマンドを隠します。同じ名前の元のIncusコマンドを実行するには、まずコマンドエイリアスを削除する必要があります。

## すべてのコマンドエイリアスを一覧表示するには

設定されたすべてのエイリアスを見るには、[`incus alias list`](incus_alias_list.md)を実行します。

```
$ incus alias list
┌──────────┬────────────┐
│  ALIAS   │   TARGET   │
├──────────┼────────────┤
│ my-alias │ image list │
└──────────┴────────────┘
```

## コマンドエイリアスを削除するには

既存のコマンドエイリアスを削除するには[`incus alias remove`](incus_alias_remove.md)にコマンドエイリアスの名前を追加して入力します。

```
$ incus alias remove my-alias
$ incus alias list
┌────────┬────────┐
│ ALIAS  │ TARGET │
└────────┴────────┘
```

## コマンドエイリアスをリネームするには

既存のコマンドエイリアスをリネームするには、[`incus alias rename`](incus_alias_rename.md)に既存のコマンドエイリアスの名前と新しいコマンドエイリアスの名前を指定して入力します。

```
$ incus alias rename my-alias my-new-alias
$ incus alias list
┌──────────────┬────────────┐
│    ALIAS     │   TARGET   │
├──────────────┼────────────┤
│ my-new-alias │ image list │
└──────────────┴────────────┘
```

## エイリアスコマンド内の引数
コマンドエイリアスを引数とともに使う場合、Incusコマンドラインクライアントはそれらの引数をエイリアスコマンドの最後に渡します。`incus alias add del "delete"`のエイリアスの場合、以下のコマンドはどちらも同じ結果になります。

```
incus delete c1 --force
incus del c1 --force
```

この挙動は`@ARGS@`という特別な文字列で変更できます。これはエイリアスの文字列内で定義された場所にすべての引数を配置します。`incus alias add create-foo "create @ARGS@ foo"`というエイリアスであれば、以下のコマンドはどちらも同じ結果になります。

```
incus create images:debian/12 foo
incus create-foo images:debian/12
```

番号付きの引数（`@ARG1@`、`@ARG2@`、…）を使うこともでき、これはエイリアスの文字列内のこれらの位置に指定した順に配置します。`incus alias add cat "exec @ARG1@ -- cat @ARG2@"`というエイリアスであれば、以下のコマンドはどちらも同じ結果になります。

```
incus exec u1 -- cat /etc/hosts
incus cat u1 /etc/hosts
```

## ビルトインの`shell`エイリアス

Incusは`shell`というビルトインのコマンドエイリアスがあります。このエイリアスは[`incus exec`](incus_exec.md)コマンドをベースにしており、`exec @ARGS@ -- su -l`を実行します。

```
$ incus alias list
┌───────┬──────────────────────┐
│ ALIAS │        TARGET        │
├───────┼──────────────────────┤
│ shell │ exec @ARGS@ -- su -l │
└───────┴──────────────────────┘
```

`incus shell myinstance`を実行すると、このコマンドは`incus exec myinstance -- su -l`に展開されます。

`--`は`-l`のようなパラメータを処理しないように指示するIncusのコマンドラインの約束事です。`--`がないと、展開された`incus exec mycontainer su -l`というコマンドはIncusコマンドクライアントが`-l`をパースしようとするため失敗します。この特定のケースでは`incus shell`に`-l`というパラメータはないため失敗します。

`su -l`コマンドは`su -`や`su --login`と同義です。
ログインシェルを`root`ユーザーでインスタンス内に起動します。
コマンドは`root`ユーザーでログインシェルを起動するために必要な設定ファイルを読みます。

`shell`エイリアスはIncusサーバーにビルトインされています。そのため、Incusクライアントでは削除できません。
削除しようとすると、エイリアスが存在しないというエラーになります。

```
$ incus alias remove shell
Error: Alias shell doesn't exist
```

`shell`という名前でエイリアスを登録すると、新しいコマンドはビルトインのコマンドエイリアスを隠すことになります。
つまり、Incusコマンドラインクライアントは新しく追加されたエイリアスを使い、代わりにビルトインのコマンドエイリアスは隠されます。追加した`shell`エイリアスを削除すると、ビルトインのエイリアスが再び現れます。

## 使い方の例
### インスタンスを削除する際に確認するようにするには
インスタンスを削除する際に必ず確認を求めるようにするには、[`incus delete`](incus_delete.md)に
常に`incus delete --interactive`を実行するようにエイリアスを作成します。

以下のコマンドは`delete`という名前でコマンドエイリアスを _追加_ し同じIncusのコマンドを`--interactive`フラグつきで実行します。

    incus alias add delete "delete --interactive"

`myinstance`と呼ばれるインスタンスを削除するために、`incus delete mycontainer`を実行した際にIncusのコマンドラインクライアントは`incus delete`を`incus delete --interactive`に置き換えて、代わりに`incus delete --interactive myinstance`を実行することに注意してください。

### インスタンス内で非rootのシェルを起動するコマンドエイリアスを使うには

いくつかのIncusイメージは以下の表に示すように非rootのユーザー名を作成するように設定されています。

| ディストリビューション          | ユーザー         | イメージ |
| :----------- | :--------------: | :----------- |
| Alpine | `alpine` | `images:alpine/edge/cloud` |
| Debian | `debian` | `images:debian/12/cloud` |
| Fedora | `fedora` | `images:fedora/42/cloud` |
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
incus alias add debian 'exec @ARGS@ -- su -l debian'
```

これで、以下の便利なコマンドでインスタンス内にシェルを起動できます:

```
$ incus debian mycontainer
debian@mycontainer:~$
```

```{note}
`su`の代わりとして、`sudo`を使いたいこともあるでしょう。その場合コマンドは以下のようになります。

     incus alias add debian `exec @ARGS@ -- sudo --login --user debian`
```

```{note}
システムコンテナや仮想マシンを起動する際に、Incusでは環境変数を指定できます。

     incus launch -c environment.MYVARIABLE=myvalue images:debian/12 myinstance

そのようなインスタンスのログインシェルではこれらの環境変数へはアクセスできません。これは`su -l`や`sudo --login`でのログインシェルのセマンティクスでは環境変数は一切維持しないからです。環境変数を維持したい場合は、代わりに`su --preserve-environment`か`sudo --preserve-env`を使う必要があります。

インスタンスに環境変数を追加する別の方法はファイルシステムの`/etc/environment`というファイルに書くことです。そうすうることでインスタンスへの新しいログインシェルはこのファイルをパースし環境変数を設定します。
```

## その他

_コマンドエイリアス_ は _イメージエイリアス_ とは違うことに注意してください。
イメージエイリアスはイメージの別名で、通常はより短いな目やそのイメージの別の一般的なニーモニックです。

イメージエイリアスはIncus APIの一部でサーバーサイドの概念ですが、コマンドエイリアスは純粋にコマンドラインツールの設定です。
