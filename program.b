

void main(){
	toPanic = False;
	error = "";	

	string s;
	s = "bla";
	int i;
	i = int(s);

	toPanic = True;
	
	if (NOT(error == "")){
		print("Неверный формат числа\n");	
	};
};

main();
