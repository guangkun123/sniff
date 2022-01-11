# sniff

sniff 是go语言编写，使用原始套接字SOCK_RAW抓包，google的gopacket解包的脚本。
输出包括 时间、sla、来源ip:port、查询语句，返回值。
看起来对类redis的nosql比较好看，mysql没法看，哈哈哈
代码只有59行，方便dba学习使用和进一步优化。
12:15:07 23 10.75.14.110:23379,*2  $4  INCR  $7  uuid_07  <===> :557051726
1、使用
./sniff 50843
2、编译。
vim sniff.go
复制代码
go build sniff.go
