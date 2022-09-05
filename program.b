void main(){
	for (int i; i = 0; i < 5; i = (i + 1)){
		string si;
		si = str(i);
		si = (("i = " + si) + "\n");
		print(si);
		for (int j; j = 0; j < 10; j = (j + 1)){
			string sj; 
			sj = str(j); 
			if (j > 5){
				break; 			
			};
			sj = (("j = " + sj) + "\n");
			print(sj);		
		};
	};
};

main();
