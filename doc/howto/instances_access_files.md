(instances-access-files)=
# インスタンス内のファイルにアクセスするには

Incus クライアントを使って、ネットワーク経由でインスタンスにアクセスする必要なしに、インスタンス内部のファイルを管理できます。
ファイルを個別に編集、削除したり、ローカルマシンからプッシュ、ローカルマシンにプルできます。
あるいは、インスタンスのファイルシステムをローカルマシン上にマウントできます。

コンテナでは、これらの操作は必ず機能し Incus で直接処理されます。
仮想マシンでは、これらの操作が機能するためには `incus-agent` プロセスが仮想マシン内部で稼働している必要があります。

## インスタンスのファイルを編集する

ローカルマシンからインスタンスのファイルを編集するには、以下のコマンドを入力します:

    incus file edit <instance_name>/<path_to_file>

たとえば、インスタンス内の `/etc/hosts` ファイルを編集するには、以下のコマンドを入力します:

    incus file edit my-container/etc/hosts

```{note}
ファイルはインスタンス上に既に存在している必要があります。
インスタンス上にファイルを作成するのに `edit` コマンドは使えません。
```

## インスタンスからファイルを削除する

インスタンスからファイルを削除するには、以下のコマンドを入力します:

    incus file delete <instance_name>/<path_to_file>

## インスタンスからローカルマシンにファイルをプルする

インスタンスからローカルマシンにファイルをプルするには、以下のコマンドを入力します:

    incus file pull <instance_name>/<path_to_file> <local_file_path>

たとえば `/etc/hosts` ファイルをカレントディレクトリーにプルするには、以下のコマンドを入力します:

    incus file pull my-instance/etc/hosts .

インスタンスのファイルをローカルマシンにプルする代わりに、標準出力にプルして別のコマンドの標準入力にパイプすることもできます。
これは、たとえば、ログファイルをチェックするのに便利です:

    incus file pull my-instance/var/log/syslog - | less

ディレクトリーのすべての中身をプルするには、以下のコマンドを入力します:

    incus file pull -r <instance_name>/<path_to_directory> <local_location>

## ローカルマシンからインスタンスにファイルをプッシュする

ローカルマシンからインスタンスにファイルをプッシュするには、以下のコマンドを入力します:

    incus file push <local_file_path> <instance_name>/<path_to_file>

ディレクトリーのすべての中身をプッシュするには、以下のコマンドを入力します:

    incus file push -r <local_location> <instance_name>/<path_to_directory>

## インスタンスのファイルシステムをマウントする

インスタンスのファイルシステムをクライアントのローカルパスにマウントできます。

そうするためには、`sshfs` がインストールされていることを確認してください。

    incus file mount <instance_name>/<path_to_directory> <local_location>

これでローカルマシンからファイルにアクセスできます。

### SSH SFTP リスナーをセットアップする

別の方法として、SSH SFTP リスナーをセットアップすることもできます。
この方法では任意の SFTP クライアントで専用のユーザー名で接続できます。

そうするには、まず以下のコマンドを入力してリスナーをセットアップします:

    incus file mount <instance_name> [--listen <address>:<port>]

たとえば、ローカルマシン上のランダムなポート（たとえば、`127.0.0.1:45467`）にリスナーをセットアップするには以下のようにします:

    incus file mount my-instance

ローカルネットワークの外側からインスタンスのファイルにアクセスしたい場合、特定のアドレスとポートを渡せます:

    incus file mount my-instance --listen 192.0.2.50:2222

```{caution}
あなたのインスタンスをリモートに公開することになるので、これを実行する際には注意してください。
```

特定のアドレスとランダムなポートでリスナーをセットアップするには以下のようにします:

    incus file mount my-instance --listen 192.0.2.50:0

コマンドは割り当てられたポートと接続に使用するユーザー名とパスワードを出力します。

```{tip}
`--auth-user` フラグを渡すとユーザ名を指定できます。
```

この情報を使ってファイルシステムにアクセスします。
たとえば `sshfs` で接続するには、以下のコマンドを入力します:

    sshfs <user_name>@<address>:<path_to_directory> <local_location> -p <port>

たとえば以下のようにします:

    sshfs xFn8ai8c@127.0.0.1:/home my-instance-files -p 35147

これでインスタンスのファイルシステムをローカルマシン上の指定の場所でアクセスできます。
