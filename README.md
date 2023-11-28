# Quan 短链接服务

本来是为免费容器服务 koyeb / zeabur / railway / replit 等开发的，不需要外部数据库，还自带备份短链接服务。但是做好之后发现：

- koyeb 免费额度 256MB 内存限制，启动不了；
- zeabur 的免费套餐只支持函数计算，完全不保证什么时候就干掉容器化；
- railway 只有一次性的 $5 试用额度，用完就要续费；
- 而 replit 一开始就是要收费的；
- ……

考虑再三，还不如使用阿里云的 FC（函数计算），一个月块把钱我还是能给的。

所以本项目暂时封存，此文档仅做留存。

## 简介

如上，golang 开发的短链接服务。

LevelDB，不依赖外部数据库，支持利用 Google Cloud Storage（每个月 5GB 免费额度，足够用）分批备份数据。

## 使用

`Dockerfile` 是项目容器化的脚本，`src/main.go` 是项目的主入口。

容器化时需要注意，项目支持下列环境变量。

- `QUAN_ADMIN_USERNAME`: 管理员用户名。管理员可以操作所有 */admin/* 的资源，例如添加、删除短连接
- `QUAN_ADMIN_PASSWORD`: 管理员密码。管理员是唯一的，http basic auth，只能通过这种方式更改
- `QUAN_PORT`: 服务端口。默认：8080
- `QUAN_BASE_URL`: 服务的基础 URL，不带结尾斜杠。默认：*http://localhost:8080*
- `QUAN_LENGTH`: 短链接 hash 长度。默认：*6*
- `QUAN_DB_FILE`: LevelDB 的位置。默认：*/data/quan/db*
- `QUAN_CHAR_RANGE`: 随机字符范围，支持 36（小写字母+数字）或 62（大小写字母+数字）。默认：*62*
- `QUAN_DEFAULT_REDIRECT_URL`: 用户访问不存在的短链时跳转到哪里。
- `QUAN_LIST_SIZE`: 管理员查看列表时，一页显示多少个。可用 url 参数覆盖。默认：*100*

以下是备份相关参数：

- `QUAN_BACKUP_BUCKET`: GCS 的 Bucket 名称。默认：*quan-backup*
- `QUAN_BACKUP_FILENAME`: 备份文件的名称。不带文件后缀名。
- `QUAN_BACKUP_BATCH_SIZE`: 分批备份，每一次备份多少条数据。默认：*100000*
- `QUAN_BACKUP_INTERVAL`: 后台多久备份一次，单位：秒。默认：*1200*
- `QUAN_BACKUP_CREDENTIAL_FILE`: Google 认证文件的完整路径。
- `QUAN_BACKUP_CREDENTIAL_CONTENT`: Google 认证文件的完整内容。和上面的选项二选一。

要使用 [Google Cloud Storage](https://console.cloud.google.com/storage)，需要在 Google Console 中启用 Google Cloud Storage 服务（以下简称 GCS）。

接着在 GCS 中创建一个 Bucket。注意，免费额度似乎仅在 `west-1` 和 `east-1` 这两个分区有效。

Bucket 创建好之后，来到 [IAM](https://console.cloud.google.com/iam-admin/iam) 中[创建一个服务账号](https://console.cloud.google.com/iam-admin/serviceaccounts/create)。创建过程中，我们要为这个服务账号赋权 GCS 的读写权限，生成并下载这个账号的密钥，最终，得到一个叫做 google credential 的认证文件（json 格式）。

这个 json 文件中的内容，就是 `QUAN_BACKUP_CREDENTIAL_FILE`。文件的内容，对应 `QUAN_BACKUP_CREDENTIAL_CONTENT`。

