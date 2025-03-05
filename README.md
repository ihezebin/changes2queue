# changes2queue

changes 可理解为 `变更` 或 `change stream` 变更流，changes2queue 主要用于监听数据库(如：MongoDB、MySQL)的数据变更，然后通过消息队列分发到其他服务实现解耦。

## 应用场景

- 数据库数据同步 ES
- 多级缓存或分布式多 POD 实例缓存更新
- 其他需要监听数据库变更的场景

## 支持数据库

目前实现和支持如下数据库：

- MongoDB：官方原生支持 ChangeStream 功能
- MySQL：通过 binlog 实现，使用 Canal（github.com/withlin/canal-go）实现变更订阅

## 支持的消息队列

目前实现和支持如下消息队列：

- Pulsar：一个多租户、高性能、多订阅模式的消息中间件
- RedisBroadcast：一个基于 Redis 的消息广播
- Redis：基于 Redis list 实现的消息队列

## 热更新

考虑可热更新任务是一个重要和实用的功能，那么具体的 change to queue 任务列表自然需要持久化存储，那么任务存储在哪里呢？

为了避免出现臃肿、杂乱无章和晦涩难懂的配置，changes2queue 选择将任务存储在其自身的 change 源中，即：

- 在 MongoDB 中 blog 库中，changes2queue 会在 blog 库下创建一个 mongo2queue 表，用于存储任务列表
- 在 MySQL 中 blog 库中，changes2queue 会在 blog 库下创建一个 mysql2queue 表，用于存储任务列表

## 配置

如在需要将 blog 库下的 article、tag、draft 表数据通过 Pulsar 消息中间件同步到 ES 的场景中。注意：这个 blog 库可以是一个 MongoDB 数据库，也可以是一个 MySQL 数据库等上述支持的数据库类型。则配置项如下：

### MongoDB 配置

```toml
mongo2queues = ["mongodb://root:root@127.0.0.1:27017/blog?authSource=admin&replicaSet=hezebin"]
```

该配置描述了一个 MongoDB 类型的任务列表，这里只配置了一个任务，这个任务的源是 blog 库，任务中的这个库具体有哪些表的 change 需要推送则通过查询 blog 库下的 mongo2queue 表来获取。

mongo2queue 表的结构如下：

```json
[
  {
    "id": "1",
    "collections": ["article", "tag"],
    "queue": {
      "type": "pulsar",
      "options": {
        "url": "pulsar://localhost:6650",
        "topic": "persistent://public/default/blog"
      }
    }
  },
  {
    "id": "2",
    "collections": ["draft"],
    "queue": {
      "type": "redis",
      "options": {
        "addrs": ["127.0.0.1:6379"],
        "password": "root",
        "topic": "blog"
      }
    }
  }
]
```

### MySQL 配置

```toml
mysql2queues = ["mysql://root:root@127.0.0.1:3306/blog?charset=utf8mb4"]
```

同理，mysql2queue 的表结构如下:

```sql
CREATE DATABASE IF NOT EXISTS `blog`;
USE `blog`;
CREATE TABLE IF NOT EXISTS `blog`.`mysql2queue` (
  `id` varchar(255) NOT NULL,
  `table` varchar(255) NOT NULL,
  `queue` TEXT NOT NULL,
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

其中 `queue` 字段为将上述 MongoDB 配置中包含 type 和 options 的对象 JSON 序列化后的字符串：

| id  | table   | queue                                                                                                          |
| --- | ------- | -------------------------------------------------------------------------------------------------------------- |
| 1   | article | {"type": "pulsar", "options": {"url": "pulsar://localhost:6650", "topic": "persistent://public/default/blog"}} |
| 2   | tag     | {"type": "pulsar", "options": {"url": "pulsar://localhost:6650", "topic": "persistent://public/default/blog"}} |
| 3   | draft   | {"type": "redis", "options": {"addrs": ["127.0.0.1:6379"], "password": "root", "topic": "blog"}}               |

## API

通过上文的配置方式，你会惊喜的发现：当你查看 blog 库下的表时，若发现了 mongo2queue 表，则可以很直观的知道该库下有表接入了 changes2queue 服务。

上述任务的添加和删除，若想避免手动去创建和维护 mongo2queue 或 mysql2queue 表，可直接通过 changes2queue 提供的 http restful api 操作，避免直接操作数据库。

> 注意仅有创建和删除，没有更新。意味着只能先删除任务，在创建任务。

### 创建

mongodb:

```bash
curl -X POST http://localhost:8080/task/mongodb/:db_name -H "Content-Type: application/json" -d '{"tables": ["article", "tag"], "queue": {"type": "pulsar", "options": {"url": "pulsar://localhost:6650", "topic": "persistent://public/default/blog"}}}'
```

mysql:

```bash
curl -X POST http://localhost:8080/task/mysql/:db_name -H "Content-Type: application/json" -d '{"tables": ["article", "tag"], "queue": {"type": "pulsar", "options": {"url": "pulsar://localhost:6650", "topic": "persistent://public/default/blog"}}}'
```

### 删除

mongodb:

```bash
curl -X DELETE http://localhost:8080/task/mongodb/:db_name/:task_id
```

mysql:

```bash
curl -X DELETE http://localhost:8080/task/mysql/:db_name/:task_id
```

## 推荐实践

考虑热更新的场景更多在于在同一库中，即同一任务只增加或减少表的监听，所以若增加任务，认为是一种较大业务变更，需更改配置文件中的 mongo2queues 或 mysql2queues 后重新启动本服务。

更为推荐的实践方式是为特定的业务使用专门的 changes2queue 服务，如上文中提到的博客业务，则使用专门的 changes2queue-blog 服务；若用户服务则使用专门的 changes2queue-user 服务。这将极少出现变更 mongo2queues 或 mysql2queues 的情况，从而避免重启服务过程中的消息丢失（RedisBroadcast 这类广播消息队列会存在该问题）。

## DDD

本服务采用 DDD 模板：https://github.com/ihezebin/go-template-ddd

## 编译打包

`make package TAG=v1.0.0` 或 `make package`
