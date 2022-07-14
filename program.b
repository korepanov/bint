void hello(){
	print("Hello!\n");
};

int factorial(int n){
	int res;

	if (0 == n){
		return 1;	
	};
	
	res = (n * factorial((n - 1)));	
	
	return res;
};

void buy(string s){
	print(s);
	print("\n");
};

void main(){
	int res;
	hello();
	res = factorial(5);
	string buf;
	buf = str(res);
	print(buf);
	print("\n");
	buy("buy!");
};

main();
