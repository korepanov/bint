void main(){
	string s;
	s = "so1me2thi3ng";

	for (int i; i = 0; i < len(s); i = (i + 1)){
		if (is_digit(s[i])){
			print(s[i]);		
		};	
	};

	print("\n"); 
};

main();
