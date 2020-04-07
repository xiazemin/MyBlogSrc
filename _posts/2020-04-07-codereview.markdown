---
title: codereview
layout: post
category: golang
author: 夏泽民
---
第一步需要注册成为一个 Go contributor 以及配置你的环境。这里有一份包含了所需步骤的清单：

步骤 0: 准备好一个你将用来给 Go 语言贡献代码的 Google 账号。在后面所有的步骤中都要使用这个账号，还有确保你的 git 已经正确配置了这个账号的邮箱地址，以便后续提交 commits。
步骤 1: 签署以及提交一个 CLA（贡献者证书协议）。
步骤 2: 给 Go Git 仓库配置好权限凭证。访问 go.googlesource.com，点击右上角的齿轮图标，接着点击 "Obtain password"，然后跟着指引操作即可。
步骤 3: 在这个页面注册一个 Gerrit 账号，它是 Go 语言团队使用的代码评审工具。CLA 的申请和 Gerrit 的注册只需要在你的账号上做一次就可以了
步骤 4: 运行 go get -u golang.org/x/review/git-codereview 命令安装 git-codereview 工具。
如果你图省事的话，可以直接用自动化工具帮你做完上面的全部步骤，只需运行：

$ go get -u golang.org/x/tools/cmd/go-contrib-init
$ cd /code/to/edit
$ go-contrib-init
这个章节的后面部分将会更加详尽地阐述上面的每一个步骤。如果你已经完成上面的所有步骤（不管是手动还是通过自动化工具），可以直接跳到贡献代码之前部分。

https://gocn.vip/topics/10185
https://golang.org/doc/gccgo_contribute.html
<!-- more -->
步骤 0: 选择一个 Google 账号
每一个提交到 Go 语言的代码贡献都是通过一个绑定了特定邮箱地址的 Google 账号来完成的。请确保你在整个流程中自始至终使用的都是同一个账号，当然，后续你提交的所有代码贡献也是如此。你可能需要想好使用哪一种邮箱，个人的还是企业的。邮箱类型的选择将决定谁拥有你编写和提交的代码的版权。在决定使用哪个账户之前，你大概要和你的雇主商议一下。

Google 账号可以是 Gmail 邮箱账号、G Suite 组织账号，或者是那些绑定了外部邮箱的账号。例如，如果你想要使用一个已存在且并不属于 G Suite 的企业邮箱，你可以创建一个绑定了外部邮箱的 Google 账号。

你还需要确保你的 Git 工具已经正确配置好你之前选定的邮箱地址，用来提交代码。你可以通过 Git 命令来进行全局配置 (所有项目都将默认使用这个配置) 或者只进行本地配置 (只指定某个特定的项目使用)。可以通过以下的命令来检查当前的配置情况：

$ git config --global user.email  # check current global config
$ git config user.email           # check current local config
修改配置好的邮箱地址：

$ git config --global user.email name@example.com   # change global config
$ git config user.email name@example.com            # change local config
步骤 1: 贡献者证书协议
在你提交第一个代码变更到 Go 语言项目之前，你必须先签署下面两种证书协议的其中之一。最后的代码版权归属于谁，将决定你应该签署哪一种协议。

如果你个人是版权持有方，你就需要同意 individual contributor license agreement 并签署，这个步骤可以在线上完成。
如果企业/组织是版权持有方，那么企业/组织就需要同意 corporate contributor license agreement 并签署。
你可以在 Google Developers Contributor License Agreements 网站上检查当前已签署的协议以及再签署新的协议。如果你代码的版权持有方之前已经在其他的 Google 开源项目上签署过这些协议了，那么就不需要再重复签署了。

如果你代码的版权持有方更改了--例如，如果你开始代表新的公司来贡献代码--请发送邮件到 golang-dev 邮件组。这样我们可以知悉情况，接着准备一份新的协议文件以及更新 作者 文件。

步骤 2: 配置 Git 认证信息
Go 语言的主仓库位于 go.googlesource.com，这是一个 Google 自建的 Git 服务器。Web 服务器上的认证信息是通过你的 Google 帐户生成的，不过你还是需要在你的个人电脑上安装配置 git 来访问它。按照以下的步骤进行：

访问 go.googlesource.com 然后点击页面右上角菜单条上的 "Generate Password" 按钮。接着你会被重定向到 accounts.google.com 去登陆。
登陆之后，你会被引导到一个标题为 "Configure Git" 的网页。这个网页包含了一段个性化的脚本代码，运行这个脚本之后会自动生成身份认证的密钥并配置到 Git 里面去。这个密钥是和另一个在远端 Server 生成并存储的密钥成对的，类似于 SSH 密钥对的工作原理。
复制这段脚本并在你的个人电脑上的终端运行一下，你的密钥认证 token 就会被保存到一个 .gitcookies 的文件里。如果你使用的是 Windows 电脑，那你应该复制并运行黄色方格里的脚本，而不是下面那个通用的脚本。
步骤 3: 创建一个 Gerrit 账号
Gerrit 是 Go 语言团队所使用的一个开源工具，用来进行讨论和代码评审。

要注册一个你自己的 Gerrit 账号，访问 go-review.googlesource.com/login/ 然后使用你上面的 Google 账号登陆一次，然后就自动注册成功了。

步骤 4: 安装 git-codereview 命令行工具
无论是谁，提交到 Go 语言源码的代码变更在被接受合并之前，必须要经过代码评审。Go 官方提供了一个叫 git-codereview 的定制化 git 命令行工具，它可以简化与 Gerrit 的交互流程。

运行下面的命令安装 git-codereview 命令行工具：

$ go get -u golang.org/x/review/git-codereview
确保 git-codereview 被正确安装到你的终端路径里，这样 git 命令才可以找到它，检查一下：

git codereview help
正确打印出帮助信息，而且没有任何错误。如果发现有错误，确保环境变量 $PATH 里有 $GOPATH/bin 这个值。

在 Windows 系统上，当使用 git-bash 的时候你必须确保 git-codereview.exe 已经存在于你的 git exec-path 上了。可以运行 git --exec-path 来找到正确的位置然后创建一个软链接指向它或者直接从 $GOPATH/bin 目录下拷贝这个可执行文件到 exec-path。

贡献代码之前
Go 语言项目欢迎提交代码补丁，但是为了确保很好地进行协调，你应该在开始提交重大代码变更之前进行必要的讨论。我们建议你把自己的意图或问题要不先提交到一个新的 GitHub issue，要不找到一个和你的问题相同或类似的 issue 跟进查看。

检查 issue 列表
不管你是已经明确了要提交什么代码，还是你正在搜寻一个想法，你都应该先到 issue 列表 搜索一下。所有 Issues 被已经分门别类以及被用来管理 Go 开发的工作流。

大多数 issues 会被标记上以下众多的工作流标签中的其中一个：

NeedsInvestigation: 该 issue 并不能被完全清晰地解读，需要更多的分析去找到问题的根源
NeedsDecision: 该 issue 已经在相当程度上被解读，但是 Go 团队还没有得出一个最好的方法去解决它。最好等 Go 团队得出了最终的结论之后才开始写代码修复它。如果你对解决这个 issue 感兴趣，而且这个 issue 已经过了很久都没得出最终结论，随时可以在该 issue 下面发表评论去"催促"维护者。
NeedsFix: 该 issue 可以被完全清晰地解读而且可以开始写代码修复它。
你可以使用 GitHub 的搜索功能去搜寻一个 issue 然后搭把手帮忙解决它。例子：

Issues that need investigation: is:issue is:open label:NeedsInvestigation
Issues that need a fix: is:issue is:open label:NeedsFix
Issues that need a fix and have a CL: is:issue is:open label:NeedsFix "golang.org/cl"
Issues that need a fix and do not have a CL: is:issue is:open label:NeedsFix NOT "golang.org/cl"
新开一个关于任何新问题的 issue
除了一些很琐碎的变更之外，所有的代码贡献都应该关联到一个已有的 issue。你随时可以新开一个 issue 来讨论你的相关计划。这个流程可以让所有人都能够参与验证代码的设计，同时帮忙减少一些重复的工作，以及确保这个想法是符合这门语言和相关工具的目标和理念的。还有就是能在真正开始写代码之前就检查这个代码设计是否合理；代码评审工具不是用来讨论高层次问题的。

在规划你的代码变更工作的时候，请知悉 Go 语言项目遵循的是 6 个月开发周期。在每一个 6 个月周期的后半部分是长达 3 个月的新功能特性冻结期：这期间我们只接受 bug 修复和文档更新相关的变更。在冻结期内还是可以提交新的变更的，但是这些变更的代码在冻结期结束之前不会被合并入主分支。

那些针对语言、标准库或者工具的重大变更必须经过变更提议流程才能被接受。

敏感性的安全相关的 issues 只能上报到 security@golang.org 邮箱！

通过 GitHub 提交一个变更
我们鼓励那些初次提交代码并且已经相当熟悉 GitHub 工作流的贡献者通过标准的 GitHub 工作流给 Go 提交代码。尽管 Go 的维护者们是使用 Gerrit 来进行代码评审，但是不用担心，会有一个叫 Gopherbot 的机器人专门来做把 GitHub PR 同步到 Gerrit 上的工作。

就像你通常情况下那样新建一个 pull request，Gopherbot 会创建一个对应的 Gerrit 变更页面然后把指向该 Gerrit 变更页面的链接发布在 GitHub PR 里面；所有 GitHub PR 的更新都会被同步更新到 Gerrit 里。当有人在 Gerrit 的代码变更页面里发表评论的时候，这些评论也会被同步更新回 GitHub PR 里，因此 PR owner 将会收到一个通知。

需要谨记于心的东西：

如果要在 GitHub PR 里进行代码更新的话，只需要把你最新的代码推送到对应的分支；你可以添加更多的 commits、或者做 rebase 和 force-push 操作（这些方式都是可以接受的）。
一旦 GitHub PR 被接受，所有的 commits 将会被合并成一条，而且最终的 commit 信息将由 PR 的标题和描述联结而成。那些单独的 commit 描述将会被丢弃掉。查看写好 Commits 信息获取更多的建议。
Gopherbot 无法逐字逐句地把代码评审的信息同步回 Github: 仅仅是 (未经格式化的) 全部评论的内容会被同步过去。请记住，你总是可以访问 Gerrit 去查看更细粒度和格式化的内容。
通过 Gerrit 提交一个变更
一般来说，我们基本不可能在 Gerrit 和 GitHub 之前完整地同步所有信息，至少在现阶段来说是这样，所以我们推荐你去学习一下 Gerrit。这个不同于 GitHub 却同样强大的工具，而且熟悉它能帮助你更好地理解我们的工作流。

概述
这是一个关于整个流程的概述：

步骤 1: 从 go.googlesource.com 克隆 Go 的源码下来，然后通过编译和测试一次确保这份源码是完整和稳定的： powershell $ git clone https://go.googlesource.com/go $ cd go/src $ ./all.bash # compile and test
步骤 2: 从 master 分支上拉出一条新分支并在这个分支上准备好你的代码变更。使用 git codereview change 来提交代码变更；这将会在这个分支上新建或者 amend 一条单独的 commit。 powershell $ git checkout -b mybranch $ [edit files...] $ git add [files...] $ git codereview change # create commit in the branch $ [edit again...] $ git add [files...] $ git codereview change # amend the existing commit with new changes $ [etc.]
步骤 3: 重跑 all.bash 脚本，测试你的代码变更。 powershell $ ./all.bash # recompile and test
步骤 4: 使用 git codereview mail 命令发送你的代码变更到 Gerrit 进行代码评审 (这个过程并不使用 e-mail，请忽略这个奇葩名字)。 powershell $ git codereview mail # send changes to Gerrit
步骤 5: 经过一轮代码评审之后，把你新的代码变更依附在同一个单独 commit 上然后再次使用 mail 命令发送到 Gerrit: powershell $ [edit files...] $ git add [files...] $ git codereview change # update same commit $ git codereview mail # send to Gerrit again
这个章节剩下的内容将会把上面的步骤进行详细的讲解。

步骤 1: 克隆 Go 语言的源码
除了你近期安装的 Go 版本，你还需要有一份从正确的远程仓库克隆下来的本地拷贝。你可以克隆 Go 语言源码到你的本地文件系统上的任意路径下，除了你的 GOPATH 环境变量对应的目录。从 go.googlesource.com 克隆下来 (不是从 Github):

$ git clone https://go.googlesource.com/go
$ cd go
步骤 2: 在新分支上准备好代码变更
每一次代码变更都必须在一条从 master 拉出来的独立分支上开发。你可以使用正常的 git 命令来新建一条分支然后把代码变更添加到暂存区：

$ git checkout -b mybranch
$ [edit files...]
$ git add [files...]
使用 git codereview change 而不是 git commit 命令来提交变更。

$ git codereview change
(open $EDITOR)
你可以像往常一样在你最喜欢的编辑器里编辑 commit 的描述信息。 git codereview change 命令会自动在靠近底部的地方添加一个唯一的 Change-Id 行。那一行是被 Gerrit 用来匹配归属于同一个变更的多次连续的上传。不要编辑或者是删除这一行。一个典型的 Change-Id 一般长的像下面这样：

Change-Id: I2fbdbffb3aab626c4b6f56348861b7909e3e8990
这个工具还会检查你是否有使用 go fmt 命令对代码进行格式化，以及你的 commit 信息是否遵循建议的格式。

如果你需要再次编辑这些文件，你可以把新的代码变更暂存到暂存区然后重跑 git codereview change : 后续每一次运行都会 amend 到现存的上一条 commit 上，同时保留同一个 Change-Id。

确保在每一条分支上都只存在一个单独的 commit，如果你不小心添加了多条 commits，你可以使用 git rebase 来把它们合并成一条。

步骤 3: 测试你的代码变更
此时，你已经写好并测试好你的代码了，但是在提交你的代码去进行代码评审之前，你还需要对整个目录树运行所有的测试来确保你的代码变更没有对其他的包或者程序造成影响/破坏：

$ cd go/src
$ ./all.bash
(如果是在 Windows 下构建，使用 all.bat ；还需要在保存 Go 语言源码树的目录下为引导编译器设置环境变量 GOROOT_BOOTSTRAP。)

在运行和打印测试输出一段时间后，这个命令在结束前打印的最后一行应该是：

ALL TESTS PASSED
你可以使用 make.bash 而不是 all.bash 来构建编译器以及标准库而不用运行整个测试套件。一旦 go 工具构建完成，一个 bin/go 可执行程序会被安装在你前面克隆下来的 Go 语言源码的根目录下，然后你可以在那个目录下直接运行那个程序。可以查看快速测试你的代码变更这个章节。

步骤 4: 提交代码变更进行代码评审
一旦代码变更准备好了而且通过完整的测试了，就可以发送代码变更去进行代码评审了。这个步骤可以通过 mail 子命令完成，当然它并没有发送任何邮件；他只是把代码变更发送到 Gerrit 上面去了：

git codereview mail
Gerrit 会给你的变更分配一个数字和 URL，通过 git codereview mail 打印出来，类似于下面的：

remote: New Changes:
remote:   https://go-review.googlesource.com/99999 math: improved Sin, Cos and Tan precision for very large arguments
如果有错误，查看 mail 命令错误大全和故障排除。

如果你的代码变更关联到一个现存的 GitHub issue 而且你也已经遵循了建议的 commit 信息格式，机器人将会在几分钟更新那个 issue：在评论区添加 Gerrit 变更页面的链接。

步骤 5: 代码评审之后修正变更
Go 语言的维护者们会在 Gerrit 上对你的代码进行 review，然后你会收到一堆邮件通知。你可以在 Gerrit 上查看详情以及发表评论，如果你更倾向于直接使用邮件回复，也没问题。

如果你需要在一轮代码评审之后更新代码，直接在你之前创建的同一条分支上编辑代码文件，接着添加这些文件进 Git 暂存区，最后通过 git codereview change amend 到上一条 commit：

$ git codereview change     # amend current commit
(open $EDITOR)
$ git codereview mail       # send new changes to Gerrit
要是你不需要更改 commit 描述信息，可以直接在编辑器保存然后退出。记得不要去碰那一行特殊的 Change-Id。

再次确保你在每一条分支上只保留了一个单独的 commit，如果你不小心添加了多条 commits，你可以使用 git rebase 来把它们合并成一条。

良好的 commit 信息
Go 语言的 commit 信息遵循一系列特定的惯例，我们将在这一章节讨论。

这是一个良好的 commit 信息的例子：

math: improve Sin, Cos and Tan precision for very large arguments

The existing implementation has poor numerical properties for
large arguments, so use the McGillicutty algorithm to improve
accuracy above 1e10.

The algorithm is described at https://wikipedia.org/wiki/McGillicutty_Algorithm

Fixes #159
首行
变更信息的第一行照惯例一般是一短行关于代码变更的概述，前缀是此次代码变更影响的主要的包名。

作为经验之谈，这一行是作为 "此次变更对 Go 的 _________ 部分进行了改动" 这一个句子的补全信息，也就是说这一行并不是一个完整的句子，因此并不需要首字母大写，仅仅只是对于代码变更的归纳总结。

紧随第一行之后的是一个空行。

主干内容
描述信息中剩下的内容会进行详尽地阐述以及会提供关于此次变更的上下文信息，而且还要解释这个变更具体做了什么。请用完整的句子以及正确的标点符号来表达，就像你在 Go 代码里的注释那样。不要使用 HTML、Markdown 或者任何其他的标记语言。

添加相关的信息，比如，如果是性能相关的改动就需要添加对应的压测数据。照惯例会使用 benchstat 工具来对压测数据进行格式化处理，以便写入变更信息里。

引用 issues
接下来那个特殊的表示法 "Fixes #12345" 把代码变更关联到了 Go issue tracker 列表里的 issue 12345。当这个代码变更最终实施之后 (也就是合入主干)，issue tracker 将会自动标记那个 issue 为"已解决"并关闭它。

如果这个代码变更只是部分解决了这个 issue 的话，请使用 "Updates #12345"，这样的话就会在那个 issue 的评论区里留下一个评论把它链接回 Gerrit 上的变更页面，但是在该代码变更被实施之后并不会关闭掉 issue。

如果你是针对一个子仓库发送的代码变更，你必须使用 GitHub 支持的完全形式的语法来确保这个代码变更是链接到主仓库的 issue 上去的，而非子仓库。主仓库的 issue tracker 会追踪所有的 issues，正确的格式是 "Fixes golang/go#159"。

代码评审流程
这个章节是对代码评审流程的详细介绍以及如何在一个变更被发送之后处理反馈。

常见的新手错误
当一个变更被发送到 Gerrit 之后，通常来说它会在几天内被分门别类。一个维护者将会查看并提供一些初始的评审，对于初次提交代码贡献者来说，这些评审通常集中在基本的修饰和常见的错误上。

内容包括诸如：

Commit 信息没有遵循建议的格式
没有链接到对应的 GitHub issue。大部分代码变更需要链接到对应的 GitHub issue，说明这次变更修复的 bug 或者实现的功能特性，而且在开始这个变更之前，issue 里应该已经达成了一致的意见。Gerrit 评审不会讨论代码变更的价值，仅仅是讨论它的具体实现。
变更如果是在开发周期的冻结阶段被发送到 Gerrit 上的，也就是说彼时 Go 代码树是不接受一般的变更的，这种情况下，一个维护者可能会在评审代码时留下一行这样的评论：R=go.1.12，意思是这个代码变更将会在下一个开发窗口期打开 Go 代码树的时候再进行评审。如果你知道那不是这个代码变更应该被评审的正确的时间范围，你可以自己加上这样的评论：R=go1.XX 来更正。
Trybots
在第一次看过你的代码变更之后，维护者会启动一些 trybots，这是一个会在不同的 CPU 架构的机器上运行完整测试套件的服务器集群。大部分 trybots 会在几分钟内执行完成，之后会有一个可以查看具体结果的链接出现在 Gerrit 变更页面上。

如果 trybot 最后执行失败了，点击链接然后查看完整的日志，看看是在哪个平台上测试失败了。尽量尝试去弄明白失败的原因，然后更新你的代码去修复它，最后重新上传你的新代码。维护者会重新启动一个新的 trybot 再跑一遍，看看问题是不是已经解决了。

有时候，Go 代码树会在某些平台上有长达数小时的执行失败；如果 trybot 上报的失败的问题看起来和你的这次代码变更无关的话，到构建面板上去查看近期内的其他 commits 在相同的平台上是不是有出现过这种一样的失败。如果有的话，你就在 Gerrit 变更页面的评论区里说明一下这个失败和你的代码变更无关，以此让维护者知悉这种情况。

评审
Go 语言社区非常重视全面的评审。你要把每一条评审的评论的当成一张罚单：你必须通过某种方式把它"关掉"，或者是你把评论里建议的修改实现一下，或者是你说服维护者那部分不需要修改。

在你更新了你的代码之后，过一遍评审页面的所有评论，确保你已经全部回复了。你可以点击 "Done" 按钮回复，这表示你已经实现了评审人建议的修改，否则的话，点击 "Reply" 按钮然后解释一下你为什么还没修改、或者是你已经做了其他地方的修改并覆盖了这一部分。

一般来说，代码评审里会经历多轮的评审，期间会有一个或者多个评审人不断地发表新的代码审查评论然后等待提交者修改更新代码之后继续评审，这是很正常的。甚至一些经验老到的代码贡献者也会经历这种循环，所以不要因此而被打击到。

投票规则
在评审人们差不多要得出结论之时，他们会对你的此次代码变更进行"投票"。Gerrit 的投票系统包含了一个在 [-2, 2] 区间的整数：

+2: 同意此次代码变更被合入到主分支。只有 Go 语言的维护者们才有权限投 +2 的票。
+1: 这个代码变更看起来没什么问题，不过要么是因为评审人还要求对代码做一些小的改动、要么是因为该评审人不是一个维护者而无法直接批准这个变更，但是该评审人支持批准这个变更。
-1: 这个代码变更并不是很合理但可能有机会做进一步的修改。如果你得到了一个 -1 票，那一定会有一个明确的解释告诉你为什么。
-2: 一个维护者否决了这个代码变更并且不同意合入主干。同样的，会有一个明确的解释来说明原因。
提交一个核准的变更
在一个代码变更被投了一个 +2 票之后，投下这票的核准人将会使用 Gerrit 的用户界面来将代码合并入主干，这个操作被称为"提交变更"。

之所以把核准和提交拆分成两步，是因为有些时候维护者们可能并不想把刚刚批准的代码变更立刻合入主干，比如，彼时可能正处于 Go 代码树的暂时冻结期。

提交一个变更将会把代码合入主仓库，代码变更的描述信息里会包含一个指向对应代码评审页面的链接，而具体代码评审页面处也会更新一个链接指向仓库里的此次代码变更 commit。把代码变更合入主干时使用的是 Git 的 "Cherry Pick" 命令，因此在主仓库里的关于此次代码变更的 commit 哈希 ID 会被这个提交操作更改。

如果你的变更已经被批准了好几天了，但是一直没有被提交到主仓库，你可以在 Gerrit 写个评论要求合入。

更多信息
除了这里的信息，Go 语言社区还维护了一个代码评审的 wiki 页面。随时欢迎你在学习相关的评审流程之时为这个页面贡献、补充新内容。

其他主题
这个章节收集了一些除了 issue/edit/code review/submit 流程之外的注解信息。

版权标头
Go 语言仓库里的文件不会保存一份作者列表，既是为了避免杂乱也是为了避免需要实时更新这份列表。相反的，你的名字将会出现在变更日志和贡献者文件里，也可能会出现在作者文件里。这些文件是定期从 commit 日志上自动生成的。作者文件定义了哪些人是 “Go 语言作者” - 版权持有者。

如果你在提交变更的时候有新添加的文件，那么应该使用标准的版权头：

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
(如果你此刻是在 2021 年或者往后的时间阅读这份文档，请使用你当前的年份。) 仓库里的文件版权生效于被添加进去的当年，不要在你变更的文件里更改版权信息里的年份。

mail 命令错误大全和故障排除
git codereview mail 命令失败的最常见原因是因为你的邮件地址和你在注册流程中使用的邮件地址不匹配。

如果你看到这样的输出信息：

remote: Processing changes: refs: 1, done
remote:
remote: ERROR:  In commit ab13517fa29487dcf8b0d48916c51639426c5ee9
remote: ERROR:  author email address XXXXXXXXXXXXXXXXXXX
remote: ERROR:  does not match your user account.
你需要在这个仓库下把 Git 用户邮箱配置为你一开始注册好的那个邮箱。更正邮箱地址以确保不会再发生这个错误：

$ git config user.email email@address.com
然后通过以下命令修改你的 commit 信息，更正里面的用户名和邮箱：

$ git commit --amend --author="Author Name <email@address.com>"
最后运行一下的命令再重试一次：

$ git codereview mail
快速测试你的代码变更
如果每一次单独的代码变更都对整个代码树运行 all.bash 脚本的话太费劲了，尽管我们极力建议你在发送代码变更之前跑一下这个脚本，然而在开发的期间你可能只想要编译和测试那些你涉及到的包。

通常来说，你可以运行 make.bash 而不是 all.bash 来只构建 Go 工具链，而不需要运行整个测试套件。或者你可以运行 run.bash 来运行整个测试套件而不构建 Go 工具链。你可以把 all.bash 看成是依次执行 make.bash 和 run.bash 。
在这个章节，我们会把你存放 Go 语言仓库的目录称为 $GODIR 。 make.bash 脚本构建的 go 工具会被安装到 $GODIR/bin/go 然后你就可以调用它来测试你的代码了。例如，如果你修改了编译器而且你想要测试看看会对你自己项目里的测试套件造成怎样的影响，直接用它运行 go test ： powershell $ cd <MYPROJECTDIR> $ $GODIR/bin/go test
如果你正在修改标准库，你可能不需要重新构建编译器：你可以直接在你正在修改的包里跑一下测试代码就可以了。你可以使用平时用的 Go 版本或者从克隆下来的源码构建而成的编译器 (有时候这个是必须的因为你正在修改的标准库代码可能会需要一个比你已经安装的稳定版更新版本的编译器) 来做这件事。 powershell $ cd $GODIR/src/hash/sha1 $ [make changes...] $ $GODIR/bin/go test .
如果你正在修改编译器本身，你可以直接重新编译 编译 工具（这是一个使用 go build 命令编译每一个单独的包之时会调用到的一个内部的二进制文件）。完成之后，你会想要编译或者运行一些代码来测试一下： powershell $ cd $GODIR/src $ [make changes...] $ $GODIR/bin/go install cmd/compile $ $GODIR/bin/go build [something...] # test the new compiler $ $GODIR/bin/go run [something...] # test the new compiler $ $GODIR/bin/go test [something...] # test the new compiler
同样的操作可以应用到 Go 工具链里的其他内部工具，像是 asm ， cover ， link 等等。直接重新编译然后使用 go install cmd/<TOOL> 命令安装，最后使用构建出来的 Go 二进制文件测试一下。

除了标准的逐包测试，在 $GODIR/test 目录下有一个顶级的测试套件，里面包含了多种黑盒和回归测试。这个测试套件是包含在 all.bash 脚本里运行的，不过你也可以手动运行它： powershell $ cd $GODIR/test $ $GODIR/bin/go run run.go
向子仓库提交贡献 (golang.org/x/...)
如果你正在向一个子仓库提交贡献，你需要使用 go get 来获取对应的 Go 包。例如，如果要向 golang.org/x/oauth2 包贡献代码，你可以通过运行以下的命令来获取代码：

$ go get -d golang.org/x/oauth2/...
紧接着，进入到包的源目录（$GOPATH/src/golang.org/x/oauth2），然后按照正常的代码贡献流程走就行了。

指定一个评审人/抄送其他人
除非有明确的说明，比如在你提交代码变更之前的讨论中，否则的话最好不要自己指定评审人。所有的代码变更都会自动抄送给 golang-codereviews@googlegroups.com 邮件组。如果这是你的第一次提交代码变更，在它出现在邮件列表之前可能会有一个审核延迟，主要是为了过滤垃圾邮件。

你可以指定一个评审人或者使用 -r/-cc 选项抄送有关各方。这两种方式都接受逗号分隔的邮件地址列表：

$ git codereview mail -r joe@golang.org -cc mabel@example.com,math-nuts@swtch.com
同步你的客户端
在你做代码变更期间，可能有其他人的变更已经先你一步被提交到主仓库里，那么为了保持你的本地分支更新，运行：

git codereview sync
（这个命令背后运行的是 git pull -r .）

其他人评审代码
评审人作为评审流程的一部分可以直接提交代码到你的变更里（就像是在 GitHub 工作流里有其他人把 commits 依附到你的 PR 上了）。你可以导入这些他人提交的变更到你的本地 Git 分支上。在 Gerrit 的评审页面，点击右上角的 "Download ▼" 链接，复制 "Checkout" 命令然后在你的本地 Git 仓库下运行它。这个命令类似如下的格式：

$ git fetch https://go.googlesource.com/review refs/changes/21/13245/1 && git checkout FETCH_HEAD
如果要撤销，切换回你之前在开发的那个分支即可。

设置 Git 别名
git codereview 相关的命令可以直接在终端键入对应的选项运行，例如：

$ git codereview sync
不过给 git codereview 子命令命令设置别名会更方便使用，上面的命令可以替换成：

$ git sync
git codereview 的子命令的名字是排除了 Git 本身的命令关键字而挑选出来的，所以不用担心设置了这些别名会和 Git 本身的命令冲突。要设置这些别名，复制下面的文本到你的 Git 配置文件里（通常是在 home 路径下的 .gitconfig 文件）：

[alias]
    change = codereview change
    gofmt = codereview gofmt
    mail = codereview mail
    pending = codereview pending
    submit = codereview submit
    sync = codereview sync
发送多个依赖的变更
老司机用户可能会想要把相关的 commits 叠加到一个单独的分支上。Gerrit 允许多个代码变更之间相互依赖，形成这样的依赖链。每一个变更需要被单独地核准和提交，但是依赖对于评审人来说是可见的。

要发送一组依赖的代码更改，请将每个变更作为不同的 commit 保存在同一分支下，然后运行：

$ git codereview mail HEAD
要确保显示地指定 HEAD ，不过这在单个变更的场景里通常是不需要指定的。

英文原文地址
https://golang.org/doc/contribute.html