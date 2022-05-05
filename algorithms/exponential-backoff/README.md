# 指数退避

根据 wiki 上对 Exponential backoff 的说明，指数退避是一种通过反馈，成倍地降低某个过程的速率，以逐渐找到合适速率的算法。

在以太网中，该算法通常用于冲突后的调度重传。根据时隙和重传尝试次数来决定延迟重传。

在 c 次碰撞后（比如请求失败），会选择 0 和 `2^c-1` 之间的随机值作为时隙的数量。

- 对于第 1 次碰撞来说，每个发送者将会等待 0 或 1 个时隙进行发送。
- 而在第 2 次碰撞后，发送者将会等待 0 到 3（ 由 `2^2-1` 计算得到）个时隙进行发送。
- 而在第 3 次碰撞后，发送者将会等待 0 到 7（ 由 `2^3-1` 计算得到）个时隙进行发送。
- 以此类推……

- 随着重传次数的增加，延迟的程度也会指数增长。

说的通俗点，每次重试的时间间隔都是上一次的两倍。

《UNIX 环境高级编程》中，使用指数退避来建立连接的示例：
```c
#include "apue.h"
#include <sys/socket.h>

#define MAXSLEEP 128

int connect_retry(int domain, int type, int protocol,
                  const struct sockaddr *addr, socklen_t alen)
{
    int numsec, fd;

    /*
    * 使用指数退避尝试连接
    */
    for (numsec = 1; numsec < MAXSLEEP; numsec <<= 1)
    {
        if (fd = socket(domain, type, protocol) < 0)
            return (-1);
        if (connect(fd, addr, alen) == 0)
        {
            /*
            * 连接接受
            */
            return (fd);
        }
        close(fd);

        /*
        * 延迟后重试
        */
        if (numsec <= MAXSLEEP / 2)
            sleep(numsec);
    }
    return (-1);
}
```

如果连接失败，进程会休眠一小段时间（numsec），然后进入下次循环再次尝试。每次循环休眠时间是上一次的 2 倍，直到最大延迟 1 分多钟，之后便不再重试。
