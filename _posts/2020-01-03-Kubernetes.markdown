---
title: Kubernetes
layout: post
category: algorithm
author: 夏泽民
---
近日，Kubernetes 1.13 正式发布，这是 2018 年发布的第四次也是最后一次大版本，该版本继续关注 Kubernetes 稳定性和可扩展性，对存储和集群生命周期的主要功能实现高可用。Kubeadm 简化了集群管理、容器存储接口（CSI）并将 CoreDNS 作为默认 DNS。

最近两年，Kubernetes 给容器战场带来了巨大冲击，Docker Swarm 也未能成为其对手，类似 AWS 的主流云供应商纷纷提供 K8s 支持。本文总结了 50 多种 Kubernetes 集群部署、监控、安全及测试等相关工具，大部分为开源项目，非常适合技术人员入门。
<!-- more -->
原生可视化与控制
1、Kubernetes Dashboard

Kubernetes Dashboard 是 Kubernetes 集群基于 Web 的通用 UI，使用本地仪表板对 K8s 集群进行故障排除和监控要容易得多，但需要在计算机和 Kubernetes API 服务器之间创建安全代理通道才能访问。原生 Kubernetes Dashboard 依赖 Heapster 数据收集器，因此也需安装在系统中。
链接:
https://github.com/kubernetes/dashboard#kubernetes-dashboard
成本: 免费

测试
2、Kube-monkey

Kube-monkey 遵循混沌工程原理，可随机删除 K8s pod 并检查服务是否具有故障恢复能力，并提高系统健康性。Kube-monkey 由 TOML 文件配置，可在其中指定要杀死的应用程序及恢复时间。
链接:  https://github.com/asobti/kube-monkey
成本: 免费

3、K8s-testsuite

K8s-testsuite 由两个 Helm 图组成，可用于网络带宽测试和单个 Kubernetes 集群负载测试。负载测试模拟带有 loadbots 的简单 Web 服务器，这些服务器作为基于 Vegeta 的 Kubernetes 微服务运行。网络测试在内部使用 iperf3 和 netperf-2.7.0 并运行三次，两组测试均会生成包含所有结果和指标的综合日志消息。
链接:  https://github.com/mrahbar/k8s-testsuite
成本: 免费

4、Test-infra

Test-infra 是 Kubernetes 测试和结果验证的工具集合，可显示历史记录、聚合故障及当前正在测试的内容。用户可通过创建测试作业增强 test-infra 套件。Test-infra 可使用 Kubetest 工具对不同提供商的完整 Kubernetes 生命周期模拟并进行端到端测试。
链接:  https://github.com/kubernetes/test-infra
成本: 免费

5、Sonobuoy

Sonobuoy 允许以可访问和非破坏方式运行一组测试了解当前 Kubernetes 集群状态。Sonobuoy 生成信息报告，其中包含有关集群性能的详细信息。Sonobuoy 支持 Kubernetes 1.8 及更高版本，Sonobuoy Scanner 是一个基于浏览器的工具，允许通过几次单击测试 Kubernetes 集群，但 CLI 版本有更多测试可用。
链接:  https://github.com/heptio/sonobuoy
成本: 免费

6、PowerfulSeal

PowerfulSeal 是一种类似 Kube-monkey 的工具，遵循混沌工程原理。PowerfulSeal 可杀死 pod 并从集群中删除或者添加 VM。与 Kube-monkey 相比，PowefulSeal 具有交互模式，允许手动中断特定集群组件。此外，PowefulSeal 除了 SSH 之外不需要外部依赖。
链接:  https://github.com/bloomberg/powerfulseal
成本: 免费

集群部署工具
7、Kubespray

Kubespray 为 Kubernetes 部署和配置提供一组 Ansible 角色。Kubespray 可使用 AWS，GCE，Azure，OpenStack 或裸机基础架构即服务（IaaS）平台。 Kubespray 是一个开源项目，具有开放的开发模型。对于已经了解 Ansible 的人来说，该工具是一个不错的选择，因为不需要使用其他工具进行配置和编排。
链接:  https://github.com/kubernetes-incubator/kubespray
成本: 免费

8、Minikube

Minikube 允许在本地安装和试用 Kubernetes，该工具是使用 Kubernetes 的良好起点，可在笔记本电脑的虚拟机（VM）中轻松启动单节点 Kubernetes 集群。Minikube 适用于 Windows、Linux 和 OSX。在短短 5 分钟内，用户就可以探索 Kubernetes 的特点，只需一个命令即可启动 Minikube 仪表板。
链接:  https://github.com/kubernetes/minikube
成本: 免费

9、Kubeadm

Kubeadm 是自 1.4 版以来的 Kubernetes 分发工具，该工具有助于在现有基础架构上引导最佳 Kubernetes 集群实践，但 Kubeadm 无法配置基础架构，其主要优点是能够在任何地方发布最小的可行 Kubernetes 集群。但是，附加组件和网络设置都不属于 Kubeadm 范围，因此需要手动或使用其他工具进行安装。
链接:  https://github.com/kubernetes/kubeadm
成本: 免费

10、Kops

Kops 可帮助用户从命令行创建、销毁、升级和维护生产级高可用 Kubernetes 集群。AWS 目前正式支持，GCE 处于测试支持状态，而 alpha vSphere 中的 VMware vSphere 以及其他平台均予以支持。Kops 可控制完整 Kubernetes 集群生命周期，从基础架构配置到集群删除。
链接:  https://github.com/kubernetes/kops
成本: 免费

11、Bootkube

CoreOS 提供自托管 Kubernetes 集群概念，自托管集群方法的中心是 Bootkube。Bootkube 可设置临时 Kubernetes 控制平面，该平面将一直运行，直到自托管控制平面能够处理请求。
链接:  https://github.com/kubernetes-incubator/bootkube
成本: 免费

12、Kubernetes on AWS (Kube-AWS)

Kube-AWS 是 CoreOS 提供的控制台工具，使用 AWS CloudFormation 部署功能齐全的 Kubernetes 集群。Kube-AWS 允许部署传统 Kubernetes 集群，并使用本机 AWS 功能（比如 ELB、S3 和 Auto Scaling 等）自动配置每个 K8s 服务。
链接:  https://github.com/kubernetes-incubator/kube-aws ]"> https://github.com/kubernetes-incubator/kube-aws ]
成本: 免费

13、SimpleKube

SimpleKube 是一个 bash 脚本，可在 Linux 服务器上部署单节点 Kubernetes 集群。虽然 Minikube 需要虚拟机管理程序（VirtualBox，KVM），但 SimpleKube 会将所有 K8s 二进制文件安装到服务器。Simplekube 在 Debian 8/9 和 Ubuntu 16.x / 17.x 上进行测试，这同样是入门 Kubernetes 的好工具。
链接:  https://github.com/valentin2105/Simplekube
成本: 免费

14、Juju

Juju 是 Canonical 提供的服务编排工具，可让用户远程操作云提供商解决方案。Juju 的工作抽象级别高于 Puppet、Ansible 和 Chef，并且管理服务而不是虚拟机。Canonical 努力在生产中提供称之为合适的“Kubernetes-core bundle”。Juju 可作为专用工具使用，具有控制台 UI 界面，也可作为 JaaS 服务。
链接:  https://jujucharms.com/
成本: 社区版免费，商业版每年 200 美元起

15、Conjure-up

Conjure-up 同样是 Canonical 的产品，允许使用简单命令部署 Kubernetes 在 Ubuntu 上的规范分布，支持 AWS、GCE、Azure、Joyent、OpenStack、VMware 和 localhost 部署。Juju、MAAS 和 LXD 是 Conjure-up 的基础。
链接:  https://conjure-up.io/
成本: 免费

监测工具
16、Kubebox

Kubebox 是 Kubernetes 集群的终端控制台，允许使用界面管理和监控集群实时状态。Kubebox 可显示 pod 资源使用情况，集群监视和容器日志等。此外，用户可轻松导航到所需的命名空间并执行到所需容器，以便快速排障或恢复。
链接:  https://github.com/astefanutti/kubebox
成本: 免费

17、Kubedash

Kubedash 为 Kubernetes 提供性能分析 UI，汇总不同来源的指标，并为管理员提供高级分析数据。Kubedash 使用 Heapster 作为数据源，默认情况下在所有 Kubernetes 集群中作为服务运行，为各个容器收集指标并分析。
链接:  https://github.com/kubernetes-retired/kubedash
成本: 免费

18、Kubernetes Operational View (Kube-ops-view)

Kube-ops-view 是一个用于多 K8s 集群的只读系统仪表板。使用 Kube-ops-view，用户可轻松在集群和监控节点之间导航，并监控 pod 健康状况。Kube-ops-view 可以动画 Kubernetes 进程，例如 pod 创建和终止，使用 Heapster 作为数据源。
链接:  https://github.com/hjacobs/kube-ops-view
成本: 免费

19、Kubetail

Kubetail 是一个小型 bash 脚本，允许将多个 pod 日志聚合到一个流中。最初的 Kubetail 没有过滤或突出显示功能，但 Github 上的额外 Kubetail 版本可使用多尾工具形成并执行日志着色。
链接: 
https://github.com/johanhaleby/kubetail https://github.com/aks/kubetail
成本: 免费

20、Kubewatch

Kubewatch 可将 K8s 活动发布到 Slack 应用。Kubewatch 作为 Kubernetes 集群内的 pod 运行，并监视系统中发生的变化，可通过编辑配置文件来指定要接收的通知。
链接:  https://github.com/bitnami-labs/kubewatch
成本: 免费

21、Weave Scope

Weave Scope 是 Docker 和 Kubernetes 集群的故障排除和监视工具，可以自动生成应用程序和基础架构拓扑，轻松识别应用程序性能瓶颈，可以将 Weave Scope 部署为本地服务器或笔记本电脑上的独立应用程序，也可以选择 Weave Cloud 上的 Weave Scope 软件即服务（SaaS）解决方案。使用 Weave Scope，用户可根据名称、标签或资源消耗轻松对容器分组、过滤或搜索。
链接:  https://www.weave.works/oss/scope/
成本: 标准模式提供 30 天免费体验，企业版每节点每月 150 美元

22、Searchlight

AppsCode 提供的 Searchlight 是满足 Icinga 的 Kubernetes 编排工具。Searchlight 会定期对 Kubernetes 集群执行检查，并在出现问题时通过电子邮件、短信或其他方式提醒。Searchlight 包含专门为 Kubernetes 编写的默认检查套件。此外，可通过外部黑匣子增强 Prometheus 监控，并在内部系统完全失效的情况下作为备份。
链接:  https://github.com/appscode/searchlight
成本: 免费

23、Heapster

Heapster 为 Kubernetes 提供容器集群监控和性能分析。Heapster 本身支持 Kubernetes，可在所有 K8s 设置上作为 pod 运行，可将 Heapster 数据推送到可配置的后端进行存储和可视化。
链接:  https://github.com/kubernetes/heapster
成本: 免费

安全
24、Trireme

Trireme 适用于所有 Kubernetes 集群，允许管理来自不同集群的 pod 间流量，主要优点是不需要任何集中策略管理，能够轻松组织部署在 Kubernetes 中的资源交互，并且没有 SDN、VLAN 标签和子网的复杂性（Trireme 使用传统 L3- 网络）。
链接:  https://github.com/aporeto-inc/trireme-kubernetes
成本: 免费

25、Aquasec

Aquasec 为 Kubernetes 部署提供完整的生命周期安全性。Aqua Security 在每个容器实例上部署专用代理，该实例用作防火墙并阻止容器中的安全漏洞，此代理与中央 Aqua Security 控制台进行通信，该控制台强制执行已定义的安全限制。Aqua Security 还有助于为云和内部部署环境组织灵活的安全交付管道。Kube-Bench 是 AquaSec 发布的开源工具，根据 CIS Kubernetes Benchmark 中的测试列表检查 Kubernetes 环境。
链接:  https://www.aquasec.com/
成本: 每次扫描 0.29 美元

26、Twistlock

Twistlock 可用作云原生应用程序防火墙，并分析容器和服务之间的网络流量。Twistlock 能够分析标准容器行为并据此生成适当规则，管理员不必手动生成。Twistlock 还支持 2.2 版本的 Kubernetes CIS 基准测试。
链接:  https://www.twistlock.com/
成本: 每年每个许可 1700 美元起 (可免费试用)

27、Sysdig Falco

Sysdig Falco 是一种行为活动监视器，旨在检测应用程序异常。Falco 基于 Sysdig 项目，这是一个开源工具（现在是商业项目），通过跟踪内核系统调用来监控容器性能。Falco 允许使用一组规则持续监视和检测容器、应用程序、主机和网络活动。
链接:  https://sysdig.com/opensource/falco/
成本: 独立工具可免费使用 
基于云: 每月 20 美元 per month (免费试用)
Pro Cloud: 每月 30 美元
Pro Software: 定制价格

28、Sysdig Secure

Sysdig Secure 是 Sysdig Container Intelligence Platform 的一部分，具有无与伦比的容器可见性并与容器编排工具深度集成，开箱即用，包括 Kubernetes、Docker、AWS ECS 和 Apache Mesos。Sysdig Secure 可实施服务感知策略，阻止攻击并分析历史记录及监控集群性能。
链接:  https://sysdig.com/product/secure/
成本: 工具免费 
Pro Cloud: 定制价格
Pro Software: 定制价格

29、 Kubesec.io

Kubesec.io 可为 Kubernetes 资源评分并提供安全功能，Kubesec.io 根据安全性最佳实践验证资源配置。因此，用户将获得有关提高整体系统安全性的全面控制和建议。该网站包含大量与容器和 Kubernetes 安全相关的链接。
成本: 免费

CLI 工具
30、Cabin

Cabin 用作移动仪表板，可远程管理 Kubernetes 集群。 借助 Cabin，用户可快速管理应用程序，扩展部署并通过 Android 或 iOS 设备对整个 K8s 集群进行故障排除。Cabin 是 K8s 集群运营商的理想工具，因为允许在发生事故时执行快速补救措施。
链接:  https://github.com/bitnami-labs/cabin
成本: 免费

31、Kubectx/Kubens

Kubectx 是一个小型开源实用工具，可以增强 Kubectl 功能，轻松切换上下文并同时连接到几个 Kubernetes 集群。Kubens 允许在 Kubernetes 命名空间之间导航，这两个工具在 bash/zsh/fish shell 上都有自动完成功能。
链接:  https://github.com/ahmetb/kubectx
成本: 免费

32、 Kube-shell

使用 kubectl 时，Kube-shell 可提高工作效率，其可通过命令自动完成部分工作。此外，Kube-shell 将提供有关已执行命令的在线文档。Kube-shell 甚至可以在错误输入时搜索和更正命令，是提高 K8s 控制台性能和工作效率的绝佳工具。
链接:  https://github.com/cloudnativelabs/kube-shell
成本: 免费

33、Kail

Kail 是 Kubernetes tail 的缩写，适用于 Kubernetes 集群。Kail 可为所有匹配 pod 添加 Docker 日志。Kail 允许按服务、部署、标签和其他功能过滤 pod。如果 Pod 符合条件，则会在启动后自动添加（或删除）到日志中。
链接:  https://github.com/boz/kail
成本: 免费

部署工具
34、Telepresence

Telepresence 通过 Kubernetes 环境的代理数据本地调试集群到本地进程的可能性。Telepresence 能够为本地代码提供对 Kubernetes 服务和 AWS/GCP 资源的访问，因为它将部署到集群。通过 Telepresence，Kubernetes 可将本地代码视为集群中的普通 pod。
链接:  https://www.telepresence.io/
成本: 免费

35、Helm

Helm 可与 Char 一起运行，Char 是构成分布式应用程序的 Kubernetes 资源归档集，可通过创建 Helm 图表共享应用程序，Helm 允许构建并轻松管理 Kubernetes 配置。
链接:  https://github.com/kubernetes/helm
成本: 免费

36、Keel

Keel 允许自动执行 Kubernetes 部署更新，并可在专用命名空间中作为 Kubernetes 服务启动。通过这种方式，Keel 为环境带来了最小负载，并显著增加稳健性。Keel 通过标签、注释和图表帮助部署 Kubernetes 服务，只需为每个部署或 Helm 版本指定更新策略。一旦新的应用程序版本在存储库中可用，Keel 将自动更新环境。
链接:  https://keel.sh/
成本: 免费

37、Apollo

Apollo 是一个开源应用程序，为团队提供自助 UI，用于创建和部署 Kubernetes 服务。Apollo 允许管理员单击一下即可查看日志并将部署恢复到任何时间点。Apollo 具有灵活的部署权限模型，每个用户只能部署需要的内容。
链接:  https://github.com/logzio/apollo
成本: 免费

38、Draft

Draft 是 Azure 团队提供的工具，可简化应用程序开发和部署到任何 Kubernetes 集群。Draft 在代码部署和提交之间创建内部循环，显著加快了变更验证过程。使用 Draft，开发人员可以准备应用程序 Dockerfiles 和 Helm 图表，并使用两个命令将应用程序部署到远程或本地 Kubernetes 集群。
链接:  https://github.com/azure/draft
成本: 免费

39、Deis Workflow

Deis Workflow 可在 Kubernetes 集群上创建额外的抽象层，这些层允许开发人员在没有特定领域认知的情况下部署或更新 Kubernetes 应用程序。Workflow 基于 Kubernetes 概念构建，提供简单且易于使用的应用程序部署方式。作为一组 Kubernetes 微服务提供，运营商可轻松安装该平台，Workflow 也可在不停机的情况下部署新版本应用。
链接:  https://deis.com/workflow/
成本: 免费

40、Kel

Kel 有助于管理 Kubernetes 应用的整个生命周期。Kel 在 Kubernetes 上提供两个用 Python 和 Go 编写的附加层。级别 0 允许配置 Kubernetes 资源，级别 1 可帮助在 K8s 上部署应用程序。
链接:  http://www.kelproject.com/
成本: 免费

CI/CD
41、Cloud 66

Cloud 66 是一个完整的 DevOps 工具链，用于生产中容器化应用程序，通过专门的 Ops 工具自动化 DevOps 大部分繁重工作。该平台目前在 Kubernetes 上运行 4,000 个客户工作负载，并管理 2,500 行配置。通过提供端到端的基础架构管理，Cloud 66 使工程师能够在任何云或服务器上构建、交付、部署和管理应用程序。
链接:  www.cloud66.com
成本: 14 天免费使用

无服务器相关工具
42、 Kubeless

Kubeless 是一个 Kubernetes 本机无服务器框架，允许部署少量代码，而无需担心底层基础架构。Kubeless 提供自动扩展、API 路由、监控和故障排除等功能。Kubeless 完全依赖 K8s 原语，因此 Kubernetes 用户也可使用原生 K8s API 服务器和 API 网关。
链接:  https://github.com/kubeless/kubeless
成本: 免费

43、Fission

Fission 是 Kubernetes 的快速无服务器框架，可任何地方的 Kubernetes 集群上运行：笔记本电脑、公有云或私有数据中心。用户可使用 Python、NodeJS、Go、C＃或 PHP 编写函数，并使用 Fission 将其部署到 K8s 集群。
链接:  https://fission.io/
成本: 免费

44、Funktion

很长一段时间，Kubernetes 只有一个功能即服务实现就是 Funktion。Funktion 是一个为 Kubernetes 设计的开源事件驱动 lambda 式编程模型。Funktion 与 fabric8 平台紧密结合。使用 Funktion，用户可创建从 200 多个事件源订阅的流来调用功能，包括大多数数据库、消息传递系统、社交媒体及其他中间件和协议。
链接:  https://github.com/funktionio/funktion
成本: 免费

45、IronFunction

IronFunctions 是一个开源无服务器平台或 FaaS 平台，可以在任何地方运行。IronFunction 是在 Golang 上编写的，并且支持任何语言函数。IronFunction 的主要优点是支持 AWS Lambda 格式，直接从 Lambda 导入函数并在任何地方运行。
链接:  https://github.com/iron-io/functions
成本: 免费

46、OpenWhisk

Apache OpenWhisk 是一个由 IBM 和 Adobe 驱动的强大开源 -FaaS 平台。OpenWhisk 可在本地内部部署或云上。Apache OpenWhisk 的设计意味着其充当异步且松散耦合的执行环境，可以针对外部触发器运行功能。OpenWhisk 在 Bluemix 上作为 SaaS 解决方案提供，或者在本地部署基于 Vagrant 的 VM。
链接:  https://console.bluemix.net/openwhisk/
成本: 免费

47、OpenFaaS

OpenFaaS 框架旨在管理 Docker Swarm 或 Kubernetes 上的无服务器功能，可收集和分析各种指标。用户可在函数内打包任何进程并使用，无需重复编码或其他任何操作。FaaS 具有 Prometheus 指标，这意味着可以根据需求自动调整功能。FaaS 本身支持基于 Web 的界面，可在其中试用功能。
链接:  https://github.com/openfaas/faas
成本: 免费

48、Nuclio

Nuclio 是一个无服务器项目，可作为独立库在内部部署设备上启动，也可在 VM 或 Docker 容器内启动。此外，Nuclio 支持开箱即用的 Kubernetes。Nuclio 提供实时数据处理，具有最大并行性和最小开销，可在页面上试用 Nuclio。
链接:  https://github.com/nuclio/nuclio
成本: 免费

49、Virtual-Kubelet

Virtual Kubelet 是一个开源 Kubernetes Kubelet 实现，伪装成一个 kubelet，用于将 Kubernetes 连接到其他 API。Virtual Kubelet 允许节点由其他服务（如 ACI、Hyper.sh 和 AWS 等）提供支持。此连接器具有可插入的体系结构，可直接使用 Kubernetes 原语，使其更容易构建。
链接: https://github.com/virtual-kubelet/virtual-kubelet
成本: 免费

50、Fnproject

Fnproject 是一个容器本机无服务器项目，几乎支持任何语言，可在任何地方运行。Fn 是在 Go 上编写的，因此具有性能和轻量级等优势。Fnproject 支持 AWS Lambda 格式，可以轻松导入 Lambda 函数并使用 Fnproject 启动。
链接:  http://fnproject.io/
成本: 免费

本地服务发现
51、CoreDNS

CoreDNS 是一组用 Go 编写的插件，用于执行 DNS，带有额外 Kubernetes 插件的 CoreDNS 可取代默认 Kube-DNS 服务，并实现为 Kubernetes 基于 DNS 服务定义的规范。CoreDNS 还可侦听通过 UDP/TCP、TLS 和 gRPC 传入的 DNS 请求。
链接:  https://coredns.io/
成本: 免费

原文链接： https://dzone.com/articles/50-useful-kubernetes-tools