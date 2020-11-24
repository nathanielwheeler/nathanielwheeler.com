cd "$GOPATH/src/nathanielwheeler.com/public/markdown"

read -p "Enter filename " filename
while true
do
	case $filename in
	* )
		break;;
	"" )
		echo "Please enter filename.";;
	esac
done

year=$(date +"%Y") # year (e.g. 2020)
if [[ ! -d $year ]]; then
	mkdir $year
fi

f=$year/$filename.md
touch $f

read -p "Enter title " title
while true
do
	case $filename in
	* )
		break;;
	"" )
		echo "Please enter title.";;
	esac
done

d=$(date +"%B %-d, %Y")

echo "---" >> $f
echo "Title: \"$title\"" >> $f
echo "Date: $d" >> $f
echo "---" >> $f
