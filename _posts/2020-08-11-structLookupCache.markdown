---
title: structLookupCache
layout: post
category: golang
author: 夏泽民
---
// The structLookupCache caches StructOf lookups.
// StructOf does not share the common lookupCache since we need to pin
// the memory associated with *structTypeFixedN.


addToCache := func(t Type) Type {
	var ts []Type
	if ti, ok := structLookupCache.m.Load(hash); ok {
		ts = ti.([]Type)
	}
	structLookupCache.m.Store(hash, append(ts, t))
	return t
}
<!-- more -->

如果通过map来获取 struct field的话，会有个问题，map是无序的，而struct 的filed顺序不一致时不同的类型，hash值不一样，比较相等也如此

因此需要对key进行排序，否则会内存泄漏
