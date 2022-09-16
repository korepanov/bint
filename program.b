int y;
stack r;

bool foo(){
	return False;
};

stack func1(){
	stack s;
	if (True){
		return s;
	};
};

stack func0(){
	int t;
	int k;
	
	return func0();  
};

void func2(float g, stack jkk){
	print("Hello world!\n");
	if ((True)AND(False)){
		int foo;
		if(False){
			stack foo;
			float r;
			print("");
		};
	}else if (5 > 4){
		int k;
		print("");	
	}else{
		float t;
		print("");
	};

	if (True){
		print("");
	}else if (False){
		print("");
	};
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
