#import "program2.b"

void main(){
	string s;
	s = "Hello";
	print((("!" + s[1:(2+1)]) + "\n"));
	print((("!" + s[0:(len(s) - 1)]) + "\n")); 
	print(((str(4) + str(5)) + "\n"));
	print(((str(index(s, "e")) + str(index(s, "l"))) + "\n")); 
	
	if (((is_letter("a"))AND(is_digit("2")))){
		print("YES\n");	
	};

	stack res;
	reg_find("l.*", s).pop(res);
	string buf;
	res.pop(buf);
	print(buf);
	print("\n"); 

	if(exists("test.txt")){
		s = s[1:(2+1)];
		print((s + "\n"));
		print("EXISTS\n");	
	}else{
		if(len(s) > 0){
			print("DOES NOT EXIST\n");
		};  
	};

	if (True){
		print((str(1) + "\n")); 
	};

	s = "Hello";

	if (True){
		if (index(s, "e") > 0){
			print("YES\n"); 		
		};
	};

	if (True){
		if (is_letter("e")){
			print("YES\n"); 		
		};
	};

	if (True){
		if (is_digit("1")){
			print("YES\n"); 		
		};
	};

	if (True){
		reg_find("l.*", s).pop(res);
		res.pop(buf);
		print(buf);
		print("\n"); 
	};

	if(True){
		if(exists("test.txt")){
			print("YES\n"); 
		};	
	};
};

main();
