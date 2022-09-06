void main(){
	int i;
	i = 0;
	
	while (i < 6){
		for (int j; j = 0; j < 4; j = (j + 1)){
			string si;
			string sj; 

			si = str(i);
			sj = str(j);
			print("i = ");
			print(si);
			print(", j = ");
			print(sj);
			print("\n");
		};
		break; 
		i = (i + 1);
	};
};

main();
