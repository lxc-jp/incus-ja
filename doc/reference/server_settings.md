(server-settings)=
# Incusの本番環境用のサーバー設定

Incusサーバーで多数のインスタンスを動かすためには、サーバーのリミットにひっかからないように以下の設定をしてください。

「値」のカラムは各パラメーターの推奨値です。

## `/etc/security/limits.conf`

| ドメイン | 種別 | 項目      | 値          | デフォルト値 | 説明                                                                                |
| :---     | :--- | :---      | :---        | :---         | :---                                                                                |
| `*`      | soft | `nofile`  | `1048576`   | 設定なし     | オープンするファイルの最大数                                                        |
| `*`      | hard | `nofile`  | `1048576`   | 設定なし     | オープンするファイルの最大数                                                        |
| `root`   | soft | `nofile`  | `1048576`   | 設定なし     | オープンするファイルの最大数                                                        |
| `root`   | hard | `nofile`  | `1048576`   | 設定なし     | オープンするファイルの最大数                                                        |
| `*`      | soft | `memlock` | `unlimited` | 設定なし     | ロックされたメモリ空間の最大値（KB）                                                |
| `*`      | hard | `memlock` | `unlimited` | 設定なし     | ロックされたメモリ空間の最大値（KB）                                                |
| `root`   | soft | `memlock` | `unlimited` | 設定なし     | ロックされたメモリ空間の最大値（KB）、`bpf`システムコールスーパービジョンでのみ必要 |
| `root`   | hard | `memlock` | `unlimited` | 設定なし     | ロックされたメモリ空間の最大値（KB）、`bpf`システムコールスーパービジョンでのみ必要 |

## `/etc/sysctl.conf`

```{note}
これらのパラメーターを変更した後はサーバーを再起動してください。
```

| パラメーター                        | 値           | デフォルト値 | 説明                                                                                                                                                                                                                                                                                                                                                         |
| :---                                | :---         | :---         | :---                                                                                                                                                                                                                                                                                                                                                         |
| `fs.aio-max-nr`                     | `524288`     | `65536`      | 同時実行可能な非同期I/Oの最大数（例えば、MySQLのようにAIOサブシステムを使う大量のワークロードがある場合はこの値を増やす必要があるかもしれません）                                                                                                                                                                                                            |
| `fs.inotify.max_queued_events`      | `1048576`    | `16384`      | 対応する`inotify`インスタンスにキューイングできるイベント数の上限（[`inotify`](https://man7.org/linux/man-pages/man7/inotify.7.html)を参照）                                                                                                                                                                                                                 |
| `fs.inotify.max_user_instances`     | `1048576`    | `128`        | 実ユーザーごとに作成できる`inotify`インスタンスの数の上限（[`inotify`](https://man7.org/linux/man-pages/man7/inotify.7.html)を参照）                                                                                                                                                                                                                         |
| `fs.inotify.max_user_watches`       | `1048576`    | `8192`       | 実ユーザーごとに作成できるwatchの数の上限（[`inotify`](https://man7.org/linux/man-pages/man7/inotify.7.html)を参照）                                                                                                                                                                                                                                         |
| `kernel.dmesg_restrict`             | `1`          | `0`          | カーネルのリングバッファー内のメッセージにコンテナーからのアクセスを許可するかどうか（これはホスト上の非rootユーザーからのアクセスも拒否することに注意）                                                                                                                                                                                                     |
| `kernel.keys.maxbytes`              | `2000000`    | `20000`      | 非ルートユーザーが使えるキーリングの最大サイズ                                                                                                                                                                                                                                                                                                               |
| `kernel.keys.maxkeys`               | `2000`       | `200`        | 非ルートユーザーが使えるキーの最大数（値はインスタンス数より多いべきです）                                                                                                                                                                                                                                                                                 |
| `net.core.bpf_jit_limit`            | `1000000000` | 環境依存     | eBPFのJIT割り当てのサイズの上限（`CONFIG_BPF_JIT_ALWAYS_ON=y`でコンパイルされた5.15より古いカーネルでは、この値は作成できるインスタンス数も制限するかもしれません）                                                                                                                                                                                          |
| `net.ipv4.neigh.default.gc_thresh3` | `8192`       | `1024`       | IPv4 ARPテーブルのエントリの最大数（1024より多くインスタンスを作成するつもりならこの値を増やしてください。でなければARPテーブルがフルになった際に`neighbour: ndisc_cache: neighbor table overflow!`エラーが出てインスタンスがネットワーク設定を取得できなくなります。[`ip-sysctl`](https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt)参照） |
| `net.ipv6.neigh.default.gc_thresh3` | `8192`       | `1024`       | IPv6 ARPテーブルのエントリの最大数（1024より多くインスタンスを作成するつもりならこの値を増やしてください。でなければARPテーブルがフルになった際に`neighbour: ndisc_cache: neighbor table overflow!`エラーが出てインスタンスがネットワーク設定を取得できなくなります。[`ip-sysctl`](https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt)参照） |
| `vm.max_map_count`                  | `262144`     | `65530`      | プロセスが持てるメモリマップエリアの最大数（メモリマップエリアは`malloc`呼び出しの副作用として、`mmap`や`mprotect`では直接、そして共有ライブラリをロードした際にも使われます）                                                                                                                                                                                     |
