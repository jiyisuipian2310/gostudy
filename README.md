# gostudy

example1:
展示了如何使用json包来解析json数据，并将其转换为结构体。
展示了如何发送http和https请求，并获取响应数据，以及如何解析json数据。
展示了如何对发送和接收的数据进行加密和解密
yaml文件读取代码
添加编译参数，程序输出版本信息

example2:
    类的抽象接口调用，即多个类有一个相同的函数接口，在其他的函数中，怎么实现相同函数的同一调用

example3:
    展示了类的继承

example4:
    go 通道：有缓冲的通道和无缓冲的通道的使用

example5:
    select case 的使用

example6:
    http 服务器的实现, 响应数据从配置文件中读取

example7:
    导入包的过程中一个典型的错误，使用其他包的函数时，函数名必须大写

example8:
    go tcp server 实现，支持多客户端连接，并发处理请求

example9:
    go 实现远程下载文件

example10:
    生成RSA公私钥对，并使用公钥加密数据，私钥解密数据；私钥签名数据，公钥验证签名

example11:
    生成后缀名为pem的私钥文件，读取私钥文件对数据加解密

example12:
    根据配置文件启动http和https服务

example13:
    自定义结构体的序列化和反序列化

example14:
    定义Animal基类，Dog和Cat继承Animal类，并实现自己的方法
    定义一个公共接口，接口参数为Animal基类，在该接口中判断传入的对象是否为Dog或Cat类型，并调用相应的方法

example15:
    缓冲 channel 和 select 语句可以提高并发程序的效率和响应能力

example16:
    实现了一个代理转发程序，通过读取来自客户端的第一条消息并解析，得到后端资源的ip，端口，域名
    第一条消息样例：proxy\x00\x3c{"dstip":"192.168.104.100", "dstport":10000, "dstdomain":""}
    然后代理和资源建立连接，然后代理创建协程读取资源数据转发给客户端，读取客户端数据转发给资源

example17:
    实现了一个代理转发程序，通过读取来自客户端的第一条消息并解析，得到后端资源的ip，端口，域名
    第一条消息样例：proxy\x00\x3c{"dstip":"192.168.104.100", "dstport":10000, "dstdomain":""}
    然后代理和资源建立连接，然后代理创建协程读取资源数据转发给客户端，读取客户端数据转发给资源

    和 example16 不同的是， example17 使用了 ctx context.Context 来控制协程之间的退出

example18:
    实现了一个代理转发程序，通过读取来自客户端的第一条消息并解析，得到后端资源的ip，端口，域名
    第一条消息样例：proxy\x00\x3c{"dstip":"192.168.104.100", "dstport":10000, "dstdomain":""}
    然后代理和资源建立连接，然后代理创建协程读取资源数据转发给客户端，读取客户端数据转发给资源

    和 example17 不同的是， example18 使用了 SetReadDeadline 感知读超时，减少了一个 goroutine 等待读超时的消耗

example19:
    go grpc 实现远程调用

example20:
    使用go操作monggo数据库，并且修改了monggo客户端驱动的源码，在登录报文中添加了Token

example21:
	interface 的 switch 使用

example22:
	反射的基本使用

example23:
	int 类型的切片，正向，反向排序练习

example24:
    插件化开发：
        统一插件接口：所有插件实现 Plugin 接口，保证一致性
        插件注册机制：使用 map[string]reflect.Type 存储插件类型，支持动态注册
        插件管理：PluginManager 负责插件加载、管理和执行
        插件示例：实现了 SQL 注入检测和 XSS 检测插件，支持自动注册和调用

example25:
    interface 的常见使用场景：依赖注入
        将依赖关系从代码中分离出来，通过将依赖关系定义为接口类型
        可以在运行时动态地替换实现，从而使得代码更加灵活、可扩展
    此例中，学生打印信息，依赖于打印店，而每个学生可以选择不同的打印店，
    所以将打印店剥离出来，定义成接口；每个打印店实现接口中的方法
    重要：这个就是设计模式中的依赖倒置原则

example26: 设计模式中的开闭原则：对修改关闭，对扩展开放

example27: 读取ini配置文件

example28: 
    使用 ReadFull 至少读取 一些数据， 超时读取不到，退出

example29: 
    通过读写锁实现一个并发安全的 map

example30:
    通过channel实现一个并发安全的 map

example31:
    实现一个http或者https服务，接收数据然后转发给其他应用程序，应用数据可以是加密的，通过AES解密出来

example32:
    实现一个http或者https服务，接收数据然后转发给其他应用程序

example33:
    实现一个http或者https服务，接收数据然后转发给其他应用程序，这个http服务可以校验客户端证书

example34:
    实现一个http服务, 接收消息，然后ssh登录到主机资源，执行ssh命令