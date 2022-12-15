void main(){
	bool res; 
	string s;

	s = "a";
	res = ((is_letter(s))AND(is_letter(s)));
	string buf;
	buf = str(res);
	print(buf);
	print("\n"); 
};

main();
