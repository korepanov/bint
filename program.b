int t;

int calc(int a){
	int res;
	string return_type;
	return_type = "Hello";
	print(return_type);
	print("\n");
	res = (a + a);
	return res;
};

void main(){
	int res;
	t = 10; 
	res = calc(t);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
