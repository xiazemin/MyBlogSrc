---
title: go 原子操作
layout: post
category: golang
author: 夏泽民
---
<!-- more -->
主要用于数值的操作，由于原子操作可有底层硬件实现，通常比操作系统层的锁机制效率要高。
一共分为5类，LoadT、StoreT、AddT、SwapT和CompareAndSwapT。
func
 LoadT(addr *T) (val T) 
func
 StoreT(addr *T) (val T)
func
 AddT(addr *T, delta T) (new T)
func
 SwapT(addr *T, new T) (old T)
func
 CompareAndSwapT(addr *T, old, new T) (swapped bool)
T代表int32、int64、uint32、uint64、unitptr、pointer。
