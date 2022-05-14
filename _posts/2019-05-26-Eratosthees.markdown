---
title: 爱拉托逊斯筛选法
layout: post
category: algorithm
author: 夏泽民
---
The Sieve of Eratosthees
爱拉托逊斯筛选法思想：对于不超过n的每个非负整数P，删除2*P, 3*P…，当处理
完所有数之后，还没有被删除的就是素数。
若用vis==1表示已被删除，则代码如下：
memset(vis, 0, sizeof(vis));
for(int i = 2; i <= 100; i++)
for(int j = i*2; j <= 100; j += i)
  vis[j] = 1;
<!-- more -->
改进的代码：
int m = sqrt(double(n+0.5));

for(int i = 2; i <= m; i++)
if(!vis[i])
{
  prime[c++] = i;
  for(int j = i*i; j <= n; j += i)
  {
   vis[j] = 1;
  }
}

为何j从i*i开始？
因为首先在i=2时，偶数都已经被删除了。
其次，“对于不超过n的每个非负整数P”， P可以限定为素数，
为什么？
因为，在 i 执行到P时，P之前所有的数的倍数都已经被删除，若P

没有被删除，则P一定是素数。
而P的倍数中，只需看：
(p-4)*p, (p-2)*p, p*p, p*(p+2), p*(p+4)
（因为P为素数，所以为奇数，而偶数已被删除，不需要考虑p*(p-1)等）
又因为(p-4)*p 已在 (p-4)的p倍中被删去，故只考虑：
p*p, p*(p+2)….即可
这也是i只需要从2到m的原因。
当然，上面 p*p, p*(p+2)…的前提是偶数都已经被删去，而代码

二若改成 j += 2*i ,则没有除去所有偶数，所以要想直接 加2*i
。只需在代码二中memset()后面加：
for(int i = 4; i <= n; i++)
     if(i % 2 == 0)
          vis = 1;
这样，i只需从3开始，而j每次可以直接加 2*i.
