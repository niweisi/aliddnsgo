原作者：https://github.com/GameBelial/AliDDNSGo

## 在原作者的基础上增加了获取本地局域网ip

## 0.简要介绍

AliDDNSGo 是基于 Golang 开发的动态 DNS 解析工具，借助于阿里云的 DNS API 来实现域名与动态 IP 的绑定功能。这样你随时就可以通过域名来访问你的设备，而不需要担心 IP 变动的问题。

## 1.使用说明

> 使用本工具的时候，请详细阅读使用说明。

### 1.1 配置说明

通过更改 ```settings.json.example``` 的内容来实现 DDNS 更新，其文件内部各个选项的说明如下：

```json
{
    "AccessIdComment": "阿里云的 Access Id。",
    "AccessId": "AccessId",
    "AccessKeyComment": "阿里云的 Access Key。",
    "AccessKey": "AccessKey",
    "MainDomainComment": "主域名。",
    "MainDomain": "example.com",
    "SubDomainsComment": "需要批量变更的子域名记录集合。",
    "SubDomains": [
        {
            "TypeComment": "子域名记录类型。",
            "Type": "A",
            "SubDomainComment": "子域名记录前缀。",
            "SubDomain": "sub1",
            "IntervalComment": "TTL 时间。",
            "Interval": 600,
            "MacComment": "要获取ip的mac地址，public为给公网ip，输入mac地址获取本地ip",
            "Mac": "08:ed:22:7e:ee:44"
        },
        {
            "Type": "A",
            "SubDomain": "sub2",
            "Interval": 600,
            "Mac": "public"
        }
    ]
}
```

其中 ```Access Id``` 与 ```Access Key``` 可以登录阿里云之后在右上角可以得到。

### 1.2 使用说明

在运行程序的时候，请建立一个新的 ```settings.json``` 文件，在里面填入配置内容，然后执行以下命令：

```shell
./AliCloudDynamicDNS
```

效果图：  
![](https://github.com/GameBelial/AliDDNSNet/blob/master/READMEPIC/Snipaste_2019-12-12_17-36-21.png)

当然如果你有其他的配置文件也可以通过指定 ```-f``` 参数来制定配置文件路径。例如：

```shell
./AliCloudDynamicDNS -f ./settings.json
```

![](https://github.com/GameBelial/AliDDNSNet/raw/master/READMEPIC/Snipaste_2019-12-12_17-38-09.png)

如果你需要启动自动周期检测的话，请通过 `-i` 参数指定执行周期，单位是秒。

```shell
./AliCloudDynamicDNS -f ./settings.json -i 3600
```

![](https://github.com/GameBelial/AliDDNSNet/raw/master/READMEPIC/Snipaste_2019-12-12_17-38-53.png)

> **注意：**
>
> **当你通过 -i 指定了周期之后，请在最末尾使用 & 符号，使应用程序在后台运行。使用群晖的同学，就不要指定 -i 参数了。**

