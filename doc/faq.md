# よく聞かれる質問（FAQ）

以下のセクションは、よくある質問への回答を提供します。
それらは一般的な問題の解決方法を説明し、より詳細な情報へと導きます。

## なぜ私のインスタンスはネットワークアクセスがないのですか？

最も可能性が高いのは、あなたのファイアウォールがインスタンスのネットワークアクセスをブロックしているためです。
問題とその修正方法についての詳細は {ref}`network-bridge-firewall` をご覧ください。

接続問題の別の一般的な原因は、Incus と Docker を同じホスト上で実行していることです。
このような問題を修正する方法については {ref}`network-incus-docker` を参照してください。

## Incus サーバーをリモートアクセス可能にするにはどうすればよいですか？

デフォルトでは、 Incus サーバーはネットワークからアクセスできません。なぜなら、それはローカルの Unix ソケットでしかリッスンしていないからです。

リモートアクセスを可能にするためには、 {ref}`server-expose` の指示に従ってください。

## `incus remote add`を行うと、トークンを求められるのはなぜですか？

リモート API にアクセスするためには、クライアントは Incus サーバーに対して認証を行わなければなりません。

トラストトークンを使用して認証する方法については {ref}`server-authenticate` を参照してください。

## なぜ特権コンテナを実行すべきではないのですか？

特権コンテナは、ホスト全体に影響を与えることができます - 例えば、`/sys`内のものを使ってネットワークカードをリセットすると、ホスト全体のそれがリセットされ、ネットワークが一時的に断線します。
詳細は {ref}`container-security` をご覧ください。

ほとんどのものは非特権コンテナで実行できます。また、NFS ファイルシステムをコンテナ内にマウントしたいなど、通常とは異なる特権を必要とするものの場合、バインドマウントを使用する必要があるかもしれません。

## ホームディレクトリーをコンテナにバインドマウントすることはできますか？

はい、それは{ref}`ディスクデバイス <devices-disk>`を使用することで可能です:

    incus config device add container-name home disk source=/home/${USER} path=/home/ubuntu

非特権コンテナの場合、コンテナ内のユーザーが適切な読み書き権限を持っていることを確認する必要があります。
そうでないと、すべてのファイルはオーバーフローUID/GID（`65536:65536`）として表示され、ワールドリーダブルでないものへのアクセスは失敗します。
必要な権限を付与するために以下の方法のいずれかを使用してください:

- [`incus config device add`](incus_config_device_add.md)の実行時に`shift=true`を指定します。これはカーネルとファイルシステムが idmapped マウントあるいは shiftfs をサポートしているかに依存します（ [`incus info`](incus_info.md)参照）。
- `raw.idmap`エントリを追加します（[User Namespace の Idmap](userns-idmap.md)参照）。
- ホームディレクトリーに再帰的な POSIX ACL を配置します。

特権コンテナはこの問題を持っていません、なぜならコンテナ内のすべての UID/GID は外部と同じだからです。
しかし、それが特権コンテナのセキュリティー問題のほとんどの原因でもあります。

## Incus コンテナの内部で Docker を実行するには？

Incus コンテナの内部で Docker を実行するには、コンテナの {config:option}`instance-security:security.nesting` プロパティを `true` にセットします:

    incus config set <container> security.nesting true

Incus コンテナはカーネルモジュールをロードできないため、 Docker の設定によっては、ホストで追加のカーネルモジュールをロードする必要があるかもしれません。
コンテナが必要とするカーネルモジュールのカンマ区切りのリストを設定することでこれを行うことができます:

    incus config set <container_name> linux.kernel_modules <modules>

さらに、コンテナ内に`/.dockerenv`ファイルを作成すると、Docker がネストした環境で実行されているために発生するいくつかのエラーを無視するのに役立ちます。

## Incus クライアント（`incus`）は設定をどこに保存しますか？

[`incus`](incus.md) コマンドはその設定を `~/.config/incus` に保存します。

様々な設定ファイルがそのディレクトリーに保存されます。例えば:

- `client.crt`：クライアント証明書（要求に応じて生成されます）
- `client.key`：クライアントキー（要求に応じて生成されます）
- `config.yml`：設定ファイル（`remotes`、`aliases`などの情報）
- `servercerts/`：`remotes`に関連するサーバー証明書が保存されているディレクトリー

## なぜ他のホストから Incus インスタンスに ping を送ることができないのですか？

多くのスイッチは MAC アドレスの変更を許可せず、不正な MAC を持つトラフィックをドロップするか、ポートを完全に無効にします。
ホストから Incus インスタンスには ping を送ることができますが、異なるホストから ping を送ることができない場合、これが原因かもしれません。

この問題を診断する方法は、アップリンク上で`tcpdump`を実行することで、``ARP Who has `xx.xx.xx.xx` tell `yy.yy.yy.yy` ``が表示され、レスポンスを送信しているにもかかわらず確認されていない、または ICMP パケットが成功裏に送受信されているにもかかわらず、他のホストには受け取られていないことを確認することです。

(faq-monitor)=
## Incusが何をしているかモニターするには？

Incus が何をしているかとどんなプロセスが稼働しているかについての詳細な情報を見るには、[`incus monitor`](incus_monitor.md)コマンドを使います。

たとえば、すべてのタイプのメッセージの出力を人間が見やすい形式で表示するには、以下のコマンドを使用します:

    incus monitor --pretty

すべてのオプションについては [`incus monitor --help`](incus_monitor.md) を、より詳しい情報は {doc}`debugging` を参照してください。

## インスタンス作成時に Incus が止まってしまうのはなぜですか？

ストレージプールの空きが無くなってないか（`incus storage info <pool_name>`を実行して）確認してください。
空きが無い場合、 Incus はイメージの展開ができず、作成しようとしているインスタンスは止まったままに見えます。

何が起きているかをより詳しく調べるには [`incus monitor`](incus_monitor.md) を実行し（{ref}`faq-monitor`参照）、`sudo dmesg`で何か I/O エラーが起きていないか確認してください。

## コンテナの起動が突然失敗するようになったのはなぜ？

コンテナの起動が cgroup 関連のエラーメッセージ（`Failed to mount "/sys/fs/cgroup"`）で失敗する場合、ホスト上で VPN クライアントが稼働しているためかもしれません。

これは [Mullvad VPN](https://github.com/mullvad/mullvadvpn-app/issues/3651) と [Private Internet Access VPN](https://github.com/pia-foss/desktop/issues/50) の両方で知られた問題ですが、他の VPN クライアントでも起きるかもしれません。
問題は VPN クライアントが（Incus が使用する）cgroiup2 上に `net_cls` cgroup1 をマウントすることです。

この問題の一番簡単な修正方法は VPN クライアントを停止し、以下のコマンドで `net_cls` cgroup1 をアンマウントすることです:

    umount /sys/fs/cgroup/net_cls

VPN クライアントを稼働したままにする必要がある場合、 `net_cls` cgroup1 を他の場所にマウントし、 VPN クライアントを適宜再設定してください。
Mullvad VPN 用の手順は [この Discourse の投稿](https://discuss.linuxcontainers.org/t/help-help-help-cgroup2-related-issue-on-ubuntu-jammy-with-mullvad-and-privateinternetaccess-vpn/14705/18) を参照してください。
