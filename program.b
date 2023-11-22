void main(){
	print("Решение квадратных уравнений вида ax^2+bx+c=0\n");	

	float a;
	float b;
	float c;
	string s;

	do{
		try{
			print("Введите a:\n");
			input(s);
			a = float(s);	
		};
		if (NOT(error == "")){
			print("Неверный формат числа\n");
			print((error + "\n")); 		
		};
	}while(NOT(error == ""));

	do{
		try{
			print("Введите b:\n");
			input(s);
			b = float(s);	
		};
		if (NOT(error == "")){
			print("Неверный формат числа\n");
			print((error + "\n")); 		
		};
	}while(NOT(error == ""));

	do{
		try{
			print("Введите c:\n");
			input(s);
			c = float(s);	
		};
		if (NOT(error == "")){
			print("Неверный формат числа\n");
			print((error + "\n")); 		
		};
	}while(NOT(error == ""));

};

main();
