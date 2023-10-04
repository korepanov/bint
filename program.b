void main(){
	string s;
	s = "False";
	bool b;
	b = bool(s); 
	s = "False";
	bool b2;
	b2 = bool(s);

	if ((b)XOR(b2)){
		print("YES\n");
	}else{
		print("NO\n"); 
	};
	
};
main();
