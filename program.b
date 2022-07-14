int factorial(int n){
	int res;

	if (0 == n){
		return 1;	
	};
	
	res = (n * factorial((n - 1)));	
	
	return res;
};

void main(){
	int res;
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};

main();
