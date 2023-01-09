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
		print("EXISTS\n");	
	}else{
		print("DOES NOT EXIST\n"); 
	};
};

main();
