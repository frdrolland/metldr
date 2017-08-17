#!/usr/bin/gawk
BEGIN {
	FS="|";
	OFS="|";
}

{
	line = "";
	for (i=1; i<=NF; i++) {
		if (1 == i) {
			line = "<datetime> ";
		} else {
			line = line OFS $i;
		}
	}
	print line;
	if (NR > 1) {
		exit 0;
	}
}
