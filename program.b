

void main(){
	
	string s;
	float a;

	do{
		error = "";
		print("Задайте число:\n");
		input(s);
		a = float(s);
		if (NOT(error == "")){
			print("Неверный формат числа\n");		
		}; 	
	}while(NOT(error == ""));

	print(s);
	print("\n");

};

main();
