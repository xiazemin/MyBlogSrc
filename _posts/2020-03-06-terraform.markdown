---
title: terraform
layout: post
category: golang
author: 夏泽民
---
https://github.com/hashicorp/terraform
在 DevOps 实践中，基础设施即代码如何落地是一个绕不开的话题。像 Chef，Puppet 等成熟的配置管理工具，都能够满足一定程度的需求，但有没有更友好的工具能够满足我们绝大多数的需求？笔者认为 Terraform 是一个很有潜力的工具，目前各大云平台也都支持的不错，尤其是使用起来简单明了。本文会简单的介绍一下 Terraform 相关的概念，然后通过一个小 demo 带大家一起进入 Terraform 的世界
<!-- more -->
Terraform 是一种安全有效地构建、更改和版本控制基础设施的工具(基础架构自动化的编排工具)。它的目标是 "Write, Plan, and create Infrastructure as Code", 基础架构即代码。Terraform 几乎可以支持所有市面上能见到的云服务。具体的说就是可以用代码来管理维护 IT 资源，把之前需要手动操作的一部分任务通过程序来自动化的完成，这样的做的结果非常明显：高效、不易出错。



Terraform 提供了对资源和提供者的灵活抽象。该模型允许表示从物理硬件、虚拟机和容器到电子邮件和 DNS 提供者的所有内容。由于这种灵活性，Terraform 可以用来解决许多不同的问题。这意味着有许多现有的工具与Terraform 的功能重叠。但是需要注意的是，Terraform 与其他系统并不相互排斥。它可以用于管理小到单个应用程序或达到整个数据中心的不同对象。

Terraform 使用配置文件描述管理的组件(小到单个应用程序，达到整个数据中心)。Terraform 生成一个执行计划，描述它将做什么来达到所需的状态，然后执行它来构建所描述的基础结构。随着配置的变化，Terraform 能够确定发生了什么变化，并创建可应用的增量执行计划。

Terraform 是用 Go 语言开发的开源项目，你可以在 github 上访问到它的源代码。

Terraform 核心功能
基础架构即代码(Infrastructure as Code)
执行计划(Execution Plans)
资源图(Resource Graph)
自动化变更(Change Automation)
基础架构即代码(Infrastructure as Code)
使用高级配置语法来描述基础架构，这样就可以对数据中心的蓝图进行版本控制，就像对待其他代码一样对待它。

执行计划(Execution Plans)
Terraform 有一个 plan 步骤，它生成一个执行计划。执行计划显示了当执行 apply 命令时 Terraform 将做什么。通过 plan 进行提前检查，可以使 Terraform 操作真正的基础结构时避免意外。

资源图(Resource Graph)
Terraform 构建的所有资源的图表，它能够并行地创建和修改任何没有相互依赖的资源。因此，Terraform 可以高效地构建基础设施，操作人员也可以通过图表深入地解其基础设施中的依赖关系。

自动化变更(Change Automation)
把复杂的变更集应用到基础设施中，而无需人工交互。通过前面提到的执行计划和资源图，我们可以确切地知道 Terraform 将会改变什么，以什么顺序改变，从而避免许多可能的人为错误。

安装 Terraform
Terraform 的安装非常简单，直接把官方提供的二进制可执行文件保存到本地就可以了。比如笔者习惯性的把它保存到 /usr/local/bin/ 目录下，当然这个目录会被添加到 PATH 环境变量中。完成后检查一下版本号：



通过 -h 选项我们可以看到 terraform 支持的所有命令：



在 Azure 上创建一个 Resource Group
要让 Terraform 访问 Azure 订阅中的资源，需要先创建 Azure service principal，Azure service principa 允许你的 Terraform 脚本在 Azure 订阅中配置资源。请参考这里创建 Azure service principal。

配置 Terraform 环境变量
若要配置 Terraform 使用 Azure service principal，需要设置以下环境变量：

ARM_SUBSCRIPTION_ID
ARM_CLIENT_ID
ARM_CLIENT_SECRET
ARM_TENANT_ID
ARM_ENVIRONMENT
这些环境变量的值都可以从前面创建 Azure service principal 的过程中获得。方便起见，我们把设置这些环境变量的步骤可以写到脚本文件 azureEnv.sh 中：

复制代码
#!/bin/sh
echo "Setting environment variables for Terraform"
export ARM_SUBSCRIPTION_ID=your_subscription_id
export ARM_CLIENT_ID=your_appId
export ARM_CLIENT_SECRET=your_password
export ARM_TENANT_ID=your_tenant_id
# Not needed for public, required for usgovernment, german, china
export ARM_ENVIRONMENT=public
复制代码
这样在执行 Terraform 命令前通过 source 命令执行该脚本就可以了！

创建 Terraform 配置文件
为了在 Azure 上创建一个 Resource Group，我们创建名称为 createrg.tf 的配置文件，并编辑内容如下：

provider "azurerm" {
}
resource "azurerm_resource_group" "rg" {
        name = "NickResourceGroup"
        location = "eastasia"
}
用 init 命令用来初始化工作目录
把当前目录切换到 createrg.tf 文件所在的目录，然后执行 init 命令：

$ terraform init 


其实就是把 createrg.tf 文件中指定的驱动程序安装到当前目录下的 .terraform 目录中：



通过 plan 命令检查配置文件
plan 命令会检查配置文件并生成执行计划，如果发现配置文件中有错误会直接报错：

$ . azureEnv.sh
$ terraform plan


通过 plan 命令的输出，我们可以清楚的看到即将在目标环境中执行的任务。

使用 graph 命令生成可视化的图表
其实 graph 命令只能生成相关图表的数据(dot 格式的数据)，我们通过 dot 命令来生成可视化的图表，先通过下面的命令安装 dot 程序：

$ sudo apt install graphviz
然后生成一个图表：

$ terraform graph | dot -Tsvg > graph.svg


上图描述了我们通过 azurerm 驱动创建了一个 Resource Group。

使用 apply 命令完成部署操作
在使用 apply 命令执行实际的部署时，默认会先执行 plan 命令并进入交互模式等待用户确认操作，我们已经执行过 plan 命令了，所以可以使用 -auto-approve 选项跳过这些步骤直接执行部署操作：

$ terraform apply -auto-approve


到 Azure 站点上检查一下，发现名称为 NickResourceGroup 的 Resource Group 已经创建成功了。

总结
Terraform 支持的平台非常多，像 AWS，Azure 等大厂自然是不用说了，一些小的厂商也可以通过提供 provider 支持 Terraform，从而让整个生态变得非常活跃。如果大家想在 DevOps 实践中引入基础设施即代码，无论是面对的是公有云还是私有云，相信 Terraform 都不会让你失望。

欢迎访问Terraform介绍指南！本指南是开始学习Terraform的最佳之处。其包含Terraform是什么，解决什么问题以及与当前已有的软件对比，并且包含使用Terraform的快速入门！

如果你已经对Terraform基础很熟悉，参考文档为所有可用功能及内部组件提供了更好的参考指南。

Terraform是什么
Terraform是一个构建、变更、和安全有效的版本化管理基础设施的工具。Terraform可以管理已存在和流行的服务提供商以及定制的内部解决方案。

配置文件为Terraform描述运行单个应用程序或你整个数据中心所需的组件。Terraform生成一个执行计划描述了它将做什么以达到预期状态，然后执行它来构建所描述的基础设施。随着配置文件的变更，Terraform可以确定有什么变更，并且创建额外可应用的执行计划。

Terraform可管理的基础设施不仅包含计算实例，存储，网络等底层组件，也包含DNS条目，SaaS服务等高级组件。

最好的Terraform工作实例，请查看用例.

Terraform的主要功能如下：

基础设施即代码
基础设施使用高级配置语法进行描述。这可以让你的数据中心蓝图像你其他代码一样进行 版本控制和管理。此外基础设施可以被 分享和重用。

执行计划
Terraform在“计划”阶段生成执行计划。执行计划展示了当你调用apply时，Terraform将做什么。这在你使用Terraform操作基础设施时避免出现任何意外。

资源图表
Terraform构建所有资源的图表，并且并行创建和修改任何无依赖的资源。因此，Terraform尽可能高效的构建基础设施，并且操作者清楚其基础设施间的依赖关系。

自动变更
复杂的变更可以在最少的人工干预下应用到你的基础设施。使用前面提到的执行计划和资源图表，你可以确切的知道Terraform将会做那些变更，以及按什么顺序，避免一些可能的人为错误。

下一步
查看Terraform用例页面，了解Terraform的多种使用方式。然后查看Terraform如何与其他软件对比了解它如何适应你现有的基础设施。最后，继续阅读入门指南来使用Terraform管理真实的基础设施并了解它如何工作。

https://helpcdn.aliyun.com/product/95817.html

Terraform与Kubernetes
看到Terraform可以替代kubectl管理k8s资源的生命周期，于是调研了下它的使用场景，并对比Terraform和Helm的区别

一.Terraform介绍
Terraform是一款开源工具，出自HashiCorp公司，著名的Vagrant、Consul也出自于该公司。其主要作用是：让用户更轻松地管理、配置任何基础架构，管理公有和私有云服务，也可以管理外部服务，如GitHub，Nomad。

区别于ansible和puppet等传统的配置管理工具，Terraform趋向于更上层的一个组装者。

Terraform使用模板来定义基础设施，通过指令来实现资源创建/更新/销毁的全生命周期管理，实现“基础设施即代码”，具体示例如下：

resource "alicloud_instance" "web" {
    # cn-beijing
    availability_zone = "cn-beijing-b"
    image_id = "ubuntu_140405_32_40G_cloudinit_20161115.vhd"

    system_disk_category = "cloud_ssd"

    instance_type = "ecs.n1.small"
    internet_charge_type = "PayByBandwidth"
    security_groups = ["${alicloud_security_group.tf_test_foo.id}"]
    instance_name = "test_foo"
    io_optimized = "optimized"
}
1
2
3
4
5
6
7
8
9
10
11
12
13
这是阿里云的一个Terraform逻辑，执行terraform apply，就可以创建一个ECS实例

Terraform AliCloud provider: terraform-provider

二.Terraform支持K8S
17年9月，Terraform官方宣布支持Kubernetes，提供Kubernetes应用程序的完整生命周期管理，包含Pod的创建、删除以及副本控制等功能（通过调用API)。

以下是操作示例：

1.安装kubernete集群
当前k8s的installer列表，已经很多了…



使用Terraform在阿里云上安装k8s集群：kubernetes-examples

2.创建应用：
1.初始化k8s-provider
因为是调用apiserver，所以需要指定k8s集群的连接方式

provider "kubernetes" {} // 默认~/.kube/config
或：
provider "kubernetes" {
  host = "https://104.196.242.174"

  client_certificate     = "${file("~/.kube/client-cert.pem")}"
  client_key             = "${file("~/.kube/client-key.pem")}"
  cluster_ca_certificate = "${file("~/.kube/cluster-ca-cert.pem")}"
}
1
2
3
4
5
6
7
8
9
10
11
2.初始化terraform
$ terraform init

Initializing provider plugins...
- Downloading plugin for provider "kubernetes"...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
1
2
3
4
5
6
7
8
9
10
11
12
13
14
3.创建RC
// Terraform不支持Deployment
// issue:https://github.com/terraform-providers/terraform-provider-kubernetes/issues/3

resource "kubernetes_replication_controller" "nginx" {
  metadata {
    name = "scalable-nginx-example"
    labels {
      App = "ScalableNginxExample"
    }
  }

  spec {
    replicas = 2
    selector {
      App = "ScalableNginxExample"
    }
    template {
      container {
        image = "nginx:1.7.8"
        name  = "example"

        port {
          container_port = 80
        }

        resources {
          limits {
            cpu    = "0.5"
            memory = "512Mi"
          }
          requests {
            cpu    = "250m"
            memory = "50Mi"
          }
        }
      }
    }
  }
}

1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
26
27
28
29
30
31
32
33
34
35
36
37
38
39
40
4.创建service
resource "kubernetes_service" "nginx" {
  metadata {
    name = "nginx-example"
  }
  spec {
    selector {
      App = "${kubernetes_replication_controller.nginx.metadata.0.labels.App}"
    }
    port {
      port = 80
      target_port = 80
    }

    type = "LoadBalancer"
  }
}
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
4.查看和执行
以上的步骤均为执行计划的定义
执行操作：terraform apply
查看当前执行几乎：terraform plan
1
2
3
三.为什么使用Terraform
1.如果你的基础设施（虚机、BLB等）是用Terraform来管理，那么你无需任何成本，可以用同样的配置语言，来管理k8s集群
2.完整的生命周期管理
3.每个执行的同步反馈
4.关系图谱：比如PVC和PV，如果PV创建失败，则不会去执行PVC的流程
四.与Helm的对比
如果是对K8S做上层的资源管理，大多数人会想到用Helm:参考

以下是Helm与Terraform都创建一个APP的操作对比：https://dzone.com/articles/terraform-vs-helm-for-kubernetes

Terraform的优势：

1.如果你的基础设施已经用了Terraform，那么k8s集群管理也可以直接用这个，没有学习成本
2.Terraform不需要在k8s集群中安装任何组件，它通过APISERVER管理资源
Terraform的缺点：

1.对K8S的支持还比较弱，而且17年9月才开始支持，项目还在初期
2.严重依赖Terraform的基础设施provider，比如外部磁盘、consul等没有支持的话，k8s中无法使用
2.不支持beta资源，这个是硬伤，如：Deployment/StatefulSet/Daemonset不支持
3.没有生态和市场的概念，比如helm中的仓库，共享大家的应用仓库
五.吐槽
对于Terraform，不支持Deployment这一条，足以让很多人放弃这个方案，而issue中对于这个的讨论，也有点不太乐观



必须在v1中的资源才会支持。对于Deployment大家只能用RC代替、或者kube exec加进去(尬

但对于kubernetes而言，beta阶段的很多资源，已经被大家广泛使用(Deployment、Daemonset)，而且新版本的Deployment已经变成了apps/v1。

k8s各种版本(v1、apps/v1)的区别：参考文章

Terraform是国际著名的开源的资源编排工具，据不完全统计，全球已有超过一百家云厂商及服务提供商支持Terraform。Terraform是HashiCorp的代码软件基础设施。它允许用户使用高级配置语言定义数据中心基础架构，从中可以创建执行计划以构建OpenStack等基础架构，或者在IBM Cloud，AWS，Microsoft Azure，Google Cloud Platform等多种云服务中构建基础架构。

2015年，当Terraform第一次进入我的视野，它就像是我的瓦尔哈拉（北欧神话中死亡之神奥丁款待阵亡将士英灵的殿堂）。Terraform旨在解决复杂基础设施的配置问题——将多个云供应商汇聚在一起——从AWS这样的巨头到Logentries这样的单一解决方案供应商。

我们的团队认为我们需要一些解决方案来处理基础设施的复杂性问题。对于一个基于Heroku和AWS的平台来说，水平伸缩到四个Terraform似乎是一个完美的解决方案。我们希望有一些东西能够让我们意识到基础设施即代码这个概念的重要性，因为对于一个DevOps团队来说，这是必须的。Terraform功能齐全，但也并非十全十美，仍然有一些问题值得我们注意。

我将列举一些比较大的“坑”，并分享我们是如何填这些坑的。最后，我将会说明，尽管存在这些挑战，但Terraform在工具领域仍然有很大的发展空间。

1. 邪恶的状态
首先，Terraform是有状态的，我个人认为这会带来两个问题：

状态必须始终与基础设施保持同步——这意味着在配置时需要全盘考虑，不能在配置工具之外进行栈的修改。

你必须在某个地方维护状态——必须是一个安全的地方，因为状态可能会带有秘钥之类的东西。

但Terraform引入状态是有原因的，状态用来维护文件中定义的资源与在云供应商平台上创建的实际资源之间的映射关系。有了这个，Terraform就可以为我们提供一些好处：

从云供应商那里处读取状态（状态同步，也称为刷新）可能非常耗时。如果我们可以100％确定状态是准确的，我们完全可以不进行同步，并立即应用变更。

能够跟踪已创建的资源，可以更轻松地进行重命名和修改结构——这些都是基本的基础设施操作。

Terraform在应用变更之前会锁定状态，这意味着我们可以确保在应用变更时，没有其他人在做同样的操作。

我认为，在考虑配置工具时，你应该权衡一下上述的几个问题，并搞清楚你的栈是像工作簿之类的东西，每次做出变更后都可以进行重建，还是像是一个生物有机体，需要在它运行过程中做出变更。

2. 难以集成已有的栈
在Terraform的早期，很多人抱怨无法将Terraform用在已有的栈上。其原因在于，Terraform无法将已有栈纳入到它的状态管理中。所幸的是，Terraform通过引入import命令解决了这个问题（至少在系统级别）。

但又出现了另一个与此紧密相关的问题——如果你的栈很大，就必须为每个资源多次使用terraform import命令。如果不使用一些自动化脚本，这就变成了一个非常耗时且令人感到沮丧的活儿。我们完全可以通过更好的方式来导入这些东西，不过这要求Terraform将资源视为树，而不是扁平的结构。在某些情况下，它是完全合理的——比如heroku_app和heroku_domain（或heroku_drain）。当然，肯定还有很大的改进空间。

3. 复杂的状态修改
在重构基础设施定义时，到最后可能就是重命名资源（修改它们的标识符）或将资源移动到模块中。遗憾的是，Terraform很难跟踪这些变化，并将它们置于一种未知的错配状态。如果再次运行apply，会重新创建资源，但这可能不是你想要的。好在Terraform提供了terraform state mv命令，可用于移动逻辑资源。不好的地方在于，在大多数情况下，你需要大量使用这个命令。

4. 古怪的条件逻辑
Terraform并非真正的命令式编程语言，有些人不喜欢这一点​​。但说实话，我并不这么认为——我认为栈的定义应该尽可能是声明性的——这样可以减少定义之间的偏差。另一方面，Terraform提供的条件逻辑有点古怪。例如，如果要定义条件资源，需要将资源定义为列表，并使用count参数来控制它：

resource "heroku_app" "some_app" {
    count = "${var.create_app}" 
    name = "some-app"
    ...
}
这已经相当具体了，而且你不需要知道if/else是怎么回事。

resource "heroku_app" "some_app" {
    count = "${var.create_app}" 
    name = "some-app"
    ...
}

resource "heroku_app" "some_app" {
    count = "${1 - var.create_app}" 
    name = "some-app"
    ...
}
不过需要注意的是，如果有可能，应该尽量避免使用这样的结构。当然，这不应该成为拒绝使用Terraform的理由，只需要知道存在这个问题即可。

在近期发布的一些版本中，可以使用resource for_each来简化这个问题。

5. 无法简单地遍历模块
模块的概念是非常棒的——它将可重用的资源集包含在可重用的工件中。我们来看一些经过简化的例子：

app/app.tf

resource "heroku_app" "app" {
  name = "${var.name}"
  region = "eu"
  organization = {
    "name" = "${var.organization}"
  }
  config_vars = ["${var.configuration}"]
}

resource "heroku_addon" "deploy_hook" {
  app = "${heroku_app.app.name}"
  plan = "deployhooks:http"
  config = {
    url = "${var.deployhook_url}"
  }
}

output "name" {
  value = "${heroku_app.app.name}"
}
可用在服务声明中：

stack/some_service.tf

module "app" {
  source = "../app"
  name = "some-service"
  configuration = {
    SOME_ENV_VARIABLE = "value"
  }
  pipeline_id = "000000-9207-490a-b050-617b01ef79f3"
  deployhook_url = "https://some.deploy.hook"
}
这对我们来说是一个重大的变化，因为我们把很多重复的资源附加到每个应用程序——监控、日志收集、部署钩子等等。但这也带来了一个问题——由于某些原因，相同的工件有时候并不代表相同的资源。也就是说，它们不支持计数参数，但计算参数在应用条件逻辑或按照每个克隆进行服务迭代时却是非常关键的。确切地说，我们不应该这样做：

stack/services.tf（这不是真的）

module "app" {
    source = "../app"
    name = "some-service-${var.clone[count.index]}"
    configuration = "${var.configuration[count.index]}"
    ...
}
我们必须重复每个克隆的定义：

stack/services.tf

module "clone1_app" {
    source = "../app"
    name = "some-service-clone1"
    configuration = "${var.clone1Configuration}"
}

module "clone2_app" {
    source = "../app"
    name = "some-service-clone2"
    configuration = "${var.clone2Configuration}"
}
这个问题也有望在可预见的未来得到解决。

6. 不稳定的资源
作为一款功能丰富的0.x版本软件，Terraform存在一大堆小问题。其中一个问题是，有些资源并不会保持稳定的状态。对我们来说，问题是SNS主题订阅策略——我们针对服务所做的所有事情，只要包含订阅了SNS的队列，Terraform就会去修改它们（无论是否有意义）。这样可能会导致混淆，尤其是对于第一次接触Terraform的人来说。虽然这个问题是供应商局部问题，并且很可能会得到修复，但在修复之前总归是要注意的。

7. 一些小细节
我们遇到的另一个小问题是无法使用依赖于某些状态的计数值（在模块中），例如：

resource "heroku_app" "sample" {
   count = "${lookup(var.some_map, "enable_sample_app", 0)}" 
   name = "sample-app"
   ...
}
如果在模块中定义了上述的内容，你将会得到一条这样的消息：无法计算count的值…这真的很烦人，特别是当你看到Hashicorp的解释，说可以使用-target来逐个初始化资源。

8. 如何处理秘钥？
Terraform文件之所以不容易保存，其中一个原因是不知道该把秘钥保存在哪里。有几种方法可以解决这个问题：

Hashicorp使用了他们的Vault——虽然这是可行的，但它使整个设置变得更加复杂，感觉有点矫枉过正。
与Vault类似，你可以使用AWS的KMS来存储秘钥——但它同样有复杂性方面的问题。
使用私有git存储库，只要电脑不被偷，就不会有什么问题。
如果你是在CD/CI服务器上运行作业，可以将秘钥保存在环境变量中，但在一个复杂的系统中，可能就难以维护太多的环境变量。
你可以将秘钥保存在本地，并使用专门用于配置的机器，但对于有一定规模的团队来说，这样做是不行的。
我们的方法是在.tf文件保存所有秘钥变量，然后使用git-secret加密。运行terraform plan和terraform apply的脚本先用git-secret显示秘钥，然后再立即隐藏。虽然这不是一个完美的解决方案，但至少足够简单，减少从本地运行Terraform可能会造成的秘钥泄露几率。
9. 托管服务
最初，Hashicorp和其他公司都没有提供Terraform托管服务。作为一款非常复杂的软件（特别是秘密管理），有必要提供一个利基，于是Terraform Enterprise出现了。可惜的是，我没有这方面的经验，所以不确定它究竟好不好。但它可能存在的一个问题的是，使用企业模式会导致供应商锁定。

那么我们应该怎么办？
很明显，在选择配置解决方案时，我们需要考虑到Terraform存在的一些问题。其中一些问题最终将得到解决，其他一些则是Hashicorp必须做出的架构选择。在文中最后，我应该给出我的观点——为什么要选择Terraform？在我看来，对于某些场景，你根本就没有其他选择，因为很难找到适合这么多云供应商的解决方案。此外，如果你的系统是一个实时系统，每天都有大量重复的基础设施小变更，那么Terraform绝对值得考虑。


