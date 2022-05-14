---
title: sparl_ml_pipline
layout: post
category: spark
author: 夏泽民
---
<div class="container">
<div class="row">
inspired by the scikit-learn project.
</div>
<div class="row">
DataFrame: This ML API uses DataFrame from Spark SQL as an ML dataset, which can hold a variety of data types. E.g., a DataFrame could have different columns storing text, feature vectors, true labels, and predictions.
</div>
<div class="row">
Transformer: A Transformer is an algorithm which can transform one DataFrame into another DataFrame. E.g., an ML model is a Transformer which transforms a DataFrame with features into a DataFrame with predictions.
</div>
<div class="row">
Estimator: An Estimator is an algorithm which can be fit on a DataFrame to produce a Transformer. E.g., a learning algorithm is an Estimator which trains on a DataFrame and produces a model.Technically, an Estimator implements a method fit(), which accepts a DataFrame and produces a Model, which is a Transformer. 
</div>
<div class="row">
Pipeline: A Pipeline chains multiple Transformers and Estimators together to specify an ML workflow.
</div>
<div class="row">
Parameter: All Transformers and Estimators now share a common API for specifying parameters.
<!-- more -->
参考：http://spark.apache.org/docs/latest/ml-pipeline.html
</div>
<div class="row">
一句话概括：管道（Pipeline）是运用数据（DataFrame）训练算法模型（Estimator）调整参数（Parameter）得到一个最优的算法模型（Transformer）,转换数据（DataFrame）的流程。
 For Transformer stages, the transform() method is called on the DataFrame. For Estimator stages, the fit() method is called to produce a Transformer (which becomes part of the PipelineModel, or fitted Pipeline), and that Transformer’s transform() method is called on the DataFrame.
</div>
</div>
