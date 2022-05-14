---
title: mapReduce
layout: post
category: php
author: 夏泽民
---
mapper.php
reducer.php
<!-- more -->

{% raw %}
#!/usr/local/bin/php  
<?php  
$word2count = array();  
// input comes from STDIN (standard input)  
// You can this code :$stdin = fopen(“php://stdin”, “r”);  
while (($line = fgets(STDIN)) !== false) {  
    // remove leading and trailing whitespace and lowercase  
    $line = strtolower(trim($line));  
    // split the line into words while removing any empty string  
    $words = preg_split('/\W/', $line, 0, PREG_SPLIT_NO_EMPTY);  
    // increase counters  
    foreach ($words as $word) {  
        $word2count[$word] += 1;  
    }  
}  
// write the results to STDOUT (standard output)  
// what we output here will be the input for the  
// Reduce step, i.e. the input for reducer.py  
foreach ($word2count as $word => $count) {  
    // tab-delimited  
    echo $word, chr(9), $count, PHP_EOL;  
}  
?>

#!/usr/local/bin/php  
<?php  
$word2count = array();  
// input comes from STDIN  
while (($line = fgets(STDIN)) !== false) {  
    // remove leading and trailing whitespace  
    $line = trim($line);  
    // parse the input we got from mapper.php  
    list($word, $count) = explode(chr(9), $line);  
    // convert count (currently a string) to int  
    $count = intval($count);  
    // sum counts  
    if ($count > 0) $word2count[$word] += $count;  
}  
// sort the words lexigraphically  
//  
// this set is NOT required, we just do it so that our  
// final output will look more like the official Hadoop  
// word count examples  
ksort($word2count);  
// write the results to STDOUT (standard output)  
foreach ($word2count as $word => $count) {  
    echo $word, chr(9), $count, PHP_EOL;  
}  
?>  
  
{% endraw %}

