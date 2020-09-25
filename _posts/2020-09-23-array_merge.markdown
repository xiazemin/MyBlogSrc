---
title: array_merge
layout: post
category: php
author: 夏泽民
---
array_merge 两个参数必须都是array，否则会返回null
即使是arrayObject也不行，必须getArrayCopy

https://stackoverflow.com/questions/19711491/merge-array-returns-null-if-one-or-more-of-arrays-is-empty

https://www.php.net/manual/zh/class.arrayobject.php

https://www.php.net/manual/en/function.array-merge.php
<!-- more -->
don't forget that numeric keys will be renumbered!

If you want to append array elements from the second array to the first array while not overwriting the elements from the first array and not re-indexing, use the + array union operator:

https://m.php.cn/manual/view/12118.html