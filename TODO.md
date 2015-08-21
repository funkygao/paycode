# TODO

- [X] 时钟不同步问题
  - 通过容错允许几分钟内的时钟不同步
  - 如果超过了容错范围，给手机客户端发push，让用户自己确认

- [X] 暴露了uid
  引入另外一套uid，仅仅给该业务使用
  同时，把该uid与real uid进行绑定

- [X] 防止码被重刷
  redis(uid): list(passwd)

- [ ] 风控
  如果密码多次不对，则suspicous user

- [ ] key的保护和传递
