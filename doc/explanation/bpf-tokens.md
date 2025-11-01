(bpf-tokens)=
# BPFトークン移譲

IncusはLinuxカーネル6.9で追加された[BPFトークン](https://docs.ebpf.io/linux/concepts/token/)を通してBPFケーパビリティーの移譲をサポートします。

{config:option}`instance-security:security.bpffs.delegate_cmds`、
{config:option}`instance-security:security.bpffs.delegate_maps`、
{config:option}`instance-security:security.bpffs.delegate_progs`、
{config:option}`instance-security:security.bpffs.delegate_attachs`のいずれかの設定オプションが設定されている場合、IncusはBPFファイルシステムを{config:option}`instance-security:security.bpffs.path`設定オプションで指定されたパスでコンテナにマウントし指定されたケーパビリティーをコンテナに移譲します。

これらのオプションに設定可能な値はカーネルのバージョンに依存し、BPFのヘッダーファイル（カーネルツリー内では`include/uapi/linux/bpf.h`、ほとんどのディストリビューションではカーネルソースをインストールしていれば`/usr/include/linux/bpf.h`）内の`enums`で定義されています:

 キー                              | カーネルの`enum`   | 除去する接頭辞
 :--                               |:--                 | :--
 `security.bpffs.delegate_cmds`    | `bpf_cmd`          | `BPF_`
 `security.bpffs.delegate_maps`    | `bpf_map_type`     | `BPF_MAP_TYPE_`
 `security.bpffs.delegate_progs`   | `bpf_prog_type`    | `BPF_PROG_TYPE_`
 `security.bpffs.delegate_attachs` | `bpf_attach_type`  | `BPF_`

これらの設定オプションはカンマ区切りリストの値をとり、さらに`any`という値を指定するとそのタイプのすべての可能な値を移譲します。

## 例

 キー                              | 値
 :--                               | :--
 `security.bpffs.delegate_cmds`    | `map_create,obj_get,link_create`
 `security.bpffs.delegate_maps`    | `hash,array,devmap,queue,stack`
 `security.bpffs.delegate_progs`   | `socket_filter,kprobe,cgroup_sysctl`
 `security.bpffs.delegate_attachs` | `any`

```bash
$ mount -t bpf
none on /sys/fs/bpf type bpf (rw,relatime,delegate_cmds=map_create:obj_get:link_create,delegate_maps=hash:array:devmap:queue:stack,delegate_progs=socket_filter:kprobe:cgroup_sysctl,delegate_attachs=any)
```
