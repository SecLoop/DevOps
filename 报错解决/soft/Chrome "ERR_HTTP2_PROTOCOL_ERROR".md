#### Chrome "ERR_HTTP2_PROTOCOL_ERROR" 解决

```

1、打开 chrome://flags/ 页面
2、找到 Block insecure private network requests. 和 Enable Trust Tokens 两项
3、将其值从 Default 改为 Enable
4、点右下角的 ReLaunch 按钮重启浏览器
5、重新打开报错的网站
6、如果打不开，在地址栏输入 chrome://restart/ 再重启一遍浏览器即可
```
#### 参考：https://www.jianshu.com/p/7f58ed7f9c0e
