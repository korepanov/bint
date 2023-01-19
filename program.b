#import "program2.b"
#import "program3.b"

void main(){
	for (int i; i = 0; i < 5; i = (i + 1)){
		print("smth\n"); 	
	};
	int i;

	do{
		print("smth2\n");
		i = (i + 1); 
	}while (i < 5);

	i = 0;
	
	while (i < 5){
		print("smth3\n"); 
		i = (i + 1);	
	};
	string s;
	s = "Hello";
	i = 0;

	for (print(""); print(""); i < len(s); i = (i + 1)){
		print("smth4\n"); 	
	};

	if(True){
		print("YES\n");	
	}else{
		print("NO\n");	
	};
	if (1 < len(str(2))){
		print(s[0:len(s[1:3])]);
	};
	print(((((s[1:3] + s[4]) + "\n") + str(999)) + "\n"));
	print((str(sum(sum(1, 2), diff(3, 2))) + "\n")); 

	if(0 < index(s, "l")){
		print("OK\n"); 	
	}; 

	if(is_letter("l")){
		print("OK2\n"); 	
	};

	if(is_digit("2")){
		print("OK3\n");	
	};

	stack st_var;
	if(True){
		reg_find("e", "Hello").pop(st_var);
	};
	string buf;
	st_var.pop(buf);
	print((buf + "\n"));

	if (exists("program2.b")){
		print("EXISTS\n"); 	
	};

	if (int("1") < 2){
		print("OK4\n"); 	
	};

	if(float("1.5") < 2){
		print("OK5\n"); 	
	};

};

main();
