void main(){
	string buf;
	buf = "Hello world!";
	int res;
	res = (index(buf, "world") + index(buf, "lo"));
	buf = str(res);
	print(buf);
	print("\n");
};

main();
