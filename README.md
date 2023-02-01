# ACDAT in Go

Aho-Corasick Automaton with Double Array Trie (Multi-patterm substitute in go)

- 多模式匹配（替换）具有很强的现实意义与实用价值：敏感词过滤，病毒特征码搜索等。
- 字典中约有一千条关键字，分为三大类：电影、音乐、电影&音乐，将输入文本中的字典记录替换为相应的类别文本。遵循贪婪原则与最先匹配原则。

## 使用说明

```bash
# Setup environment: create ramdisk and copy input files
# Assume you are using Mac
$ make setup
util/ramdisk.sh && mkdir -p /ramdisk/vonng
Started erase on disk2
Unmounting disk
Erasing
Initialized /dev/rdisk2 as a 1024 MB case-insensitive HFS Plus volume
Mounting disk
Finished erase on disk2 ramdisk
tar -xf data/dict.txt.tgz -C /ramdisk/vonng
cat data/xa* | tar -x -C /ramdisk/vonng

# build and run
$ make
GOGC=off ./ac -p=profile
time: 2.497193491s sig: 5d76461b53079d20c08eb0b33c46b7cd
Round 0: 2.497955425s
Round 1: 2.504176103s
Round 2: 2.518718803s
Round 3: 2.506447912s
Round 4: 2.531897625s
Round 5: 2.529866993s
Round 6: 2.524811658s
Round 7: 2.531879072s
Round 8: 2.532776425s
Round 9: 2.549934714s
Avg: 2.522571823s

# run with profile
make runp

# see profile
# type 'web' to see graphviz call graph
make prof
```

## 题目

### 1. 目的

- 庆新年找点乐子。
- 切磋一下高性能编程。

### 2. 比赛时间

提交程序请在这个时间：2018-01–02 ～ 2018-01-03

### 3. 比赛奖励

- 总排名第一名 ，语言不限， 奖品 Kindle Oasis，价值 2400+。 为了减弱 C/C++的天然优势，把 C/C++的时间 x2 计算。
- 主流编程语言内的第一名（Go，Java，OC，C/C++，Python，Shell，PHP，SQL），奖品罗技（Logitech）G502 炫光自适应游戏鼠标， 价值 400+

### 4. 题目：基于词典的字符串替换。

#### 输入数据

一个词典文件 dict.txt：内容如下所示。

```
让我爱你	-*音乐*-
超人奥特曼	-*电影*-
圣战奇兵	-*电影*-
火线勇气	-*电影*-
龙在江湖	-*电影*-
一个人	-*音乐*-
听海	-*音乐*-
双面特工	-*电影*-
花开自在	-*音乐*-
穿牛仔裤的十字军	-*电影*-
蜻蜓	-*电影*-
X档案	-*电影*-
金钱帝国	-*电影*-
```

此文件每一行分两个字段。用跳格分开。第一列为`key string`；第二列为`value string`。这个文件有 1000+行。

一个输入文件`video_title.txt`，此文件 1000 万行。 请对这个文件中的每一行做如下加工：

#### 功能

将所有命中 key string 的字符串替换成 value string。例如

`从前，东北一家人生活在伤心太平洋。` 替换完成后变成 `从前，-*电影*-生活在-*音乐&电影*-。`

替换后的结果请输出到`result.txt`文件。

#### 匹配原则

为了加强结果的稳定性，规定一些原则：

1. 从前向后匹配。例如如果`寂寞的季节`， `节日`都是 key string， 一条输入`寂寞的季节日日悲催`只会匹配到`寂寞的季节`， 不会再匹配`节日`
2. 贪心匹配原则，匹配最长的那个 key string。例如`独立`，`独立日` 都是 key string 的话。一条输入`听说独立日这部电影是美国人拍的`，只会匹配上`独立日`，不会匹配到`独立`。如果没有比`独立日`更长的匹配的话。

#### 其他要求

1. 必须提交源代码。
2. 比的是算法与编程细节。必须单机运算，禁止分布式，禁止多核并发。
3. 可以对词典文件进行预处理，不允许对`video_title.txt`文件预处理。
4. 不允许把结果文件存下来直接输出。
5. 内存占用必须小于 500M。
6. 发现抄袭两人成绩同时作废。
7. 除了语言自带的标准库之外，禁止引入其他三方标准库。
8. 业余时间完成，不允许耽误工作。

## 比赛结果

![](image/result2.png)

结果文件 MD5：`5d76461b53079d20c08eb0b33c46b7cd`

本机测试平均为 2.5 秒，标准服务器测试结果为 3.2 秒。

本机重复执行十次的 profile：

![](image/pprof.png)

### 服务器 bench 脚本

```bash
#!/bin/bash
# REQUIRES SUDO
# Benchmark runner

repeats=20
output_file='/ramdisk/bench/results.csv'
command_to_run='echo 1'

run_tests() {
    # --------------------------------------------------------------------------
    # Benchmark loop
    # --------------------------------------------------------------------------
    echo 'Benchmarking ' $command_to_run '...';
    # Indicate the command we just run in the csv file
    echo '======' $command_to_run '======' >> $output_file;

    # Run the given command [repeats] times
    for (( i = 1; i <= $repeats ; i++ ))
    do
        # percentage completion
        p=$(( $i * 100 / $repeats))
        # indicator of progress
        l=$(seq -s "+" $i | sed 's/[0-9]//g')

        rm -f /ramdisk/bench/tmp/*
        # runs time function for the called script, output in a comma seperated
        # format output file specified with -o command and -a specifies append
        # /usr/bin/time -f "%E,%U,%S" -o ${output_file} -a ${command_to_run} > /dev/null 2>&1
        GOGC=off chrt -f 99 /usr/bin/time -f "%e" -o ${output_file} -a ${command_to_run} > /dev/null 2>&1

        # Clear the HDD cache (I hope?)
        sync && echo 3 > /proc/sys/vm/drop_caches

        echo -ne ${l}' ('${p}'%) \r'
    done;
    md5sum -c <<<"5d76461b53079d20c08eb0b33c46b7cd  /ramdisk/bench/tmp/result"

    echo -ne '\n'

    # Convenience seperator for file
    echo '--------------------------' >> $output_file
}

# Option parsing
while getopts n:c:o: OPT
do
    case "$OPT" in
        n)
            repeats=$OPTARG
            ;;
        o)
            output_file=$OPTARG
            ;;
        c)
            command_to_run=$OPTARG
            run_tests
            ;;
        \?)
            echo 'No arguments supplied'
            exit 1
            ;;
    esac
done

shift `expr $OPTIND - 1`

```

## 第一名感想总结

- 过早优化是万恶之源
  - 先保证正确性，在正确性的前提下渐进式优化。
- 以 Profiling 为依据，以 Benchmark 为准绳
  - 先优化算法，再做场景特定优化，然后优化 IO，最后优化细节。

### 1. 选型

#### 语言

运算密集型应用，不考虑 Python，Javascript，SQL，Shell 等脚本语言。

Go 的性能损失绝对不到 C 的一倍，因此一定选 Go。但如果比绝对时长，那么选 C。

#### 算法

直觉：状态机，Trie，KMP

搜索发现：AC 自动机，是融合了 KMP 思想的 Trie 树。

进一步了解：使用数组代替指针的进阶实现：Double Array Trie

### 2. 分治

将问题划分为几个子问题：

- 生成 ACDAT
  - 构建一颗 Trie 数
  - 构建 Fail 状态数组与 Info 状态数组
  - Info 数组要保留该状态对应的重要信息：激发的关键词长度，激发的关键词类型。
- 使用 ACDAT 逐行处理
  - 先找出所有匹配：《开始位置，结束位置，关键词长度，关键词类型》
  - 根据匹配进行字符串替换
- IO 使用 Bufio，以行处理的方式进行。

```go
BuildDict()
for err = nil; err != io.EOF; line, err = reader.ReadSlice('\n') {
    writer.Write(HandleLine(line))
}
```

### 3. 优化

算法层面的优化不细说了，大家用的基本上都是一样的方法。主要谈一谈工程上的优化。

工程上的一些技巧，从本机 5 秒优化至 2.6 秒

#### 处理粒度: 1s

粒度是输入的字符宽度，对实现有决定性的影响。

本题中可以使用三种粒度：字节(byte)，双字节(int16)，四字节(int)。

本例中文本均为 UTF8 编码，UTF8 为变长编码。

- 使用`int32`，即`rune`为单位
  - 最为通用的做法，可以处理所有 Unicode 字符，包括 Emoji 等。（在本题中并没有）
  - 需要较少的状态转移判断，但存在 rune 解码编码开销
  - 状态数组长度约 10W。比较宽，总大小在 1.6M 左右。
- 使用 int16，即双字节编码
  - 可选 UTF16，GBK 等编码，但这是开历史倒车的行为。
  - 需要对输入文件做转码，违背题目约束，如果运行时转码总用时不见得会更快。
  - 如果输入已经转码，处理会很快，尤其对于 Java，C#这类语言。
- 使用 int8，按照字节处理
  - 可以处理二进制数据
  - **可以避免[]Byte、String、[]Rune 相互转化的开销。**
  - 状态数组非常更小，缓存友好。
  - 无效判断大大增多，状态转移判断数翻倍（2KW 到 4KW）。

最终选择了 `int32`的做法，这是实际 bench 中性能最好的方案，也具有良好的通用性。

###

### 字符串、字符数组、字节数组转换： 0.5s

- [深入 Go 文本类型](https://vonng.com/blog/go-text-types/)，[]byte 转为 string 是**相对安全**的操作，因为 string 相比[]byte 只是少了一个 cap 字段。

```go
for _, c := range *(*string)(unsafe.Pointer(&input))
```

### 编码与解码： 0.3s

使用自定义的 Rune 解码函数代替标准库的 WriteRune，有 0.3s 的性能提升。

```go
func WriteRune(r rune) {
	switch i := uint32(r); {
	case i <= 127:
		Buf[BSP] = byte(r)
		BSP++
	case i <= 2047:
		Buf[BSP] = 0xC0 | byte(r>>6)
		BSP++
		Buf[BSP] = 0x80 | byte(r)&0x3F
		BSP++
	case i <= 65535:
		Buf[BSP] = 0xE0 | byte(r>>12)
		BSP++
		Buf[BSP] = 0x80 | byte(r>>6)&0x3F
		BSP++
		Buf[BSP] = 0x80 | byte(r)&0x3F
		BSP++
	default:
		Buf[BSP] = 0xF0 | byte(r>>18)
		BSP++
		Buf[BSP] = 0x80 | byte(r>>12)&0x3F
		BSP++
		Buf[BSP] = 0x80 | byte(r>>6)&0x3F
		BSP++
		Buf[BSP] = 0x80 | byte(r)&0x3F
		BSP++
	}
}
```

### 缓冲区大小

操作系统一个 Page 大小为 4k，服务器默认是宝存的 PCI-E SSD，一次写入的单位是 32K。所以在普通闪存盘上 32k 的缓冲区表现已经很不错了。不过后来换成了 Ramdisk，就又需要修改了。

![](image/buf-bench.png)

最后 IO 的 Buffer 设置为 128k，表现还不错。

我也试过用系统调用一次性读完，一次性写入，有一定提升，但内存消耗就变成 O(n)的了。所以就没用。

出乎意料的是，使用 ramdisk 并没有对性能产生质的提升。一个原因是 IO 的 pattern 是顺序读取与顺序写入。对于 PCI-E SSD 而言，内存盘并没有绝对性的优势。

### 条件分支重排： 0.3s

通过检查分支命中的概率，重排、组合条件分支，能有 0.3 秒的性能提升。

### 全局变量 vs 本地变量: 0.1

使用全局变量带来了 0.1s 的优化，但一部分变量改为全局变量反而会影响性能。

### 字典常量 vs 动态创建 0.15s

存在只读段的数组，访问起来竟然比本地读取重新构建还要慢，不禁陷入深思。
