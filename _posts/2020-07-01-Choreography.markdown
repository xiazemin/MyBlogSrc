---
title: Choreography choreography
layout: post
category: web
author: 夏泽民
---
SOA中的两个概念：编制（orchestration）和编排（choreography）

<!-- more -->
以下是摘自《Understanding SOA with Web Services》（中文版）关于两个概念的解释：
编制（orchestration）和编排（choreography）是常用于描述“合成Web服务的两种方式”的术语。虽然它们有共同之处，但还是有些区别的。Web服务编制（Web Services Orchestration，WSO）指为业务流程（business processes）而进行Web服务合成，而Web服务编排（Web Services Choreography，WSC）指为业务协作（business collaborations）而进行Web服务合成。

WSO关注于以一种说明性的（declarative）方式（而不是编程的方式）创建合成服务。WSO定义了组成编制（orchestration）的服务，以及这些服务的执行顺序（比如并行活动、条件分支逻辑等）。因此，可以将编制（orchestration）视为一种简单的流程，这种流程自身也是一个Web服务。WSO流通常包括分支控制点、并行处理选择、人类响应步骤以及各种类型的预定义步骤（例如转换、适配器、电子邮件及Web服务等）。

WSC关注于定义多方如何在一个更大的业务事务中进行协作。WSC通过“各方描述自己如何与其他Web服务进行公共消息交换”来定义业务交互，而不是像WSO中那样描述一方是如何执行某个具体业务流程的。
在用WSC来定义业务交互时，需要一个对“业务流程在交互过程中所使用的消息交换协议”的正式描述，对在“有状态的、长期运行的、涉及多方的流程”中的对等的（peer-to-peer）消息交换（同步的或异步的）进行建模。

WSO与WSC的关键区别在于：WSC是一种对等模型（peer-to-peer model），业务流程中会有很多协作方；而WSO是一种层次化的请求者/提供者模型（hierarchical requester/provider model），WSO仅定义了应调用什么服务以及应该何时调用，没有定义多方如何进行协作。
https://blog.csdn.net/villasy/article/details/83839126

