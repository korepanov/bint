float find_descr(float a, float b, float c){
	return (((b^2.0)-((4.0*a)*c)));
};

void solve(float a, float b, float c, float d){
	float x1;
	float x2;

	if (d < 0){
		print("Нет решений\n");	
	}else{
		x1 = ((((-1.0)*b)-(d^0.5)) / (2.0*a));	
		x2 = ((((-1.0)*b)+(d^0.5)) / (2.0*a));
		print("Решение:\n"); 		
		print((str(x1) + "\n"));
		print((str(x2) + "\n")); 	
	};  
};

void main(){
	print("Решение квадратных уравнений вида ax^2+bx+c=0\n");	
	string next; 

	do{
		float a; 
		float b;
		float c;
		float d; 

		do{
			print("Введите a:\n");
			string s; 
			input(s); 
			try{
				a = float(s);
			};		
			
			if (NOT(error == "")){
				print("Ошибка! Неверный формат числа\n");
				print((error + "\n")); 
			};
		}while(NOT(error == ""));

		do{
			print("Введите b:\n");
			string s; 
			input(s); 
			try{
				b = float(s);
			};		
			
			if (NOT(error == "")){
				print("Ошибка! Неверный формат числа\n");
				print((error + "\n")); 
			};
		}while(NOT(error == ""));

		do{
			print("Введите c:\n");
			string s; 
			input(s); 
			try{
				c = float(s);
			};		
			
			if (NOT(error == "")){
				print("Ошибка! Неверный формат числа\n");
				print((error + "\n")); 
			};
		}while(NOT(error == ""));

		d = find_descr(a, b, c); 
		
		solve(a, b, c, d);

		print("Решить еще одно уравнение? д/н\n");
		input(next);
	}while (next == "д"); 
};

main();
