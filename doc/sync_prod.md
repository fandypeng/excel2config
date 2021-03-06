# 同步正式环境



点击同步正式环境的按钮，将测试环境的配置导入正式环境，同步界面会有差异比对，左边显示的是正式环境的版本，右边显示的是测试环境的版本。

![img](./images/sync_compare.jpg?lastModify=1620962878)



## 正式环境和测试环境



创建项目的时候，会自动创建一个测试环境和一个正式环境，两者是自动镜像的存在。在测试环境添加、删除表格的时候，正式环境会同步创建和删除，测试环境表格里新创建的sheet不会自动同步到正式环境，需要点击同步正式环境按钮手动同步。正式环境的Excel只可以查看，无法直接修改，发布正式环境必须按照如下的业务流程执行。



正式环境和测试环境可以独立配置数据仓库，独立一套发布记录，正常业务的使用流程是：

1. 测试环境添加配置表，编辑好表内容，发布到测试环境测试
2. 测试完成之后，需要点击同步正式环境，正式环境的Excel里会自动同步到sheet内容
3. 切换到正式环境之后，点击发布即可发布到正式环境的数据库，至此完成正式环境发布过程



注意：同步正式环境需要校验权限，用户必须是该项目正式环境的成员才可以执行同步操作，同步操作会将测试环境的配置数据覆盖到正式环境。

