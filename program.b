void main(){
	bool b;
	string s;
	string s2;
	string s3;
	s = "/home/slava/Go";
	s2 = "/bint/";
	s3 = "s.txt";

	b = exists(((s + s2) + "s.txt"));
	print((str(b) + "\n"));
	
	
};
main();
