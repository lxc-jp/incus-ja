(instances-console)=
# コンソールにアクセスするには

インスタンスのコンソールにアタッチするには [`incus console`](incus_console.md) コマンドを使います。
コンソールは起動時に既に利用可能になり、必要なら、ブートメッセージを見て、コンテナや仮想マシンの起動時の問題をデバッグするのに使えます。

インタラクティブなコンソールに接続するには、以下のコマンドを入力します:

    incus console <instance_name>

ログ出力を見るには `--show-log` フラグを渡します:

    incus console <instance_name> --show-log

インスタンスが起動したらすぐにコンソールにアタッチできます:

    incus start <instance_name> --console
    incus start <instance_name> --console=vga

## グラフィカルなコンソールにアタッチする（仮想マシンの場合）

仮想マシンでは、コンソールにログオンしてグラフィカルな出力を見ることができます。
コンソールを使えば、たとえば、グラフィカルなインターフェースを使ってオペレーティングシステムをインストールしたりデスクトップ環境を実行できます。

さらなる利点は `incus-agent` プロセスが実行していなくても、コンソールは利用可能です。
これは `incus-agent` が起動する前や `incus-agent` が全く利用可能でない場合にもコンソール経由で仮想マシンにアクセスできることを意味します。

仮想マシンにグラフィカルなアウトプットを持つ VGA コンソールを開始するには、 SPICE クライアント (たとえば、`virt-viewer` または `spice-gtk-client`) をインストールする必要があります。
次に以下のコマンドを入力します:

    incus console <vm_name> --type vga
