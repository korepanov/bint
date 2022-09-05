void main(){
	print("");
	for (int i; i = 0; i < 5; i = (i + 1)){
		print("");
		string si;
		si = str(i);
		si = (("i = " + si) + "\n");
		print(si);
		for (int j; j = 0; j < 10; j = (j + 1)){
			print("");
			string sj; 
			sj = str(j); 
			sj = (("j = " + sj) + "\n");
			print(sj);		
		};
	};
};

main();
