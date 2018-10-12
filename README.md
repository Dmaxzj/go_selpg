# 使用golang实现selpg

selpg是一个Linux命令行使用程序, 参考[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)和其[C语言源码](https://www.ibm.com/developerworks/cn/linux/shell/clutil/selpg.c), 实现了go语言版本的selpg

参考资料:  
go文档: https://godoc.org/  
Golang之使用Flag和Pflag: https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/




## 测试

测试文本 test.txt

```vim
line1
line2^L
line3
line4^L
line5
line6^L
line7
line8^L
line9
line10^L
good bye
```

### 测试输入

#### 测试1

```bash
$ selpg -s1 -e5 test.txt
line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
good bye
```

#### 测试2

重定向标准输入

```bash
$ ./bin/selpg.exe -s1 -e5 < test.txt
line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
good bye
```

#### 测试3

测试管道

```bash
$ cat test.txt | selpg -s1 -e5
line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
good bye
```

### 测试输出

#### 测试4

```bash
$ selpg -s1 -e5 test.txt > output.txt
$ cat output.txt
line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
good bye
```

#### 测试5

selpg的标准输出被透明地传递至cat的标准输入。

```bash
$ selpg -s1 -e5 test.txt | cat
line1
line2
line3
line4
line5
line6
line7
line8
line9
line10
good bye
```

### 测试错误输出

#### 测试6

```bash
$ selpg.exe -s0 -e5 test.txt
invalid start page 0

USAGE: -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]
```

#### 测试7

将错误输出重定向到error.txt

```bash
$ selpg -s0 -e5 test.txt 2>error.txt
$ cat error.txt
invalid start page 0

USAGE: -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]
```

### 测试-l, -f参数

-l 和 -f 同时使用会产生错误

```bash
$ selpg -s1 -e2 -l3 -f test.txt
-l and -f can not be used together

USAGE: -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]
```

#### 测试8

读取test.txt第一到第二页, 每页三行

```bash
$ selpg -s1 -e2 -l3 test.txt
line1
line2
line3
line4
line5
line6
```

#### 测试9

读取test.txt第一到第二页, 以'/f'分页

```bash
$ selpg -s1 -e2 -f test.txt
line1
line2
line3
line4
```