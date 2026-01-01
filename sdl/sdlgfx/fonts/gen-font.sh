#!/bin/sh

set -e

O=fonts.go

echo -e "package sdlgfx\n" > $O

for i in *.fnt
do
	cat $i | gobin2hex `basename Font$i .fnt` >> $O
	echo -e "\n" >> $O
done

mv $O ../
