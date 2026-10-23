# OAuth 刷新截止时间边界

`OAuthRefreshAPI.RefreshIfNeeded` 负责加锁、重读账号和条件写入；
`TokenRefreshService` 负责后台重试及账号状态同步。两条路径均须在首次写入前
检查取消状态和显式 deadline，不能只依赖异步取消通知。恰好到达截止时间也视为超时。

父级取消或到期直接结束本轮，不据此设置账号永久错误或临时冷却。
单次尝试超时保持原有暂时失败分类。Grok CAS 结果不确定时保持 provider
containment；已确认写入成功后，仅清理越过单次尝试 deadline，仍同步缓存并返回成功。
此检查不提供跨数据库事务的原子截止时间承诺，也不变更 token 轮换协议。

定向回归（Go 版本跟随 `backend/go.mod`，从 `backend/` 执行）：

```sh
go test -mod=readonly -tags=unit -p 1 -parallel 1 -count=1 ./internal/service -run 'Deadline|LateSuccess|ParentCancellation'
```

测试为两个对象注入同一私有时钟，直接推进到 deadline 前、恰好到达、之后及
持久化后清理阶段；检查迟到时零首次写入、父级错误分类、缓存同步和无错误重试。
生产默认使用 `time.Now`，没有新增配置项。
