---
title: APIExtensionServer
layout: post
category: k8s
author: 夏泽民
---
APIExtensionServer：主要处理 CustomResourceDefinition（CRD）和 CustomResource（CR）的 REST 请求，也是 Delegation 的最后一环，如果对应 CR 不能被处理的话则会返回 404。
Aggregator 和 APIExtensionsServer 对应两种主要扩展 APIServer 资源的方式，即分别是 AA 和 CRD。

https://www.bookstack.cn/read/source-code-reading-notes/kubernetes-kube_apiserver.md

registry 层
实现各种资源对象的存储逻辑

kubernetes/pkg/registry负责k8s内置的资源对象存储逻辑实现
k8s.io/apiextensions-apiserver/pkg/registry负责crd和cr资源对象存储逻辑实现

https://qiankunli.github.io/2019/01/05/kubernetes_source_apiserver.html
<!-- more -->
kube-apiserver作为整个Kubernetes集群操作etcd的唯一入口，负责Kubernetes各资源的认证&鉴权，校验以及CRUD等操作，提供RESTful APIs

kube-apiserver包含三种APIServer：

aggregatorServer：负责处理 apiregistration.k8s.io 组下的APIService资源请求，同时将来自用户的请求拦截转发给aggregated server(AA)
kubeAPIServer：负责对请求的一些通用处理，包括：认证、鉴权以及各个内建资源(pod, deployment，service and etc)的REST服务等
apiExtensionsServer：负责CustomResourceDefinition（CRD）apiResources以及apiVersions的注册，同时处理CRD以及相应CustomResource（CR）的REST请求(如果对应CR不能被处理的话则会返回404)，也是apiserver Delegation的最后一环

https://www.cnblogs.com/tencent-cloud-native/p/14301277.html

https://www.jianshu.com/p/7100880a8858
kube-apiserver-autoregistration
在aggregator启动的时候，即在createAggregatorServer()方法中，crd和apiserver中定义的资源组(GroupVersion)，会通过kube-apiserver-autoregistration poststarthook，被转换成APIService，然后注册进aggregator中，比如将GroupVersion{Group: "apps", Version: "v1"}转成APIService{Spec: v1.APIServiceSpec{Group: "apps", Version" "v1"}}，存进数据库中。apiserver中的对象资源因为是k8s内置的，是固定的，所以只需要在启动的时候，注册一次就可以了，但是CRD中的资源，是用户自定义的，可能随时增删改，所以需要不断的进行更新同步

这个转换过程，主要是通过两个Controller来实现的: crdRegistrationController和autoRegistrationController，这里就体现了Kubernetes中非常核心的设计模式，Controller-Loop模式，即不断从API中获取对象定义，然后按照API对象的定义，执行对应的操作，确保API对象定义和实际的效果是相符的，这种API也叫做declarative api，即申明式API。

autoRegistrationController中定义了一个队列，用来保存添加进来的APIService对象，这些APIService，可能是KubeAPIServer或者APIExtensions APIServer转换过来的，也可能是通过APIService的API直接添加进来的，然后在kube-apiserver-autoregistration PostStartHook中，启动这个Controller，通过不断轮询，将队列中的APIService取出，然后调用apiservice对应的API，将他们添加或者更新到etcd数据库中，固化下来。

crdRegistrationController则是将APIExtensions APIServer中定义的CRD对象转换成APIService，注册到autoRegistrationController的队列中，然后在kube-apiserver-autoregistration PostStartHook中，启动这个Controller，通过不断轮询CRD的API，将CRD中的GroupVersion，转换成APIService，添加到autoRegistrationController的队列中。除了APIExtensions APIServer，还有KubeAPIServer，因为它里面的对象资源是内置的，不会发生变化，所以，在apiServicesToRegister()方法中只进行一次转换，然后注册进autoRegistrationController队列中。

https://hackerain.me/2020/10/08/kubernetes/kube-apiserver-extensions.html

https://insujang.github.io/2020-02-11/kubernetes-custom-resource/

