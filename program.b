
void main(){
	for (int i; i = 0; i < 5; i = (i + 1)){
		print("smth\n"); 	
	};
	int i;

	do{
		print("smth2\n");
		i = (i + 1); 
	}while (i < 5);

	i = 0;
	
	while (i < 5){
		print("smth3\n"); 
		i = (i + 1);	
	};
	string s;
	s = "Hello";
	i = 0;

	for (print(""); print(""); i < len(s); i = (i + 1)){
		print("smth4\n"); 	
	};

	if(True){
		print("YES\n");	
	}else{
		print("NO\n");	
	};

	print((((s[1:3] + s[4]) + "\n") + str(999))); 
};

main();
