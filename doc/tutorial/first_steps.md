(first-steps)=
# Incusを使う最初のステップ

このチュートリアルは Incus を使う最初のステップのガイドです。
Incus のインストールと初期化、インスタンスの作成と設定、インスタンスの操作、スナップショットの作成を取り扱います。

これらのステップを経験した後、どのように Incus を使うかのアイデアがあれば、より高度な使用例について探索を開始できることでしょう！

## Incusのインストールと初期化

1. Incus パッケージのインストール

Incus はほとんどのよくある Linux ディストリビューションで利用できます。

ディストリビューションごとの詳細な手順は{ref}`installing`を参照してください。

1. あなたのユーザーに Incus を制御する許可を与えます。

   上記のパッケージに含まれる Incus へのアクセスは 2 つのグループで制御されます。

   - `incus`は基本的なユーザーアクセスを許可します。設定はできずすべてのアクションはユーザーごとのプロジェクトに限定されます。
   - `incus-admin`は Incus の完全なコントロールを許可します。

   すべてのコマンドを root で実行することなく Incus を制御するには、あなた自身を`incus-admin`グループに追加してください。

       sudo adduser YOUR-USERNAME incus-admin
       newgrp incus-admin

   `newgrp`の手順はあなたの端末セッションを再起動しないままで Incus を利用する場合に必要です（訳注：端末を起動し直す場合は不要です）。

1. Incus の初期化

   ```{note}
   既存の LXD 環境からマイグレートする場合、このステップはスキップして代わりに{ref}`server-migrate-lxd`を参照してください。
   ```

   Incus はネットワークとストレージでいくつかの初期設定が必要です。これは以下のコマンドでインタラクティブにできます:

       incus admin init

   あるいは以下のコマンドで基本的な設定を自動で設定できます:

       incus admin init --minimal

   初期化オプションをチューニングしたい場合、詳細は{ref}`initialize`を参照してください。


## インスタンスの起動と調査

Incus はイメージベースです。そしてさまざまなイメージサーバーからイメージをロードできます。
このチュートリアルでは、[公式イメージサーバ](https://images.linuxcontainers.org/)イメージサーバーを使います。

このイメージサーバーで利用可能なすべてのイメージを一覧表示するには以下のようにします。

    incus image list images:

Incus が使用するイメージについてのより詳細情報は{ref}`images`を参照してください。

では、いくつかインスタンスを起動してみましょう。
コンテナまたは仮想マシンのことを*インスタンス*と呼びます。
2 つのインスタンスタイプの違いについての情報は{ref}`containers-and-vms`を参照してください。

インスタンスを管理するには、Incus のコマンドラインクライアント`incus`を使います。

1. Ubuntu 22.04 イメージを使って`first`という名前のコンテナを起動します。

       incus launch images:ubuntu/22.04 first

   ```{note}
   最初はイメージをダウンロードして展開しなければならないため、コンテナの起動には少し時間がかかることに注意してください。
   ```

1. 同じイメージを使って`second`という名前のコンテナを起動します。

       incus launch images:ubuntu/22.04 second

   ```{note}
   イメージを取得済みなので、最初の（first）コンテナの起動に比べると早く起動します。
   ```

1. 最初の（first）コンテナを`third`という名前のコンテナとしてコピーします。

       incus copy first third

1. Ubuntu 22.04 イメージを使って`ubuntu-vm`という名前の仮想マシンを起動します。

       incus launch images:ubuntu/22.04 ubuntu-vm --vm

   ```{note}
   インスタンスの起動に同じイメージ名を使っていますが、Incusは仮想マシンに適した少し異なるイメージをダウンロードします。
   ```

1. 起動したインスタンスのリストをチェックします。

       incus list

   3 番目のコンテナ以外は稼働中になっています。
   3 つ目のコンテナ以外が起動していることが確認できるでしょう。これは、3 つ目のコンテナを最初の（first）コンテナからコピーして作成はしたものの、起動処理を実行していないからです。

   3 つ目のインスタンスを次のように起動できます。

       incus start third

1. それぞれのコンテナの情報をもう少し詳しく見ることができます。

       incus info first
       incus info second
       incus info third
       incus info ubuntu-vm

1. チュートリアルではこの後、これらのインスタンスすべては必要ありませんので、不要なインスタンスを消しましょう。

   1. 2 つ目のコンテナを停止します。

          incus stop second

   1. 2 つ目のコンテナを削除します。

          incus delete second

   1. 3 つ目のコンテナを削除します。

          incus delete third

      このコンテナはまだ実行中なので、最初に停止しないとエラーメッセージが出るでしょう。その代わりに、強制的に削除できます。

          incus delete third --force

詳細は{ref}`instances-create`と{ref}`instances-manage`を参照してください。

## インスタンスの設定

インスタンスに設定できる制限や設定オプションがいくつか存在します。その概要については{ref}`instance-options`を参照してください。

リソース制限を持つインスタンスを 1 つ作ってみましょう。

1. コンテナを起動し、1vCPU と 192MiB メモリーの制限を設定します。

       incus launch images:ubuntu/22.04 limited --config limits.cpu=1 --config limits.memory=192MiB

1. 現在の設定を確認し、制限が設定されていない最初の（first）コンテナの設定と比べてみましょう。

       incus config show limited
       incus config show first

1. 親環境のシステムと 2 つのコンテナで空きメモリーと使用済メモリーの量をチェックしましょう。

       free -m
       incus exec first -- free -m
       incus exec limited -- free -m

   ```{note}
   デフォルトでは、コンテナは親環境からリソースを継承するため、親環境と最初の（first）インスタンスではメモリの総量が同じであることに注意してください。一方で、制限を設定したインスタンスは192MiBだけが使用できます。
   ```

1. 親環境と 2 つのインスタンスで使用できる CPU の数をチェックしましょう。

       nproc
       incus exec first -- nproc
       incus exec limited -- nproc

   ```{note}
   ふたたび、親環境と最初の（first）インスタンスのCPU数は同じで、制限を設定したインスタンスでは減少していることに注意してください。
   ```

1. 実行中のインスタンスの設定を更新することもできます。

   1. インスタンスのメモリー制限を設定する。

          incus config set limited limits.memory=128MiB

   1. 適用した設定をチェックする。

          incus config show limited

   1. コンテナで使用できるメモリー量をチェックする。

          incus exec limited -- free -m

      数値が変わっていることを確認してください。

1. 使用するインスタンスタイプとストレージドライバーによっては、より多くの設定を指定できます。
   たとえば、仮想マシンの root ディスクデバイスのサイズを指定できます。

   1. Ubuntu 仮想マシンの root ディスクデバイスの現在のサイズをチェックします。

      ```{terminal}
      :input: incus exec ubuntu-vm -- df -h

      Filesystem      Size  Used Avail Use% Mounted on
      /dev/root       9.6G  1.4G  8.2G  15% /
      tmpfs           483M     0  483M   0% /dev/shm
      tmpfs           193M  604K  193M   1% /run
      tmpfs           5.0M     0  5.0M   0% /run/lock
      tmpfs            50M   14M   37M  27% /run/incus_agent
      /dev/sda15      105M  6.1M   99M   6% /boot/efi
      ```

   1. root ディスクデバイスのサイズを上書きします。

          incus config device override ubuntu-vm root size=30GiB

   1. 仮想マシンを再起動します。

          incus restart ubuntu-vm

   1. ふたたび、root ディスクデバイスのサイズをチェックします。

       ```{terminal}
       :input: incus exec ubuntu-vm -- df -h

       Filesystem      Size  Used Avail Use% Mounted on
       /dev/root        29G  1.4G   28G   5% /
       tmpfs           483M     0  483M   0% /dev/shm
       tmpfs           193M  588K  193M   1% /run
       tmpfs           5.0M     0  5.0M   0% /run/lock
       tmpfs            50M   14M   37M  27% /run/incus_agent
       /dev/sda15      105M  6.1M   99M   6% /boot/efi
       ```

より詳細な情報は`instances-configure`と{ref}`instance-config`を参照してください。

## インスタンスの操作

インスタンス内で（インタラクティブなシェルを含む）コマンドを実行したりインスタンス内のファイルにアクセスできます。

まずインスタンス内でインタラクティブなシェルを起動しましょう。

1. コンテナ内で`bash`コマンドを実行します。

       incus exec first -- bash

1. たとえば以下のコマンドを入力するとオペレーティングシステムについての情報が表示されます。

       cat /etc/*release

1. インタラクティブなシェルを抜けます。

       exit

インスタンスにログインしてコマンドを実行する代わりに、ホストから直接コマンドを実行できます。

たとえば、コンテナ上にコマンドラインツールをインストールし、それを実行できます。

    incus exec first -- apt-get update
    incus exec first -- apt-get install sl -y
    incus exec first -- /usr/games/sl

詳細は{ref}`run-commands`を参照してください。

インスタンスのファイルにアクセスしたり、ファイルを操作できます。

1. インスタンスからファイルを取得します。

       incus file pull first/etc/hosts .

1. ファイルにエントリーを追加します。

       echo "1.2.3.4 my-example" >> hosts

1. インスタンスにファイルを戻します。

       incus file push hosts first/etc/hosts

1. ログファイルにアクセスするために同じメカニズムを使います。

       incus file pull first/var/log/syslog - | less

   ```{note}
   Press `q` to exit the `less` command.
   ```

詳細は{ref}`instances-access-files`を参照してください。

## スナップショットの管理

インスタンスのスナップショットを作成したり、スナップショットからリストアしたりできます。

1. "clean"という名前のスナップショットを作ります。

       incus snapshot create first clean

1. スナップショットが作られたことを確認します。

       incus list first
       incus info first

   ```{note}
   `incus list`はスナップショットの数を表示します。
   `incus info`はそれぞれのスナップショットについての情報を表示します。
   ```

1. コンテナを破壊します。

       incus exec first -- rm /usr/bin/bash

1. 壊れたことを確認します。

       incus exec first -- bash

   ```{note}
   `bash`コマンドを削除したので、シェルが実行できないことに注意してください。
   ```

1. スナップショットの状態にインスタンスをリストアします。

       incus snapshot restore first clean

1. すべて通常状態に戻ったことを確認します。

       incus exec first -- bash
       exit

1. スナップショットを削除します。

       incus delete first/clean

詳細は{ref}`instances-snapshots`を参照してください。
