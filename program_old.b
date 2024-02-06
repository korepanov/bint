// Программа вычисления таблицы простых чисел 

int mod(int a, int b){
	float fa;
	float fb;

	fa = float(a);
	fb = float(b); 

	int c; 

	c = int((fa / fb)); 

	return (a-((a@b)*b));
};

// вычисляет простое число, следующее за n 
int next_prime(int n){
	
	while(True){	
		n = (n + 1);
		bool is_prime; 
		is_prime = True; 

		for (int i; i = 2; i < n; i = (i + 1)){
			if (mod(n, i) == 0){
				is_prime = False;
				break;			
			};		
		};	
		
		if (is_prime){
			return n;		
		};	

	};

	return -1; 
};

int main(){
	print("Программа вычисления таблицы простых чисел\n");
	print("Задайте количество чисел (не более 200)\n");
	string s;
	int n;

	
	do{
		try{
			input(s);
			n = int(s);
		};

		if (NOT(error == "")){
			print("Ошибка! Неправильный формат числа\n");			
			print(error);
			print("\n");	
			print("Повторите ввод\n");	
		}else if (n < 0){
			error = "Введено отрицательное число"; 
			print((("Ошибка! " + error) + "\n")); 
			print("Повторите ввод\n");		
		}else if (n > 200){
			error = "Слишком большое количество чисел";
			print((("Ошибка! " + error) + "\n")); 
			print("Повторите ввод\n");	
		};

	}while(NOT(error==""));

	if (n == 0){
		return 0;	
	};
	
	int p;
	p = 2;
 
	if (n == 1){
		print(str(p));
		print("\n");
		return 0;
	};

	
	//print((str(p) + ", ")); 
	
	for (int i; i = 0; i < (n - 1); i = (i + 1)){
		p = next_prime(p);
		print(str(i));
		print("\n");
		//print(str(p)); 
		if (NOT(i == (n - 2))){
			//print(", ");		
		};
	};
	print("\n");
	
	return 0;
};

int res;
res = main();
