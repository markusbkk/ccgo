// Problem with conftest.c generated by ./configure by GNU autotools.
// Note: strtoul actually returns unsigned long.

extern int strtoul();

int main() {
	char *term, *string = "0";
	exit(strtoul(string,&term,0) != 0 || term != string+1);
}
