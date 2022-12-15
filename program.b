void main(){
	bool res; 
	string s;

	s = "2";
	string s2;
	s2 = "3";
	res = ((is_digit(s))AND(is_digit(s2)));
	string buf;
	buf = str(res);
	print(buf);
	print("\n"); 
};

main();
