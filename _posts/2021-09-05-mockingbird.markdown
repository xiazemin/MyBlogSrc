---
title: mockingbird
layout: post
category: python
author: 夏泽民
---
https://github.com/babysor/MockingBird

pip3 install -r requirements.txt


MockingBird如何使用
MockingBird的安装要求如下：

首先，MockingBird需要Python 3.7 或更高版本

安装 PyTorch

安装 ffmpeg。

运行pip install -r requirements.txt 来安装剩余的必要包。

安装 webrtcvad 用 pip install webrtcvad-wheels。

接着，你需要使用数据集训练合成器：

下载 数据集并解压：确保您可以访问 train 文件夹中的所有音频文件(如.wav)

使用音频和梅尔频谱图进行预处理：python synthesizer_preprocess_audio.py 可以传入参数 --dataset {dataset} 支持 adatatang_200zh, magicdata, aishell3

预处理嵌入：python synthesizer_preprocess_embeds.py /SV2TTS/synthesizer

训练合成器：python synthesizer_train.py mandarin /SV2TTS/synthesizer

当你在训练文件夹 synthesizer/saved_models/ 中看到注意线显示和损失满足您的需要时，请转到下一步。
<!-- more -->
https://www.easemob.com/news/7090

https://developer.51cto.com/art/202108/680019.htm

https://mp.weixin.qq.com/s/oPnDKf8XdYC0FnfGf4Lz8Q

No module named pathlib

 pip3 install pathlib
 
 https://docs.python.org/zh-cn/3/library/pathlib.html
 python3 demo_toolbox.py
 
 https://stackoverflow.com/questions/60788709/i-get-importerror-no-module-named-pathlib-even-after-installing-pathlib-with-p
 
  pip3 install torch
  
python3 demo_toolbox.py

% python3 demo_toolbox.py
Arguments:
    datasets_root:    None
    enc_models_dir:   encoder/saved_models
    syn_models_dir:   synthesizer/saved_models
    voc_models_dir:   vocoder/saved_models
    cpu:              False
    seed:             None
    no_mp3_support:   False

Warning: you did not pass a root directory for datasets as argument.
The recognized datasets are:
	LibriSpeech/dev-clean
	LibriSpeech/dev-other
	LibriSpeech/test-clean
	LibriSpeech/test-other
	LibriSpeech/train-clean-100
	LibriSpeech/train-clean-360
	LibriSpeech/train-other-500
	LibriTTS/dev-clean
	LibriTTS/dev-other
	LibriTTS/test-clean
	LibriTTS/test-other
	LibriTTS/train-clean-100
	LibriTTS/train-clean-360
	LibriTTS/train-other-500
	LJSpeech-1.1
	VoxCeleb1/wav
	VoxCeleb1/test_wav
	VoxCeleb2/dev/aac
	VoxCeleb2/test/aac
	VCTK-Corpus/wav48
	aidatatang_200zh/corpus/dev
	aidatatang_200zh/corpus/test
	aishell3/test/wav
	magicdata/train
Feel free to add your own. You can still use the toolbox by recording samples yourself.
Traceback (most recent call last):
  File "/Users/xiazemin/MockingBird/demo_toolbox.py", line 43, in <module>
    Toolbox(**vars(args))
  File "/Users/xiazemin/MockingBird/toolbox/__init__.py", line 77, in __init__
    self.reset_ui(enc_models_dir, syn_models_dir, voc_models_dir, seed)
  File "/Users/xiazemin/MockingBird/toolbox/__init__.py", line 145, in reset_ui
    self.ui.populate_models(encoder_models_dir, synthesizer_models_dir, vocoder_models_dir)
  File "/Users/xiazemin/MockingBird/toolbox/ui.py", line 339, in populate_models
    raise Exception("No synthesizer models found in %s" % synthesizer_models_dir)
Exception: No synthesizer models found in synthesizer/saved_models
  
 
 
