#!/bin/bash
if read -t 20 -p "please input your blog name:"
then
  echo "hello $REPLY, welcome to come back here"
  prefix=`date '+%Y-%m-%d'`
  name='_posts/'$prefix'-'$REPLY'.markdown'
  echo $name' opened'
  path=`pwd`'/'
  file=$path$name
if [ -f "$file" ]
then
  echo "open existed file"
else
  template=$path'head.markdown'
  sed -E "s/title:.*/title: $REPLY/" $template  > $file 
fi 
 /Applications/MacDown.app/Contents/MacOS/MacDown  $file

else
  echo "sorry , you are too slow "
fi
