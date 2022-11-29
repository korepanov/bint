bool calc(int a){
	string returntype;
	returntype = "Hello\n";
	print(returntype);
	if (returntype == "Hello\n"){
		print("hello2\n");	
	};
	return (10 == a);
};

void main(){
	bool res;
	res = calc(10);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
