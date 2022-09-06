void main(){
	int i;
	i = 0;
	
	while (i < 6){
		print("");
		int j;
		j = 0;
		while (j < 4){
			print("");
			string si;
			string sj; 

			si = str(i);
			sj = str(j);
			print("i = ");
			print(si);
			print("; j = ");
			print(sj);
			print("\n");
			j = (j + 1);
		};
		i = (i + 1);
	};
};

main();
