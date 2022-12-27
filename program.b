void println(string s){
	print(s);
	print("\n");
};

void main(){
	string s;
	s = "Hello";
	println(s[len(s[(len(s) - 2): len(s)])]); 	
};

main();
