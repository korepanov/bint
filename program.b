int fibonacci(int n){
	if (n < 1){
		return -1; 	
	};

	if (n == 1){
		print(str(0));
		print("\n");
		return 0;	
	};

	if (n == 2){
		print(((str(0) + ", ") + str(1)));
		print("\n");
		return 0; 	
	};

	if (n == 3){
		print(((str(0) + ", ") + str(1)));
		print((", " + str(1)));
		print("\n");
		return 0;
	};

	int a1;
	int a2; 
	int a3;

	a1 = 0;
	a2 = 1; 
	
	print((str(a1) + ", "));
	print((str(a2) + ", "));

	for (int i; i = 0; i < (n - 2); i = (i + 1)){
		
		a3 = (a1 + a2); 
		print(str(a3));
		if (i < (n - 3)){
			print(", ");
		};	

		a1 = a2;
		a2 = a3;
	};

	print("\n");

	return 0; 
};



void main(){
	print("Программа вычисления чисел Фибоначчи\n");
	print("Задайте количество чисел\n");
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
		};

		if (n < 0){
			error = "Введено отрицательное число"; 
			print((("Ошибка! " + error) + "\n")); 
			print("Повторите ввод\n");		
		};
	}while(NOT(error==""));

	int res;	
	print("\n");
	print("Результат:\n");
	res = fibonacci(n); 
	
};

main();
