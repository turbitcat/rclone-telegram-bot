package main

var download_and_unzip_sh = `downloadurl=$1
tempdir=$2
mkdir "$tempdir"
cd "$tempdir"
wget "$downloadurl"
file=` + "`ls`" + `
extname=${file##*.}
if [ "$extname" == "zip" ];
then
echo "unzip"
unzip "$file" && rm "$file"
fi

if [ "$extname" == "7z" ];
then
echo "7z"
7z x "$file" && rm "$file"
fi

if [ "$extname" == "rar" ];
then
echo "unrar"
unrar x "$file" && rm "$file"
fi`
