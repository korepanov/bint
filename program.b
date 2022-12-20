void main(){
	string s;
	string s2;
	string re;
	stack r;

	re="lo.*";
	s = "Hello world!\n";
	reg_find(re, s).pop(r);
	r.pop(s2);
	print(s2);
	print("\n");
	reg_find(re, s).pop(r);
	r.pop(s2);
	print(s2);
	print("\n");
};

main();
