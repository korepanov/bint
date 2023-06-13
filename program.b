string print1(string s){
	return s;
};

string print2(string s){
	return s;
};

void println(string s){
	print(s);
	print("\n");
};

void main(){
	println(print2(print1("Привет мир!"))); 
};

main();
