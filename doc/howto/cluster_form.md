(cluster-form)=
# クラスタを形成するには

Incus クラスタを形成するときはブートストラップサーバーから始めます。
このブートストラップサーバーは既存の Incus サーバーでもよいですし新しくインストールしたものでもよいです。

ブートストラップサーバーを初期化した後、クラスタに追加のサーバーをジョインできます。
詳細は {ref}`clustering-members` を参照してください。

Incus クラスタを形成するために初期化プロセス中に設定をインタラクティブに指定することもできますし、完全な設定を含むプリシードファイルを使うこともできます。

## クラスタをインタラクティブに設定する

クラスタを形成するには、まずブートストラップサーバー上で `incus admin init` を実行する必要があります。その後クラスタにジョインさせたい他のサーバー上でもそのコマンドを実行します。

クラスタをインタラクティブに形成する際、クラスタを設定するために `incus admin init` のプロンプトの質問に回答します。

### ブートストラップサーバーを初期化する

ブートストラップサーバーを初期化するには、 `incus admin init` を実行して希望の設定に応じて質問に回答します。

ほとんどの質問はデフォルト値を受け入れることができますが、以下の質問には適切に答えるようにしてください:

- `Would you like to use Incus clustering?`

  **yes** を選択。
- `What IP address or DNS name should be used to reach this server?`

  他のサーバーがアクセスできる IP または DNS のアドレスを確実に使用してください。
- `Are you joining an existing cluster?`

  **no** を選択。

<details>
<summary>ブートストラップ上での <code>incus admin init</code> の完全な例を見るには展開してください</summary>

```{terminal}
:input: incus admin init

Would you like to use Incus clustering? (yes/no) [default=no]: yes
What IP address or DNS name should be used to reach this server? [default=192.0.2.101]:
Are you joining an existing cluster? (yes/no) [default=no]: no
What member name should be used to identify this server in the cluster? [default=server1]:
Do you want to configure a new local storage pool? (yes/no) [default=yes]:
Name of the storage backend to use (btrfs, dir, lvm, zfs) [default=zfs]:
Create a new ZFS pool? (yes/no) [default=yes]:
Would you like to use an existing empty block device (e.g. a disk or partition)? (yes/no) [default=no]:
Size in GiB of the new loop device (1GiB minimum) [default=9GiB]:
Do you want to configure a new remote storage pool? (yes/no) [default=no]:
Would you like to configure Incus to use an existing bridge or host interface? (yes/no) [default=no]:
Would you like stale cached images to be updated automatically? (yes/no) [default=yes]:
Would you like a YAML "incus admin init" preseed to be printed? (yes/no) [default=no]:
```

</details>

初期化プロセスが終了したら、最初のクラスタメンバーが起動してネットワーク上で利用可能になるはずです。
これは [`incus cluster list`](incus_cluster_list.md) で確認できます。

### 追加のサーバーをジョインさせる

これでクラスタに追加のサーバーをジョインできるようになりました。

```{note}
追加するサーバーは新規にインストールした Incus サーバーにするほうがよいです。
既存のサーバーを使う場合、既存のデータは消失するので、ジョインする前にデータを確実にクリアしてください。
```

クラスタにサーバーをジョインさせるには、クラスタ上で `incus admin init` を実行します。
既存のクラスタにジョインするには root 権限が必要ですので、コマンドを root で実行するか `sudo` をつけて実行するのを忘れないでください。

基本的に、初期化プロセスは以下のステップからなります:

1. 既存のクラスタにジョインをリクエストする。

   `incus admin init` の最初の質問に適切に回答します:

   - `Would you like to use Incus clustering?`

     **yes** を選択。
   - `What IP address or DNS name should be used to reach this server?`

     他のサーバーがアクセスできる IP または DNS のアドレスを確実に使用してください。
   - `Are you joining an existing cluster?`

     **yes** を選択。

1. クラスタで認証する。

   ブートスラップサーバーを設定する際に選んだ認証方法に応じて 2 つの方法があります。

   `````{tabs}

   ````{group-tab} 認証トークン
   {ref}`認証トークン <authentication-token>` を使うようにクラスタを設定した場合、新メンバーごとにジョイントークンを生成する必要があります。
   そのためには、既存のクラスタメンバー（たとえば、ブートストラップサーバー）で以下のコマンドを実行します:

       incus cluster add <new_member_name>

   このコマンドは設定（{config:option}`server-cluster:cluster.join_token_expiry`参照）時に有効な一回限りのジョイントークンを返します。
   `incus admin init` のプロンプトでジョイントークンを求められたときにこのトークンを入力してください。

   ジョイントークンは既存のオンラインメンバーのアドレス、一回限りのシークレットとクラスタ証明書のフィンガープリントを含みます。
   ジョイントークンがこれらの質問に自動で回答できるので、 `incus admin init` 中に回答が必要な質問の量を減らすことができます。
   ````

   `````

1. クラスタにジョインする際サーバーのすべてのローカルデータが消失することを確認します。
1. サーバー固有の設定を行います（詳細は {ref}`clustering-member-config` を参照）。

   デフォルト値を受け入れることもできますし、各サーバーにカスタム値を指定することもできます。

<details>
<summary>追加のサーバー上で <code>incus admin init</code> を実行する完全な例を見るには展開してください</summary>

`````{tabs}

````{group-tab} 認証トークン

```{terminal}
:input: sudo incus admin init

Would you like to use Incus clustering? (yes/no) [default=no]: yes
What IP address or DNS name should be used to reach this server? [default=192.0.2.102]:
Are you joining an existing cluster? (yes/no) [default=no]: yes
Do you have a join token? (yes/no/[token]) [default=no]: yes
Please provide join token: eyJzZXJ2ZXJfbmFtZSI6InJwaTAxIiwiZmluZ2VycHJpbnQiOiIyNjZjZmExZDk0ZDZiMjk2Nzk0YjU0YzJlYzdjOTMwNDA5ZjIzNjdmNmM1YjRhZWVjOGM0YjAxYTc2NjU0MjgxIiwiYWRkcmVzc2VzIjpbIjE3Mi4xNy4zMC4xODM6ODQ0MyJdLCJzZWNyZXQiOiJmZGI1OTgyNjgxNTQ2ZGQyNGE2ZGE0Mzg5MTUyOGM1ZGUxNWNmYmQ5M2M3OTU3ODNkNGI5OGU4MTQ4MWMzNmUwIn0=
All existing data is lost when joining a cluster, continue? (yes/no) [default=no] yes
Choose "size" property for storage pool "local":
Choose "source" property for storage pool "local":
Choose "zfs.pool_name" property for storage pool "local":
Would you like a YAML "incus admin init" preseed to be printed? (yes/no) [default=no]:
```

````
`````

</details>

初期化プロセスが終わった後、サーバーが新しいクラスタメンバーとして追加されます。
これは [`incus cluster list`](incus_cluster_list.md) で確認できます。

## クラスタをプリシードファイルで設定する

クラスタを形成するには、まずブートストラップサーバー上で `incus admin init` を実行します。
その後、クラスタにジョインさせたい他のサーバーでもこのコマンドを実行します。

`incus admin init` の質問にインタラクティブに回答する代わりに、プリシードファイルを使って必要な情報を提供できます。
以下のコマンドを使って `incus admin init` にファイルをフィードできます:

    cat <preseed-file> | incus admin init --preseed

サーバーごとに異なるプリシードファイルが必要です。

### ブートストラップサーバーを初期化する

`````{tabs}

````{group-tab} 認証トークン
クラスタリングを有効にするには、ブートストラップサーバー用のプリシードファイルは以下のフィールドを含む必要があります:

```yaml
config:
  core.https_address: <IP_address_and_port>
cluster:
  server_name: <server_name>
  enabled: true
```

ブートストラップサーバー用のプリシードファイルの例を以下に示します:

```yaml
config:
  core.https_address: 192.0.2.101:8443
  images.auto_update_interval: 15
storage_pools:
- name: default
  driver: dir
- name: my-pool
  driver: zfs
networks:
- name: incusbr0
  type: bridge
profiles:
- name: default
  devices:
    root:
      path: /
      pool: my-pool
      type: disk
    eth0:
      name: eth0
      nictype: bridged
      parent: incusbr0
      type: nic
cluster:
  server_name: server1
  enabled: true
```

````

`````

### 追加のサーバーをジョインさせる

新しいクラスタメンバー用のプリシードファイルは参加するサーバーに固有のデータと設定値を含む `cluster` セクションのみが必要です。

`````{tabs}

````{group-tab} 認証トークン
追加のサーバーのプリシードファイルは以下の項目を含む必要があります:

```yaml
cluster:
  enabled: true
  server_address: <IP_address_of_server>
  cluster_token: <join_token>
```

新しいクラスタメンバー用のプリシードファイルの例を以下に示します:

```yaml
cluster:
  enabled: true
  server_address: 192.0.2.102:8443
  cluster_token: eyJzZXJ2ZXJfbmFtZSI6Im5vZGUyIiwiZmluZ2VycHJpbnQiOiJjZjlmNmVhMWIzYjhiNjgxNzQ1YTY1NTY2YjM3ZGUwOTUzNjRmM2MxMDAwMGNjZWQyOTk5NDU5YzY2MGIxNWQ4IiwiYWRkcmVzc2VzIjpbIjE3Mi4xNy4zMC4xODM6ODQ0MyJdLCJzZWNyZXQiOiIxNGJmY2EzMDhkOTNhY2E3MGJmYThkMzE0NWM4NWY3YmE0ZmU1YmYyNmJiNDhmMmUwNzhhOGZhMDczZDc0YTFiIn0=
  member_config:
  - entity: storage-pool
    name: default
    key: source
    value: ""
  - entity: storage-pool
    name: my-pool
    key: source
    value: ""
  - entity: storage-pool
    name: my-pool
    key: driver
    value: "zfs"

```

````

`````
