int factorial(int n){
	if (0 == n){
		return 1;	
	};
	return (n * factorial(n - 1));
};

void main(){
	int res;
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
};
