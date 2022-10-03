int y;
stack r;

int a(){
	return 1;
};

int aba(){
	return 3;
};

int a_ba(){
	return 4; 
};

int b(int a, int c, int k){
	return 2; 
};

int c(){
	return (((((a() + b(1, 2, 3))^a())*b(1, 2, 3)) / aba()) - a_ba());
};

bool foo(){
	return False;
};

stack func1(){
	stack s;
	stack d;
	s = d;
	
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
			si = ((str(55) + str(6.7)) + str("lala"));
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
