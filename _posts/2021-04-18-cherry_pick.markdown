---
title: cherry pick
layout: post
category: web
author: 夏泽民
---
这时分两种情况。一种情况是，你需要另一个分支的所有代码变动，那么就采用合并（git merge）。另一种情况是，你只需要部分代码变动（某几个提交），这时可以采用 Cherry pick。

git cherry-pick <commitHash>

 git cherry-pick <HashA> <HashB>
 
 
<!-- more -->

https://www.ruanyifeng.com/blog/2020/04/git-cherry-pick.html

https://backlog.com/git-tutorial/cn/stepup/stepup7_4.html