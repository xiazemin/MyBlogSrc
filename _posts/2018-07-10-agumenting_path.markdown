---
title: agumenting_path
layout: post
category: algorithm
author: 夏泽民
---
若P是图G中一条连通两个未匹配顶点的路径，并且属于M的边和不属于M的边(即已匹配和待匹配的边)在P上交替出现，则称P为相对于M的一条增广路径（举例来说，有A、B集合，增广路由A中一个点通向B中一个点，再由B中这个点通向A中一个点……交替进行）。
由增广路的定义可以推出下述五个结论：
1－P的路径长度必定为奇数，第一条边和最后一条边都不属于M。
2－不断寻找增广路可以得到一个更大的匹配M’，直到找不到更多的增广路。
3－M为G的最大匹配当且仅当不存在M的增广路径。
4－最大匹配数M+最大独立数N=总的结点数
5 -- 二分图的最小路径覆盖数 = 原图点数 - 最大匹配数
增广路主要应用于匈牙利算法中，用于求二分图最大匹配。
<!-- more -->
<div class="container">
	<div class="row">
	<img src="{{site.url}}{{site.baseurl}}/img/jupyterSlider.png"/>
	</div>
	<div class="row">
	</div>
</div>
