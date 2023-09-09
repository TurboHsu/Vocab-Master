# Vocab Master
## 词达人辅助工具

<img src="icon/VocabMaster.svg" width="120" alt="logo">

### WARN
我自己用不上了，只提供有限的维护。

### Tip
使用```-proxy=false```参数可以使程序在打开时不对系统代理进行操作

程序内提供词达人处理器开关，可打开/关闭词达人注入。

### 原理
使用MITM Proxy监听&注入词达人网页，实现辅助功能。

### 注意事项
1. 本项目旨在为了让词汇已经过关的人，不再被这个东西所困扰。
2. 打开程序会自动覆盖系统代理
3. **维护者目前只接触到CET4的一些课堂任务，只见过部分题型，（题型一样TopicMode也可能不同），后续遇到相关题型我会继续更新。**

### 如何使用？
#### 学习任务
1. 下载release中的压缩包并且解压
2. 运行vocab-master，他会在对应平台的应用数据存储区中生成cert文件夹，只需要将其中的.cer文件添加到系统的信任根证书中。这个文件夹可以用程序中的一个按钮打开。
3. 对于Windows，你需要把Release中的fonts拷贝到程序数据文件夹和cert同级的目录中，不然你可能看不到中文。（打算重构UI，这个问题将会得到解决。）
4. 打开一个课堂任务，```重新选词```,and feel great ;)

#### 测试任务
** 测试任务需要您获取全部词库，参见数据库控制台 **
1. 在打开词达人之前打开程序，勾选```Javascript Hijack```，以绕过手机端检测。
2. 导入词库。
3. You nailed it.

### 关于数据控制台
#### 导入数据
点击导入数据，选一个数据文件打开即可。
#### 抓取数据
1. 抓取词库时，首先勾选```Fetch Identify```，启用Cookie/Header抓取。
2. 输入Course ID。他可能长这样：```CET6_hx```
3. 点击开始抓取，为了保证不达到单个```Token```达到使用限制，请您在程序抓取时**不断操作词达人**，你只需要到处点点，浏览浏览词汇。注意控制台输出。
4. 结束后可导出。存到你喜欢的地方即可。

### 我要编译
这是个go module，go ahead!

### 鸣谢
- [fyne.io](https://fyne.io)
- [github.com/Trisia/gosysproxy](https://github.com/Trisia/gosysproxy)
- [github.com/andybalholm/brotli](https://github.com/andybalholm/brotli)
- [github.com/lqqyt2423/go-mitmproxy](https://github.com/lqqyt2423/go-mitmproxy)
