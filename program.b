int y;
stack r;

int func0(){
	return 1; 
};

void func1(float g, stack jkk){
	print("Hello world!\n");
};

int factorial(int n){
	if (0 == n){
		return 1;	
	};
	return (factorial((n - 1)) * n);
};

void main(){
	int i;
	i = 0;
	
	while (i < 6){
		print("");
		if (True){
			print("");
			print("YES!\n");		
		};
		int j;
		j = 0;
		do{
			print("");
			string si;
			string sj; 

			si = str(i);
			sj = str(j);
			print("i = ");
			print(si);
			print(", j = ");
			print(sj);
			print("\n");
			j = (j + 1); 
		}while(j < 7);
		
		i = (i + 1);
	};

	do{
		print("");
		print("Hello!\n");
	}while(False);
};

main();
