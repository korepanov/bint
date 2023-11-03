void main(){
	string s;
	s = "something";
	for (int i; i = 0; i < 5; i = (i+1)){
		string s;
		s = "nothing";
		if (i < 3){
			continue;		
		};
		print("Hello world!\n");	
	};

	print((s + "\n"));
};

main();
