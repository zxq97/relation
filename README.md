# relation
关系链服务

### 功能模块
关注/取关

读取关注/粉丝 数量/列表

判断关系 无/单方关注/互相关注

### 架构设计

BFF: relation

关系链唯一对外grpc接口 关系链业务的服务编排
比如 判断关注双方是否有拉黑行为等

Service: relation_svc
服务层 专注在关系链功能的API实现上 比如 关注/取关 判断关系等

Task: relation_task
关系链内部异步消费任务 比如 异步添加/删除关系链 粉丝列表缓存重建

Admin: relation_admin
关系链后台

    grpc -> relation (grpc) -> relation_svc (kafka) -> relation_task (canal) -> relation_admin
