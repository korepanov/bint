void main(){
	string s;
	s = "banana";
	if ("bana" == (s[0:2] + s[(len(s) - 2):len(s)])){
		print("YES\n");	
	}else{
		print("NO\n");
	}; 
	
};
main();
